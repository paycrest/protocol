package schema

import (
	"context"

	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	gen "github.com/paycrest/paycrest-protocol/ent"
	"github.com/paycrest/paycrest-protocol/ent/hook"
	"golang.org/x/crypto/bcrypt"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Mixin of the User.
func (User) Mixin() []ent.Mixin {
	return []ent.Mixin{
		TimeMixin{},
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("first_name").MaxLen(80),
		field.String("last_name").MaxLen(80),
		field.String("email").
			Unique(),
		field.String("password").Sensitive(),
		field.Bool("is_verified").
			Default(true),
	}
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("api_keys", APIKey.Type),
	}
}

// Indexes of the User.
func (User) Indexes() []ent.Index {
	return nil
}

// Hooks of the User.
func (User) Hooks() []ent.Hook {
	return []ent.Hook{
		hook.On(hashPasswordHook(), ent.OpCreate),
	}
}

// hashPasswordHook is a hook that hashes the password before saving the User entity.
func hashPasswordHook() ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return hook.UserFunc(func(ctx context.Context, m *gen.UserMutation) (ent.Value, error) {
			// Hash the password if it's set in the mutation.
			if password, ok := m.Field("password"); ok {
				hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password.(string)), 14)
				if err != nil {
					return nil, err
				}
				m.SetField("password", string(hashedPassword))
			}
			return next.Mutate(ctx, m)
		})
	}
}
