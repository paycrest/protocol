// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/paycrest/aggregator/ent/providerordertoken"
	"github.com/paycrest/aggregator/ent/providerprofile"
	"github.com/shopspring/decimal"
)

// ProviderOrderToken is the model entity for the ProviderOrderToken schema.
type ProviderOrderToken struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Symbol holds the value of the "symbol" field.
	Symbol string `json:"symbol,omitempty"`
	// FixedConversionRate holds the value of the "fixed_conversion_rate" field.
	FixedConversionRate decimal.Decimal `json:"fixed_conversion_rate,omitempty"`
	// FloatingConversionRate holds the value of the "floating_conversion_rate" field.
	FloatingConversionRate decimal.Decimal `json:"floating_conversion_rate,omitempty"`
	// ConversionRateType holds the value of the "conversion_rate_type" field.
	ConversionRateType providerordertoken.ConversionRateType `json:"conversion_rate_type,omitempty"`
	// MaxOrderAmount holds the value of the "max_order_amount" field.
	MaxOrderAmount decimal.Decimal `json:"max_order_amount,omitempty"`
	// MinOrderAmount holds the value of the "min_order_amount" field.
	MinOrderAmount decimal.Decimal `json:"min_order_amount,omitempty"`
	// Addresses holds the value of the "addresses" field.
	Addresses []struct {
		Address string "json:\"address\""
		Network string "json:\"network\""
	} `json:"addresses,omitempty"`
	// RateSlippage holds the value of the "rate_slippage" field.
	RateSlippage decimal.Decimal `json:"rate_slippage,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ProviderOrderTokenQuery when eager-loading is set.
	Edges                         ProviderOrderTokenEdges `json:"edges"`
	provider_profile_order_tokens *string
	selectValues                  sql.SelectValues
}

// ProviderOrderTokenEdges holds the relations/edges for other nodes in the graph.
type ProviderOrderTokenEdges struct {
	// Provider holds the value of the provider edge.
	Provider *ProviderProfile `json:"provider,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// ProviderOrErr returns the Provider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ProviderOrderTokenEdges) ProviderOrErr() (*ProviderProfile, error) {
	if e.Provider != nil {
		return e.Provider, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: providerprofile.Label}
	}
	return nil, &NotLoadedError{edge: "provider"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ProviderOrderToken) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case providerordertoken.FieldAddresses:
			values[i] = new([]byte)
		case providerordertoken.FieldFixedConversionRate, providerordertoken.FieldFloatingConversionRate, providerordertoken.FieldMaxOrderAmount, providerordertoken.FieldMinOrderAmount, providerordertoken.FieldRateSlippage:
			values[i] = new(decimal.Decimal)
		case providerordertoken.FieldID:
			values[i] = new(sql.NullInt64)
		case providerordertoken.FieldSymbol, providerordertoken.FieldConversionRateType:
			values[i] = new(sql.NullString)
		case providerordertoken.FieldCreatedAt, providerordertoken.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case providerordertoken.ForeignKeys[0]: // provider_profile_order_tokens
			values[i] = new(sql.NullString)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ProviderOrderToken fields.
