package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// FiatCurrency holds the schema definition for the FiatCurrency entity.
type FiatCurrency struct {
	ent.Schema
}

// Mixin of the FiatCurrency.
func (FiatCurrency) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the FiatCurrency.
func (FiatCurrency) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("code").Unique(),
		field.String("short_name").Unique(),
		field.Int("decimals").Default(2),
		field.String("symbol"),
		field.String("name"),
	}
}

// Edges of the FiatCurrency.
func (FiatCurrency) Edges() []ent.Edge {
	return nil
}
