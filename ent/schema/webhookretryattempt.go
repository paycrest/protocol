package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// WebhookRetryAttempt holds the schema definition for the WebhookRetryAttempt entity.
type WebhookRetryAttempt struct {
	ent.Schema
}

// Mixin of the WebhookRetryAttempt.
func (WebhookRetryAttempt) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the WebhookRetryAttempt.
func (WebhookRetryAttempt) Fields() []ent.Field {
	return []ent.Field{
		field.Int("attempt_number"),
		field.Time("next_retry_time").
			Default(time.Now),
		field.JSON("payload", map[string]interface{}{}),
		field.String("signature").
			Optional(),
		field.String("webhook_url"),
		field.Enum("status").
			Values("success", "failed", "expired").
			Default("failed"),
	}
}

// Edges of the WebhookRetryAttempt.
func (WebhookRetryAttempt) Edges() []ent.Edge {
	return nil
}
