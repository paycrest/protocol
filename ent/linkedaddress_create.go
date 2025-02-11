// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/paycrest/aggregator/ent/linkedaddress"
	"github.com/paycrest/aggregator/ent/paymentorder"
)

// LinkedAddressCreate is the builder for creating a LinkedAddress entity.
type LinkedAddressCreate struct {
	config
	mutation *LinkedAddressMutation
	hooks    []Hook
}

// SetCreatedAt sets the "created_at" field.
func (lac *LinkedAddressCreate) SetCreatedAt(t time.Time) *LinkedAddressCreate {
	lac.mutation.SetCreatedAt(t)
	return lac
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (lac *LinkedAddressCreate) SetNillableCreatedAt(t *time.Time) *LinkedAddressCreate {
	if t != nil {
		lac.SetCreatedAt(*t)
	}
	return lac
}

// SetUpdatedAt sets the "updated_at" field.
func (lac *LinkedAddressCreate) SetUpdatedAt(t time.Time) *LinkedAddressCreate {
	lac.mutation.SetUpdatedAt(t)
	return lac
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (lac *LinkedAddressCreate) SetNillableUpdatedAt(t *time.Time) *LinkedAddressCreate {
	if t != nil {
		lac.SetUpdatedAt(*t)
	}
	return lac
}

// SetAddress sets the "address" field.
func (lac *LinkedAddressCreate) SetAddress(s string) *LinkedAddressCreate {
	lac.mutation.SetAddress(s)
	return lac
}

// SetSalt sets the "salt" field.
func (lac *LinkedAddressCreate) SetSalt(b []byte) *LinkedAddressCreate {
	lac.mutation.SetSalt(b)
	return lac
}

// SetInstitution sets the "institution" field.
func (lac *LinkedAddressCreate) SetInstitution(s string) *LinkedAddressCreate {
	lac.mutation.SetInstitution(s)
	return lac
}

// SetAccountIdentifier sets the "account_identifier" field.
func (lac *LinkedAddressCreate) SetAccountIdentifier(s string) *LinkedAddressCreate {
	lac.mutation.SetAccountIdentifier(s)
	return lac
}

// SetAccountName sets the "account_name" field.
func (lac *LinkedAddressCreate) SetAccountName(s string) *LinkedAddressCreate {
	lac.mutation.SetAccountName(s)
	return lac
}

// SetOwnerAddress sets the "owner_address" field.
func (lac *LinkedAddressCreate) SetOwnerAddress(s string) *LinkedAddressCreate {
	lac.mutation.SetOwnerAddress(s)
	return lac
}

// SetLastIndexedBlock sets the "last_indexed_block" field.
func (lac *LinkedAddressCreate) SetLastIndexedBlock(i int64) *LinkedAddressCreate {
	lac.mutation.SetLastIndexedBlock(i)
	return lac
}

// SetNillableLastIndexedBlock sets the "last_indexed_block" field if the given value is not nil.
func (lac *LinkedAddressCreate) SetNillableLastIndexedBlock(i *int64) *LinkedAddressCreate {
	if i != nil {
		lac.SetLastIndexedBlock(*i)
	}
	return lac
}

// SetTxHash sets the "tx_hash" field.
func (lac *LinkedAddressCreate) SetTxHash(s string) *LinkedAddressCreate {
	lac.mutation.SetTxHash(s)
	return lac
}

// SetNillableTxHash sets the "tx_hash" field if the given value is not nil.
func (lac *LinkedAddressCreate) SetNillableTxHash(s *string) *LinkedAddressCreate {
	if s != nil {
		lac.SetTxHash(*s)
	}
	return lac
}

// AddPaymentOrderIDs adds the "payment_orders" edge to the PaymentOrder entity by IDs.
func (lac *LinkedAddressCreate) AddPaymentOrderIDs(ids ...uuid.UUID) *LinkedAddressCreate {
	lac.mutation.AddPaymentOrderIDs(ids...)
	return lac
}

// AddPaymentOrders adds the "payment_orders" edges to the PaymentOrder entity.
func (lac *LinkedAddressCreate) AddPaymentOrders(p ...*PaymentOrder) *LinkedAddressCreate {
	ids := make([]uuid.UUID, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return lac.AddPaymentOrderIDs(ids...)
}

// Mutation returns the LinkedAddressMutation object of the builder.
func (lac *LinkedAddressCreate) Mutation() *LinkedAddressMutation {
	return lac.mutation
}

// Save creates the LinkedAddress in the database.
func (lac *LinkedAddressCreate) Save(ctx context.Context) (*LinkedAddress, error) {
	lac.defaults()
	return withHooks(ctx, lac.sqlSave, lac.mutation, lac.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (lac *LinkedAddressCreate) SaveX(ctx context.Context) *LinkedAddress {
	v, err := lac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (lac *LinkedAddressCreate) Exec(ctx context.Context) error {
	_, err := lac.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lac *LinkedAddressCreate) ExecX(ctx context.Context) {
	if err := lac.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (lac *LinkedAddressCreate) defaults() {
	if _, ok := lac.mutation.CreatedAt(); !ok {
		v := linkedaddress.DefaultCreatedAt()
		lac.mutation.SetCreatedAt(v)
	}
	if _, ok := lac.mutation.UpdatedAt(); !ok {
		v := linkedaddress.DefaultUpdatedAt()
		lac.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (lac *LinkedAddressCreate) check() error {
	if _, ok := lac.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "LinkedAddress.created_at"`)}
	}
	if _, ok := lac.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "LinkedAddress.updated_at"`)}
	}
	if _, ok := lac.mutation.Address(); !ok {
		return &ValidationError{Name: "address", err: errors.New(`ent: missing required field "LinkedAddress.address"`)}
	}
	if _, ok := lac.mutation.Salt(); !ok {
		return &ValidationError{Name: "salt", err: errors.New(`ent: missing required field "LinkedAddress.salt"`)}
	}
	if _, ok := lac.mutation.Institution(); !ok {
		return &ValidationError{Name: "institution", err: errors.New(`ent: missing required field "LinkedAddress.institution"`)}
	}
	if _, ok := lac.mutation.AccountIdentifier(); !ok {
		return &ValidationError{Name: "account_identifier", err: errors.New(`ent: missing required field "LinkedAddress.account_identifier"`)}
	}
	if _, ok := lac.mutation.AccountName(); !ok {
		return &ValidationError{Name: "account_name", err: errors.New(`ent: missing required field "LinkedAddress.account_name"`)}
	}
	if _, ok := lac.mutation.OwnerAddress(); !ok {
		return &ValidationError{Name: "owner_address", err: errors.New(`ent: missing required field "LinkedAddress.owner_address"`)}
	}
	if v, ok := lac.mutation.TxHash(); ok {
		if err := linkedaddress.TxHashValidator(v); err != nil {
			return &ValidationError{Name: "tx_hash", err: fmt.Errorf(`ent: validator failed for field "LinkedAddress.tx_hash": %w`, err)}
		}
	}
	return nil
}

func (lac *LinkedAddressCreate) sqlSave(ctx context.Context) (*LinkedAddress, error) {
	if err := lac.check(); err != nil {
		return nil, err
	}
	_node, _spec := lac.createSpec()
	if err := sqlgraph.CreateNode(ctx, lac.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	lac.mutation.id = &_node.ID
	lac.mutation.done = true
	return _node, nil
}

func (lac *LinkedAddressCreate) createSpec() (*LinkedAddress, *sqlgraph.CreateSpec) {
	var (
		_node = &LinkedAddress{config: lac.config}
		_spec = sqlgraph.NewCreateSpec(linkedaddress.Table, sqlgraph.NewFieldSpec(linkedaddress.FieldID, field.TypeInt))
	)
	if value, ok := lac.mutation.CreatedAt(); ok {
		_spec.SetField(linkedaddress.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := lac.mutation.UpdatedAt(); ok {
		_spec.SetField(linkedaddress.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := lac.mutation.Address(); ok {
		_spec.SetField(linkedaddress.FieldAddress, field.TypeString, value)
		_node.Address = value
	}
	if value, ok := lac.mutation.Salt(); ok {
		_spec.SetField(linkedaddress.FieldSalt, field.TypeBytes, value)
		_node.Salt = value
	}
	if value, ok := lac.mutation.Institution(); ok {
		_spec.SetField(linkedaddress.FieldInstitution, field.TypeString, value)
		_node.Institution = value
	}
	if value, ok := lac.mutation.AccountIdentifier(); ok {
		_spec.SetField(linkedaddress.FieldAccountIdentifier, field.TypeString, value)
		_node.AccountIdentifier = value
	}
	if value, ok := lac.mutation.AccountName(); ok {
		_spec.SetField(linkedaddress.FieldAccountName, field.TypeString, value)
		_node.AccountName = value
	}
	if value, ok := lac.mutation.OwnerAddress(); ok {
		_spec.SetField(linkedaddress.FieldOwnerAddress, field.TypeString, value)
		_node.OwnerAddress = value
	}
	if value, ok := lac.mutation.LastIndexedBlock(); ok {
		_spec.SetField(linkedaddress.FieldLastIndexedBlock, field.TypeInt64, value)
		_node.LastIndexedBlock = value
	}
	if value, ok := lac.mutation.TxHash(); ok {
		_spec.SetField(linkedaddress.FieldTxHash, field.TypeString, value)
		_node.TxHash = value
	}
	if nodes := lac.mutation.PaymentOrdersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   linkedaddress.PaymentOrdersTable,
			Columns: []string{linkedaddress.PaymentOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(paymentorder.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// LinkedAddressCreateBulk is the builder for creating many LinkedAddress entities in bulk.
type LinkedAddressCreateBulk struct {
	config
	err      error
	builders []*LinkedAddressCreate
}

// Save creates the LinkedAddress entities in the database.
func (lacb *LinkedAddressCreateBulk) Save(ctx context.Context) ([]*LinkedAddress, error) {
	if lacb.err != nil {
		return nil, lacb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(lacb.builders))
	nodes := make([]*LinkedAddress, len(lacb.builders))
	mutators := make([]Mutator, len(lacb.builders))
	for i := range lacb.builders {
		func(i int, root context.Context) {
			builder := lacb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*LinkedAddressMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, lacb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, lacb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				if specs[i].ID.Value != nil {
					id := specs[i].ID.Value.(int64)
					nodes[i].ID = int(id)
				}
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, lacb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (lacb *LinkedAddressCreateBulk) SaveX(ctx context.Context) []*LinkedAddress {
	v, err := lacb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (lacb *LinkedAddressCreateBulk) Exec(ctx context.Context) error {
	_, err := lacb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lacb *LinkedAddressCreateBulk) ExecX(ctx context.Context) {
	if err := lacb.Exec(ctx); err != nil {
		panic(err)
	}
}
