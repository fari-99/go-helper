// Package activitylog emits user-activity audit logs for the Radiance service.
//
// Every audited data mutation is rendered as a single line of JSON on stdout
// (consumed by Fluentd per the "User Activity Logs" PRD). Two entry points are
// provided:
//
//   - LogActivity — a low-coupling helper any controller/service can call
//     directly with pre-computed field maps.
//   - RegisterCallbacks (see callback.go) — wires gorm Create/Update/Delete
//     callbacks so mutations are logged automatically, without every call site
//     remembering to log.
//
// Immutability (PRD NFR-03): these logs are strictly WRITE-ONLY. They are
// emitted to stdout only and are NEVER persisted to an application table, so no
// HTTP controller or other application code path can update or delete a past
// entry. Do not add a database destination for these logs.
package activitylog

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	"radiance/helpers/sentry"
)

// Action names emitted in the "action" field of an ActivityLog. The BULK_*
// variants are used by the gorm callbacks when an operation targets rows via a
// WHERE condition rather than a single primary key (see callback.go).
const (
	ActionCreate     = "CREATE"
	ActionUpdate     = "UPDATE"
	ActionDelete     = "DELETE"
	ActionBulkCreate = "BULK_CREATE"
	ActionBulkUpdate = "BULK_UPDATE"
	ActionBulkDelete = "BULK_DELETE"
)

// ActivityLog is the JSON payload emitted (one line) per audited mutation.
//
// old_data / new_data are intentionally map[string]interface{} to avoid coupling
// the logger to any concrete model. For updates a caller may pass ONLY the
// changed fields (partial maps); this package never diffs objects — it outputs
// exactly what it is given, after stripping sensitive columns.
type ActivityLog struct {
	Timestamp   string                 `json:"timestamp"`
	ActorUserID interface{}            `json:"actor_user_id"`
	Action      string                 `json:"action"`
	Entity      string                 `json:"entity"`
	EntityID    interface{}            `json:"entity_id"`
	OldData     map[string]interface{} `json:"old_data"`
	NewData     map[string]interface{} `json:"new_data"`
}

// sensitiveColumns lists column/field-name fragments whose values must never
// reach the log sink. Matching is case-insensitive and substring-based, so
// "password" also hides "password_hash". Extend as new sensitive fields appear.
var sensitiveColumns = []string{
	"password",
	"salt",
	"token",
	"secret",
	"otp",
	"pin",
	"credential",
	"private_key",
}

// LogActivity builds an ActivityLog from the supplied fields and emits it as a
// single line of JSON to stdout. The timestamp is generated here in UTC and
// formatted as RFC3339 (e.g. 2026-06-23T10:00:00Z) — callers never pass it.
//
// Passing nil (or a partial map) for oldData/newData is safe: nil becomes a JSON
// null, and only the fields provided are emitted. Sensitive columns (see
// sensitiveColumns) are stripped from both maps.
//
// It is defensive by contract: any marshal/write failure — or panic — is
// reported to Sentry (using the optional ctx args, which helpers/sentry accepts
// as iris.Context / context.Context) and swallowed, so activity logging can
// never return an error that would block the caller's business operation.
//
// Usage:
//
//	activitylog.LogActivity(actorUserID, activitylog.ActionUpdate, "users", user.ID,
//	    map[string]interface{}{"email": oldEmail},
//	    map[string]interface{}{"email": newEmail},
//	    ctx)
func LogActivity(actorUserID interface{}, action, entity string, entityID interface{}, oldData, newData map[string]interface{}, ctx ...interface{}) {
	emit(buildEntry(actorUserID, action, entity, entityID, oldData, newData), ctx...)
}

// buildEntry assembles an ActivityLog, stamping the current UTC time (RFC3339)
// and sanitizing the data maps. Split out from LogActivity so it can be unit
// tested without capturing stdout.
func buildEntry(actorUserID interface{}, action, entity string, entityID interface{}, oldData, newData map[string]interface{}) ActivityLog {
	return ActivityLog{
		Timestamp:   time.Now().UTC().Format(time.RFC3339),
		ActorUserID: actorUserID,
		Action:      action,
		Entity:      entity,
		EntityID:    entityID,
		OldData:     RemoveSensitive(oldData),
		NewData:     RemoveSensitive(newData),
	}
}

// emit marshals an ActivityLog to one JSON line and writes it to stdout. On any
// failure (marshal, write, or panic) it falls back to Sentry and returns — it
// never panics and never returns an error.
func emit(entry ActivityLog, ctx ...interface{}) {
	defer func() {
		if r := recover(); r != nil {
			sentry.SentryLogger(fmt.Sprintf("activitylog: recovered from panic while emitting log: %v", r), ctx...)
		}
	}()

	_ = emitStructuredLog(entry.Action, entry)
}

func emitStructuredLog(logName string, data interface{}) error {
	now := time.Now().UTC()

	dataPayload := map[string]interface{}{
		"created_at": now,
		"updated_at": now,
	}

	if payloadData, ok := data.(map[string]interface{}); ok {
		for k, v := range payloadData {
			dataPayload[k] = v
		}

		for k, v := range dataPayload {
			dataPayload[k] = normalizeStructuredLogDataField(v)
		}
	} else {
		dataPayload["payload"] = data
	}

	payload := struct {
		LogName string                 `json:"log_name"`
		Data    map[string]interface{} `json:"data"`
	}{
		LogName: logName,
		Data:    dataPayload,
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	// Emit pure compact JSON so each log event stays on a single line.
	log.Print(string(payloadBytes))
	return nil
}

// normalizeStructuredLogDataField attempts to decode JSON-encoded strings into objects/arrays.
func normalizeStructuredLogDataField(value interface{}) interface{} {
	s, ok := value.(string)
	if !ok {
		return value
	}

	var v interface{}
	if err := json.Unmarshal([]byte(s), &v); err != nil {
		return value
	}

	switch v.(type) {
	case map[string]interface{}, []interface{}:
		return v
	default:
		b, err := json.Marshal(value)
		if err != nil {
			return value
		}

		var nestedMap map[string]interface{}
		if err := json.Unmarshal(b, &nestedMap); err == nil && nestedMap != nil {
			return RemoveSensitive(nestedMap)
		}

		var nestedSlice []interface{}
		if err := json.Unmarshal(b, &nestedSlice); err == nil {
			return sanitizeValue(nestedSlice)
		}

		return value
	}
}

// RemoveSensitive returns a shallow copy of data with any sensitive columns
// (see sensitiveColumns) dropped, so secrets like passwords never reach the log
// sink. It is nil-safe: a nil map returns nil.
func RemoveSensitive(data map[string]interface{}) map[string]interface{} {
	if data == nil {
		return nil
	}
	cleaned := make(map[string]interface{}, len(data))
	for key, value := range data {
		if isSensitive(key) {
			continue
		}
		cleaned[key] = sanitizeValue(value)
	}
	return cleaned
}

func sanitizeValue(value interface{}) interface{} {
	switch v := value.(type) {
	case map[string]interface{}:
		return RemoveSensitive(v)
	case []interface{}:
		cleaned := make([]interface{}, len(v))
		for i, item := range v {
			cleaned[i] = sanitizeValue(item)
		}
		return cleaned
	default:
		return value
	}
}

// isSensitive reports whether a column/field name matches any sensitiveColumns
// fragment (case-insensitive, substring).
func isSensitive(column string) bool {
	lower := strings.ToLower(column)
	for _, sensitive := range sensitiveColumns {
		if strings.Contains(lower, sensitive) {
			return true
		}
	}
	return false
}
