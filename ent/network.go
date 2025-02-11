// Code generated by ent, DO NOT EDIT.

package ent

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/paycrest/aggregator/ent/network"
	"github.com/shopspring/decimal"
)

// Network is the model entity for the Network schema.
type Network struct {
	config `json:"-"`
	// ID of the ent.
	ID int `json:"id,omitempty"`
	// CreatedAt holds the value of the "created_at" field.
	CreatedAt time.Time `json:"created_at,omitempty"`
	// UpdatedAt holds the value of the "updated_at" field.
	UpdatedAt time.Time `json:"updated_at,omitempty"`
	// ChainID holds the value of the "chain_id" field.
	ChainID int64 `json:"chain_id,omitempty"`
	// ChainIDHex holds the value of the "chain_id_hex" field.
	ChainIDHex string `json:"chain_id_hex,omitempty"`
	// Identifier holds the value of the "identifier" field.
	Identifier string `json:"identifier,omitempty"`
	// RPCEndpoint holds the value of the "rpc_endpoint" field.
	RPCEndpoint string `json:"rpc_endpoint,omitempty"`
	// GatewayContractAddress holds the value of the "gateway_contract_address" field.
	GatewayContractAddress string `json:"gateway_contract_address,omitempty"`
	// IsTestnet holds the value of the "is_testnet" field.
	IsTestnet bool `json:"is_testnet,omitempty"`
	// Fee holds the value of the "fee" field.
	Fee decimal.Decimal `json:"fee,omitempty"`
	// IsEnabled holds the value of the "is_enabled" field.
	IsEnabled bool `json:"is_enabled,omitempty"`
	// Edges holds the relations/edges for other nodes in the graph.
	// The values are being populated by the NetworkQuery when eager-loading is set.
	Edges        NetworkEdges `json:"edges"`
	selectValues sql.SelectValues
}

// NetworkEdges holds the relations/edges for other nodes in the graph.
type NetworkEdges struct {
	// Tokens holds the value of the tokens edge.
	Tokens []*Token `json:"tokens,omitempty"`
	// loadedTypes holds the information for reporting if a
	// type was loaded (or requested) in eager-loading or not.
	loadedTypes [1]bool
}

// TokensOrErr returns the Tokens value or an error if the edge
// was not loaded in eager-loading.
func (e NetworkEdges) TokensOrErr() ([]*Token, error) {
	if e.loadedTypes[0] {
		return e.Tokens, nil
	}
	return nil, &NotLoadedError{edge: "tokens"}
}

// scanValues returns the types for scanning values from sql.Rows.
func (*Network) scanValues(columns []string) ([]any, error) {
	values := make([]any, len(columns))
	for i := range columns {
		switch columns[i] {
		case network.FieldFee:
			values[i] = new(decimal.Decimal)
		case network.FieldIsTestnet, network.FieldIsEnabled:
			values[i] = new(sql.NullBool)
		case network.FieldID, network.FieldChainID:
			values[i] = new(sql.NullInt64)
		case network.FieldChainIDHex, network.FieldIdentifier, network.FieldRPCEndpoint, network.FieldGatewayContractAddress:
			values[i] = new(sql.NullString)
		case network.FieldCreatedAt, network.FieldUpdatedAt:
			values[i] = new(sql.NullTime)
		default:
			values[i] = new(sql.UnknownType)
		}
	}
	return values, nil
}

