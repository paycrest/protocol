// Code generated by ent, DO NOT EDIT.

package fiatcurrency

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/google/uuid"
)

const (
	// Label holds the string label denoting the fiatcurrency type in the database.
	Label = "fiat_currency"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldCode holds the string denoting the code field in the database.
	FieldCode = "code"
	// FieldShortName holds the string denoting the short_name field in the database.
	FieldShortName = "short_name"
	// FieldDecimals holds the string denoting the decimals field in the database.
	FieldDecimals = "decimals"
	// FieldSymbol holds the string denoting the symbol field in the database.
	FieldSymbol = "symbol"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldMarketRate holds the string denoting the market_rate field in the database.
	FieldMarketRate = "market_rate"
	// FieldIsEnabled holds the string denoting the is_enabled field in the database.
	FieldIsEnabled = "is_enabled"
	// EdgeProviders holds the string denoting the providers edge name in mutations.
	EdgeProviders = "providers"
	// EdgeProvisionBuckets holds the string denoting the provision_buckets edge name in mutations.
	EdgeProvisionBuckets = "provision_buckets"
	// EdgeInstitutions holds the string denoting the institutions edge name in mutations.
	EdgeInstitutions = "institutions"
	// EdgeProviderSettings holds the string denoting the provider_settings edge name in mutations.
	EdgeProviderSettings = "provider_settings"
	// Table holds the table name of the fiatcurrency in the database.
	Table = "fiat_currencies"
	// ProvidersTable is the table that holds the providers relation/edge. The primary key declared below.
	ProvidersTable = "fiat_currency_providers"
	// ProvidersInverseTable is the table name for the ProviderProfile entity.
	// It exists in this package in order to avoid circular dependency with the "providerprofile" package.
	ProvidersInverseTable = "provider_profiles"
	// ProvisionBucketsTable is the table that holds the provision_buckets relation/edge.
	ProvisionBucketsTable = "provision_buckets"
	// ProvisionBucketsInverseTable is the table name for the ProvisionBucket entity.
	// It exists in this package in order to avoid circular dependency with the "provisionbucket" package.
	ProvisionBucketsInverseTable = "provision_buckets"
	// ProvisionBucketsColumn is the table column denoting the provision_buckets relation/edge.
	ProvisionBucketsColumn = "fiat_currency_provision_buckets"
	// InstitutionsTable is the table that holds the institutions relation/edge.
	InstitutionsTable = "institutions"
	// InstitutionsInverseTable is the table name for the Institution entity.
	// It exists in this package in order to avoid circular dependency with the "institution" package.
	InstitutionsInverseTable = "institutions"
	// InstitutionsColumn is the table column denoting the institutions relation/edge.
	InstitutionsColumn = "fiat_currency_institutions"
	// ProviderSettingsTable is the table that holds the provider_settings relation/edge.
	ProviderSettingsTable = "provider_order_tokens"
	// ProviderSettingsInverseTable is the table name for the ProviderOrderToken entity.
	// It exists in this package in order to avoid circular dependency with the "providerordertoken" package.
	ProviderSettingsInverseTable = "provider_order_tokens"
	// ProviderSettingsColumn is the table column denoting the provider_settings relation/edge.
	ProviderSettingsColumn = "fiat_currency_provider_settings"
)

// Columns holds all SQL columns for fiatcurrency fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldCode,
	FieldShortName,
	FieldDecimals,
	FieldSymbol,
	FieldName,
	FieldMarketRate,
	FieldIsEnabled,
}

var (
	// ProvidersPrimaryKey and ProvidersColumn2 are the table columns denoting the
	// primary key for the providers relation (M2M).
	ProvidersPrimaryKey = []string{"fiat_currency_id", "provider_profile_id"}
)

// ValidColumn reports if the column name is valid (part of the table columns).
func ValidColumn(column string) bool {
	for i := range Columns {
		if column == Columns[i] {
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
	// DefaultDecimals holds the default value on creation for the "decimals" field.
	DefaultDecimals int
	// DefaultIsEnabled holds the default value on creation for the "is_enabled" field.
	DefaultIsEnabled bool
	// DefaultID holds the default value on creation for the "id" field.
	DefaultID func() uuid.UUID
)

// OrderOption defines the ordering options for the FiatCurrency queries.
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

// ByCode orders the results by the code field.
func ByCode(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldCode, opts...).ToFunc()
}

// ByShortName orders the results by the short_name field.
func ByShortName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldShortName, opts...).ToFunc()
}

// ByDecimals orders the results by the decimals field.
func ByDecimals(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldDecimals, opts...).ToFunc()
}

// BySymbol orders the results by the symbol field.
func BySymbol(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldSymbol, opts...).ToFunc()
}

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByMarketRate orders the results by the market_rate field.
func ByMarketRate(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldMarketRate, opts...).ToFunc()
}

// ByIsEnabled orders the results by the is_enabled field.
func ByIsEnabled(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldIsEnabled, opts...).ToFunc()
}

// ByProvidersCount orders the results by providers count.
func ByProvidersCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newProvidersStep(), opts...)
	}
}

// ByProviders orders the results by providers terms.
func ByProviders(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProvidersStep(), append([]sql.OrderTerm{term}, terms...)...)
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

// ByInstitutionsCount orders the results by institutions count.
func ByInstitutionsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newInstitutionsStep(), opts...)
	}
}

// ByInstitutions orders the results by institutions terms.
func ByInstitutions(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newInstitutionsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}

// ByProviderSettingsCount orders the results by provider_settings count.
func ByProviderSettingsCount(opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborsCount(s, newProviderSettingsStep(), opts...)
	}
}

// ByProviderSettings orders the results by provider_settings terms.
func ByProviderSettings(term sql.OrderTerm, terms ...sql.OrderTerm) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newProviderSettingsStep(), append([]sql.OrderTerm{term}, terms...)...)
	}
}
func newProvidersStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProvidersInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2M, false, ProvidersTable, ProvidersPrimaryKey...),
	)
}
func newProvisionBucketsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProvisionBucketsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ProvisionBucketsTable, ProvisionBucketsColumn),
	)
}
func newInstitutionsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(InstitutionsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, InstitutionsTable, InstitutionsColumn),
	)
}
func newProviderSettingsStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(ProviderSettingsInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.O2M, false, ProviderSettingsTable, ProviderSettingsColumn),
	)
}
