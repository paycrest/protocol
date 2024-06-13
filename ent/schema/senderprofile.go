package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// SenderProfile holds the schema definition for the SenderProfile entity.
type SenderProfile struct {
	ent.Schema
}

// Fields of the SenderProfile.
func (SenderProfile) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("webhook_url").Optional(),
		field.Strings("domain_whitelist").
			Default([]string{}),
		field.Bool("is_partner").Default(false),
		field.Bool("is_active").
			Default(false),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.JSON("addresses", []struct {
			Token           string          `json:"token"`
			Address         string          `json:"address"`
			Network         string          `json:"network"`
			FeeAddress      string          `json:"feeAddress"`
			RefundAddress   string          `json:"refundAddress"`
			FeePerTokenUnit decimal.Decimal `json:"fee_per_token_unit"`
		}{}),
	}
}

// Edges of the SenderProfile.
func (SenderProfile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Ref("sender_profile").
			Unique().
			Required().
			Immutable(),
		edge.To("api_key", APIKey.Type).
			Unique().
			Annotations(entsql.OnDelete(entsql.Cascade)),
		edge.To("payment_orders", PaymentOrder.Type).
			Annotations(entsql.OnDelete(entsql.SetNull)),
	}
}
