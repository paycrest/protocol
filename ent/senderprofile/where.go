// Code generated by ent, DO NOT EDIT.

package senderprofile

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/paycrest/protocol/ent/predicate"
	"github.com/shopspring/decimal"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldLTE(FieldID, id))
}

// WebhookURL applies equality check predicate on the "webhook_url" field. It's identical to WebhookURLEQ.
func WebhookURL(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldWebhookURL, v))
}

// FeePerTokenUnit applies equality check predicate on the "fee_per_token_unit" field. It's identical to FeePerTokenUnitEQ.
func FeePerTokenUnit(v decimal.Decimal) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldFeePerTokenUnit, v))
}

// FeeAddress applies equality check predicate on the "fee_address" field. It's identical to FeeAddressEQ.
func FeeAddress(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldFeeAddress, v))
}

// RefundAddress applies equality check predicate on the "refund_address" field. It's identical to RefundAddressEQ.
func RefundAddress(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldRefundAddress, v))
}

// IsPartner applies equality check predicate on the "is_partner" field. It's identical to IsPartnerEQ.
func IsPartner(v bool) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldIsPartner, v))
}

// IsActive applies equality check predicate on the "is_active" field. It's identical to IsActiveEQ.
func IsActive(v bool) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldIsActive, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldUpdatedAt, v))
}

// WebhookURLEQ applies the EQ predicate on the "webhook_url" field.
func WebhookURLEQ(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldWebhookURL, v))
}

// WebhookURLNEQ applies the NEQ predicate on the "webhook_url" field.
func WebhookURLNEQ(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNEQ(FieldWebhookURL, v))
}

// WebhookURLIn applies the In predicate on the "webhook_url" field.
func WebhookURLIn(vs ...string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldIn(FieldWebhookURL, vs...))
}

// WebhookURLNotIn applies the NotIn predicate on the "webhook_url" field.
func WebhookURLNotIn(vs ...string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNotIn(FieldWebhookURL, vs...))
}

// WebhookURLGT applies the GT predicate on the "webhook_url" field.
func WebhookURLGT(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldGT(FieldWebhookURL, v))
}

// WebhookURLGTE applies the GTE predicate on the "webhook_url" field.
func WebhookURLGTE(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldGTE(FieldWebhookURL, v))
}

// WebhookURLLT applies the LT predicate on the "webhook_url" field.
func WebhookURLLT(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldLT(FieldWebhookURL, v))
}

// WebhookURLLTE applies the LTE predicate on the "webhook_url" field.
func WebhookURLLTE(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldLTE(FieldWebhookURL, v))
}

// WebhookURLContains applies the Contains predicate on the "webhook_url" field.
func WebhookURLContains(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldContains(FieldWebhookURL, v))
}

// WebhookURLHasPrefix applies the HasPrefix predicate on the "webhook_url" field.
func WebhookURLHasPrefix(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldHasPrefix(FieldWebhookURL, v))
}

// WebhookURLHasSuffix applies the HasSuffix predicate on the "webhook_url" field.
func WebhookURLHasSuffix(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldHasSuffix(FieldWebhookURL, v))
}

// WebhookURLIsNil applies the IsNil predicate on the "webhook_url" field.
func WebhookURLIsNil() predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldIsNull(FieldWebhookURL))
}

// WebhookURLNotNil applies the NotNil predicate on the "webhook_url" field.
func WebhookURLNotNil() predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNotNull(FieldWebhookURL))
}

// WebhookURLEqualFold applies the EqualFold predicate on the "webhook_url" field.
func WebhookURLEqualFold(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEqualFold(FieldWebhookURL, v))
}

// WebhookURLContainsFold applies the ContainsFold predicate on the "webhook_url" field.
func WebhookURLContainsFold(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldContainsFold(FieldWebhookURL, v))
}

// FeePerTokenUnitEQ applies the EQ predicate on the "fee_per_token_unit" field.
func FeePerTokenUnitEQ(v decimal.Decimal) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldFeePerTokenUnit, v))
}

// FeePerTokenUnitNEQ applies the NEQ predicate on the "fee_per_token_unit" field.
func FeePerTokenUnitNEQ(v decimal.Decimal) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNEQ(FieldFeePerTokenUnit, v))
}

