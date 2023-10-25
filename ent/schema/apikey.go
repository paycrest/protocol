package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// APIKey holds the schema definition for the APIKey entity.
type APIKey struct {
	ent.Schema
}

// Fields of the APIKey.
func (APIKey) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("secret").
			NotEmpty().
			Unique(),
	}
}

// Edges of the APIKey.
func (APIKey) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("sender_profile", SenderProfile.Type).
			Ref("api_key").
			Unique().
			Immutable(),
		edge.From("provider_profile", ProviderProfile.Type).
			Ref("api_key").
			Unique().
			Immutable(),
		edge.From("validator_profile", ValidatorProfile.Type).
			Ref("api_key").
			Unique().
			Immutable(),
		edge.To("payment_orders", PaymentOrder.Type).
			Annotations(entsql.OnDelete(entsql.SetNull)),
	}
}
