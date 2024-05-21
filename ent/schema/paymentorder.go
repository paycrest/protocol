package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
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
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.Float("amount").GoType(decimal.Decimal{}),
		field.Float("amount_paid").GoType(decimal.Decimal{}),
		field.Float("amount_returned").GoType(decimal.Decimal{}),
		field.Float("percent_settled").GoType(decimal.Decimal{}),
		field.Float("sender_fee").GoType(decimal.Decimal{}),
		field.Float("network_fee").GoType(decimal.Decimal{}),
		field.Float("protocol_fee").GoType(decimal.Decimal{}),
		field.Float("rate").GoType(decimal.Decimal{}),
		field.String("tx_hash").
			MaxLen(70).
			Optional(),
		field.Int64("block_number").Default(0),
		field.String("from_address").
			MaxLen(60).
			Optional(),
		field.String("return_address").
			MaxLen(60).
			Optional(),
		field.String("receive_address_text").
			MaxLen(60),
		field.Float("fee_per_token_unit").GoType(decimal.Decimal{}),
		field.String("fee_address").Optional(),
		field.String("gateway_id").Optional(),
		field.Enum("status").
			Values("initiated", "reverted", "pending", "expired", "settled", "refunded").
			Default("initiated"),
	}
}

// Edges of the PaymentOrder.
func (PaymentOrder) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("sender_profile", SenderProfile.Type).
			Ref("payment_orders").
			Unique().
			Required(),
		edge.From("token", Token.Type).
			Ref("payment_orders").
			Unique().
			Required(),
		edge.To("receive_address", ReceiveAddress.Type).
			Unique().
			Annotations(entsql.OnDelete(entsql.SetNull)),
		edge.To("recipients", PaymentOrderRecipient.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
