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
	"github.com/paycrest/protocol/ent/paymentorder"
	"github.com/paycrest/protocol/ent/paymentorderrecipient"
	"github.com/paycrest/protocol/ent/receiveaddress"
	"github.com/paycrest/protocol/ent/senderprofile"
	"github.com/paycrest/protocol/ent/token"
	"github.com/shopspring/decimal"
)

// PaymentOrderCreate is the builder for creating a PaymentOrder entity.
type PaymentOrderCreate struct {
	config
	mutation *PaymentOrderMutation
	hooks    []Hook
}

// SetCreatedAt sets the "created_at" field.
func (poc *PaymentOrderCreate) SetCreatedAt(t time.Time) *PaymentOrderCreate {
	poc.mutation.SetCreatedAt(t)
	return poc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (poc *PaymentOrderCreate) SetNillableCreatedAt(t *time.Time) *PaymentOrderCreate {
	if t != nil {
		poc.SetCreatedAt(*t)
	}
	return poc
}

// SetUpdatedAt sets the "updated_at" field.
func (poc *PaymentOrderCreate) SetUpdatedAt(t time.Time) *PaymentOrderCreate {
	poc.mutation.SetUpdatedAt(t)
	return poc
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (poc *PaymentOrderCreate) SetNillableUpdatedAt(t *time.Time) *PaymentOrderCreate {
	if t != nil {
		poc.SetUpdatedAt(*t)
	}
	return poc
}

// SetAmount sets the "amount" field.
func (poc *PaymentOrderCreate) SetAmount(d decimal.Decimal) *PaymentOrderCreate {
	poc.mutation.SetAmount(d)
	return poc
}

// SetAmountPaid sets the "amount_paid" field.
func (poc *PaymentOrderCreate) SetAmountPaid(d decimal.Decimal) *PaymentOrderCreate {
	poc.mutation.SetAmountPaid(d)
	return poc
}

// SetRate sets the "rate" field.
func (poc *PaymentOrderCreate) SetRate(d decimal.Decimal) *PaymentOrderCreate {
	poc.mutation.SetRate(d)
	return poc
}

// SetTxHash sets the "tx_hash" field.
func (poc *PaymentOrderCreate) SetTxHash(s string) *PaymentOrderCreate {
	poc.mutation.SetTxHash(s)
	return poc
}

// SetNillableTxHash sets the "tx_hash" field if the given value is not nil.
func (poc *PaymentOrderCreate) SetNillableTxHash(s *string) *PaymentOrderCreate {
	if s != nil {
		poc.SetTxHash(*s)
	}
	return poc
}

// SetReceiveAddressText sets the "receive_address_text" field.
func (poc *PaymentOrderCreate) SetReceiveAddressText(s string) *PaymentOrderCreate {
	poc.mutation.SetReceiveAddressText(s)
	return poc
}

// SetLabel sets the "label" field.
func (poc *PaymentOrderCreate) SetLabel(s string) *PaymentOrderCreate {
	poc.mutation.SetLabel(s)
	return poc
}

// SetStatus sets the "status" field.
func (poc *PaymentOrderCreate) SetStatus(pa paymentorder.Status) *PaymentOrderCreate {
	poc.mutation.SetStatus(pa)
	return poc
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (poc *PaymentOrderCreate) SetNillableStatus(pa *paymentorder.Status) *PaymentOrderCreate {
	if pa != nil {
		poc.SetStatus(*pa)
	}
	return poc
}

// SetID sets the "id" field.
func (poc *PaymentOrderCreate) SetID(u uuid.UUID) *PaymentOrderCreate {
	poc.mutation.SetID(u)
	return poc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (poc *PaymentOrderCreate) SetNillableID(u *uuid.UUID) *PaymentOrderCreate {
	if u != nil {
		poc.SetID(*u)
	}
	return poc
}

// SetSenderProfileID sets the "sender_profile" edge to the SenderProfile entity by ID.
func (poc *PaymentOrderCreate) SetSenderProfileID(id uuid.UUID) *PaymentOrderCreate {
	poc.mutation.SetSenderProfileID(id)
	return poc
}

// SetSenderProfile sets the "sender_profile" edge to the SenderProfile entity.
func (poc *PaymentOrderCreate) SetSenderProfile(s *SenderProfile) *PaymentOrderCreate {
	return poc.SetSenderProfileID(s.ID)
}

// SetTokenID sets the "token" edge to the Token entity by ID.
func (poc *PaymentOrderCreate) SetTokenID(id int) *PaymentOrderCreate {
	poc.mutation.SetTokenID(id)
	return poc
}

// SetToken sets the "token" edge to the Token entity.
func (poc *PaymentOrderCreate) SetToken(t *Token) *PaymentOrderCreate {
	return poc.SetTokenID(t.ID)
}

// SetReceiveAddressID sets the "receive_address" edge to the ReceiveAddress entity by ID.
func (poc *PaymentOrderCreate) SetReceiveAddressID(id int) *PaymentOrderCreate {
	poc.mutation.SetReceiveAddressID(id)
	return poc
}

// SetNillableReceiveAddressID sets the "receive_address" edge to the ReceiveAddress entity by ID if the given value is not nil.
func (poc *PaymentOrderCreate) SetNillableReceiveAddressID(id *int) *PaymentOrderCreate {
	if id != nil {
		poc = poc.SetReceiveAddressID(*id)
	}
	return poc
}

// SetReceiveAddress sets the "receive_address" edge to the ReceiveAddress entity.
func (poc *PaymentOrderCreate) SetReceiveAddress(r *ReceiveAddress) *PaymentOrderCreate {
	return poc.SetReceiveAddressID(r.ID)
}

// SetRecipientID sets the "recipient" edge to the PaymentOrderRecipient entity by ID.
func (poc *PaymentOrderCreate) SetRecipientID(id int) *PaymentOrderCreate {
	poc.mutation.SetRecipientID(id)
	return poc
}

// SetNillableRecipientID sets the "recipient" edge to the PaymentOrderRecipient entity by ID if the given value is not nil.
func (poc *PaymentOrderCreate) SetNillableRecipientID(id *int) *PaymentOrderCreate {
	if id != nil {
		poc = poc.SetRecipientID(*id)
	}
	return poc
}

// SetRecipient sets the "recipient" edge to the PaymentOrderRecipient entity.
func (poc *PaymentOrderCreate) SetRecipient(p *PaymentOrderRecipient) *PaymentOrderCreate {
	return poc.SetRecipientID(p.ID)
}

// Mutation returns the PaymentOrderMutation object of the builder.
func (poc *PaymentOrderCreate) Mutation() *PaymentOrderMutation {
	return poc.mutation
}

// Save creates the PaymentOrder in the database.
func (poc *PaymentOrderCreate) Save(ctx context.Context) (*PaymentOrder, error) {
	poc.defaults()
	return withHooks(ctx, poc.sqlSave, poc.mutation, poc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (poc *PaymentOrderCreate) SaveX(ctx context.Context) *PaymentOrder {
	v, err := poc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (poc *PaymentOrderCreate) Exec(ctx context.Context) error {
	_, err := poc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (poc *PaymentOrderCreate) ExecX(ctx context.Context) {
	if err := poc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (poc *PaymentOrderCreate) defaults() {
	if _, ok := poc.mutation.CreatedAt(); !ok {
		v := paymentorder.DefaultCreatedAt()
		poc.mutation.SetCreatedAt(v)
	}
	if _, ok := poc.mutation.UpdatedAt(); !ok {
		v := paymentorder.DefaultUpdatedAt()
		poc.mutation.SetUpdatedAt(v)
	}
	if _, ok := poc.mutation.Status(); !ok {
		v := paymentorder.DefaultStatus
		poc.mutation.SetStatus(v)
	}
	if _, ok := poc.mutation.ID(); !ok {
		v := paymentorder.DefaultID()
		poc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (poc *PaymentOrderCreate) check() error {
	if _, ok := poc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "PaymentOrder.created_at"`)}
	}
	if _, ok := poc.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "PaymentOrder.updated_at"`)}
	}
	if _, ok := poc.mutation.Amount(); !ok {
		return &ValidationError{Name: "amount", err: errors.New(`ent: missing required field "PaymentOrder.amount"`)}
	}
	if _, ok := poc.mutation.AmountPaid(); !ok {
		return &ValidationError{Name: "amount_paid", err: errors.New(`ent: missing required field "PaymentOrder.amount_paid"`)}
	}
	if _, ok := poc.mutation.Rate(); !ok {
		return &ValidationError{Name: "rate", err: errors.New(`ent: missing required field "PaymentOrder.rate"`)}
	}
	if v, ok := poc.mutation.TxHash(); ok {
		if err := paymentorder.TxHashValidator(v); err != nil {
			return &ValidationError{Name: "tx_hash", err: fmt.Errorf(`ent: validator failed for field "PaymentOrder.tx_hash": %w`, err)}
		}
	}
	if _, ok := poc.mutation.ReceiveAddressText(); !ok {
		return &ValidationError{Name: "receive_address_text", err: errors.New(`ent: missing required field "PaymentOrder.receive_address_text"`)}
	}
	if v, ok := poc.mutation.ReceiveAddressText(); ok {
		if err := paymentorder.ReceiveAddressTextValidator(v); err != nil {
			return &ValidationError{Name: "receive_address_text", err: fmt.Errorf(`ent: validator failed for field "PaymentOrder.receive_address_text": %w`, err)}
		}
	}
	if _, ok := poc.mutation.Label(); !ok {
		return &ValidationError{Name: "label", err: errors.New(`ent: missing required field "PaymentOrder.label"`)}
	}
	if _, ok := poc.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "PaymentOrder.status"`)}
	}
	if v, ok := poc.mutation.Status(); ok {
		if err := paymentorder.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "PaymentOrder.status": %w`, err)}
		}
	}
	if _, ok := poc.mutation.SenderProfileID(); !ok {
		return &ValidationError{Name: "sender_profile", err: errors.New(`ent: missing required edge "PaymentOrder.sender_profile"`)}
	}
	if _, ok := poc.mutation.TokenID(); !ok {
		return &ValidationError{Name: "token", err: errors.New(`ent: missing required edge "PaymentOrder.token"`)}
	}
	return nil
}

func (poc *PaymentOrderCreate) sqlSave(ctx context.Context) (*PaymentOrder, error) {
	if err := poc.check(); err != nil {
		return nil, err
	}
	_node, _spec := poc.createSpec()
	if err := sqlgraph.CreateNode(ctx, poc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	poc.mutation.id = &_node.ID
	poc.mutation.done = true
	return _node, nil
}

func (poc *PaymentOrderCreate) createSpec() (*PaymentOrder, *sqlgraph.CreateSpec) {
	var (
		_node = &PaymentOrder{config: poc.config}
		_spec = sqlgraph.NewCreateSpec(paymentorder.Table, sqlgraph.NewFieldSpec(paymentorder.FieldID, field.TypeUUID))
	)
	if id, ok := poc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := poc.mutation.CreatedAt(); ok {
		_spec.SetField(paymentorder.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := poc.mutation.UpdatedAt(); ok {
		_spec.SetField(paymentorder.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := poc.mutation.Amount(); ok {
		_spec.SetField(paymentorder.FieldAmount, field.TypeFloat64, value)
		_node.Amount = value
	}
	if value, ok := poc.mutation.AmountPaid(); ok {
		_spec.SetField(paymentorder.FieldAmountPaid, field.TypeFloat64, value)
		_node.AmountPaid = value
	}
	if value, ok := poc.mutation.Rate(); ok {
		_spec.SetField(paymentorder.FieldRate, field.TypeFloat64, value)
		_node.Rate = value
	}
	if value, ok := poc.mutation.TxHash(); ok {
		_spec.SetField(paymentorder.FieldTxHash, field.TypeString, value)
		_node.TxHash = value
	}
	if value, ok := poc.mutation.ReceiveAddressText(); ok {
		_spec.SetField(paymentorder.FieldReceiveAddressText, field.TypeString, value)
		_node.ReceiveAddressText = value
	}
	if value, ok := poc.mutation.Label(); ok {
		_spec.SetField(paymentorder.FieldLabel, field.TypeString, value)
		_node.Label = value
	}
	if value, ok := poc.mutation.Status(); ok {
		_spec.SetField(paymentorder.FieldStatus, field.TypeEnum, value)
		_node.Status = value
	}
	if nodes := poc.mutation.SenderProfileIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   paymentorder.SenderProfileTable,
			Columns: []string{paymentorder.SenderProfileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(senderprofile.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.sender_profile_payment_orders = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := poc.mutation.TokenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   paymentorder.TokenTable,
			Columns: []string{paymentorder.TokenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(token.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.token_payment_orders = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := poc.mutation.ReceiveAddressIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   paymentorder.ReceiveAddressTable,
			Columns: []string{paymentorder.ReceiveAddressColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(receiveaddress.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := poc.mutation.RecipientIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   paymentorder.RecipientTable,
			Columns: []string{paymentorder.RecipientColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(paymentorderrecipient.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// PaymentOrderCreateBulk is the builder for creating many PaymentOrder entities in bulk.
type PaymentOrderCreateBulk struct {
	config
	builders []*PaymentOrderCreate
}

// Save creates the PaymentOrder entities in the database.
func (pocb *PaymentOrderCreateBulk) Save(ctx context.Context) ([]*PaymentOrder, error) {
	specs := make([]*sqlgraph.CreateSpec, len(pocb.builders))
	nodes := make([]*PaymentOrder, len(pocb.builders))
	mutators := make([]Mutator, len(pocb.builders))
	for i := range pocb.builders {
		func(i int, root context.Context) {
			builder := pocb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*PaymentOrderMutation)
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
					_, err = mutators[i+1].Mutate(root, pocb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, pocb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
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
		if _, err := mutators[0].Mutate(ctx, pocb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (pocb *PaymentOrderCreateBulk) SaveX(ctx context.Context) []*PaymentOrder {
	v, err := pocb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (pocb *PaymentOrderCreateBulk) Exec(ctx context.Context) error {
	_, err := pocb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pocb *PaymentOrderCreateBulk) ExecX(ctx context.Context) {
	if err := pocb.Exec(ctx); err != nil {
		panic(err)
	}
}