func (pot *ProviderOrderToken) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case providerordertoken.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			pot.ID = int(value.Int64)
		case providerordertoken.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				pot.CreatedAt = value.Time
			}
		case providerordertoken.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				pot.UpdatedAt = value.Time
			}
		case providerordertoken.FieldSymbol:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field symbol", values[i])
			} else if value.Valid {
				pot.Symbol = value.String
			}
		case providerordertoken.FieldFixedConversionRate:
			if value, ok := values[i].(*decimal.Decimal); !ok {
				return fmt.Errorf("unexpected type %T for field fixed_conversion_rate", values[i])
			} else if value != nil {
				pot.FixedConversionRate = *value
			}
		case providerordertoken.FieldFloatingConversionRate:
			if value, ok := values[i].(*decimal.Decimal); !ok {
				return fmt.Errorf("unexpected type %T for field floating_conversion_rate", values[i])
			} else if value != nil {
				pot.FloatingConversionRate = *value
			}
		case providerordertoken.FieldConversionRateType:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field conversion_rate_type", values[i])
			} else if value.Valid {
				pot.ConversionRateType = providerordertoken.ConversionRateType(value.String)
			}
		case providerordertoken.FieldMaxOrderAmount:
			if value, ok := values[i].(*decimal.Decimal); !ok {
				return fmt.Errorf("unexpected type %T for field max_order_amount", values[i])
			} else if value != nil {
				pot.MaxOrderAmount = *value
			}
		case providerordertoken.FieldMinOrderAmount:
			if value, ok := values[i].(*decimal.Decimal); !ok {
				return fmt.Errorf("unexpected type %T for field min_order_amount", values[i])
			} else if value != nil {
				pot.MinOrderAmount = *value
			}
		case providerordertoken.FieldAddresses:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field addresses", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &pot.Addresses); err != nil {
					return fmt.Errorf("unmarshal field addresses: %w", err)
				}
			}
		case providerordertoken.FieldRateSlippage:
			if value, ok := values[i].(*decimal.Decimal); !ok {
				return fmt.Errorf("unexpected type %T for field rate_slippage", values[i])
			} else if value != nil {
				pot.RateSlippage = *value
			}
		case providerordertoken.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field provider_profile_order_tokens", values[i])
			} else if value.Valid {
				pot.provider_profile_order_tokens = new(string)
				*pot.provider_profile_order_tokens = value.String
			}
		default:
			pot.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ProviderOrderToken.
// This includes values selected through modifiers, order, etc.
func (pot *ProviderOrderToken) Value(name string) (ent.Value, error) {
	return pot.selectValues.Get(name)
}

// QueryProvider queries the "provider" edge of the ProviderOrderToken entity.
func (pot *ProviderOrderToken) QueryProvider() *ProviderProfileQuery {
	return NewProviderOrderTokenClient(pot.config).QueryProvider(pot)
}

// Update returns a builder for updating this ProviderOrderToken.
// Note that you need to call ProviderOrderToken.Unwrap() before calling this method if this ProviderOrderToken
// was returned from a transaction, and the transaction was committed or rolled back.
func (pot *ProviderOrderToken) Update() *ProviderOrderTokenUpdateOne {
	return NewProviderOrderTokenClient(pot.config).UpdateOne(pot)
}

// Unwrap unwraps the ProviderOrderToken entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (pot *ProviderOrderToken) Unwrap() *ProviderOrderToken {
	_tx, ok := pot.config.driver.(*txDriver)
	if !ok {
		panic("ent: ProviderOrderToken is not a transactional entity")
	}
	pot.config.driver = _tx.drv
	return pot
}

// String implements the fmt.Stringer.
func (pot *ProviderOrderToken) String() string {
	var builder strings.Builder
	builder.WriteString("ProviderOrderToken(")
	builder.WriteString(fmt.Sprintf("id=%v, ", pot.ID))
	builder.WriteString("created_at=")
	builder.WriteString(pot.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(pot.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("symbol=")
	builder.WriteString(pot.Symbol)
	builder.WriteString(", ")
	builder.WriteString("fixed_conversion_rate=")
	builder.WriteString(fmt.Sprintf("%v", pot.FixedConversionRate))
	builder.WriteString(", ")
	builder.WriteString("floating_conversion_rate=")
	builder.WriteString(fmt.Sprintf("%v", pot.FloatingConversionRate))
	builder.WriteString(", ")
	builder.WriteString("conversion_rate_type=")
	builder.WriteString(fmt.Sprintf("%v", pot.ConversionRateType))
	builder.WriteString(", ")
	builder.WriteString("max_order_amount=")
	builder.WriteString(fmt.Sprintf("%v", pot.MaxOrderAmount))
	builder.WriteString(", ")
	builder.WriteString("min_order_amount=")
	builder.WriteString(fmt.Sprintf("%v", pot.MinOrderAmount))
	builder.WriteString(", ")
	builder.WriteString("addresses=")
	builder.WriteString(fmt.Sprintf("%v", pot.Addresses))
	builder.WriteString(", ")
	builder.WriteString("rate_slippage=")
	builder.WriteString(fmt.Sprintf("%v", pot.RateSlippage))
	builder.WriteByte(')')
	return builder.String()
}

// ProviderOrderTokens is a parsable slice of ProviderOrderToken.
type ProviderOrderTokens []*ProviderOrderToken
