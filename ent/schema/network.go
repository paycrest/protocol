package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Network holds the schema definition for the Network entity.
type Network struct {
	ent.Schema
}

// Mixin of the Network.
func (Network) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the Network.
func (Network) Fields() []ent.Field {
	return []ent.Field{
		field.Int64("chain_id"),
		// e.g "bnb-smart-chain", "polygon", "tron", "polygon-mumbai", "tron-shasta"
		field.String("identifier").
			Unique(),
		field.String("rpc_endpoint"),
		field.Bool("is_testnet"),
	}
}

// Edges of the Network.
func (Network) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("tokens", Token.Type).
			Annotations(entsql.OnDelete(entsql.Cascade)),
	}
}
