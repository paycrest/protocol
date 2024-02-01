// Code generated by ent, DO NOT EDIT.

package senderprofile

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the senderprofile type in the database.
	Label = "sender_profile"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldWebhookURL holds the string denoting the webhook_url field in the database.
	FieldWebhookURL = "webhook_url"
	// FieldFeePerTokenUnit holds the string denoting the fee_per_token_unit field in the database.
	FieldFeePerTokenUnit = "fee_per_token_unit"
	// FieldFeeAddress holds the string denoting the fee_address field in the database.
	FieldFeeAddress = "fee_address"
	// FieldRefundAddress holds the string denoting the refund_address field in the database.
	FieldRefundAddress = "refund_address"
	// FieldDomainWhitelist holds the string denoting the domain_whitelist field in the database.
	FieldDomainWhitelist = "domain_whitelist"
	// FieldIsPartner holds the string denoting the is_partner field in the database.
	FieldIsPartner = "is_partner"
	// FieldIsActive holds the string denoting the is_active field in the database.
	FieldIsActive = "is_active"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// EdgeAPIKey holds the string denoting the api_key edge name in mutations.
	EdgeAPIKey = "api_key"
	// EdgePaymentOrders holds the string denoting the payment_orders edge name in mutations.
	EdgePaymentOrders = "payment_orders"
	// Table holds the table name of the senderprofile in the database.
	Table = "sender_profiles"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "sender_profiles"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_sender_profile"
	// APIKeyTable is the table that holds the api_key relation/edge.
	APIKeyTable = "api_keys"
	// APIKeyInverseTable is the table name for the APIKey entity.
	// It exists in this package in order to avoid circular dependency with the "apikey" package.
	APIKeyInverseTable = "api_keys"
	// APIKeyColumn is the table column denoting the api_key relation/edge.
	APIKeyColumn = "sender_profile_api_key"
	// PaymentOrdersTable is the table that holds the payment_orders relation/edge.
	PaymentOrdersTable = "payment_orders"
	// PaymentOrdersInverseTable is the table name for the PaymentOrder entity.
	// It exists in this package in order to avoid circular dependency with the "paymentorder" package.
	PaymentOrdersInverseTable = "payment_orders"
	// PaymentOrdersColumn is the table column denoting the payment_orders relation/edge.
	PaymentOrdersColumn = "sender_profile_payment_orders"
)

// Columns holds all SQL columns for senderprofile fields.
var Columns = []string{
	FieldID,
	FieldWebhookURL,
	FieldFeePerTokenUnit,
	FieldFeeAddress,
	FieldRefundAddress,
	FieldDomainWhitelist,
	FieldIsPartner,
	FieldIsActive,
	FieldUpdatedAt,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "sender_profiles"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"user_sender_profile",
}

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
			return true
		}
	}
	for i := range ForeignKeys {
		if column == ForeignKeys[i] {
			return true
		}
	}
	return false
}

var (
	// DefaultDomainWhitelist holds the default value on creation for the "domain_whitelist" field.
	DefaultDomainWhitelist []string
	// DefaultIsPartner holds the default value on creation for the "is_partner" field.
	DefaultIsPartner bool
	// DefaultIsActive holds the default value on creation for the "is_active" field.
	DefaultIsActive bool
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the SenderProfile queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByWebhookURL orders the results by the webhook_url field.
func ByWebhookURL(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldWebhookURL, opts...).ToFunc()
}

// ByFeePerTokenUnit orders the results by the fee_per_token_unit field.
func ByFeePerTokenUnit(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFeePerTokenUnit, opts...).ToFunc()
}

// ByFeeAddress orders the results by the fee_address field.
func ByFeeAddress(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldFeeAddress, opts...).ToFunc()
}

// ByRefundAddress orders the results by the refund_address field.
func ByRefundAddress(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldRefundAddress, opts...).ToFunc()
}

// ByIsPartner orders the results by the is_partner field.
func ByIsPartner(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsPartner, opts...).ToFunc()
}

// ByIsActive orders the results by the is_active field.
func ByIsActive(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsActive, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByUserField orders the results by user field.
func ByUserField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newUserStep(), sql.OrderByField(field, opts...))
	}
}

// ByAPIKeyField orders the results by api_key field.
func ByAPIKeyField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newAPIKeyStep(), sql.OrderByField(field, opts...))
	}
}

// ByPaymentOrdersCount orders the results by payment_orders count.
func ByPaymentOrdersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newPaymentOrdersStep(), opts...)
	}
}

// ByPaymentOrders orders the results by payment_orders terms.
func ByPaymentOrders(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newPaymentOrdersStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newUserStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(UserInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2O, true, UserTable, UserColumn),
	)
}
func newAPIKeyStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(APIKeyInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2O, false, APIKeyTable, APIKeyColumn),
	)
}
func newPaymentOrdersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(PaymentOrdersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, PaymentOrdersTable, PaymentOrdersColumn),
	)
}
