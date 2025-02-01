// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/paycrest/aggregator/ent/fiatcurrency"
	"github.com/paycrest/aggregator/ent/institution"
)

// Institution is the model entity for the Institution schema.
type Institution struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Code holds the value of the "code" field.
	Code string `json:"code,omitempty"`
	// Name holds the value of the "name" field.
	Name string `json:"name,omitempty"`
	// Type holds the value of the "type" field.
	Type institution.Type `json:"type,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the InstitutionQuery when eager-loading is set.
	Edges                      InstitutionEdges `json:"edges"`
	fiat_currency_institutions *uuid.UUID
	selectValues               sql.SelectValues
}

// InstitutionEdges holds the relations/edges for other nodes in the graph.
type InstitutionEdges struct {
	// FiatCurrency holds the value of the fiat_currency edge.
	FiatCurrency *FiatCurrency `json:"fiat_currency,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// FiatCurrencyOrErr returns the FiatCurrency value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e InstitutionEdges) FiatCurrencyOrErr() (*FiatCurrency, error) {
	if e.FiatCurrency != nil {
		return e.FiatCurrency, nil
	} else if e.loadedTypes[0] {
		return nil, &NotFoundError{label: fiatcurrency.Label}
	}
	return nil, &NotLoadedError{edge: "fiat_currency"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Institution) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case institution.FieldID:
			values[i] = new(sql.NullInt64)
		case institution.FieldCode, institution.FieldName, institution.FieldType:
			values[i] = new(sql.NullString)
		case institution.FieldCreatedAt, institution.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case institution.ForeignKeys[0]: // fiat_currency_institutions
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Institution fields.
func (i *Institution) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for j := range columns {
		switch columns[j] {
		case institution.FieldID:
			value, ok := values[j].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			i.ID = int(value.Int64)
		case institution.FieldCreatedAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[j])
			} else if value.Valid {
				i.CreatedAt = value.Time
			}
		case institution.FieldUpdatedAt:
			if value, ok := values[j].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[j])
			} else if value.Valid {
				i.UpdatedAt = value.Time
			}
		case institution.FieldCode:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field code", values[j])
			} else if value.Valid {
				i.Code = value.String
			}
		case institution.FieldName:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field name", values[j])
			} else if value.Valid {
				i.Name = value.String
			}
		case institution.FieldType:
			if value, ok := values[j].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field type", values[j])
			} else if value.Valid {
				i.Type = institution.Type(value.String)
			}
		case institution.ForeignKeys[0]:
			if value, ok := values[j].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field fiat_currency_institutions", values[j])
			} else if value.Valid {
				i.fiat_currency_institutions = new(uuid.UUID)
				*i.fiat_currency_institutions = *value.S.(*uuid.UUID)
			}
		default:
			i.selectValues.Set(columns[j], values[j])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Institution.
// This includes values selected through modifiers, order, etc.
func (i *Institution) Value(name string) (ent.Value, error) {
	return i.selectValues.Get(name)
}

// QueryFiatCurrency queries the "fiat_currency" edge of the Institution entity.
func (i *Institution) QueryFiatCurrency() *FiatCurrencyQuery {
	return NewInstitutionClient(i.config).QueryFiatCurrency(i)
}

// Update returns a builder for updating this Institution.
// Note that you need to call Institution.Unwrap() before calling this method if this Institution
// was returned from a transaction, and the transaction was committed or rolled back.
func (i *Institution) Update() *InstitutionUpdateOne {
	return NewInstitutionClient(i.config).UpdateOne(i)
}

// Unwrap unwraps the Institution entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (i *Institution) Unwrap() *Institution {
	_tx, ok := i.config.driver.(*txDriver)
	if !ok {
		panic("ent: Institution is not a transactional entity")
	}
	i.config.driver = _tx.drv
	return i
}

// String implements the fmt.Stringer.
func (i *Institution) String() string {
	var builder strings.Builder
	builder.WriteString("Institution(")
	builder.WriteString(fmt.Sprintf("id=%v, ", i.ID))
	builder.WriteString("created_at=")
	builder.WriteString(i.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(i.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("code=")
	builder.WriteString(i.Code)
	builder.WriteString(", ")
	builder.WriteString("name=")
	builder.WriteString(i.Name)
	builder.WriteString(", ")
	builder.WriteString("type=")
	builder.WriteString(fmt.Sprintf("%v", i.Type))
	builder.WriteByte(')')
	return builder.String()
}

// Institutions is a parsable slice of Institution.
type Institutions []*Institution
