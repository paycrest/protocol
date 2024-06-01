package schema

import (
	"time"

	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// TransactionLog holds the schema definition for the TransactionLog entity.
type TransactionLog struct {
	ent.Schema
}

// Fields of the TransactionLog.
func (TransactionLog) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).
			Default(uuid.New).Immutable(),
		field.String("sender_id").Optional(),
		field.String("provider_id").Optional(),
		field.String("gateway_id").Optional(),
		field.Enum("status").
			Values("unset", "crypto_deposited", "order_created", "order_settled", "order_refunded", "order_reverted", "gas_prefunded", "gateway_approved").
			Default("unset").Immutable(),
		field.String("network").Optional(),
		field.String("transaction_hash").Optional(),
		field.JSON("metadata", map[string]interface{}{}),
		field.Time("created_at").Default(time.Now).Immutable(),
	}
}

// Edges of the TransactionLog.
func (TransactionLog) Edges() []ent.Edge {
	return nil
}