// FeePerTokenUnitIn applies the In predicate on the "fee_per_token_unit" field.
func FeePerTokenUnitIn(vs ...decimal.Decimal) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldIn(FieldFeePerTokenUnit, vs...))
}

// FeePerTokenUnitNotIn applies the NotIn predicate on the "fee_per_token_unit" field.
func FeePerTokenUnitNotIn(vs ...decimal.Decimal) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNotIn(FieldFeePerTokenUnit, vs...))
}

// FeePerTokenUnitGT applies the GT predicate on the "fee_per_token_unit" field.
func FeePerTokenUnitGT(v decimal.Decimal) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldGT(FieldFeePerTokenUnit, v))
}

// FeePerTokenUnitGTE applies the GTE predicate on the "fee_per_token_unit" field.
func FeePerTokenUnitGTE(v decimal.Decimal) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldGTE(FieldFeePerTokenUnit, v))
}

// FeePerTokenUnitLT applies the LT predicate on the "fee_per_token_unit" field.
func FeePerTokenUnitLT(v decimal.Decimal) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldLT(FieldFeePerTokenUnit, v))
}

// FeePerTokenUnitLTE applies the LTE predicate on the "fee_per_token_unit" field.
func FeePerTokenUnitLTE(v decimal.Decimal) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldLTE(FieldFeePerTokenUnit, v))
}

// FeeAddressEQ applies the EQ predicate on the "fee_address" field.
func FeeAddressEQ(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldFeeAddress, v))
}

// FeeAddressNEQ applies the NEQ predicate on the "fee_address" field.
func FeeAddressNEQ(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNEQ(FieldFeeAddress, v))
}

// FeeAddressIn applies the In predicate on the "fee_address" field.
func FeeAddressIn(vs ...string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldIn(FieldFeeAddress, vs...))
}

// FeeAddressNotIn applies the NotIn predicate on the "fee_address" field.
func FeeAddressNotIn(vs ...string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNotIn(FieldFeeAddress, vs...))
}

// FeeAddressGT applies the GT predicate on the "fee_address" field.
func FeeAddressGT(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldGT(FieldFeeAddress, v))
}

// FeeAddressGTE applies the GTE predicate on the "fee_address" field.
func FeeAddressGTE(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldGTE(FieldFeeAddress, v))
}

// FeeAddressLT applies the LT predicate on the "fee_address" field.
func FeeAddressLT(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldLT(FieldFeeAddress, v))
}

// FeeAddressLTE applies the LTE predicate on the "fee_address" field.
func FeeAddressLTE(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldLTE(FieldFeeAddress, v))
}

// FeeAddressContains applies the Contains predicate on the "fee_address" field.
func FeeAddressContains(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldContains(FieldFeeAddress, v))
}

// FeeAddressHasPrefix applies the HasPrefix predicate on the "fee_address" field.
func FeeAddressHasPrefix(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldHasPrefix(FieldFeeAddress, v))
}

// FeeAddressHasSuffix applies the HasSuffix predicate on the "fee_address" field.
func FeeAddressHasSuffix(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldHasSuffix(FieldFeeAddress, v))
}

// FeeAddressIsNil applies the IsNil predicate on the "fee_address" field.
func FeeAddressIsNil() predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldIsNull(FieldFeeAddress))
}

// FeeAddressNotNil applies the NotNil predicate on the "fee_address" field.
func FeeAddressNotNil() predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNotNull(FieldFeeAddress))
}

// FeeAddressEqualFold applies the EqualFold predicate on the "fee_address" field.
func FeeAddressEqualFold(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEqualFold(FieldFeeAddress, v))
}

// FeeAddressContainsFold applies the ContainsFold predicate on the "fee_address" field.
func FeeAddressContainsFold(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldContainsFold(FieldFeeAddress, v))
}

// RefundAddressEQ applies the EQ predicate on the "refund_address" field.
func RefundAddressEQ(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldRefundAddress, v))
}

// RefundAddressNEQ applies the NEQ predicate on the "refund_address" field.
func RefundAddressNEQ(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNEQ(FieldRefundAddress, v))
}

// RefundAddressIn applies the In predicate on the "refund_address" field.
func RefundAddressIn(vs ...string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldIn(FieldRefundAddress, vs...))
}

