package schema

import "entgo.io/ent"

// ProviderOrderTokenAddress holds the schema definition for the ProviderOrderTokenAddress entity.
type ProviderOrderTokenAddress struct {
	ent.Schema
}

// Fields of the ProviderOrderTokenAddress.
func (ProviderOrderTokenAddress) Fields() []ent.Field {
	return []ent.Field{,
		field.Enum("network").
			Values("BNB Smart Chain (BEP20)", "Polygon", "TRON (TRC20)", "Polygon Mumbai", "Tron Shasta"),
		field.String("address").MaxLen(50),
	}
}

// Edges of the ProviderOrderTokenAddress.
func (ProviderOrderTokenAddress) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("providerordertoken", ProviderOrderToken.Type).
			Ref("addresses").
			Required(),
	}
}
