package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// LockOrderFulfillment holds the schema definition for the LockOrderFulfillment entity.
type LockOrderFulfillment struct {
	ent.Schema
}

// Mixin of the LockOrderFulfillment.
func (LockOrderFulfillment) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the LockOrderFulfillment.
func (LockOrderFulfillment) Fields() []ent.Field {
	return []ent.Field{
		field.String("tx_id"),
		field.String("tx_receipt_image"),
		field.Int("confirmations").
			Default(0),
	}
}

// Edges of the LockOrderFulfillment.
func (LockOrderFulfillment) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("order", LockPaymentOrder.Type).
			Ref("fulfillment").
			Unique().
			Required(),
	}
}
