// Code generated by ent, DO NOT EDIT.

package lockpaymentorder

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent/predicate"
	"github.com/shopspring/decimal"
)

// ID filters vertices based on their ID field.
func ID(id uuid.UUID) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id uuid.UUID) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id uuid.UUID) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...uuid.UUID) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...uuid.UUID) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id uuid.UUID) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id uuid.UUID) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id uuid.UUID) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id uuid.UUID) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLTE(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldUpdatedAt, v))
}

// OrderID applies equality check predicate on the "order_id" field. It's identical to OrderIDEQ.
func OrderID(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldOrderID, v))
}

// Amount applies equality check predicate on the "amount" field. It's identical to AmountEQ.
func Amount(v decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldAmount, v))
}

// Rate applies equality check predicate on the "rate" field. It's identical to RateEQ.
func Rate(v decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldRate, v))
}

// TxHash applies equality check predicate on the "tx_hash" field. It's identical to TxHashEQ.
func TxHash(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldTxHash, v))
}

// BlockNumber applies equality check predicate on the "block_number" field. It's identical to BlockNumberEQ.
func BlockNumber(v int64) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldBlockNumber, v))
}

// Institution applies equality check predicate on the "institution" field. It's identical to InstitutionEQ.
func Institution(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldInstitution, v))
}

// AccountIdentifier applies equality check predicate on the "account_identifier" field. It's identical to AccountIdentifierEQ.
func AccountIdentifier(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldAccountIdentifier, v))
}

// AccountName applies equality check predicate on the "account_name" field. It's identical to AccountNameEQ.
func AccountName(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldAccountName, v))
}

// CancellationCount applies equality check predicate on the "cancellation_count" field. It's identical to CancellationCountEQ.
func CancellationCount(v int) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldCancellationCount, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLTE(FieldUpdatedAt, v))
}

// OrderIDEQ applies the EQ predicate on the "order_id" field.
func OrderIDEQ(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldOrderID, v))
}

// OrderIDNEQ applies the NEQ predicate on the "order_id" field.
func OrderIDNEQ(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNEQ(FieldOrderID, v))
}

// OrderIDIn applies the In predicate on the "order_id" field.
func OrderIDIn(vs ...string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldIn(FieldOrderID, vs...))
}

// OrderIDNotIn applies the NotIn predicate on the "order_id" field.
func OrderIDNotIn(vs ...string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNotIn(FieldOrderID, vs...))
}

// OrderIDGT applies the GT predicate on the "order_id" field.
func OrderIDGT(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGT(FieldOrderID, v))
}

// OrderIDGTE applies the GTE predicate on the "order_id" field.
func OrderIDGTE(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGTE(FieldOrderID, v))
}

// OrderIDLT applies the LT predicate on the "order_id" field.
func OrderIDLT(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLT(FieldOrderID, v))
}

// OrderIDLTE applies the LTE predicate on the "order_id" field.
func OrderIDLTE(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLTE(FieldOrderID, v))
}

// OrderIDContains applies the Contains predicate on the "order_id" field.
func OrderIDContains(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldContains(FieldOrderID, v))
}

// OrderIDHasPrefix applies the HasPrefix predicate on the "order_id" field.
func OrderIDHasPrefix(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldHasPrefix(FieldOrderID, v))
}

// OrderIDHasSuffix applies the HasSuffix predicate on the "order_id" field.
func OrderIDHasSuffix(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldHasSuffix(FieldOrderID, v))
}

// OrderIDEqualFold applies the EqualFold predicate on the "order_id" field.
func OrderIDEqualFold(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEqualFold(FieldOrderID, v))
}

// OrderIDContainsFold applies the ContainsFold predicate on the "order_id" field.
func OrderIDContainsFold(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldContainsFold(FieldOrderID, v))
}

// AmountEQ applies the EQ predicate on the "amount" field.
func AmountEQ(v decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldAmount, v))
}

// AmountNEQ applies the NEQ predicate on the "amount" field.
func AmountNEQ(v decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNEQ(FieldAmount, v))
}

// AmountIn applies the In predicate on the "amount" field.
func AmountIn(vs ...decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldIn(FieldAmount, vs...))
}

// AmountNotIn applies the NotIn predicate on the "amount" field.
func AmountNotIn(vs ...decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNotIn(FieldAmount, vs...))
}

// AmountGT applies the GT predicate on the "amount" field.
func AmountGT(v decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGT(FieldAmount, v))
}

// AmountGTE applies the GTE predicate on the "amount" field.
func AmountGTE(v decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGTE(FieldAmount, v))
}

// AmountLT applies the LT predicate on the "amount" field.
func AmountLT(v decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLT(FieldAmount, v))
}

// AmountLTE applies the LTE predicate on the "amount" field.
func AmountLTE(v decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLTE(FieldAmount, v))
}

