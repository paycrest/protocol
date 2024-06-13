// Code generated by ent, DO NOT EDIT.

package institution

import (
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
)

const (
	// Label holds the string label denoting the institution type in the database.
	Label = "institution"
	// FieldID holds the string denoting the id field in the database.
	FieldID = "id"
	// FieldCreatedAt holds the string denoting the created_at field in the database.
	FieldCreatedAt = "created_at"
	// FieldUpdatedAt holds the string denoting the updated_at field in the database.
	FieldUpdatedAt = "updated_at"
	// FieldCode holds the string denoting the code field in the database.
	FieldCode = "code"
	// FieldName holds the string denoting the name field in the database.
	FieldName = "name"
	// FieldType holds the string denoting the type field in the database.
	FieldType = "type"
	// EdgeFiatCurrency holds the string denoting the fiat_currency edge name in mutations.
	EdgeFiatCurrency = "fiat_currency"
	// Table holds the table name of the institution in the database.
	Table = "institutions"
	// FiatCurrencyTable is the table that holds the fiat_currency relation/edge.
	FiatCurrencyTable = "institutions"
	// FiatCurrencyInverseTable is the table name for the FiatCurrency entity.
	// It exists in this package in order to avoid circular dependency with the "fiatcurrency" package.
	FiatCurrencyInverseTable = "fiat_currencies"
	// FiatCurrencyColumn is the table column denoting the fiat_currency relation/edge.
	FiatCurrencyColumn = "fiat_currency_institutions"
)

// Columns holds all SQL columns for institution fields.
var Columns = []string{
	FieldID,
	FieldCreatedAt,
	FieldUpdatedAt,
	FieldCode,
	FieldName,
	FieldType,
}

// ForeignKeys holds the SQL foreign-keys that are owned by the "institutions"
// table and are not defined as standalone fields in the schema.
var ForeignKeys = []string{
	"fiat_currency_institutions",
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
	// DefaultCreatedAt holds the default value on creation for the "created_at" field.
	DefaultCreatedAt func() time.Time
	// DefaultUpdatedAt holds the default value on creation for the "updated_at" field.
	DefaultUpdatedAt func() time.Time
	// UpdateDefaultUpdatedAt holds the default value on update for the "updated_at" field.
	UpdateDefaultUpdatedAt func() time.Time
)

// Type defines the type for the "type" enum field.
type Type string

// TypeBank is the default value of the Type enum.
const DefaultType = TypeBank

// Type values.
const (
	TypeBank        Type = "bank"
	TypeMobileMoney Type = "mobile_money"
)

func (_type Type) String() string {
	return string(_type)
}

// TypeValidator is a validator for the "type" field enum values. It is called by the builders before save.
func TypeValidator(_type Type) error {
	switch _type {
	case TypeBank, TypeMobileMoney:
		return nil
	default:
		return fmt.Errorf("institution: invalid enum value for type field: %q", _type)
	}
}

// OrderOption defines the ordering options for the Institution queries.
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

// ByName orders the results by the name field.
func ByName(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldName, opts...).ToFunc()
}

// ByType orders the results by the type field.
func ByType(opts ...sql.OrderTermOption) OrderOption {
	return sql.OrderByField(FieldType, opts...).ToFunc()
}

// ByFiatCurrencyField orders the results by fiat_currency field.
func ByFiatCurrencyField(field string, opts ...sql.OrderTermOption) OrderOption {
	return func(s *sql.Selector) {
		sqlgraph.OrderByNeighborTerms(s, newFiatCurrencyStep(), sql.OrderByField(field, opts...))
	}
}
func newFiatCurrencyStep() *sqlgraph.Step {
	return sqlgraph.NewStep(
		sqlgraph.From(Table, FieldID),
		sqlgraph.To(FiatCurrencyInverseTable, FieldID),
		sqlgraph.Edge(sqlgraph.M2O, true, FiatCurrencyTable, FiatCurrencyColumn),
	)
}
