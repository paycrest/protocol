package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// ValidatorProfile holds the schema definition for the ValidatorProfile entity.
type ValidatorProfile struct {
	ent.Schema
}

// Mixin of the ValidatorProfile.
func (ValidatorProfile) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the ValidatorProfile.
func (ValidatorProfile) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("wallet_address").MaxLen(80),
	}
}

// Edges of the ValidatorProfile.
func (ValidatorProfile) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("api_key", APIKey.Type).
			Ref("validator_profile").
			Unique().
			Required(),
		edge.From("validated_fulfillments", LockOrderFulfillment.Type).
			Ref("validators"),
	}
}
