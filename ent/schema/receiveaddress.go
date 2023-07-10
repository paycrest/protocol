package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ReceiveAddress holds the schema definition for the ReceiveAddress entity.
type ReceiveAddress struct {
	ent.Schema
}

// Mixin of the ReceiveAddress.
func (ReceiveAddress) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the ReceiveAddress.
func (ReceiveAddress) Fields() []ent.Field {
	return []ent.Field{
		field.String("address").Unique(),
		field.Int("account_index"),
		field.Enum("status").Values("unused", "partial", "used", "expired").Default("unused"),
		field.Int64("last_indexed_block").Optional(),
		field.Time("last_used").Optional(),
	}
}

// Edges of the ReceiveAddress.
func (ReceiveAddress) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("payment_order", PaymentOrder.Type).
			Ref("receive_address_fk").
			Unique(),
	}
}
