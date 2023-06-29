package schema

import (
	"entgo.io/ent"
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
		field.Int("accountIndex"),
		field.Enum("status").Values("unused", "partial", "used", "expired").Default("unused"),
		field.Time("last_used").Optional(),
	}
}

// Edges of the ReceiveAddress.
func (ReceiveAddress) Edges() []ent.Edge {
	return nil
}
