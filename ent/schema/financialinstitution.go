package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// FinancialInstitution holds the schema definition for the FinancialInstitution entity.
type FinancialInstitution struct {
	ent.Schema
}

// Fields of the FinancialInstitution.
func (FinancialInstitution) Fields() []ent.Field {
	return []ent.Field{
		field.String("code").Unique().Immutable(),
		field.String("name").Optional(),
		field.String("Type").Optional(),
		field.Time("created_at").Default(time.Now).Immutable(),
		field.Time("updated_at").Default(time.Now),
	}
}

// Edges of the FinancialInstitution.
func (FinancialInstitution) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("fiat", FiatCurrency.Type).
			Ref("financialInstitutions").
			Unique(),
	}
}
