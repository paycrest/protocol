// Code generated by ent, DO NOT EDIT.

package providerprofile

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the providerprofile type in the database.
	Label = "provider_profile"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldTradingName holds the string denoting the trading_name field in the database.
	FieldTradingName = "trading_name"
	// FieldHostIdentifier holds the string denoting the host_identifier field in the database.
	FieldHostIdentifier = "host_identifier"
	// FieldProvisionMode holds the string denoting the provision_mode field in the database.
	FieldProvisionMode = "provision_mode"
	// FieldIsPartner holds the string denoting the is_partner field in the database.
	FieldIsPartner = "is_partner"
	// EdgeUser holds the string denoting the user edge name in mutations.
	EdgeUser = "user"
	// EdgeAPIKey holds the string denoting the api_key edge name in mutations.
	EdgeAPIKey = "api_key"
	// EdgeCurrency holds the string denoting the currency edge name in mutations.
	EdgeCurrency = "currency"
	// EdgeProvisionBuckets holds the string denoting the provision_buckets edge name in mutations.
	EdgeProvisionBuckets = "provision_buckets"
	// EdgeOrderTokens holds the string denoting the order_tokens edge name in mutations.
	EdgeOrderTokens = "order_tokens"
	// EdgeAvailability holds the string denoting the availability edge name in mutations.
	EdgeAvailability = "availability"
	// EdgeProviderRating holds the string denoting the provider_rating edge name in mutations.
	EdgeProviderRating = "provider_rating"
	// EdgeAssignedOrders holds the string denoting the assigned_orders edge name in mutations.
	EdgeAssignedOrders = "assigned_orders"
	// Table holds the table name of the providerprofile in the database.
	Table = "provider_profiles"
	// UserTable is the table that holds the user relation/edge.
	UserTable = "provider_profiles"
	// UserInverseTable is the table name for the User entity.
	// It exists in this package in order to avoid circular dependency with the "user" package.
	UserInverseTable = "users"
	// UserColumn is the table column denoting the user relation/edge.
	UserColumn = "user_provider_profile"
	// APIKeyTable is the table that holds the api_key relation/edge.
	APIKeyTable = "api_keys"
	// APIKeyInverseTable is the table name for the APIKey entity.
	// It exists in this package in order to avoid circular dependency with the "apikey" package.
	APIKeyInverseTable = "api_keys"
	// APIKeyColumn is the table column denoting the api_key relation/edge.
	APIKeyColumn = "provider_profile_api_key"
	// CurrencyTable is the table that holds the currency relation/edge.
	CurrencyTable = "provider_profiles"
	// CurrencyInverseTable is the table name for the FiatCurrency entity.
	// It exists in this package in order to avoid circular dependency with the "fiatcurrency" package.
	CurrencyInverseTable = "fiat_currencies"
	// CurrencyColumn is the table column denoting the currency relation/edge.
	CurrencyColumn = "fiat_currency_provider"
	// ProvisionBucketsTable is the table that holds the provision_buckets relation/edge. The primary key declared below.
	ProvisionBucketsTable = "provision_bucket_provider_profiles"
	// ProvisionBucketsInverseTable is the table name for the ProvisionBucket entity.
	// It exists in this package in order to avoid circular dependency with the "provisionbucket" package.
	ProvisionBucketsInverseTable = "provision_buckets"
	// OrderTokensTable is the table that holds the order_tokens relation/edge.
	OrderTokensTable = "provider_order_tokens"
	// OrderTokensInverseTable is the table name for the ProviderOrderToken entity.
	// It exists in this package in order to avoid circular dependency with the "providerordertoken" package.
	OrderTokensInverseTable = "provider_order_tokens"
	// OrderTokensColumn is the table column denoting the order_tokens relation/edge.
	OrderTokensColumn = "provider_profile_order_tokens"
	// AvailabilityTable is the table that holds the availability relation/edge.
	AvailabilityTable = "provider_availabilities"
	// AvailabilityInverseTable is the table name for the ProviderAvailability entity.
	// It exists in this package in order to avoid circular dependency with the "provideravailability" package.
	AvailabilityInverseTable = "provider_availabilities"
	// AvailabilityColumn is the table column denoting the availability relation/edge.
	AvailabilityColumn = "provider_profile_availability"
	// ProviderRatingTable is the table that holds the provider_rating relation/edge.
	ProviderRatingTable = "provider_ratings"
	// ProviderRatingInverseTable is the table name for the ProviderRating entity.
	// It exists in this package in order to avoid circular dependency with the "providerrating" package.
	ProviderRatingInverseTable = "provider_ratings"
	// ProviderRatingColumn is the table column denoting the provider_rating relation/edge.
	ProviderRatingColumn = "provider_profile_provider_rating"
	// AssignedOrdersTable is the table that holds the assigned_orders relation/edge.
	AssignedOrdersTable = "lock_payment_orders"
	// AssignedOrdersInverseTable is the table name for the LockPaymentOrder entity.
	// It exists in this package in order to avoid circular dependency with the "lockpaymentorder" package.
	AssignedOrdersInverseTable = "lock_payment_orders"
	// AssignedOrdersColumn is the table column denoting the assigned_orders relation/edge.
	AssignedOrdersColumn = "provider_profile_assigned_orders"
)

