package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// ProviderOrderTokenAddress holds the schema definition for the ProviderOrderTokenAddress entity.
type ProviderOrderTokenAddress struct {
	ent.Schema
}

// Fields of the ProviderOrderTokenAddress.
func (ProviderOrderTokenAddress) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("network").
			Values("bnb-smart-chain", "polygon", "tron", "polygon-mumbai", "tron-shasta"),
		field.String("address").MaxLen(50),
	}
}

// Edges of the ProviderOrderTokenAddress.
func (ProviderOrderTokenAddress) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("providerordertoken", ProviderOrderToken.Type).
			Ref("addresses").
			Unique(),
	}
}
