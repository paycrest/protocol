package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ProviderAvailability holds the schema definition for the ProviderAvailability entity.
type ProviderAvailability struct {
	ent.Schema
}

// Fields of the ProviderAvailability.
func (ProviderAvailability) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("cadence").
			Values("24/7", "weekdays", "weekends"),
		field.Time("start_time"),
		field.Time("end_time"),
	}
}

// Edges of the ProviderAvailability.
func (ProviderAvailability) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("provider", ProviderProfile.Type).
			Ref("availability").
			Unique().
			Required(),
	}
}