// Columns holds all SQL columns for providerprofile fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldTradingName,
	FieldHostIdentifier,
	FieldProvisionMode,
	FieldIsPartner,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "provider_profiles"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"fiat_currency_provider",
	"user_provider_profile",
}

var (
	// ProvisionBucketsPrimaryKey and ProvisionBucketsColumn2 are the table columns denoting the
	// primary key for the provision_buckets relation (M2M).
	ProvisionBucketsPrimaryKey = []string{"provision_bucket_id", "provider_profile_id"}
)

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
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
	// TradingNameValidator is a validator for the "trading_name" field. It is called by the builders before save.
	TradingNameValidator func(string) error
	// DefaultIsPartner holds the default value on creation for the "is_partner" field.
	DefaultIsPartner bool
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() string
)

// ProvisionMode defines the type for the "provision_mode" enum field.
type ProvisionMode string

// ProvisionModeAuto is the default value of the ProvisionMode enum.
const DefaultProvisionMode = ProvisionModeAuto

// ProvisionMode values.
const (
	ProvisionModeManual ProvisionMode = "manual"
	ProvisionModeAuto   ProvisionMode = "auto"
)

func (pm ProvisionMode) String() string {
	return string(pm)
}

// ProvisionModeValidator is a validator for the "provision_mode" field enum values. It is called by the builders before save.
func ProvisionModeValidator(pm ProvisionMode) error {
	switch pm {
	case ProvisionModeManual, ProvisionModeAuto:
		return nil
	default:
		return fmt.Errorf("providerprofile: invalid enum value for provision_mode field: %q", pm)
	}
}

// OrderOption defines the ordering options for the ProviderProfile queries.
type OrderOption func(*sql.Selector)

// ByID orders the results by the id field.
func ByID(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldID, opts...).ToFunc()
}

// ByCreatedAt orders the results by the created_at field.
func ByCreatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCreatedAt, opts...).ToFunc()
}

// ByUpdatedAt orders the results by the updated_at field.
func ByUpdatedAt(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldUpdatedAt, opts...).ToFunc()
}

// ByTradingName orders the results by the trading_name field.
func ByTradingName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldTradingName, opts...).ToFunc()
}

// ByHostIdentifier orders the results by the host_identifier field.
func ByHostIdentifier(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldHostIdentifier, opts...).ToFunc()
}

// ByProvisionMode orders the results by the provision_mode field.
func ByProvisionMode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldProvisionMode, opts...).ToFunc()
}

// ByIsPartner orders the results by the is_partner field.
func ByIsPartner(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsPartner, opts...).ToFunc()
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

// ByCurrencyField orders the results by currency field.
func ByCurrencyField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newCurrencyStep(), sql.OrderByField(field, opts...))
	}
}

// ByProvisionBucketsCount orders the results by provision_buckets count.
func ByProvisionBucketsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newProvisionBucketsStep(), opts...)
	}
}

// ByProvisionBuckets orders the results by provision_buckets terms.
func ByProvisionBuckets(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProvisionBucketsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByOrderTokensCount orders the results by order_tokens count.
func ByOrderTokensCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newOrderTokensStep(), opts...)
	}
}

// ByOrderTokens orders the results by order_tokens terms.
func ByOrderTokens(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newOrderTokensStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByAvailabilityField orders the results by availability field.
func ByAvailabilityField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newAvailabilityStep(), sql.OrderByField(field, opts...))
	}
}

// ByProviderRatingField orders the results by provider_rating field.
func ByProviderRatingField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProviderRatingStep(), sql.OrderByField(field, opts...))
	}
}

// ByAssignedOrdersCount orders the results by assigned_orders count.
func ByAssignedOrdersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newAssignedOrdersStep(), opts...)
	}
}

// ByAssignedOrders orders the results by assigned_orders terms.
func ByAssignedOrders(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newAssignedOrdersStep(), append([]sql.OrderTerm{term}, terms...)...)
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
func newCurrencyStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(CurrencyInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2O, true, CurrencyTable, CurrencyColumn),
	)
}
func newProvisionBucketsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProvisionBucketsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, true, ProvisionBucketsTable, ProvisionBucketsPrimaryKey...),
	)
}
func newOrderTokensStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(OrderTokensInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, OrderTokensTable, OrderTokensColumn),
	)
}
func newAvailabilityStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(AvailabilityInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2O, false, AvailabilityTable, AvailabilityColumn),
	)
}
func newProviderRatingStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProviderRatingInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2O, false, ProviderRatingTable, ProviderRatingColumn),
	)
}
func newAssignedOrdersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(AssignedOrdersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, AssignedOrdersTable, AssignedOrdersColumn),
	)
}