// RateEQ applies the EQ predicate on the "rate" field.
func RateEQ(v decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldRate, v))
}

// RateNEQ applies the NEQ predicate on the "rate" field.
func RateNEQ(v decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNEQ(FieldRate, v))
}

// RateIn applies the In predicate on the "rate" field.
func RateIn(vs ...decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldIn(FieldRate, vs...))
}

// RateNotIn applies the NotIn predicate on the "rate" field.
func RateNotIn(vs ...decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNotIn(FieldRate, vs...))
}

// RateGT applies the GT predicate on the "rate" field.
func RateGT(v decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGT(FieldRate, v))
}

// RateGTE applies the GTE predicate on the "rate" field.
func RateGTE(v decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGTE(FieldRate, v))
}

// RateLT applies the LT predicate on the "rate" field.
func RateLT(v decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLT(FieldRate, v))
}

// RateLTE applies the LTE predicate on the "rate" field.
func RateLTE(v decimal.Decimal) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLTE(FieldRate, v))
}

// TxHashEQ applies the EQ predicate on the "tx_hash" field.
func TxHashEQ(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldTxHash, v))
}

// TxHashNEQ applies the NEQ predicate on the "tx_hash" field.
func TxHashNEQ(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNEQ(FieldTxHash, v))
}

// TxHashIn applies the In predicate on the "tx_hash" field.
func TxHashIn(vs ...string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldIn(FieldTxHash, vs...))
}

// TxHashNotIn applies the NotIn predicate on the "tx_hash" field.
func TxHashNotIn(vs ...string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNotIn(FieldTxHash, vs...))
}

// TxHashGT applies the GT predicate on the "tx_hash" field.
func TxHashGT(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGT(FieldTxHash, v))
}

// TxHashGTE applies the GTE predicate on the "tx_hash" field.
func TxHashGTE(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGTE(FieldTxHash, v))
}

// TxHashLT applies the LT predicate on the "tx_hash" field.
func TxHashLT(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLT(FieldTxHash, v))
}

// TxHashLTE applies the LTE predicate on the "tx_hash" field.
func TxHashLTE(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLTE(FieldTxHash, v))
}

// TxHashContains applies the Contains predicate on the "tx_hash" field.
func TxHashContains(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldContains(FieldTxHash, v))
}

// TxHashHasPrefix applies the HasPrefix predicate on the "tx_hash" field.
func TxHashHasPrefix(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldHasPrefix(FieldTxHash, v))
}

// TxHashHasSuffix applies the HasSuffix predicate on the "tx_hash" field.
func TxHashHasSuffix(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldHasSuffix(FieldTxHash, v))
}

// TxHashIsNil applies the IsNil predicate on the "tx_hash" field.
func TxHashIsNil() predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldIsNull(FieldTxHash))
}

// TxHashNotNil applies the NotNil predicate on the "tx_hash" field.
func TxHashNotNil() predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNotNull(FieldTxHash))
}

// TxHashEqualFold applies the EqualFold predicate on the "tx_hash" field.
func TxHashEqualFold(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEqualFold(FieldTxHash, v))
}

// TxHashContainsFold applies the ContainsFold predicate on the "tx_hash" field.
func TxHashContainsFold(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldContainsFold(FieldTxHash, v))
}

// StatusEQ applies the EQ predicate on the "status" field.
func StatusEQ(v Status) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldStatus, v))
}

// StatusNEQ applies the NEQ predicate on the "status" field.
func StatusNEQ(v Status) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNEQ(FieldStatus, v))
}

// StatusIn applies the In predicate on the "status" field.
func StatusIn(vs ...Status) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldIn(FieldStatus, vs...))
}

// StatusNotIn applies the NotIn predicate on the "status" field.
func StatusNotIn(vs ...Status) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNotIn(FieldStatus, vs...))
}

// BlockNumberEQ applies the EQ predicate on the "block_number" field.
func BlockNumberEQ(v int64) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldBlockNumber, v))
}

// BlockNumberNEQ applies the NEQ predicate on the "block_number" field.
func BlockNumberNEQ(v int64) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNEQ(FieldBlockNumber, v))
}

// BlockNumberIn applies the In predicate on the "block_number" field.
func BlockNumberIn(vs ...int64) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldIn(FieldBlockNumber, vs...))
}

// BlockNumberNotIn applies the NotIn predicate on the "block_number" field.
func BlockNumberNotIn(vs ...int64) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNotIn(FieldBlockNumber, vs...))
}

// BlockNumberGT applies the GT predicate on the "block_number" field.
func BlockNumberGT(v int64) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGT(FieldBlockNumber, v))
}

// BlockNumberGTE applies the GTE predicate on the "block_number" field.
func BlockNumberGTE(v int64) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGTE(FieldBlockNumber, v))
}

