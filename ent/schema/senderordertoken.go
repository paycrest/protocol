package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/shopspring/decimal"
)

// SenderOrderToken holds the schema definition for the SenderOrderToken entity.
type SenderOrderToken struct {
	ent.Schema
}

// Fields of the SenderOrderToken.
func (SenderOrderToken) Fields() []ent.Field {
	return []ent.Field{
		field.String("symbol"),
		field.Float("fee_per_token_unit").
			GoType(decimal.Decimal{}),

		// Even though this is what was approved, I would recommend promoting Network from a JSON field to an Actual DB field
		// and making a unique index with both symbol and network that way the JSON field is not an array
		// and direct queries can be made to find address by Network or by token
		field.JSON("addresses", []struct {
			IsDisabled    bool   `json:"isDisabled"` // addition field to disable a network
			FeeAddress    string `json:"feeAddress"`
			RefundAddress string `json:"refundAddress"`
			Network       string `json:"network"`
		}{}),
	}
}

// Edges of the SenderOrderToken.
func (SenderOrderToken) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("sender", SenderProfile.Type).
			Ref("order_tokens").
			Unique(),
	}
}
