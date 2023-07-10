package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
)

// PaymentOrder holds the schema definition for the PaymentOrder entity.
type PaymentOrder struct {
	ent.Schema
}

// Mixin of the PaymentOrder.
func (PaymentOrder) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the PaymentOrder.
func (PaymentOrder) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("token").
			Values("USDT", "USDC", "BUSD"),
		field.Float("amount").
			GoType(decimal.Decimal{}),
		field.Float("amount_paid").
			GoType(decimal.Decimal{}),
		field.Enum("network").
			Values("bnb-smart-chain", "polygon", "tron", "polygon-mumbai", "tron-shasta"),
		field.String("tx_hash").
			MaxLen(70).
			Optional(),
		field.String("receive_address"),
		field.Enum("status").
			Values("initiated", "pending", "settled", "cancelled", "failed", "refunded").
			Default("initiated"),
		field.Time("last_used").Optional(),
	}
}

// Edges of the PaymentOrder.
func (PaymentOrder) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("api_key", APIKey.Type).
			Ref("payment_orders").
			Unique(),
		edge.To("receive_address_fk", ReceiveAddress.Type).
			Unique(),
		edge.To("recipient", PaymentOrderRecipient.Type).
			Unique(),
	}
}
