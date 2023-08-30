package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// LockPaymentOrder holds the schema definition for the LockPaymentOrder entity.
type LockPaymentOrder struct {
	ent.Schema
}

// Mixin of the LockPaymentOrder.
func (LockPaymentOrder) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the LockPaymentOrder.
func (LockPaymentOrder) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("order_id").
			Unique(),
		field.Float("amount").
			GoType(decimal.Decimal{}),
		field.Float("amount_paid").
			GoType(decimal.Decimal{}),
		field.Float("rate").
			GoType(decimal.Decimal{}),
		field.String("tx_hash").
			MaxLen(70).
			Optional(),
		field.Enum("status").
			Values("pending", "processing", "cancelled", "fulfilled", "validated", "settled").
			Default("pending"),
		field.Int64("block_number"),
		field.String("institution"),
		field.String("account_identifier"),
		field.String("account_name"),
		field.String("provider_id").
			Optional(),
	}
}

// Edges of the LockPaymentOrder.
func (LockPaymentOrder) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("token", Token.Type).
			Ref("lock_payment_orders").
			Unique().
			Required(),
	}
}
