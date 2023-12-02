// Code generated by ent, DO NOT EDIT.

package ent

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/google/uuid"
	"github.com/paycrest/protocol/ent/lockorderfulfillment"
	"github.com/paycrest/protocol/ent/lockpaymentorder"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/provisionbucket"
	"github.com/paycrest/protocol/ent/token"
	"github.com/shopspring/decimal"
)

// LockPaymentOrder is the model entity for the LockPaymentOrder schema.
type LockPaymentOrder struct {
	config `json:"-"`
	// ID of the ent.
	ID uuid.UUID `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// OrderID holds the value of the "order_id" field.
	OrderID string `json:"order_id,omitempty"`
	// Amount holds the value of the "amount" field.
	Amount decimal.Decimal `json:"amount,omitempty"`
	// Rate holds the value of the "rate" field.
	Rate decimal.Decimal `json:"rate,omitempty"`
	// OrderPercent holds the value of the "order_percent" field.
	OrderPercent decimal.Decimal `json:"order_percent,omitempty"`
	// TxHash holds the value of the "tx_hash" field.
	TxHash string `json:"tx_hash,omitempty"`
	// Label holds the value of the "label" field.
	Label string `json:"label,omitempty"`
	// Status holds the value of the "status" field.
	Status lockpaymentorder.Status `json:"status,omitempty"`
	// BlockNumber holds the value of the "block_number" field.
	BlockNumber int64 `json:"block_number,omitempty"`
	// Institution holds the value of the "institution" field.
	Institution string `json:"institution,omitempty"`
	// AccountIdentifier holds the value of the "account_identifier" field.
	AccountIdentifier string `json:"account_identifier,omitempty"`
	// AccountName holds the value of the "account_name" field.
	AccountName string `json:"account_name,omitempty"`
	// Memo holds the value of the "memo" field.
	Memo string `json:"memo,omitempty"`
	// CancellationCount holds the value of the "cancellation_count" field.
	CancellationCount int `json:"cancellation_count,omitempty"`
	// CancellationReasons holds the value of the "cancellation_reasons" field.
	CancellationReasons []string `json:"cancellation_reasons,omitempty"`
	// IsRefunded holds the value of the "is_refunded" field.
	IsRefunded bool `json:"is_refunded,omitempty"`
	// RefundTxHash holds the value of the "refund_tx_hash" field.
	RefundTxHash string `json:"refund_tx_hash,omitempty"`
	// IsRefundConfirmed holds the value of the "is_refund_confirmed" field.
	IsRefundConfirmed bool `json:"is_refund_confirmed,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the LockPaymentOrderQuery when eager-loading is set.
	Edges                                LockPaymentOrderEdges `json:"edges"`
	provider_profile_assigned_orders     *string
	provision_bucket_lock_payment_orders *int
	token_lock_payment_orders            *int
	selectValues                         sql.SelectValues
}

