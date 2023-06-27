package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// ReceiveAddress holds the schema definition for the ReceiveAddress entity.
type ReceiveAddress struct {
	ent.Schema
}

// Fields of the ReceiveAddress.
func (ReceiveAddress) Fields() []ent.Field {
	return []ent.Field{
		field.String("address").Unique(),
		field.Int("accountIndex"),
		field.Enum("status").Values("active", "inactive").Default("active"),
		// TODO: add timestamps from mixin
		// TODO: add "last_used" datetime
	}
}

// Edges of the ReceiveAddress.
func (ReceiveAddress) Edges() []ent.Edge {
	return nil
}
