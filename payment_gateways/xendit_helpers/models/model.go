package xendit_helpers

import validation "github.com/go-ozzo/ozzo-validation/v4"

const HeaderXCallbackToken = "x-callback-token"
const HeaderWebhookID = "webhook-id"

type CallbackHeaders struct {
	XCallbackToken string `json:"x_callback_token"`
	WebhookID      string `json:"webhook_id"`
}

func (model CallbackHeaders) Validate() error {
	return validation.ValidateStruct(&model,
		validation.Field(&model.XCallbackToken, validation.Required),
		validation.Field(&model.WebhookID, validation.Required),
	)
}
