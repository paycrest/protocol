package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
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
		field.String("name"),
		field.Enum("scope").
			Values("sender", "provider", "tx_validator"),
		field.String("secret").
			NotEmpty().
			Unique(),
		field.Bool("is_active").
			Default(true),
		field.Time("created_at").
			Immutable().
			Default(time.Now),
	}
}

// Edges of the APIKey.
func (APIKey) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("owner", User.Type).
			Ref("api_keys").
			Required(),
		edge.To("provider_profile", ProviderProfile.Type),
	}
}

// Indexes of the APIKey.
func (APIKey) Indexes() []ent.Index {
	return []ent.Index{
		index.Fields("scope").
			Edges("owner").
			Unique(),
	}
}
