// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/paycrest/protocol/ent/paymentorder"
	"github.com/paycrest/protocol/ent/paymentorderrecipient"
	"github.com/paycrest/protocol/ent/predicate"
	"github.com/paycrest/protocol/ent/receiveaddress"
	"github.com/paycrest/protocol/ent/senderprofile"
	"github.com/paycrest/protocol/ent/token"
	"github.com/shopspring/decimal"
)

// PaymentOrderUpdate is the builder for updating PaymentOrder entities.
type PaymentOrderUpdate struct {
	config
	hooks    []Hook
	mutation *PaymentOrderMutation
}

// Where appends a list predicates to the PaymentOrderUpdate builder.
func (pou *PaymentOrderUpdate) Where(ps ...predicate.PaymentOrder) *PaymentOrderUpdate {
	pou.mutation.Where(ps...)
	return pou
}

// SetUpdatedAt sets the "updated_at" field.
func (pou *PaymentOrderUpdate) SetUpdatedAt(t time.Time) *PaymentOrderUpdate {
	pou.mutation.SetUpdatedAt(t)
	return pou
}

// SetAmount sets the "amount" field.
func (pou *PaymentOrderUpdate) SetAmount(d decimal.Decimal) *PaymentOrderUpdate {
	pou.mutation.ResetAmount()
	pou.mutation.SetAmount(d)
	return pou
}

// AddAmount adds d to the "amount" field.
func (pou *PaymentOrderUpdate) AddAmount(d decimal.Decimal) *PaymentOrderUpdate {
	pou.mutation.AddAmount(d)
	return pou
}

// SetAmountPaid sets the "amount_paid" field.
func (pou *PaymentOrderUpdate) SetAmountPaid(d decimal.Decimal) *PaymentOrderUpdate {
	pou.mutation.ResetAmountPaid()
	pou.mutation.SetAmountPaid(d)
	return pou
}

// AddAmountPaid adds d to the "amount_paid" field.
func (pou *PaymentOrderUpdate) AddAmountPaid(d decimal.Decimal) *PaymentOrderUpdate {
	pou.mutation.AddAmountPaid(d)
	return pou
}

// SetRate sets the "rate" field.
func (pou *PaymentOrderUpdate) SetRate(d decimal.Decimal) *PaymentOrderUpdate {
	pou.mutation.ResetRate()
	pou.mutation.SetRate(d)
	return pou
}

// AddRate adds d to the "rate" field.
func (pou *PaymentOrderUpdate) AddRate(d decimal.Decimal) *PaymentOrderUpdate {
	pou.mutation.AddRate(d)
	return pou
}

// SetTxHash sets the "tx_hash" field.
func (pou *PaymentOrderUpdate) SetTxHash(s string) *PaymentOrderUpdate {
	pou.mutation.SetTxHash(s)
	return pou
}

// SetNillableTxHash sets the "tx_hash" field if the given value is not nil.
func (pou *PaymentOrderUpdate) SetNillableTxHash(s *string) *PaymentOrderUpdate {
	if s != nil {
		pou.SetTxHash(*s)
	}
	return pou
}

// ClearTxHash clears the value of the "tx_hash" field.
func (pou *PaymentOrderUpdate) ClearTxHash() *PaymentOrderUpdate {
	pou.mutation.ClearTxHash()
	return pou
}

// SetReceiveAddressText sets the "receive_address_text" field.
func (pou *PaymentOrderUpdate) SetReceiveAddressText(s string) *PaymentOrderUpdate {
	pou.mutation.SetReceiveAddressText(s)
	return pou
}

// SetLabel sets the "label" field.
func (pou *PaymentOrderUpdate) SetLabel(s string) *PaymentOrderUpdate {
	pou.mutation.SetLabel(s)
	return pou
}