// BlockNumberLT applies the LT predicate on the "block_number" field.
func BlockNumberLT(v int64) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLT(FieldBlockNumber, v))
}

// BlockNumberLTE applies the LTE predicate on the "block_number" field.
func BlockNumberLTE(v int64) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLTE(FieldBlockNumber, v))
}

// InstitutionEQ applies the EQ predicate on the "institution" field.
func InstitutionEQ(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldInstitution, v))
}

// InstitutionNEQ applies the NEQ predicate on the "institution" field.
func InstitutionNEQ(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNEQ(FieldInstitution, v))
}

// InstitutionIn applies the In predicate on the "institution" field.
func InstitutionIn(vs ...string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldIn(FieldInstitution, vs...))
}

// InstitutionNotIn applies the NotIn predicate on the "institution" field.
func InstitutionNotIn(vs ...string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNotIn(FieldInstitution, vs...))
}

// InstitutionGT applies the GT predicate on the "institution" field.
func InstitutionGT(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGT(FieldInstitution, v))
}

// InstitutionGTE applies the GTE predicate on the "institution" field.
func InstitutionGTE(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGTE(FieldInstitution, v))
}

// InstitutionLT applies the LT predicate on the "institution" field.
func InstitutionLT(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLT(FieldInstitution, v))
}

// InstitutionLTE applies the LTE predicate on the "institution" field.
func InstitutionLTE(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLTE(FieldInstitution, v))
}

// InstitutionContains applies the Contains predicate on the "institution" field.
func InstitutionContains(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldContains(FieldInstitution, v))
}

// InstitutionHasPrefix applies the HasPrefix predicate on the "institution" field.
func InstitutionHasPrefix(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldHasPrefix(FieldInstitution, v))
}

// InstitutionHasSuffix applies the HasSuffix predicate on the "institution" field.
func InstitutionHasSuffix(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldHasSuffix(FieldInstitution, v))
}

// InstitutionEqualFold applies the EqualFold predicate on the "institution" field.
func InstitutionEqualFold(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEqualFold(FieldInstitution, v))
}

// InstitutionContainsFold applies the ContainsFold predicate on the "institution" field.
func InstitutionContainsFold(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldContainsFold(FieldInstitution, v))
}

// AccountIdentifierEQ applies the EQ predicate on the "account_identifier" field.
func AccountIdentifierEQ(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldAccountIdentifier, v))
}

// AccountIdentifierNEQ applies the NEQ predicate on the "account_identifier" field.
func AccountIdentifierNEQ(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNEQ(FieldAccountIdentifier, v))
}

// AccountIdentifierIn applies the In predicate on the "account_identifier" field.
func AccountIdentifierIn(vs ...string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldIn(FieldAccountIdentifier, vs...))
}

// AccountIdentifierNotIn applies the NotIn predicate on the "account_identifier" field.
func AccountIdentifierNotIn(vs ...string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNotIn(FieldAccountIdentifier, vs...))
}

// AccountIdentifierGT applies the GT predicate on the "account_identifier" field.
func AccountIdentifierGT(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGT(FieldAccountIdentifier, v))
}

// AccountIdentifierGTE applies the GTE predicate on the "account_identifier" field.
func AccountIdentifierGTE(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGTE(FieldAccountIdentifier, v))
}

// AccountIdentifierLT applies the LT predicate on the "account_identifier" field.
func AccountIdentifierLT(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLT(FieldAccountIdentifier, v))
}

// AccountIdentifierLTE applies the LTE predicate on the "account_identifier" field.
func AccountIdentifierLTE(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLTE(FieldAccountIdentifier, v))
}

// AccountIdentifierContains applies the Contains predicate on the "account_identifier" field.
func AccountIdentifierContains(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldContains(FieldAccountIdentifier, v))
}

// AccountIdentifierHasPrefix applies the HasPrefix predicate on the "account_identifier" field.
func AccountIdentifierHasPrefix(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldHasPrefix(FieldAccountIdentifier, v))
}

// AccountIdentifierHasSuffix applies the HasSuffix predicate on the "account_identifier" field.
func AccountIdentifierHasSuffix(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldHasSuffix(FieldAccountIdentifier, v))
}

// AccountIdentifierEqualFold applies the EqualFold predicate on the "account_identifier" field.
func AccountIdentifierEqualFold(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEqualFold(FieldAccountIdentifier, v))
}

// AccountIdentifierContainsFold applies the ContainsFold predicate on the "account_identifier" field.
func AccountIdentifierContainsFold(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldContainsFold(FieldAccountIdentifier, v))
}

// AccountNameEQ applies the EQ predicate on the "account_name" field.
func AccountNameEQ(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldAccountName, v))
}

// AccountNameNEQ applies the NEQ predicate on the "account_name" field.
func AccountNameNEQ(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNEQ(FieldAccountName, v))
}