// RefundAddressNotIn applies the NotIn predicate on the "refund_address" field.
func RefundAddressNotIn(vs ...string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNotIn(FieldRefundAddress, vs...))
}

// RefundAddressGT applies the GT predicate on the "refund_address" field.
func RefundAddressGT(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldGT(FieldRefundAddress, v))
}

// RefundAddressGTE applies the GTE predicate on the "refund_address" field.
func RefundAddressGTE(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldGTE(FieldRefundAddress, v))
}

// RefundAddressLT applies the LT predicate on the "refund_address" field.
func RefundAddressLT(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldLT(FieldRefundAddress, v))
}

// RefundAddressLTE applies the LTE predicate on the "refund_address" field.
func RefundAddressLTE(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldLTE(FieldRefundAddress, v))
}

// RefundAddressContains applies the Contains predicate on the "refund_address" field.
func RefundAddressContains(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldContains(FieldRefundAddress, v))
}

// RefundAddressHasPrefix applies the HasPrefix predicate on the "refund_address" field.
func RefundAddressHasPrefix(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldHasPrefix(FieldRefundAddress, v))
}

// RefundAddressHasSuffix applies the HasSuffix predicate on the "refund_address" field.
func RefundAddressHasSuffix(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldHasSuffix(FieldRefundAddress, v))
}

// RefundAddressIsNil applies the IsNil predicate on the "refund_address" field.
func RefundAddressIsNil() predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldIsNull(FieldRefundAddress))
}

// RefundAddressNotNil applies the NotNil predicate on the "refund_address" field.
func RefundAddressNotNil() predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNotNull(FieldRefundAddress))
}

// RefundAddressEqualFold applies the EqualFold predicate on the "refund_address" field.
func RefundAddressEqualFold(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEqualFold(FieldRefundAddress, v))
}

// RefundAddressContainsFold applies the ContainsFold predicate on the "refund_address" field.
func RefundAddressContainsFold(v string) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldContainsFold(FieldRefundAddress, v))
}

// IsPartnerEQ applies the EQ predicate on the "is_partner" field.
func IsPartnerEQ(v bool) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldIsPartner, v))
}

// IsPartnerNEQ applies the NEQ predicate on the "is_partner" field.
func IsPartnerNEQ(v bool) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNEQ(FieldIsPartner, v))
}

// IsActiveEQ applies the EQ predicate on the "is_active" field.
func IsActiveEQ(v bool) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldIsActive, v))
}

// IsActiveNEQ applies the NEQ predicate on the "is_active" field.
func IsActiveNEQ(v bool) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNEQ(FieldIsActive, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.SenderProfile {
	return predicate.SenderProfile(sql.FieldLTE(FieldUpdatedAt, v))
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.SenderProfile {
	return predicate.SenderProfile(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.SenderProfile {
	return predicate.SenderProfile(func(s *sql.Selector) {
		step := newUserStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasAPIKey applies the HasEdge predicate on the "api_key" edge.
func HasAPIKey() predicate.SenderProfile {
	return predicate.SenderProfile(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, APIKeyTable, APIKeyColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasAPIKeyWith applies the HasEdge predicate on the "api_key" edge with a given conditions (other predicates).
func HasAPIKeyWith(preds ...predicate.APIKey) predicate.SenderProfile {
	return predicate.SenderProfile(func(s *sql.Selector) {
		step := newAPIKeyStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasPaymentOrders applies the HasEdge predicate on the "payment_orders" edge.
func HasPaymentOrders() predicate.SenderProfile {
	return predicate.SenderProfile(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, PaymentOrdersTable, PaymentOrdersColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasPaymentOrdersWith applies the HasEdge predicate on the "payment_orders" edge with a given conditions (other predicates).
func HasPaymentOrdersWith(preds ...predicate.PaymentOrder) predicate.SenderProfile {
	return predicate.SenderProfile(func(s *sql.Selector) {
		step := newPaymentOrdersStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.SenderProfile) predicate.SenderProfile {
	return predicate.SenderProfile(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.SenderProfile) predicate.SenderProfile {
	return predicate.SenderProfile(func(s *sql.Selector) {
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
func Not(p predicate.SenderProfile) predicate.SenderProfile {
	return predicate.SenderProfile(func(s *sql.Selector) {
		p(s.Not())
	})
}