// assignValues assigns the values that were returned from sql.Rows (after scanning)
// to the Network fields.
func (n *Network) assignValues(columns []string, values []any) error {
	if m, n := len(values), len(columns); m < n {
		return fmt.Errorf("mismatch number of scan values: %d != %d", m, n)
	}
	for i := range columns {
		switch columns[i] {
		case network.FieldID:
			value, ok := values[i].(*sql.NullInt64)
			if !ok {
				return fmt.Errorf("unexpected type %T for field id", value)
			}
			n.ID = int(value.Int64)
		case network.FieldCreatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field created_at", values[i])
			} else if value.Valid {
				n.CreatedAt = value.Time
			}
		case network.FieldUpdatedAt:
			if value, ok := values[i].(*sql.NullTime); !ok {
				return fmt.Errorf("unexpected type %T for field updated_at", values[i])
			} else if value.Valid {
				n.UpdatedAt = value.Time
			}
		case network.FieldChainID:
			if value, ok := values[i].(*sql.NullInt64); !ok {
				return fmt.Errorf("unexpected type %T for field chain_id", values[i])
			} else if value.Valid {
				n.ChainID = value.Int64
			}
		case network.FieldChainIDHex:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field chain_id_hex", values[i])
			} else if value.Valid {
				n.ChainIDHex = value.String
			}
		case network.FieldIdentifier:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field identifier", values[i])
			} else if value.Valid {
				n.Identifier = value.String
			}
		case network.FieldRPCEndpoint:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field rpc_endpoint", values[i])
			} else if value.Valid {
				n.RPCEndpoint = value.String
			}
		case network.FieldGatewayContractAddress:
			if value, ok := values[i].(*sql.NullString); !ok {
				return fmt.Errorf("unexpected type %T for field gateway_contract_address", values[i])
			} else if value.Valid {
				n.GatewayContractAddress = value.String
			}
		case network.FieldIsTestnet:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_testnet", values[i])
			} else if value.Valid {
				n.IsTestnet = value.Bool
			}
		case network.FieldFee:
			if value, ok := values[i].(*decimal.Decimal); !ok {
				return fmt.Errorf("unexpected type %T for field fee", values[i])
			} else if value != nil {
				n.Fee = *value
			}
		case network.FieldIsEnabled:
			if value, ok := values[i].(*sql.NullBool); !ok {
				return fmt.Errorf("unexpected type %T for field is_enabled", values[i])
			} else if value.Valid {
				n.IsEnabled = value.Bool
			}
		default:
			n.selectValues.Set(columns[i], values[i])
		}
	}
	return nil
}

// Value returns the ent.Value that was dynamically selected and assigned to the Network.
// This includes values selected through modifiers, order, etc.
func (n *Network) Value(name string) (ent.Value, error) {
	return n.selectValues.Get(name)
}

// QueryTokens queries the "tokens" edge of the Network entity.
func (n *Network) QueryTokens() *TokenQuery {
	return NewNetworkClient(n.config).QueryTokens(n)
}

// Update returns a builder for updating this Network.
// Note that you need to call Network.Unwrap() before calling this method if this Network
// was returned from a transaction, and the transaction was committed or rolled back.
func (n *Network) Update() *NetworkUpdateOne {
	return NewNetworkClient(n.config).UpdateOne(n)
}

// Unwrap unwraps the Network entity that was returned from a transaction after it was closed,
// so that all future queries will be executed through the driver which created the transaction.
func (n *Network) Unwrap() *Network {
	_tx, ok := n.config.driver.(*txDriver)
	if !ok {
		panic("ent: Network is not a transactional entity")
	}
	n.config.driver = _tx.drv
	return n
}

// String implements the fmt.Stringer.
func (n *Network) String() string {
	var builder strings.Builder
	builder.WriteString("Network(")
	builder.WriteString(fmt.Sprintf("id=%v, ", n.ID))
	builder.WriteString("created_at=")
	builder.WriteString(n.CreatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("updated_at=")
	builder.WriteString(n.UpdatedAt.Format(time.ANSIC))
	builder.WriteString(", ")
	builder.WriteString("chain_id=")
	builder.WriteString(fmt.Sprintf("%v", n.ChainID))
	builder.WriteString(", ")
	builder.WriteString("chain_id_hex=")
	builder.WriteString(n.ChainIDHex)
	builder.WriteString(", ")
	builder.WriteString("identifier=")
	builder.WriteString(n.Identifier)
	builder.WriteString(", ")
	builder.WriteString("rpc_endpoint=")
	builder.WriteString(n.RPCEndpoint)
	builder.WriteString(", ")
	builder.WriteString("gateway_contract_address=")
	builder.WriteString(n.GatewayContractAddress)
	builder.WriteString(", ")
	builder.WriteString("is_testnet=")
	builder.WriteString(fmt.Sprintf("%v", n.IsTestnet))
	builder.WriteString(", ")
	builder.WriteString("fee=")
	builder.WriteString(fmt.Sprintf("%v", n.Fee))
	builder.WriteString(", ")
	builder.WriteString("is_enabled=")
	builder.WriteString(fmt.Sprintf("%v", n.IsEnabled))
	builder.WriteByte(')')
	return builder.String()
}

// Networks is a parsable slice of Network.
type Networks []*Network
