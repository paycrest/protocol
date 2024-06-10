// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/paycrest/protocol/ent/fiatcurrency"
	"github.com/shopspring/decimal"
)

// FiatCurrency is the model entity for the FiatCurrency schema.
type FiatCurrency struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Code holds the value of the "code" field.
	Code string `json:"code,omitempty"`
	// ShortName holds the value of the "short_name" field.
	ShortName string `json:"short_name,omitempty"`
	// Decimals holds the value of the "decimals" field.
	Decimals int `json:"decimals,omitempty"`
	// Symbol holds the value of the "symbol" field.
	Symbol string `json:"symbol,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// MarketRate holds the value of the "market_rate" field.
	MarketRate decimal.Decimal `json:"market_rate,omitempty"`
	// IsEnabled holds the value of the "is_enabled" field.
	IsEnabled bool `json:"is_enabled,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the FiatCurrencyQuery when eager-loading is set.
	Edges        FiatCurrencyEdges `json:"edges"`
	selectValues sql.SelectValues
}

// FiatCurrencyEdges holds the relations/edges for other nodes in the graph.
type FiatCurrencyEdges struct {
	// Providers holds the value of the providers edge.
	Providers []*ProviderProfile `json:"providers,omitempty"`
	// ProvisionBuckets holds the value of the provision_buckets edge.
	ProvisionBuckets []*ProvisionBucket `json:"provision_buckets,omitempty"`
	// FinancialInstitutions holds the value of the financialInstitutions edge.
	FinancialInstitutions []*FinancialInstitution `json:"financialInstitutions,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [3]bool
}

// ProvidersOrErr returns the Providers value or an error if the edge
// was not loaded in eager-loading.
func (e FiatCurrencyEdges) ProvidersOrErr() ([]*ProviderProfile, error) {
	if e.loadedTypes[0] {
		return e.Providers, nil
	}
	return nil, &NotLoadedError{edge: "providers"}
}

// ProvisionBucketsOrErr returns the ProvisionBuckets value or an error if the edge
// was not loaded in eager-loading.
func (e FiatCurrencyEdges) ProvisionBucketsOrErr() ([]*ProvisionBucket, error) {
	if e.loadedTypes[1] {
		return e.ProvisionBuckets, nil
	}
	return nil, &NotLoadedError{edge: "provision_buckets"}
}

// FinancialInstitutionsOrErr returns the FinancialInstitutions value or an error if the edge
// was not loaded in eager-loading.
func (e FiatCurrencyEdges) FinancialInstitutionsOrErr() ([]*FinancialInstitution, error) {
	if e.loadedTypes[2] {
		return e.FinancialInstitutions, nil
	}
	return nil, &NotLoadedError{edge: "financialInstitutions"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*FiatCurrency) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case fiatcurrency.FieldMarketRate:
			values[i] = new(decimal.Decimal)
		case fiatcurrency.FieldIsEnabled:
			values[i] = new(sql.NullBool)
		case fiatcurrency.FieldDecimals:
			values[i] = new(sql.NullInt64)
		case fiatcurrency.FieldCode, fiatcurrency.FieldShortName, fiatcurrency.FieldSymbol, fiatcurrency.FieldName:
			values[i] = new(sql.NullString)
		case fiatcurrency.FieldCreatedAt, fiatcurrency.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case fiatcurrency.FieldID:
			values[i] = new(uuid.UUID)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the FiatCurrency fields.
func (fc *FiatCurrency) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case fiatcurrency.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				fc.ID = *value
			}
		case fiatcurrency.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				fc.CreatedAt = value.Time
			}
		case fiatcurrency.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				fc.UpdatedAt = value.Time
			}
		case fiatcurrency.FieldCode:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field code", values[i])
			} else if value.Valid {
				fc.Code = value.String
			}
		case fiatcurrency.FieldShortName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field short_name", values[i])
			} else if value.Valid {
				fc.ShortName = value.String
			}
		case fiatcurrency.FieldDecimals:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field decimals", values[i])
			} else if value.Valid {
				fc.Decimals = int(value.Int64)
			}
		case fiatcurrency.FieldSymbol:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field symbol", values[i])
			} else if value.Valid {
				fc.Symbol = value.String
			}
		case fiatcurrency.FieldName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[i])
			} else if value.Valid {
				fc.Name = value.String
			}
		case fiatcurrency.FieldMarketRate:
			if value, ok := values[i].(*decimal.Decimal); !ok {
				return fmt.Errorf("unexpected type %T for field market_rate", values[i])
			} else if value != nil {
				fc.MarketRate = *value
			}
		case fiatcurrency.FieldIsEnabled:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_enabled", values[i])
			} else if value.Valid {
				fc.IsEnabled = value.Bool
			}
		default:
			fc.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the FiatCurrency.
// This includes values selected through modifiers, order, etc.
func (fc *FiatCurrency) Value(name string) (ent.Value, error) {
	return fc.selectValues.Get(name)
}

// QueryProviders queries the "providers" edge of the FiatCurrency entity.
func (fc *FiatCurrency) QueryProviders() *ProviderProfileQuery {
	return NewFiatCurrencyClient(fc.config).QueryProviders(fc)
}

// QueryProvisionBuckets queries the "provision_buckets" edge of the FiatCurrency entity.
func (fc *FiatCurrency) QueryProvisionBuckets() *ProvisionBucketQuery {
	return NewFiatCurrencyClient(fc.config).QueryProvisionBuckets(fc)
}

// QueryFinancialInstitutions queries the "financialInstitutions" edge of the FiatCurrency entity.
func (fc *FiatCurrency) QueryFinancialInstitutions() *FinancialInstitutionQuery {
	return NewFiatCurrencyClient(fc.config).QueryFinancialInstitutions(fc)
}

// Update returns a builder for updating this FiatCurrency.
// Note that you need to call FiatCurrency.Unwrap() before calling this method if this FiatCurrency
// was returned from a transaction, and the transaction was committed or rolled back.
func (fc *FiatCurrency) Update() *FiatCurrencyUpdateOne {
	return NewFiatCurrencyClient(fc.config).UpdateOne(fc)
}

// Unwrap unwraps the FiatCurrency entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (fc *FiatCurrency) Unwrap() *FiatCurrency {
	_tx, ok := fc.config.driver.(*txDriver)
	if !ok {
		panic("ent: FiatCurrency is not a transactional entity")
	}
	fc.config.driver = _tx.drv
	return fc
}

// String implements the fmt.Stringer.
func (fc *FiatCurrency) String() string {
	var builder strings.Builder
	builder.WriteString("FiatCurrency(")
	builder.WriteString(fmt.Sprintf("id=%v, ", fc.ID))
	builder.WriteString("created_at=")
	builder.WriteString(fc.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(fc.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("code=")
	builder.WriteString(fc.Code)
	builder.WriteString(", ")
	builder.WriteString("short_name=")
	builder.WriteString(fc.ShortName)
	builder.WriteString(", ")
	builder.WriteString("decimals=")
	builder.WriteString(fmt.Sprintf("%v", fc.Decimals))
	builder.WriteString(", ")
	builder.WriteString("symbol=")
	builder.WriteString(fc.Symbol)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(fc.Name)
	builder.WriteString(", ")
	builder.WriteString("market_rate=")
	builder.WriteString(fmt.Sprintf("%v", fc.MarketRate))
	builder.WriteString(", ")
	builder.WriteString("is_enabled=")
	builder.WriteString(fmt.Sprintf("%v", fc.IsEnabled))
	builder.WriteByte(')')
	return builder.String()
}

// FiatCurrencies is a parsable slice of FiatCurrency.
type FiatCurrencies []*FiatCurrency