// LockPaymentOrderEdges holds the relations/edges for other nodes in the graph.
type LockPaymentOrderEdges struct {
	// Token holds the value of the token edge.
	Token *Token `json:"token,omitempty"`
	// ProvisionBucket holds the value of the provision_bucket edge.
	ProvisionBucket *ProvisionBucket `json:"provision_bucket,omitempty"`
	// Provider holds the value of the provider edge.
	Provider *ProviderProfile `json:"provider,omitempty"`
	// Fulfillment holds the value of the fulfillment edge.
	Fulfillment *LockOrderFulfillment `json:"fulfillment,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [4]bool
}

// TokenOrErr returns the Token value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LockPaymentOrderEdges) TokenOrErr() (*Token, error) {
	if e.loadedTypes[0] {
		if e.Token == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: token.Label}
		}
		return e.Token, nil
	}
	return nil, &NotLoadedError{edge: "token"}
}

// ProvisionBucketOrErr returns the ProvisionBucket value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LockPaymentOrderEdges) ProvisionBucketOrErr() (*ProvisionBucket, error) {
	if e.loadedTypes[1] {
		if e.ProvisionBucket == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: provisionbucket.Label}
		}
		return e.ProvisionBucket, nil
	}
	return nil, &NotLoadedError{edge: "provision_bucket"}
}

// ProviderOrErr returns the Provider value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LockPaymentOrderEdges) ProviderOrErr() (*ProviderProfile, error) {
	if e.loadedTypes[2] {
		if e.Provider == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: providerprofile.Label}
		}
		return e.Provider, nil
	}
	return nil, &NotLoadedError{edge: "provider"}
}

// FulfillmentOrErr returns the Fulfillment value or an error if the edge
// was not loaded in eager-loading, or loaded but was not found.
func (e LockPaymentOrderEdges) FulfillmentOrErr() (*LockOrderFulfillment, error) {
	if e.loadedTypes[3] {
		if e.Fulfillment == nil {
			// Edge was loaded but was not found.
			return nil, &NotFoundError{label: lockorderfulfillment.Label}
		}
		return e.Fulfillment, nil
	}
	return nil, &NotLoadedError{edge: "fulfillment"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*LockPaymentOrder) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case lockpaymentorder.FieldCancellationReasons:
			values[i] = new([]byte)
		case lockpaymentorder.FieldAmount, lockpaymentorder.FieldRate, lockpaymentorder.FieldOrderPercent:
			values[i] = new(decimal.Decimal)
		case lockpaymentorder.FieldIsRefunded, lockpaymentorder.FieldIsRefundConfirmed:
			values[i] = new(sql.NullBool)
		case lockpaymentorder.FieldBlockNumber, lockpaymentorder.FieldCancellationCount:
			values[i] = new(sql.NullInt64)
		case lockpaymentorder.FieldOrderID, lockpaymentorder.FieldTxHash, lockpaymentorder.FieldLabel, lockpaymentorder.FieldStatus, lockpaymentorder.FieldInstitution, lockpaymentorder.FieldAccountIdentifier, lockpaymentorder.FieldAccountName, lockpaymentorder.FieldMemo, lockpaymentorder.FieldRefundTxHash:
			values[i] = new(sql.NullString)
		case lockpaymentorder.FieldCreatedAt, lockpaymentorder.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		case lockpaymentorder.FieldID:
			values[i] = new(uuid.UUID)
		case lockpaymentorder.ForeignKeys[0]: // provider_profile_assigned_orders
			values[i] = new(sql.NullString)
		case lockpaymentorder.ForeignKeys[1]: // provision_bucket_lock_payment_orders
			values[i] = new(sql.NullInt64)
		case lockpaymentorder.ForeignKeys[2]: // token_lock_payment_orders
			values[i] = new(sql.NullInt64)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the LockPaymentOrder fields.
func (lpo *LockPaymentOrder) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case lockpaymentorder.FieldID:
			if value, ok := values[i].(*uuid.UUID); !ok {
				return fmt.Errorf("unexpected type %T for field id", values[i])
			} else if value != nil {
				lpo.ID = *value
			}
		case lockpaymentorder.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				lpo.CreatedAt = value.Time
			}
		case lockpaymentorder.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				lpo.UpdatedAt = value.Time
			}
		case lockpaymentorder.FieldOrderID:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field order_id", values[i])
			} else if value.Valid {
				lpo.OrderID = value.String
			}
		case lockpaymentorder.FieldAmount:
			if value, ok := values[i].(*decimal.Decimal); !ok {
				return fmt.Errorf("unexpected type %T for field amount", values[i])
			} else if value != nil {
				lpo.Amount = *value
			}
		case lockpaymentorder.FieldRate:
			if value, ok := values[i].(*decimal.Decimal); !ok {
				return fmt.Errorf("unexpected type %T for field rate", values[i])
			} else if value != nil {
				lpo.Rate = *value
			}
		case lockpaymentorder.FieldOrderPercent:
			if value, ok := values[i].(*decimal.Decimal); !ok {
				return fmt.Errorf("unexpected type %T for field order_percent", values[i])
			} else if value != nil {
				lpo.OrderPercent = *value
			}
		case lockpaymentorder.FieldTxHash:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field tx_hash", values[i])
			} else if value.Valid {
				lpo.TxHash = value.String
			}
		case lockpaymentorder.FieldLabel:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field label", values[i])
			} else if value.Valid {
				lpo.Label = value.String
			}
		case lockpaymentorder.FieldStatus:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field status", values[i])
			} else if value.Valid {
				lpo.Status = lockpaymentorder.Status(value.String)
			}
		case lockpaymentorder.FieldBlockNumber:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field block_number", values[i])
			} else if value.Valid {
				lpo.BlockNumber = value.Int64
			}
		case lockpaymentorder.FieldInstitution:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field institution", values[i])
			} else if value.Valid {
				lpo.Institution = value.String
			}
		case lockpaymentorder.FieldAccountIdentifier:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field account_identifier", values[i])
			} else if value.Valid {
				lpo.AccountIdentifier = value.String
			}
		case lockpaymentorder.FieldAccountName:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field account_name", values[i])
			} else if value.Valid {
				lpo.AccountName = value.String
			}
		case lockpaymentorder.FieldMemo:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field memo", values[i])
			} else if value.Valid {
				lpo.Memo = value.String
			}
		case lockpaymentorder.FieldCancellationCount:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field cancellation_count", values[i])
			} else if value.Valid {
				lpo.CancellationCount = int(value.Int64)
			}
		case lockpaymentorder.FieldCancellationReasons:
			if value, ok := values[i].(*[]byte); !ok {
				return fmt.Errorf("unexpected type %T for field cancellation_reasons", values[i])
			} else if value != nil && len(*value) > 0 {
				if err := json.Unmarshal(*value, &lpo.CancellationReasons); err != nil {
					return fmt.Errorf("unmarshal field cancellation_reasons: %w", err)
				}
			}
		case lockpaymentorder.FieldIsRefunded:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_refunded", values[i])
			} else if value.Valid {
				lpo.IsRefunded = value.Bool
			}
		case lockpaymentorder.FieldRefundTxHash:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field refund_tx_hash", values[i])
			} else if value.Valid {
				lpo.RefundTxHash = value.String
			}
		case lockpaymentorder.FieldIsRefundConfirmed:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_refund_confirmed", values[i])
			} else if value.Valid {
				lpo.IsRefundConfirmed = value.Bool
			}
		case lockpaymentorder.ForeignKeys[0]:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field provider_profile_assigned_orders", values[i])
			} else if value.Valid {
				lpo.provider_profile_assigned_orders = new(string)
				*lpo.provider_profile_assigned_orders = value.String
			}
		case lockpaymentorder.ForeignKeys[1]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field provision_bucket_lock_payment_orders", value)
			} else if value.Valid {
				lpo.provision_bucket_lock_payment_orders = new(int)
				*lpo.provision_bucket_lock_payment_orders = int(value.Int64)
			}
		case lockpaymentorder.ForeignKeys[2]:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for edge-field token_lock_payment_orders", value)
			} else if value.Valid {
				lpo.token_lock_payment_orders = new(int)
				*lpo.token_lock_payment_orders = int(value.Int64)
			}
		default:
			lpo.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the LockPaymentOrder.
// This includes values selected through modifiers, order, etc.
func (lpo *LockPaymentOrder) Value(name string) (ent.Value, error) {
	return lpo.selectValues.Get(name)
}

// QueryToken queries the "token" edge of the LockPaymentOrder entity.
func (lpo *LockPaymentOrder) QueryToken() *TokenQuery {
	return NewLockPaymentOrderClient(lpo.config).QueryToken(lpo)
}

// QueryProvisionBucket queries the "provision_bucket" edge of the LockPaymentOrder entity.
func (lpo *LockPaymentOrder) QueryProvisionBucket() *ProvisionBucketQuery {
	return NewLockPaymentOrderClient(lpo.config).QueryProvisionBucket(lpo)
}

// QueryProvider queries the "provider" edge of the LockPaymentOrder entity.
func (lpo *LockPaymentOrder) QueryProvider() *ProviderProfileQuery {
	return NewLockPaymentOrderClient(lpo.config).QueryProvider(lpo)
}

// QueryFulfillment queries the "fulfillment" edge of the LockPaymentOrder entity.
func (lpo *LockPaymentOrder) QueryFulfillment() *LockOrderFulfillmentQuery {
	return NewLockPaymentOrderClient(lpo.config).QueryFulfillment(lpo)
}

// Update returns a builder for updating this LockPaymentOrder.
// Note that you need to call LockPaymentOrder.Unwrap() before calling this method if this LockPaymentOrder
// was returned from a transaction, and the transaction was committed or rolled back.
func (lpo *LockPaymentOrder) Update() *LockPaymentOrderUpdateOne {
	return NewLockPaymentOrderClient(lpo.config).UpdateOne(lpo)
}

// Unwrap unwraps the LockPaymentOrder entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (lpo *LockPaymentOrder) Unwrap() *LockPaymentOrder {
	_tx, ok := lpo.config.driver.(*txDriver)
	if !ok {
		panic("ent: LockPaymentOrder is not a transactional entity")
	}
	lpo.config.driver = _tx.drv
	return lpo
}

// String implements the fmt.Stringer.
func (lpo *LockPaymentOrder) String() string {
	var builder strings.Builder
	builder.WriteString("LockPaymentOrder(")
	builder.WriteString(fmt.Sprintf("id=%v, ", lpo.ID))
	builder.WriteString("created_at=")
	builder.WriteString(lpo.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(lpo.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("order_id=")
	builder.WriteString(lpo.OrderID)
	builder.WriteString(", ")
	builder.WriteString("amount=")
	builder.WriteString(fmt.Sprintf("%v", lpo.Amount))
	builder.WriteString(", ")
	builder.WriteString("rate=")
	builder.WriteString(fmt.Sprintf("%v", lpo.Rate))
	builder.WriteString(", ")
	builder.WriteString("order_percent=")
	builder.WriteString(fmt.Sprintf("%v", lpo.OrderPercent))
	builder.WriteString(", ")
	builder.WriteString("tx_hash=")
	builder.WriteString(lpo.TxHash)
	builder.WriteString(", ")
	builder.WriteString("label=")
	builder.WriteString(lpo.Label)
	builder.WriteString(", ")
	builder.WriteString("status=")
	builder.WriteString(fmt.Sprintf("%v", lpo.Status))
	builder.WriteString(", ")
	builder.WriteString("block_number=")
	builder.WriteString(fmt.Sprintf("%v", lpo.BlockNumber))
	builder.WriteString(", ")
	builder.WriteString("institution=")
	builder.WriteString(lpo.Institution)
	builder.WriteString(", ")
	builder.WriteString("account_identifier=")
	builder.WriteString(lpo.AccountIdentifier)
	builder.WriteString(", ")
	builder.WriteString("account_name=")
	builder.WriteString(lpo.AccountName)
	builder.WriteString(", ")
	builder.WriteString("memo=")
	builder.WriteString(lpo.Memo)
	builder.WriteString(", ")
	builder.WriteString("cancellation_count=")
	builder.WriteString(fmt.Sprintf("%v", lpo.CancellationCount))
	builder.WriteString(", ")
	builder.WriteString("cancellation_reasons=")
	builder.WriteString(fmt.Sprintf("%v", lpo.CancellationReasons))
	builder.WriteString(", ")
	builder.WriteString("is_refunded=")
	builder.WriteString(fmt.Sprintf("%v", lpo.IsRefunded))
	builder.WriteString(", ")
	builder.WriteString("refund_tx_hash=")
	builder.WriteString(lpo.RefundTxHash)
	builder.WriteString(", ")
	builder.WriteString("is_refund_confirmed=")
	builder.WriteString(fmt.Sprintf("%v", lpo.IsRefundConfirmed))
	builder.WriteByte(')')
	return builder.String()
}

// LockPaymentOrders is a parsable slice of LockPaymentOrder.
type LockPaymentOrders []*LockPaymentOrder
