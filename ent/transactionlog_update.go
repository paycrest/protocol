// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/paycrest/protocol/ent/predicate"
	"github.com/paycrest/protocol/ent/transactionlog"
)

// TransactionLogUpdate is the builder for updating TransactionLog entities.
type TransactionLogUpdate struct {
	config
	hooks    []Hook
	mutation *TransactionLogMutation
}

// Where appends a list predicates to the TransactionLogUpdate builder.
func (tlu *TransactionLogUpdate) Where(ps ...predicate.TransactionLog) *TransactionLogUpdate {
	tlu.mutation.Where(ps...)
	return tlu
}

// SetSenderAddress sets the "sender_address" field.
func (tlu *TransactionLogUpdate) SetSenderAddress(s string) *TransactionLogUpdate {
	tlu.mutation.SetSenderAddress(s)
	return tlu
}

// SetNillableSenderAddress sets the "sender_address" field if the given value is not nil.
func (tlu *TransactionLogUpdate) SetNillableSenderAddress(s *string) *TransactionLogUpdate {
	if s != nil {
		tlu.SetSenderAddress(*s)
	}
	return tlu
}

// ClearSenderAddress clears the value of the "sender_address" field.
func (tlu *TransactionLogUpdate) ClearSenderAddress() *TransactionLogUpdate {
	tlu.mutation.ClearSenderAddress()
	return tlu
}

// SetProviderAddress sets the "provider_address" field.
func (tlu *TransactionLogUpdate) SetProviderAddress(s string) *TransactionLogUpdate {
	tlu.mutation.SetProviderAddress(s)
	return tlu
}

// SetNillableProviderAddress sets the "provider_address" field if the given value is not nil.
func (tlu *TransactionLogUpdate) SetNillableProviderAddress(s *string) *TransactionLogUpdate {
	if s != nil {
		tlu.SetProviderAddress(*s)
	}
	return tlu
}

// ClearProviderAddress clears the value of the "provider_address" field.
func (tlu *TransactionLogUpdate) ClearProviderAddress() *TransactionLogUpdate {
	tlu.mutation.ClearProviderAddress()
	return tlu
}

// SetGatewayID sets the "gateway_id" field.
func (tlu *TransactionLogUpdate) SetGatewayID(s string) *TransactionLogUpdate {
	tlu.mutation.SetGatewayID(s)
	return tlu
}

// SetNillableGatewayID sets the "gateway_id" field if the given value is not nil.
func (tlu *TransactionLogUpdate) SetNillableGatewayID(s *string) *TransactionLogUpdate {
	if s != nil {
		tlu.SetGatewayID(*s)
	}
	return tlu
}

// ClearGatewayID clears the value of the "gateway_id" field.
func (tlu *TransactionLogUpdate) ClearGatewayID() *TransactionLogUpdate {
	tlu.mutation.ClearGatewayID()
	return tlu
}

// SetNetwork sets the "network" field.
func (tlu *TransactionLogUpdate) SetNetwork(s string) *TransactionLogUpdate {
	tlu.mutation.SetNetwork(s)
	return tlu
}

// SetNillableNetwork sets the "network" field if the given value is not nil.
func (tlu *TransactionLogUpdate) SetNillableNetwork(s *string) *TransactionLogUpdate {
	if s != nil {
		tlu.SetNetwork(*s)
	}
	return tlu
}

// ClearNetwork clears the value of the "network" field.
func (tlu *TransactionLogUpdate) ClearNetwork() *TransactionLogUpdate {
	tlu.mutation.ClearNetwork()
	return tlu
}

// SetTransactionHash sets the "transaction_hash" field.
func (tlu *TransactionLogUpdate) SetTransactionHash(s string) *TransactionLogUpdate {
	tlu.mutation.SetTransactionHash(s)
	return tlu
}

// SetNillableTransactionHash sets the "transaction_hash" field if the given value is not nil.
func (tlu *TransactionLogUpdate) SetNillableTransactionHash(s *string) *TransactionLogUpdate {
	if s != nil {
		tlu.SetTransactionHash(*s)
	}
	return tlu
}

// ClearTransactionHash clears the value of the "transaction_hash" field.
func (tlu *TransactionLogUpdate) ClearTransactionHash() *TransactionLogUpdate {
	tlu.mutation.ClearTransactionHash()
	return tlu
}

// SetMetadata sets the "metadata" field.
func (tlu *TransactionLogUpdate) SetMetadata(m map[string]interface{}) *TransactionLogUpdate {
	tlu.mutation.SetMetadata(m)
	return tlu
}

