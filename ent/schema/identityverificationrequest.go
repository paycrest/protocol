package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// IdentityVerificationRequest holds the schema definition for the IdentityVerificationRequest entity.
type IdentityVerificationRequest struct {
	ent.Schema
}

// Fields of the IdentityVerificationRequest.
func (IdentityVerificationRequest) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New),
		field.String("wallet_address").
			Unique(),
		field.String("wallet_signature").
			Unique(),
		field.Enum("platform").
			Values("smile_id", "metamap", "sumsub", "synaps"),
		field.String("platform_ref"),
		field.String("verification_url"),
		field.Enum("status").
			Values("pending", "success", "failed").
			Default("pending"),
		field.Bool("fee_reclaimed").Default(false),
		field.Time("updated_at").
			Default(time.Now).
			UpdateDefault(time.Now),
		field.Time("last_url_created_at").
			Default(time.Now),
	}
}

// Edges of the IdentityVerificationRequest.
func (IdentityVerificationRequest) Edges() []ent.Edge {
	return nil
}

// Indexes of the IdentityVerificationRequest.
func (IdentityVerificationRequest) Indexes() []ent.Index {
	return nil
}
