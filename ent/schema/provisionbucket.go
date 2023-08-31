package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
)

// ProvisionBucket holds the schema definition for the ProvisionBucket entity.
type ProvisionBucket struct {
	ent.Schema
}

// Fields of the ProvisionBucket.
func (ProvisionBucket) Fields() []ent.Field {
	return []ent.Field{
		field.Float("min_amount").
			GoType(decimal.Decimal{}),
		field.Float("max_amount").
			GoType(decimal.Decimal{}),
		field.String("currency").
			MaxLen(3), // https://en.wikipedia.org/wiki/ISO_4217
		field.Time("created_at").
			Immutable().
			Default(time.Now),
	}
}

// Edges of the ProvisionBucket.
func (ProvisionBucket) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("lock_payment_orders", LockPaymentOrder.Type).
			Annotations(entsql.OnDelete(entsql.SetNull)),
		edge.To("provider_profiles", ProviderProfile.Type),
	}
}