// AccountNameIn applies the In predicate on the "account_name" field.
func AccountNameIn(vs ...string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldIn(FieldAccountName, vs...))
}

// AccountNameNotIn applies the NotIn predicate on the "account_name" field.
func AccountNameNotIn(vs ...string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNotIn(FieldAccountName, vs...))
}

// AccountNameGT applies the GT predicate on the "account_name" field.
func AccountNameGT(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGT(FieldAccountName, v))
}

// AccountNameGTE applies the GTE predicate on the "account_name" field.
func AccountNameGTE(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGTE(FieldAccountName, v))
}

// AccountNameLT applies the LT predicate on the "account_name" field.
func AccountNameLT(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLT(FieldAccountName, v))
}

// AccountNameLTE applies the LTE predicate on the "account_name" field.
func AccountNameLTE(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLTE(FieldAccountName, v))
}

// AccountNameContains applies the Contains predicate on the "account_name" field.
func AccountNameContains(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldContains(FieldAccountName, v))
}

// AccountNameHasPrefix applies the HasPrefix predicate on the "account_name" field.
func AccountNameHasPrefix(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldHasPrefix(FieldAccountName, v))
}

// AccountNameHasSuffix applies the HasSuffix predicate on the "account_name" field.
func AccountNameHasSuffix(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldHasSuffix(FieldAccountName, v))
}

// AccountNameEqualFold applies the EqualFold predicate on the "account_name" field.
func AccountNameEqualFold(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEqualFold(FieldAccountName, v))
}

// AccountNameContainsFold applies the ContainsFold predicate on the "account_name" field.
func AccountNameContainsFold(v string) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldContainsFold(FieldAccountName, v))
}

// CancellationCountEQ applies the EQ predicate on the "cancellation_count" field.
func CancellationCountEQ(v int) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldEQ(FieldCancellationCount, v))
}

// CancellationCountNEQ applies the NEQ predicate on the "cancellation_count" field.
func CancellationCountNEQ(v int) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNEQ(FieldCancellationCount, v))
}

// CancellationCountIn applies the In predicate on the "cancellation_count" field.
func CancellationCountIn(vs ...int) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldIn(FieldCancellationCount, vs...))
}

// CancellationCountNotIn applies the NotIn predicate on the "cancellation_count" field.
func CancellationCountNotIn(vs ...int) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldNotIn(FieldCancellationCount, vs...))
}

// CancellationCountGT applies the GT predicate on the "cancellation_count" field.
func CancellationCountGT(v int) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGT(FieldCancellationCount, v))
}

// CancellationCountGTE applies the GTE predicate on the "cancellation_count" field.
func CancellationCountGTE(v int) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldGTE(FieldCancellationCount, v))
}

// CancellationCountLT applies the LT predicate on the "cancellation_count" field.
func CancellationCountLT(v int) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLT(FieldCancellationCount, v))
}

// CancellationCountLTE applies the LTE predicate on the "cancellation_count" field.
func CancellationCountLTE(v int) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(sql.FieldLTE(FieldCancellationCount, v))
}

// HasToken applies the HasEdge predicate on the "token" edge.
func HasToken() predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, TokenTable, TokenColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasTokenWith applies the HasEdge predicate on the "token" edge with a given conditions (other predicates).
func HasTokenWith(preds ...predicate.Token) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(func(s *sql.Selector) {
		step := newTokenStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProvisionBucket applies the HasEdge predicate on the "provision_bucket" edge.
func HasProvisionBucket() predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProvisionBucketTable, ProvisionBucketColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProvisionBucketWith applies the HasEdge predicate on the "provision_bucket" edge with a given conditions (other predicates).
func HasProvisionBucketWith(preds ...predicate.ProvisionBucket) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(func(s *sql.Selector) {
		step := newProvisionBucketStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProvider applies the HasEdge predicate on the "provider" edge.
func HasProvider() predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, ProviderTable, ProviderColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProviderWith applies the HasEdge predicate on the "provider" edge with a given conditions (other predicates).
func HasProviderWith(preds ...predicate.ProviderProfile) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(func(s *sql.Selector) {
		step := newProviderStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasFulfillment applies the HasEdge predicate on the "fulfillment" edge.
func HasFulfillment() predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, FulfillmentTable, FulfillmentColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasFulfillmentWith applies the HasEdge predicate on the "fulfillment" edge with a given conditions (other predicates).
func HasFulfillmentWith(preds ...predicate.LockOrderFulfillment) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(func(s *sql.Selector) {
		step := newFulfillmentStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.LockPaymentOrder) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.LockPaymentOrder) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(func(s *sql.Selector) {
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
func Not(p predicate.LockPaymentOrder) predicate.LockPaymentOrder {
	return predicate.LockPaymentOrder(func(s *sql.Selector) {
		p(s.Not())
	})
}
