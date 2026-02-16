package go_platon

import (
	"net/url"

	"github.com/stremovskyy/go-platon/platon"
)

// ParseWebhookForm parses a Platon callback payload sent as
// application/x-www-form-urlencoded.
func ParseWebhookForm(data []byte) (*platon.WebhookForm, error) {
	return platon.ParseWebhookForm(data)
}

// ParseWebhookValues maps decoded callback form values to WebhookForm.
func ParseWebhookValues(values url.Values) *platon.WebhookForm {
	return platon.ParseWebhookValues(values)
}