// SetStatus sets the "status" field.
func (pou *PaymentOrderUpdate) SetStatus(pa paymentorder.Status) *PaymentOrderUpdate {
	pou.mutation.SetStatus(pa)
	return pou
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (pou *PaymentOrderUpdate) SetNillableStatus(pa *paymentorder.Status) *PaymentOrderUpdate {
	if pa != nil {
		pou.SetStatus(*pa)
	}
	return pou
}

// SetSenderProfileID sets the "sender_profile" edge to the SenderProfile entity by ID.
func (pou *PaymentOrderUpdate) SetSenderProfileID(id uuid.UUID) *PaymentOrderUpdate {
	pou.mutation.SetSenderProfileID(id)
	return pou
}

// SetSenderProfile sets the "sender_profile" edge to the SenderProfile entity.
func (pou *PaymentOrderUpdate) SetSenderProfile(s *SenderProfile) *PaymentOrderUpdate {
	return pou.SetSenderProfileID(s.ID)
}

// SetTokenID sets the "token" edge to the Token entity by ID.
func (pou *PaymentOrderUpdate) SetTokenID(id int) *PaymentOrderUpdate {
	pou.mutation.SetTokenID(id)
	return pou
}

// SetToken sets the "token" edge to the Token entity.
func (pou *PaymentOrderUpdate) SetToken(t *Token) *PaymentOrderUpdate {
	return pou.SetTokenID(t.ID)
}

// SetReceiveAddressID sets the "receive_address" edge to the ReceiveAddress entity by ID.
func (pou *PaymentOrderUpdate) SetReceiveAddressID(id int) *PaymentOrderUpdate {
	pou.mutation.SetReceiveAddressID(id)
	return pou
}

// SetNillableReceiveAddressID sets the "receive_address" edge to the ReceiveAddress entity by ID if the given value is not nil.
func (pou *PaymentOrderUpdate) SetNillableReceiveAddressID(id *int) *PaymentOrderUpdate {
	if id != nil {
		pou = pou.SetReceiveAddressID(*id)
	}
	return pou
}

// SetReceiveAddress sets the "receive_address" edge to the ReceiveAddress entity.
func (pou *PaymentOrderUpdate) SetReceiveAddress(r *ReceiveAddress) *PaymentOrderUpdate {
	return pou.SetReceiveAddressID(r.ID)
}

// SetRecipientID sets the "recipient" edge to the PaymentOrderRecipient entity by ID.
func (pou *PaymentOrderUpdate) SetRecipientID(id int) *PaymentOrderUpdate {
	pou.mutation.SetRecipientID(id)
	return pou
}

// SetNillableRecipientID sets the "recipient" edge to the PaymentOrderRecipient entity by ID if the given value is not nil.
func (pou *PaymentOrderUpdate) SetNillableRecipientID(id *int) *PaymentOrderUpdate {
	if id != nil {
		pou = pou.SetRecipientID(*id)
	}
	return pou
}

// SetRecipient sets the "recipient" edge to the PaymentOrderRecipient entity.
func (pou *PaymentOrderUpdate) SetRecipient(p *PaymentOrderRecipient) *PaymentOrderUpdate {
	return pou.SetRecipientID(p.ID)
}

// Mutation returns the PaymentOrderMutation object of the builder.
func (pou *PaymentOrderUpdate) Mutation() *PaymentOrderMutation {
	return pou.mutation
}

// ClearSenderProfile clears the "sender_profile" edge to the SenderProfile entity.
func (pou *PaymentOrderUpdate) ClearSenderProfile() *PaymentOrderUpdate {
	pou.mutation.ClearSenderProfile()
	return pou
}

// ClearToken clears the "token" edge to the Token entity.
func (pou *PaymentOrderUpdate) ClearToken() *PaymentOrderUpdate {
	pou.mutation.ClearToken()
	return pou
}

// ClearReceiveAddress clears the "receive_address" edge to the ReceiveAddress entity.
func (pou *PaymentOrderUpdate) ClearReceiveAddress() *PaymentOrderUpdate {
	pou.mutation.ClearReceiveAddress()
	return pou
}

// ClearRecipient clears the "recipient" edge to the PaymentOrderRecipient entity.
func (pou *PaymentOrderUpdate) ClearRecipient() *PaymentOrderUpdate {
	pou.mutation.ClearRecipient()
	return pou
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (pou *PaymentOrderUpdate) Save(ctx context.Context) (int, error) {
	pou.defaults()
	return withHooks(ctx, pou.sqlSave, pou.mutation, pou.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pou *PaymentOrderUpdate) SaveX(ctx context.Context) int {
	affected, err := pou.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (pou *PaymentOrderUpdate) Exec(ctx context.Context) error {
	_, err := pou.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pou *PaymentOrderUpdate) ExecX(ctx context.Context) {
	if err := pou.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pou *PaymentOrderUpdate) defaults() {
	if _, ok := pou.mutation.UpdatedAt(); !ok {
		v := paymentorder.UpdateDefaultUpdatedAt()
		pou.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pou *PaymentOrderUpdate) check() error {
	if v, ok := pou.mutation.TxHash(); ok {
		if err := paymentorder.TxHashValidator(v); err != nil {
			return &ValidationError{Name: "tx_hash", err: fmt.Errorf(`ent: validator failed for field "PaymentOrder.tx_hash": %w`, err)}
		}
	}
	if v, ok := pou.mutation.ReceiveAddressText(); ok {
		if err := paymentorder.ReceiveAddressTextValidator(v); err != nil {
			return &ValidationError{Name: "receive_address_text", err: fmt.Errorf(`ent: validator failed for field "PaymentOrder.receive_address_text": %w`, err)}
		}
	}
	if v, ok := pou.mutation.Status(); ok {
		if err := paymentorder.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "PaymentOrder.status": %w`, err)}
		}
	}
	if _, ok := pou.mutation.SenderProfileID(); pou.mutation.SenderProfileCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "PaymentOrder.sender_profile"`)
	}
	if _, ok := pou.mutation.TokenID(); pou.mutation.TokenCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "PaymentOrder.token"`)
	}
	return nil
}

func (pou *PaymentOrderUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := pou.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(paymentorder.Table, paymentorder.Columns, sqlgraph.NewFieldSpec(paymentorder.FieldID, field.TypeUUID))
	if ps := pou.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pou.mutation.UpdatedAt(); ok {
		_spec.SetField(paymentorder.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := pou.mutation.Amount(); ok {
		_spec.SetField(paymentorder.FieldAmount, field.TypeFloat64, value)
	}
	if value, ok := pou.mutation.AddedAmount(); ok {
		_spec.AddField(paymentorder.FieldAmount, field.TypeFloat64, value)
	}
	if value, ok := pou.mutation.AmountPaid(); ok {
		_spec.SetField(paymentorder.FieldAmountPaid, field.TypeFloat64, value)
	}
	if value, ok := pou.mutation.AddedAmountPaid(); ok {
		_spec.AddField(paymentorder.FieldAmountPaid, field.TypeFloat64, value)
	}
	if value, ok := pou.mutation.Rate(); ok {
		_spec.SetField(paymentorder.FieldRate, field.TypeFloat64, value)
	}
	if value, ok := pou.mutation.AddedRate(); ok {
		_spec.AddField(paymentorder.FieldRate, field.TypeFloat64, value)
	}
	if value, ok := pou.mutation.TxHash(); ok {
		_spec.SetField(paymentorder.FieldTxHash, field.TypeString, value)
	}
	if pou.mutation.TxHashCleared() {
		_spec.ClearField(paymentorder.FieldTxHash, field.TypeString)
	}
	if value, ok := pou.mutation.ReceiveAddressText(); ok {
		_spec.SetField(paymentorder.FieldReceiveAddressText, field.TypeString, value)
	}
	if value, ok := pou.mutation.Label(); ok {
		_spec.SetField(paymentorder.FieldLabel, field.TypeString, value)
	}
	if value, ok := pou.mutation.Status(); ok {
		_spec.SetField(paymentorder.FieldStatus, field.TypeEnum, value)
	}
	if pou.mutation.SenderProfileCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pou.mutation.SenderProfileIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pou.mutation.TokenCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pou.mutation.TokenIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pou.mutation.ReceiveAddressCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pou.mutation.ReceiveAddressIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pou.mutation.RecipientCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pou.mutation.RecipientIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, pou.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{paymentorder.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	pou.mutation.done = true
	return n, nil
}

// PaymentOrderUpdateOne is the builder for updating a single PaymentOrder entity.
type PaymentOrderUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *PaymentOrderMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (pouo *PaymentOrderUpdateOne) SetUpdatedAt(t time.Time) *PaymentOrderUpdateOne {
	pouo.mutation.SetUpdatedAt(t)
	return pouo
}

// SetAmount sets the "amount" field.
func (pouo *PaymentOrderUpdateOne) SetAmount(d decimal.Decimal) *PaymentOrderUpdateOne {
	pouo.mutation.ResetAmount()
	pouo.mutation.SetAmount(d)
	return pouo
}

// AddAmount adds d to the "amount" field.
func (pouo *PaymentOrderUpdateOne) AddAmount(d decimal.Decimal) *PaymentOrderUpdateOne {
	pouo.mutation.AddAmount(d)
	return pouo
}

// SetAmountPaid sets the "amount_paid" field.
func (pouo *PaymentOrderUpdateOne) SetAmountPaid(d decimal.Decimal) *PaymentOrderUpdateOne {
	pouo.mutation.ResetAmountPaid()
	pouo.mutation.SetAmountPaid(d)
	return pouo
}

// AddAmountPaid adds d to the "amount_paid" field.
func (pouo *PaymentOrderUpdateOne) AddAmountPaid(d decimal.Decimal) *PaymentOrderUpdateOne {
	pouo.mutation.AddAmountPaid(d)
	return pouo
}

// SetRate sets the "rate" field.
func (pouo *PaymentOrderUpdateOne) SetRate(d decimal.Decimal) *PaymentOrderUpdateOne {
	pouo.mutation.ResetRate()
	pouo.mutation.SetRate(d)
	return pouo
}

// AddRate adds d to the "rate" field.
func (pouo *PaymentOrderUpdateOne) AddRate(d decimal.Decimal) *PaymentOrderUpdateOne {
	pouo.mutation.AddRate(d)
	return pouo
}

// SetTxHash sets the "tx_hash" field.
func (pouo *PaymentOrderUpdateOne) SetTxHash(s string) *PaymentOrderUpdateOne {
	pouo.mutation.SetTxHash(s)
	return pouo
}

// SetNillableTxHash sets the "tx_hash" field if the given value is not nil.
func (pouo *PaymentOrderUpdateOne) SetNillableTxHash(s *string) *PaymentOrderUpdateOne {
	if s != nil {
		pouo.SetTxHash(*s)
	}
	return pouo
}

// ClearTxHash clears the value of the "tx_hash" field.
func (pouo *PaymentOrderUpdateOne) ClearTxHash() *PaymentOrderUpdateOne {
	pouo.mutation.ClearTxHash()
	return pouo
}

// SetReceiveAddressText sets the "receive_address_text" field.
func (pouo *PaymentOrderUpdateOne) SetReceiveAddressText(s string) *PaymentOrderUpdateOne {
	pouo.mutation.SetReceiveAddressText(s)
	return pouo
}

// SetLabel sets the "label" field.
func (pouo *PaymentOrderUpdateOne) SetLabel(s string) *PaymentOrderUpdateOne {
	pouo.mutation.SetLabel(s)
	return pouo
}

// SetStatus sets the "status" field.
func (pouo *PaymentOrderUpdateOne) SetStatus(pa paymentorder.Status) *PaymentOrderUpdateOne {
	pouo.mutation.SetStatus(pa)
	return pouo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (pouo *PaymentOrderUpdateOne) SetNillableStatus(pa *paymentorder.Status) *PaymentOrderUpdateOne {
	if pa != nil {
		pouo.SetStatus(*pa)
	}
	return pouo
}

// SetSenderProfileID sets the "sender_profile" edge to the SenderProfile entity by ID.
func (pouo *PaymentOrderUpdateOne) SetSenderProfileID(id uuid.UUID) *PaymentOrderUpdateOne {
	pouo.mutation.SetSenderProfileID(id)
	return pouo
}

// SetSenderProfile sets the "sender_profile" edge to the SenderProfile entity.
func (pouo *PaymentOrderUpdateOne) SetSenderProfile(s *SenderProfile) *PaymentOrderUpdateOne {
	return pouo.SetSenderProfileID(s.ID)
}

// SetTokenID sets the "token" edge to the Token entity by ID.
func (pouo *PaymentOrderUpdateOne) SetTokenID(id int) *PaymentOrderUpdateOne {
	pouo.mutation.SetTokenID(id)
	return pouo
}

// SetToken sets the "token" edge to the Token entity.
func (pouo *PaymentOrderUpdateOne) SetToken(t *Token) *PaymentOrderUpdateOne {
	return pouo.SetTokenID(t.ID)
}

// SetReceiveAddressID sets the "receive_address" edge to the ReceiveAddress entity by ID.
func (pouo *PaymentOrderUpdateOne) SetReceiveAddressID(id int) *PaymentOrderUpdateOne {
	pouo.mutation.SetReceiveAddressID(id)
	return pouo
}

// SetNillableReceiveAddressID sets the "receive_address" edge to the ReceiveAddress entity by ID if the given value is not nil.
func (pouo *PaymentOrderUpdateOne) SetNillableReceiveAddressID(id *int) *PaymentOrderUpdateOne {
	if id != nil {
		pouo = pouo.SetReceiveAddressID(*id)
	}
	return pouo
}

// SetReceiveAddress sets the "receive_address" edge to the ReceiveAddress entity.
func (pouo *PaymentOrderUpdateOne) SetReceiveAddress(r *ReceiveAddress) *PaymentOrderUpdateOne {
	return pouo.SetReceiveAddressID(r.ID)
}

// SetRecipientID sets the "recipient" edge to the PaymentOrderRecipient entity by ID.
func (pouo *PaymentOrderUpdateOne) SetRecipientID(id int) *PaymentOrderUpdateOne {
	pouo.mutation.SetRecipientID(id)
	return pouo
}

// SetNillableRecipientID sets the "recipient" edge to the PaymentOrderRecipient entity by ID if the given value is not nil.
func (pouo *PaymentOrderUpdateOne) SetNillableRecipientID(id *int) *PaymentOrderUpdateOne {
	if id != nil {
		pouo = pouo.SetRecipientID(*id)
	}
	return pouo
}

// SetRecipient sets the "recipient" edge to the PaymentOrderRecipient entity.
func (pouo *PaymentOrderUpdateOne) SetRecipient(p *PaymentOrderRecipient) *PaymentOrderUpdateOne {
	return pouo.SetRecipientID(p.ID)
}

// Mutation returns the PaymentOrderMutation object of the builder.
func (pouo *PaymentOrderUpdateOne) Mutation() *PaymentOrderMutation {
	return pouo.mutation
}

// ClearSenderProfile clears the "sender_profile" edge to the SenderProfile entity.
func (pouo *PaymentOrderUpdateOne) ClearSenderProfile() *PaymentOrderUpdateOne {
	pouo.mutation.ClearSenderProfile()
	return pouo
}

// ClearToken clears the "token" edge to the Token entity.
func (pouo *PaymentOrderUpdateOne) ClearToken() *PaymentOrderUpdateOne {
	pouo.mutation.ClearToken()
	return pouo
}

// ClearReceiveAddress clears the "receive_address" edge to the ReceiveAddress entity.
func (pouo *PaymentOrderUpdateOne) ClearReceiveAddress() *PaymentOrderUpdateOne {
	pouo.mutation.ClearReceiveAddress()
	return pouo
}

// ClearRecipient clears the "recipient" edge to the PaymentOrderRecipient entity.
func (pouo *PaymentOrderUpdateOne) ClearRecipient() *PaymentOrderUpdateOne {
	pouo.mutation.ClearRecipient()
	return pouo
}

// Where appends a list predicates to the PaymentOrderUpdate builder.
func (pouo *PaymentOrderUpdateOne) Where(ps ...predicate.PaymentOrder) *PaymentOrderUpdateOne {
	pouo.mutation.Where(ps...)
	return pouo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (pouo *PaymentOrderUpdateOne) Select(field string, fields ...string) *PaymentOrderUpdateOne {
	pouo.fields = append([]string{field}, fields...)
	return pouo
}

// Save executes the query and returns the updated PaymentOrder entity.
func (pouo *PaymentOrderUpdateOne) Save(ctx context.Context) (*PaymentOrder, error) {
	pouo.defaults()
	return withHooks(ctx, pouo.sqlSave, pouo.mutation, pouo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (pouo *PaymentOrderUpdateOne) SaveX(ctx context.Context) *PaymentOrder {
	node, err := pouo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (pouo *PaymentOrderUpdateOne) Exec(ctx context.Context) error {
	_, err := pouo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (pouo *PaymentOrderUpdateOne) ExecX(ctx context.Context) {
	if err := pouo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (pouo *PaymentOrderUpdateOne) defaults() {
	if _, ok := pouo.mutation.UpdatedAt(); !ok {
		v := paymentorder.UpdateDefaultUpdatedAt()
		pouo.mutation.SetUpdatedAt(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (pouo *PaymentOrderUpdateOne) check() error {
	if v, ok := pouo.mutation.TxHash(); ok {
		if err := paymentorder.TxHashValidator(v); err != nil {
			return &ValidationError{Name: "tx_hash", err: fmt.Errorf(`ent: validator failed for field "PaymentOrder.tx_hash": %w`, err)}
		}
	}
	if v, ok := pouo.mutation.ReceiveAddressText(); ok {
		if err := paymentorder.ReceiveAddressTextValidator(v); err != nil {
			return &ValidationError{Name: "receive_address_text", err: fmt.Errorf(`ent: validator failed for field "PaymentOrder.receive_address_text": %w`, err)}
		}
	}
	if v, ok := pouo.mutation.Status(); ok {
		if err := paymentorder.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "PaymentOrder.status": %w`, err)}
		}
	}
	if _, ok := pouo.mutation.SenderProfileID(); pouo.mutation.SenderProfileCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "PaymentOrder.sender_profile"`)
	}
	if _, ok := pouo.mutation.TokenID(); pouo.mutation.TokenCleared() && !ok {
		return errors.New(`ent: clearing a required unique edge "PaymentOrder.token"`)
	}
	return nil
}

func (pouo *PaymentOrderUpdateOne) sqlSave(ctx context.Context) (_node *PaymentOrder, err error) {
	if err := pouo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(paymentorder.Table, paymentorder.Columns, sqlgraph.NewFieldSpec(paymentorder.FieldID, field.TypeUUID))
	id, ok := pouo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "PaymentOrder.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := pouo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, paymentorder.FieldID)
		for _, f := range fields {
			if !paymentorder.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != paymentorder.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := pouo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := pouo.mutation.UpdatedAt(); ok {
		_spec.SetField(paymentorder.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := pouo.mutation.Amount(); ok {
		_spec.SetField(paymentorder.FieldAmount, field.TypeFloat64, value)
	}
	if value, ok := pouo.mutation.AddedAmount(); ok {
		_spec.AddField(paymentorder.FieldAmount, field.TypeFloat64, value)
	}
	if value, ok := pouo.mutation.AmountPaid(); ok {
		_spec.SetField(paymentorder.FieldAmountPaid, field.TypeFloat64, value)
	}
	if value, ok := pouo.mutation.AddedAmountPaid(); ok {
		_spec.AddField(paymentorder.FieldAmountPaid, field.TypeFloat64, value)
	}
	if value, ok := pouo.mutation.Rate(); ok {
		_spec.SetField(paymentorder.FieldRate, field.TypeFloat64, value)
	}
	if value, ok := pouo.mutation.AddedRate(); ok {
		_spec.AddField(paymentorder.FieldRate, field.TypeFloat64, value)
	}
	if value, ok := pouo.mutation.TxHash(); ok {
		_spec.SetField(paymentorder.FieldTxHash, field.TypeString, value)
	}
	if pouo.mutation.TxHashCleared() {
		_spec.ClearField(paymentorder.FieldTxHash, field.TypeString)
	}
	if value, ok := pouo.mutation.ReceiveAddressText(); ok {
		_spec.SetField(paymentorder.FieldReceiveAddressText, field.TypeString, value)
	}
	if value, ok := pouo.mutation.Label(); ok {
		_spec.SetField(paymentorder.FieldLabel, field.TypeString, value)
	}
	if value, ok := pouo.mutation.Status(); ok {
		_spec.SetField(paymentorder.FieldStatus, field.TypeEnum, value)
	}
	if pouo.mutation.SenderProfileCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pouo.mutation.SenderProfileIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pouo.mutation.TokenCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pouo.mutation.TokenIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pouo.mutation.ReceiveAddressCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pouo.mutation.ReceiveAddressIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if pouo.mutation.RecipientCleared() {
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
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := pouo.mutation.RecipientIDs(); len(nodes) > 0 {
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
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &PaymentOrder{config: pouo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, pouo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{paymentorder.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	pouo.mutation.done = true
	return _node, nil
}