// Mutation returns the TransactionLogMutation object of the builder.
func (tlu *TransactionLogUpdate) Mutation() *TransactionLogMutation {
	return tlu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (tlu *TransactionLogUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, tlu.sqlSave, tlu.mutation, tlu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tlu *TransactionLogUpdate) SaveX(ctx context.Context) int {
	affected, err := tlu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (tlu *TransactionLogUpdate) Exec(ctx context.Context) error {
	_, err := tlu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tlu *TransactionLogUpdate) ExecX(ctx context.Context) {
	if err := tlu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tlu *TransactionLogUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(transactionlog.Table, transactionlog.Columns, sqlgraph.NewFieldSpec(transactionlog.FieldID, field.TypeUUID))
	if ps := tlu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tlu.mutation.SenderAddress(); ok {
		_spec.SetField(transactionlog.FieldSenderAddress, field.TypeString, value)
	}
	if tlu.mutation.SenderAddressCleared() {
		_spec.ClearField(transactionlog.FieldSenderAddress, field.TypeString)
	}
	if value, ok := tlu.mutation.ProviderAddress(); ok {
		_spec.SetField(transactionlog.FieldProviderAddress, field.TypeString, value)
	}
	if tlu.mutation.ProviderAddressCleared() {
		_spec.ClearField(transactionlog.FieldProviderAddress, field.TypeString)
	}
	if value, ok := tlu.mutation.GatewayID(); ok {
		_spec.SetField(transactionlog.FieldGatewayID, field.TypeString, value)
	}
	if tlu.mutation.GatewayIDCleared() {
		_spec.ClearField(transactionlog.FieldGatewayID, field.TypeString)
	}
	if value, ok := tlu.mutation.Network(); ok {
		_spec.SetField(transactionlog.FieldNetwork, field.TypeString, value)
	}
	if tlu.mutation.NetworkCleared() {
		_spec.ClearField(transactionlog.FieldNetwork, field.TypeString)
	}
	if value, ok := tlu.mutation.TransactionHash(); ok {
		_spec.SetField(transactionlog.FieldTransactionHash, field.TypeString, value)
	}
	if tlu.mutation.TransactionHashCleared() {
		_spec.ClearField(transactionlog.FieldTransactionHash, field.TypeString)
	}
	if value, ok := tlu.mutation.Metadata(); ok {
		_spec.SetField(transactionlog.FieldMetadata, field.TypeJSON, value)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, tlu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{transactionlog.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	tlu.mutation.done = true
	return n, nil
}

// TransactionLogUpdateOne is the builder for updating a single TransactionLog entity.
type TransactionLogUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *TransactionLogMutation
}

// SetSenderAddress sets the "sender_address" field.
func (tluo *TransactionLogUpdateOne) SetSenderAddress(s string) *TransactionLogUpdateOne {
	tluo.mutation.SetSenderAddress(s)
	return tluo
}

// SetNillableSenderAddress sets the "sender_address" field if the given value is not nil.
func (tluo *TransactionLogUpdateOne) SetNillableSenderAddress(s *string) *TransactionLogUpdateOne {
	if s != nil {
		tluo.SetSenderAddress(*s)
	}
	return tluo
}

// ClearSenderAddress clears the value of the "sender_address" field.
func (tluo *TransactionLogUpdateOne) ClearSenderAddress() *TransactionLogUpdateOne {
	tluo.mutation.ClearSenderAddress()
	return tluo
}

// SetProviderAddress sets the "provider_address" field.
func (tluo *TransactionLogUpdateOne) SetProviderAddress(s string) *TransactionLogUpdateOne {
	tluo.mutation.SetProviderAddress(s)
	return tluo
}

// SetNillableProviderAddress sets the "provider_address" field if the given value is not nil.
func (tluo *TransactionLogUpdateOne) SetNillableProviderAddress(s *string) *TransactionLogUpdateOne {
	if s != nil {
		tluo.SetProviderAddress(*s)
	}
	return tluo
}

// ClearProviderAddress clears the value of the "provider_address" field.
func (tluo *TransactionLogUpdateOne) ClearProviderAddress() *TransactionLogUpdateOne {
	tluo.mutation.ClearProviderAddress()
	return tluo
}

// SetGatewayID sets the "gateway_id" field.
func (tluo *TransactionLogUpdateOne) SetGatewayID(s string) *TransactionLogUpdateOne {
	tluo.mutation.SetGatewayID(s)
	return tluo
}

// SetNillableGatewayID sets the "gateway_id" field if the given value is not nil.
func (tluo *TransactionLogUpdateOne) SetNillableGatewayID(s *string) *TransactionLogUpdateOne {
	if s != nil {
		tluo.SetGatewayID(*s)
	}
	return tluo
}

// ClearGatewayID clears the value of the "gateway_id" field.
func (tluo *TransactionLogUpdateOne) ClearGatewayID() *TransactionLogUpdateOne {
	tluo.mutation.ClearGatewayID()
	return tluo
}

// SetNetwork sets the "network" field.
func (tluo *TransactionLogUpdateOne) SetNetwork(s string) *TransactionLogUpdateOne {
	tluo.mutation.SetNetwork(s)
	return tluo
}

// SetNillableNetwork sets the "network" field if the given value is not nil.
func (tluo *TransactionLogUpdateOne) SetNillableNetwork(s *string) *TransactionLogUpdateOne {
	if s != nil {
		tluo.SetNetwork(*s)
	}
	return tluo
}

// ClearNetwork clears the value of the "network" field.
func (tluo *TransactionLogUpdateOne) ClearNetwork() *TransactionLogUpdateOne {
	tluo.mutation.ClearNetwork()
	return tluo
}

// SetTransactionHash sets the "transaction_hash" field.
func (tluo *TransactionLogUpdateOne) SetTransactionHash(s string) *TransactionLogUpdateOne {
	tluo.mutation.SetTransactionHash(s)
	return tluo
}

// SetNillableTransactionHash sets the "transaction_hash" field if the given value is not nil.
func (tluo *TransactionLogUpdateOne) SetNillableTransactionHash(s *string) *TransactionLogUpdateOne {
	if s != nil {
		tluo.SetTransactionHash(*s)
	}
	return tluo
}

// ClearTransactionHash clears the value of the "transaction_hash" field.
func (tluo *TransactionLogUpdateOne) ClearTransactionHash() *TransactionLogUpdateOne {
	tluo.mutation.ClearTransactionHash()
	return tluo
}

// SetMetadata sets the "metadata" field.
func (tluo *TransactionLogUpdateOne) SetMetadata(m map[string]interface{}) *TransactionLogUpdateOne {
	tluo.mutation.SetMetadata(m)
	return tluo
}

// Mutation returns the TransactionLogMutation object of the builder.
func (tluo *TransactionLogUpdateOne) Mutation() *TransactionLogMutation {
	return tluo.mutation
}

// Where appends a list predicates to the TransactionLogUpdate builder.
func (tluo *TransactionLogUpdateOne) Where(ps ...predicate.TransactionLog) *TransactionLogUpdateOne {
	tluo.mutation.Where(ps...)
	return tluo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (tluo *TransactionLogUpdateOne) Select(field string, fields ...string) *TransactionLogUpdateOne {
	tluo.fields = append([]string{field}, fields...)
	return tluo
}

// Save executes the query and returns the updated TransactionLog entity.
func (tluo *TransactionLogUpdateOne) Save(ctx context.Context) (*TransactionLog, error) {
	return withHooks(ctx, tluo.sqlSave, tluo.mutation, tluo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (tluo *TransactionLogUpdateOne) SaveX(ctx context.Context) *TransactionLog {
	node, err := tluo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (tluo *TransactionLogUpdateOne) Exec(ctx context.Context) error {
	_, err := tluo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (tluo *TransactionLogUpdateOne) ExecX(ctx context.Context) {
	if err := tluo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (tluo *TransactionLogUpdateOne) sqlSave(ctx context.Context) (_node *TransactionLog, err error) {
	_spec := sqlgraph.NewUpdateSpec(transactionlog.Table, transactionlog.Columns, sqlgraph.NewFieldSpec(transactionlog.FieldID, field.TypeUUID))
	id, ok := tluo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "TransactionLog.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := tluo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, transactionlog.FieldID)
		for _, f := range fields {
			if !transactionlog.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != transactionlog.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := tluo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := tluo.mutation.SenderAddress(); ok {
		_spec.SetField(transactionlog.FieldSenderAddress, field.TypeString, value)
	}
	if tluo.mutation.SenderAddressCleared() {
		_spec.ClearField(transactionlog.FieldSenderAddress, field.TypeString)
	}
	if value, ok := tluo.mutation.ProviderAddress(); ok {
		_spec.SetField(transactionlog.FieldProviderAddress, field.TypeString, value)
	}
	if tluo.mutation.ProviderAddressCleared() {
		_spec.ClearField(transactionlog.FieldProviderAddress, field.TypeString)
	}
	if value, ok := tluo.mutation.GatewayID(); ok {
		_spec.SetField(transactionlog.FieldGatewayID, field.TypeString, value)
	}
	if tluo.mutation.GatewayIDCleared() {
		_spec.ClearField(transactionlog.FieldGatewayID, field.TypeString)
	}
	if value, ok := tluo.mutation.Network(); ok {
		_spec.SetField(transactionlog.FieldNetwork, field.TypeString, value)
	}
	if tluo.mutation.NetworkCleared() {
		_spec.ClearField(transactionlog.FieldNetwork, field.TypeString)
	}
	if value, ok := tluo.mutation.TransactionHash(); ok {
		_spec.SetField(transactionlog.FieldTransactionHash, field.TypeString, value)
	}
	if tluo.mutation.TransactionHashCleared() {
		_spec.ClearField(transactionlog.FieldTransactionHash, field.TypeString)
	}
	if value, ok := tluo.mutation.Metadata(); ok {
		_spec.SetField(transactionlog.FieldMetadata, field.TypeJSON, value)
	}
	_node = &TransactionLog{config: tluo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, tluo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{transactionlog.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	tluo.mutation.done = true
	return _node, nil
}
