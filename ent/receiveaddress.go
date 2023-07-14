// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent/paymentorder"
	"github.com/paycrest/paycrest-protocol/ent/receiveaddress"
)

// ReceiveAddress is the model entity for the ReceiveAddress schema.
type ReceiveAddress struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// Address holds the value of the "address" field.
	Address string `json:"address,omitempty"`
	// AccountIndex holds the value of the "account_index" field.
	AccountIndex int `json:"account_index,omitempty"`
	// Status holds the value of the "status" field.
	Status receiveaddress.Status `json:"status,omitempty"`
	// LastIndexedBlock holds the value of the "last_indexed_block" field.
	LastIndexedBlock int64 `json:"last_indexed_block,omitempty"`
	// LastUsed holds the value of the "last_used" field.
	LastUsed time.Time `json:"last_used,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the ReceiveAddressQuery when eager-loading is set.
	Edges                         ReceiveAddressEdges `json:"edges"`
	payment_order_receive_address *uuid.UUID
	selectValues                  sql.SelectValues
}

// ReceiveAddressEdges holds the relations/edges for other nodes in the graph.
type ReceiveAddressEdges struct {
	// PaymentOrder holds the value of the payment_order edge.
	PaymentOrder *PaymentOrder `json:"payment_order,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// PaymentOrderOrErr returns the PaymentOrder value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e ReceiveAddressEdges) PaymentOrderOrErr() (*PaymentOrder, error) {
	if e.loadedTypes[0] {
		if e.PaymentOrder == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: paymentorder.Label}
		}
		return e.PaymentOrder, nil
	}
	return nil, &NotLoadedError{edge: "payment_order"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*ReceiveAddress) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case receiveaddress.FieldID, receiveaddress.FieldAccountIndex, receiveaddress.FieldLastIndexedBlock:
			values[i] = new(sql.NullInt64)
		case receiveaddress.FieldAddress, receiveaddress.FieldStatus:
			values[i] = new(sql.NullString)
		case receiveaddress.FieldCreatedAt, receiveaddress.FieldUpdatedAt, receiveaddress.FieldLastUsed:
			values[i] = new(sql.NullTime)
		case receiveaddress.ForeignKeys[0]: // payment_order_receive_address
			values[i] = &sql.NullScanner{S: new(uuid.UUID)}
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the ReceiveAddress fields.
func (ra *ReceiveAddress) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case receiveaddress.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			ra.ID = int(value.Int64)
		case receiveaddress.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				ra.CreatedAt = value.Time
			}
		case receiveaddress.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				ra.UpdatedAt = value.Time
			}
		case receiveaddress.FieldAddress:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field address", values[i])
			} else if value.Valid {
				ra.Address = value.String
			}
		case receiveaddress.FieldAccountIndex:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field account_index", values[i])
			} else if value.Valid {
				ra.AccountIndex = int(value.Int64)
			}
		case receiveaddress.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				ra.Status = receiveaddress.Status(value.String)
			}
		case receiveaddress.FieldLastIndexedBlock:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field last_indexed_block", values[i])
			} else if value.Valid {
				ra.LastIndexedBlock = value.Int64
			}
		case receiveaddress.FieldLastUsed:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field last_used", values[i])
			} else if value.Valid {
				ra.LastUsed = value.Time
			}
		case receiveaddress.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullScanner); !ok {
				return fmt.Errorf("unexpected type %T for field payment_order_receive_address", values[i])
			} else if value.Valid {
				ra.payment_order_receive_address = new(uuid.UUID)
				*ra.payment_order_receive_address = *value.S.(*uuid.UUID)
			}
		default:
			ra.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the ReceiveAddress.
// This includes values selected through modifiers, order, etc.
func (ra *ReceiveAddress) Value(name string) (ent.Value, error) {
	return ra.selectValues.Get(name)
}

// QueryPaymentOrder queries the "payment_order" edge of the ReceiveAddress entity.
func (ra *ReceiveAddress) QueryPaymentOrder() *PaymentOrderQuery {
	return NewReceiveAddressClient(ra.config).QueryPaymentOrder(ra)
}

// Update returns a builder for updating this ReceiveAddress.
// Note that you need to call ReceiveAddress.Unwrap() before calling this method if this ReceiveAddress
// was returned from a transaction, and the transaction was committed or rolled back.
func (ra *ReceiveAddress) Update() *ReceiveAddressUpdateOne {
	return NewReceiveAddressClient(ra.config).UpdateOne(ra)
}

// Unwrap unwraps the ReceiveAddress entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (ra *ReceiveAddress) Unwrap() *ReceiveAddress {
	_tx, ok := ra.config.driver.(*txDriver)
	if !ok {
		panic("ent: ReceiveAddress is not a transactional entity")
	}
	ra.config.driver = _tx.drv
	return ra
}

// String implements the fmt.Stringer.
func (ra *ReceiveAddress) String() string {
	var builder strings.Builder
	builder.WriteString("ReceiveAddress(")
	builder.WriteString(fmt.Sprintf("id=%v, ", ra.ID))
	builder.WriteString("created_at=")
	builder.WriteString(ra.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(ra.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("address=")
	builder.WriteString(ra.Address)
	builder.WriteString(", ")
	builder.WriteString("account_index=")
	builder.WriteString(fmt.Sprintf("%v", ra.AccountIndex))
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", ra.Status))
	builder.WriteString(", ")
	builder.WriteString("last_indexed_block=")
	builder.WriteString(fmt.Sprintf("%v", ra.LastIndexedBlock))
	builder.WriteString(", ")
	builder.WriteString("last_used=")
	builder.WriteString(ra.LastUsed.Format(time.ANSIC))
	builder.WriteByte(')')
	return builder.String()
}

// ReceiveAddresses is a parsable slice of ReceiveAddress.
type ReceiveAddresses []*ReceiveAddress
