// Code generated by ent, DO NOT EDIT.

package validatorprofile

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldEQ(FieldUpdatedAt, v))
}

// WalletAddress applies equality check predicate on the "wallet_address" field. It's identical to WalletAddressEQ.
func WalletAddress(v string) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldEQ(FieldWalletAddress, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldLTE(FieldUpdatedAt, v))
}

// WalletAddressEQ applies the EQ predicate on the "wallet_address" field.
func WalletAddressEQ(v string) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldEQ(FieldWalletAddress, v))
}

// WalletAddressNEQ applies the NEQ predicate on the "wallet_address" field.
func WalletAddressNEQ(v string) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldNEQ(FieldWalletAddress, v))
}

// WalletAddressIn applies the In predicate on the "wallet_address" field.
func WalletAddressIn(vs ...string) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldIn(FieldWalletAddress, vs...))
}

// WalletAddressNotIn applies the NotIn predicate on the "wallet_address" field.
func WalletAddressNotIn(vs ...string) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldNotIn(FieldWalletAddress, vs...))
}

// WalletAddressGT applies the GT predicate on the "wallet_address" field.
func WalletAddressGT(v string) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldGT(FieldWalletAddress, v))
}

// WalletAddressGTE applies the GTE predicate on the "wallet_address" field.
func WalletAddressGTE(v string) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldGTE(FieldWalletAddress, v))
}

// WalletAddressLT applies the LT predicate on the "wallet_address" field.
func WalletAddressLT(v string) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldLT(FieldWalletAddress, v))
}

// WalletAddressLTE applies the LTE predicate on the "wallet_address" field.
func WalletAddressLTE(v string) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldLTE(FieldWalletAddress, v))
}

// WalletAddressContains applies the Contains predicate on the "wallet_address" field.
func WalletAddressContains(v string) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldContains(FieldWalletAddress, v))
}

// WalletAddressHasPrefix applies the HasPrefix predicate on the "wallet_address" field.
func WalletAddressHasPrefix(v string) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldHasPrefix(FieldWalletAddress, v))
}

// WalletAddressHasSuffix applies the HasSuffix predicate on the "wallet_address" field.
func WalletAddressHasSuffix(v string) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldHasSuffix(FieldWalletAddress, v))
}

// WalletAddressEqualFold applies the EqualFold predicate on the "wallet_address" field.
func WalletAddressEqualFold(v string) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldEqualFold(FieldWalletAddress, v))
}

// WalletAddressContainsFold applies the ContainsFold predicate on the "wallet_address" field.
func WalletAddressContainsFold(v string) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(sql.FieldContainsFold(FieldWalletAddress, v))
}

// HasAPIKey applies the HasEdge predicate on the "api_key" edge.
func HasAPIKey() predicate.ValidatorProfile {
	return predicate.ValidatorProfile(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, APIKeyTable, APIKeyColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasAPIKeyWith applies the HasEdge predicate on the "api_key" edge with a given conditions (other predicates).
func HasAPIKeyWith(preds ...predicate.APIKey) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(func(s *sql.Selector) {
		step := newAPIKeyStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasValidatedFulfillments applies the HasEdge predicate on the "validated_fulfillments" edge.
func HasValidatedFulfillments() predicate.ValidatorProfile {
	return predicate.ValidatorProfile(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, ValidatedFulfillmentsTable, ValidatedFulfillmentsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasValidatedFulfillmentsWith applies the HasEdge predicate on the "validated_fulfillments" edge with a given conditions (other predicates).
func HasValidatedFulfillmentsWith(preds ...predicate.LockOrderFulfillment) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(func(s *sql.Selector) {
		step := newValidatedFulfillmentsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ValidatorProfile) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ValidatorProfile) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.ValidatorProfile) predicate.ValidatorProfile {
	return predicate.ValidatorProfile(func(s *sql.Selector) {
		p(s.Not())
	})
}
