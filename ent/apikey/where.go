// Code generated by ent, DO NOT EDIT.

package apikey

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.APIKey {
	return predicate.APIKey(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.APIKey {
	return predicate.APIKey(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.APIKey {
	return predicate.APIKey(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.APIKey {
	return predicate.APIKey(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.APIKey {
	return predicate.APIKey(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.APIKey {
	return predicate.APIKey(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.APIKey {
	return predicate.APIKey(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.APIKey {
	return predicate.APIKey(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.APIKey {
	return predicate.APIKey(sql.FieldLTE(FieldID, id))
}

// Secret applies equality check predicate on the "secret" field. It's identical to SecretEQ.
func Secret(v string) predicate.APIKey {
	return predicate.APIKey(sql.FieldEQ(FieldSecret, v))
}

// IsActive applies equality check predicate on the "is_active" field. It's identical to IsActiveEQ.
func IsActive(v bool) predicate.APIKey {
	return predicate.APIKey(sql.FieldEQ(FieldIsActive, v))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.APIKey {
	return predicate.APIKey(sql.FieldEQ(FieldCreatedAt, v))
}

// SecretEQ applies the EQ predicate on the "secret" field.
func SecretEQ(v string) predicate.APIKey {
	return predicate.APIKey(sql.FieldEQ(FieldSecret, v))
}

// SecretNEQ applies the NEQ predicate on the "secret" field.
func SecretNEQ(v string) predicate.APIKey {
	return predicate.APIKey(sql.FieldNEQ(FieldSecret, v))
}

// SecretIn applies the In predicate on the "secret" field.
func SecretIn(vs ...string) predicate.APIKey {
	return predicate.APIKey(sql.FieldIn(FieldSecret, vs...))
}

// SecretNotIn applies the NotIn predicate on the "secret" field.
func SecretNotIn(vs ...string) predicate.APIKey {
	return predicate.APIKey(sql.FieldNotIn(FieldSecret, vs...))
}

// SecretGT applies the GT predicate on the "secret" field.
func SecretGT(v string) predicate.APIKey {
	return predicate.APIKey(sql.FieldGT(FieldSecret, v))
}

// SecretGTE applies the GTE predicate on the "secret" field.
func SecretGTE(v string) predicate.APIKey {
	return predicate.APIKey(sql.FieldGTE(FieldSecret, v))
}

// SecretLT applies the LT predicate on the "secret" field.
func SecretLT(v string) predicate.APIKey {
	return predicate.APIKey(sql.FieldLT(FieldSecret, v))
}

// SecretLTE applies the LTE predicate on the "secret" field.
func SecretLTE(v string) predicate.APIKey {
	return predicate.APIKey(sql.FieldLTE(FieldSecret, v))
}

// SecretContains applies the Contains predicate on the "secret" field.
func SecretContains(v string) predicate.APIKey {
	return predicate.APIKey(sql.FieldContains(FieldSecret, v))
}

// SecretHasPrefix applies the HasPrefix predicate on the "secret" field.
func SecretHasPrefix(v string) predicate.APIKey {
	return predicate.APIKey(sql.FieldHasPrefix(FieldSecret, v))
}

// SecretHasSuffix applies the HasSuffix predicate on the "secret" field.
func SecretHasSuffix(v string) predicate.APIKey {
	return predicate.APIKey(sql.FieldHasSuffix(FieldSecret, v))
}

// SecretEqualFold applies the EqualFold predicate on the "secret" field.
func SecretEqualFold(v string) predicate.APIKey {
	return predicate.APIKey(sql.FieldEqualFold(FieldSecret, v))
}

// SecretContainsFold applies the ContainsFold predicate on the "secret" field.
func SecretContainsFold(v string) predicate.APIKey {
	return predicate.APIKey(sql.FieldContainsFold(FieldSecret, v))
}

// IsActiveEQ applies the EQ predicate on the "is_active" field.
func IsActiveEQ(v bool) predicate.APIKey {
	return predicate.APIKey(sql.FieldEQ(FieldIsActive, v))
}

// IsActiveNEQ applies the NEQ predicate on the "is_active" field.
func IsActiveNEQ(v bool) predicate.APIKey {
	return predicate.APIKey(sql.FieldNEQ(FieldIsActive, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.APIKey {
	return predicate.APIKey(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.APIKey {
	return predicate.APIKey(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.APIKey {
	return predicate.APIKey(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.APIKey {
	return predicate.APIKey(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.APIKey {
	return predicate.APIKey(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.APIKey {
	return predicate.APIKey(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.APIKey {
	return predicate.APIKey(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.APIKey {
	return predicate.APIKey(sql.FieldLTE(FieldCreatedAt, v))
}

// HasSenderProfile applies the HasEdge predicate on the "sender_profile" edge.
func HasSenderProfile() predicate.APIKey {
	return predicate.APIKey(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, SenderProfileTable, SenderProfileColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasSenderProfileWith applies the HasEdge predicate on the "sender_profile" edge with a given conditions (other predicates).
func HasSenderProfileWith(preds ...predicate.SenderProfile) predicate.APIKey {
	return predicate.APIKey(func(s *sql.Selector) {
		step := newSenderProfileStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProviderProfile applies the HasEdge predicate on the "provider_profile" edge.
func HasProviderProfile() predicate.APIKey {
	return predicate.APIKey(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, ProviderProfileTable, ProviderProfileColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProviderProfileWith applies the HasEdge predicate on the "provider_profile" edge with a given conditions (other predicates).
func HasProviderProfileWith(preds ...predicate.ProviderProfile) predicate.APIKey {
	return predicate.APIKey(func(s *sql.Selector) {
		step := newProviderProfileStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasValidatorProfile applies the HasEdge predicate on the "validator_profile" edge.
func HasValidatorProfile() predicate.APIKey {
	return predicate.APIKey(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, ValidatorProfileTable, ValidatorProfileColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasValidatorProfileWith applies the HasEdge predicate on the "validator_profile" edge with a given conditions (other predicates).
func HasValidatorProfileWith(preds ...predicate.ValidatorProfile) predicate.APIKey {
	return predicate.APIKey(func(s *sql.Selector) {
		step := newValidatorProfileStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPaymentOrders applies the HasEdge predicate on the "payment_orders" edge.
func HasPaymentOrders() predicate.APIKey {
	return predicate.APIKey(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PaymentOrdersTable, PaymentOrdersColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPaymentOrdersWith applies the HasEdge predicate on the "payment_orders" edge with a given conditions (other predicates).
func HasPaymentOrdersWith(preds ...predicate.PaymentOrder) predicate.APIKey {
	return predicate.APIKey(func(s *sql.Selector) {
		step := newPaymentOrdersStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.APIKey) predicate.APIKey {
	return predicate.APIKey(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.APIKey) predicate.APIKey {
	return predicate.APIKey(func(s *sql.Selector) {
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
func Not(p predicate.APIKey) predicate.APIKey {
	return predicate.APIKey(func(s *sql.Selector) {
		p(s.Not())
	})
}
