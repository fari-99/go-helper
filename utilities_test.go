package gohelper

import (
	"encoding/json"
	"testing"

	"golang.org/x/exp/slices"
)

type TestStruct struct {
	ID    int
	Value string
}

func TestInArray(t *testing.T) {
	t.Run("testing Int", testInt)
	t.Run("testing Float", testFloat)
	t.Run("testing Struct", testStruct)
}

func testInt(t *testing.T) {
	// simple int exists
	var haystackIntArray = []any{int(0), int(1), int(2), int(3), int(4), int(5)}
	var needleIntArray int = 3

	isExistsIntArray, _, _ := InArray(needleIntArray, haystackIntArray)
	if !isExistsIntArray {
		t.Log("can't find that int value")
		t.FailNow()
	}

	// simple int not exists
	haystackIntArray = []any{int(0), int(1), int(2), int(3), int(4), int(5)}
	needleIntArray = 6

	isExistsIntArray, _, _ = InArray(needleIntArray, haystackIntArray)
	if isExistsIntArray {
		t.Log("shouldn't have find that int value")
		t.FailNow()
	}

	// simple int64 exists
	var haystackInt64Array = []any{int64(0), int64(1), int64(2), int64(3), int64(4), int64(5)}
	var needleInt64Array int64 = 3

	isExistsInt64Array, _, _ := InArray(needleInt64Array, haystackInt64Array)
	if !isExistsInt64Array {
		t.Log("can't find that int64 value")
		t.FailNow()
	}

	// simple int int64 exists
	haystackInt64Array = []any{int64(0), int64(1), int64(2), int64(3), int64(4), int64(5)}
	needleInt64Array = 6

	isExistsInt64Array, _, _ = InArray(needleInt64Array, haystackInt64Array)
	if isExistsInt64Array {
		t.Log("shouldn't have find that int64 value")
		t.FailNow()
	}
}

func testFloat(t *testing.T) {
	// simple float32 exists
	var haystackFloat32Array = []any{float32(0), float32(1), float32(2), float32(3), float32(4), float32(5)}
	var needleFloat32Array float32 = 3

	isExistsFloat32Array, _, _ := InArray(needleFloat32Array, haystackFloat32Array)
	if !isExistsFloat32Array {
		t.Log("can't find that Float32 value")
		t.FailNow()
	}

	// simple float32 not exists
	haystackFloat32Array = []any{float32(0), float32(1), float32(2), float32(3), float32(4), float32(5)}
	needleFloat32Array = float32(6)

	isExistsFloat32Array, _, _ = InArray(needleFloat32Array, haystackFloat32Array)
	if isExistsFloat32Array {
		t.Log("shouldn't have find that Float32 value")
		t.FailNow()
	}

	// simple float64 exists
	var haystackFloat64Array = []any{float64(0), float64(1), float64(2), float64(3), float64(4), float64(5)}
	var needleFloat64Array float64 = float64(3)

	isExistsFloat64Array, _, _ := InArray(needleFloat64Array, haystackFloat64Array)
	if !isExistsFloat64Array {
		t.Log("can't find that Float64 value")
		t.FailNow()
	}

	// simple float64 exists
	haystackFloat64Array = []any{float64(0), float64(1), float64(2), float64(3), float64(4), float64(5)}
	needleFloat64Array = float64(6)

	isExistsFloat64Array, _, _ = InArray(needleFloat64Array, haystackFloat64Array)
	if isExistsFloat64Array {
		t.Log("shouldn't have find that Float64 value")
		t.FailNow()
	}
}

func testStruct(t *testing.T) {
	var haystackStruct = []any{
		TestStruct{ID: 0, Value: "0"},
		TestStruct{ID: 1, Value: "1"},
		TestStruct{ID: 2, Value: "2"},
		TestStruct{ID: 3, Value: "3"},
		TestStruct{ID: 4, Value: "4"},
		TestStruct{ID: 5, Value: "5"},
	}

	needleExistsFunction := func(data any) bool {
		var testStruct TestStruct
		dataMarshal, _ := json.Marshal(data)
		_ = json.Unmarshal(dataMarshal, &testStruct)

		return testStruct.ID == 3
	}

	idxExists := slices.IndexFunc(haystackStruct, needleExistsFunction)
	isExistsStruct, idxExistsStruct, _ := InArray(needleExistsFunction, haystackStruct)
	if !isExistsStruct {
		t.Log("can't find that value")
		t.Fail()
	}

	if idxExistsStruct != idxExists {
		t.Log("index not the same")
		t.Fail()
	}

	needleNotExistsFunction := func(data any) bool {
		var testStruct TestStruct
		dataMarshal, _ := json.Marshal(data)
		_ = json.Unmarshal(dataMarshal, &testStruct)

		return testStruct.ID == 6
	}

	idxNotExists := slices.IndexFunc(haystackStruct, needleNotExistsFunction)
	isNotExistsStruct, idxNotExistsStruct, _ := InArray(needleNotExistsFunction, haystackStruct)
	if isNotExistsStruct {
		t.Log("shouldn't find the value in the struct")
		t.Fail()
	}

	if idxNotExistsStruct != -1 && idxNotExists != -1 {
		t.Log("idx shouldn't be find")
		t.Fail()
	}
}
