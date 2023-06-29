package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/shopspring/decimal"
)

// ProviderOrderToken holds the schema definition for the ProviderOrderToken entity.
type ProviderOrderToken struct {
	ent.Schema
}

// Mixin of the ProviderOrderToken.
func (ProviderOrderToken) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the ProviderOrderToken.
func (ProviderOrderToken) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("name").
			Values("USDT", "USDC", "BUSD"),
		field.Float("fixed_conversion_rate").
			GoType(decimal.Decimal{}),
		field.Float("floating_conversion_rate").
			GoType(decimal.Decimal{}),
		field.Enum("conversion_rate_type").
			Values("fixed", "floating"),
		field.String("max_order_amount").
			GoType(decimal.Decimal{}),
		field.String("min_order_amount").
			GoType(decimal.Decimal{}),
	}
}

// Edges of the ProviderOrderToken.
func (ProviderOrderToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("provider", ProviderProfile.Type).
			Ref("order_tokens").
			Unique(),
		edge.To("addresses", ProviderOrderTokenAddress.Type),
	}
}

// Indexes of the ProviderOrderToken.
func (ProviderOrderToken) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("name").
			Edges("provider").
			Unique(),
	}
}
