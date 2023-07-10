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
	"github.com/paycrest/paycrest-protocol/ent/apikey"
	"github.com/paycrest/paycrest-protocol/ent/paymentorder"
	"github.com/paycrest/paycrest-protocol/ent/providerprofile"
	"github.com/paycrest/paycrest-protocol/ent/user"
)

// APIKeyCreate is the builder for creating a APIKey entity.
type APIKeyCreate struct {
	config
	mutation *APIKeyMutation
	hooks    []Hook
}

// SetName sets the "name" field.
func (akc *APIKeyCreate) SetName(s string) *APIKeyCreate {
	akc.mutation.SetName(s)
	return akc
}

// SetScope sets the "scope" field.
func (akc *APIKeyCreate) SetScope(a apikey.Scope) *APIKeyCreate {
	akc.mutation.SetScope(a)
	return akc
}

// SetSecret sets the "secret" field.
func (akc *APIKeyCreate) SetSecret(s string) *APIKeyCreate {
	akc.mutation.SetSecret(s)
	return akc
}

// SetIsActive sets the "is_active" field.
func (akc *APIKeyCreate) SetIsActive(b bool) *APIKeyCreate {
	akc.mutation.SetIsActive(b)
	return akc
}

// SetNillableIsActive sets the "is_active" field if the given value is not nil.
func (akc *APIKeyCreate) SetNillableIsActive(b *bool) *APIKeyCreate {
	if b != nil {
		akc.SetIsActive(*b)
	}
	return akc
}

// SetCreatedAt sets the "created_at" field.
func (akc *APIKeyCreate) SetCreatedAt(t time.Time) *APIKeyCreate {
	akc.mutation.SetCreatedAt(t)
	return akc
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (akc *APIKeyCreate) SetNillableCreatedAt(t *time.Time) *APIKeyCreate {
	if t != nil {
		akc.SetCreatedAt(*t)
	}
	return akc
}

// SetID sets the "id" field.
func (akc *APIKeyCreate) SetID(u uuid.UUID) *APIKeyCreate {
	akc.mutation.SetID(u)
	return akc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (akc *APIKeyCreate) SetNillableID(u *uuid.UUID) *APIKeyCreate {
	if u != nil {
		akc.SetID(*u)
	}
	return akc
}

// SetOwnerID sets the "owner" edge to the User entity by ID.
func (akc *APIKeyCreate) SetOwnerID(id uuid.UUID) *APIKeyCreate {
	akc.mutation.SetOwnerID(id)
	return akc
}

// SetNillableOwnerID sets the "owner" edge to the User entity by ID if the given value is not nil.
func (akc *APIKeyCreate) SetNillableOwnerID(id *uuid.UUID) *APIKeyCreate {
	if id != nil {
		akc = akc.SetOwnerID(*id)
	}
	return akc
}

// SetOwner sets the "owner" edge to the User entity.
func (akc *APIKeyCreate) SetOwner(u *User) *APIKeyCreate {
	return akc.SetOwnerID(u.ID)
}

// SetProviderProfileID sets the "provider_profile" edge to the ProviderProfile entity by ID.
func (akc *APIKeyCreate) SetProviderProfileID(id string) *APIKeyCreate {
	akc.mutation.SetProviderProfileID(id)
	return akc
}

// SetNillableProviderProfileID sets the "provider_profile" edge to the ProviderProfile entity by ID if the given value is not nil.
func (akc *APIKeyCreate) SetNillableProviderProfileID(id *string) *APIKeyCreate {
	if id != nil {
		akc = akc.SetProviderProfileID(*id)
	}
	return akc
}

// SetProviderProfile sets the "provider_profile" edge to the ProviderProfile entity.
func (akc *APIKeyCreate) SetProviderProfile(p *ProviderProfile) *APIKeyCreate {
	return akc.SetProviderProfileID(p.ID)
}

// AddPaymentOrderIDs adds the "payment_orders" edge to the PaymentOrder entity by IDs.
func (akc *APIKeyCreate) AddPaymentOrderIDs(ids ...int) *APIKeyCreate {
	akc.mutation.AddPaymentOrderIDs(ids...)
	return akc
}

// AddPaymentOrders adds the "payment_orders" edges to the PaymentOrder entity.
func (akc *APIKeyCreate) AddPaymentOrders(p ...*PaymentOrder) *APIKeyCreate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return akc.AddPaymentOrderIDs(ids...)
}

// Mutation returns the APIKeyMutation object of the builder.
func (akc *APIKeyCreate) Mutation() *APIKeyMutation {
	return akc.mutation
}

// Save creates the APIKey in the database.
func (akc *APIKeyCreate) Save(ctx context.Context) (*APIKey, error) {
	akc.defaults()
	return withHooks(ctx, akc.sqlSave, akc.mutation, akc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (akc *APIKeyCreate) SaveX(ctx context.Context) *APIKey {
	v, err := akc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (akc *APIKeyCreate) Exec(ctx context.Context) error {
	_, err := akc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (akc *APIKeyCreate) ExecX(ctx context.Context) {
	if err := akc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (akc *APIKeyCreate) defaults() {
	if _, ok := akc.mutation.IsActive(); !ok {
		v := apikey.DefaultIsActive
		akc.mutation.SetIsActive(v)
	}
	if _, ok := akc.mutation.CreatedAt(); !ok {
		v := apikey.DefaultCreatedAt()
		akc.mutation.SetCreatedAt(v)
	}
	if _, ok := akc.mutation.ID(); !ok {
		v := apikey.DefaultID()
		akc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (akc *APIKeyCreate) check() error {
	if _, ok := akc.mutation.Name(); !ok {
		return &ValidationError{Name: "name", err: errors.New(`ent: missing required field "APIKey.name"`)}
	}
	if _, ok := akc.mutation.Scope(); !ok {
		return &ValidationError{Name: "scope", err: errors.New(`ent: missing required field "APIKey.scope"`)}
	}
	if v, ok := akc.mutation.Scope(); ok {
		if err := apikey.ScopeValidator(v); err != nil {
			return &ValidationError{Name: "scope", err: fmt.Errorf(`ent: validator failed for field "APIKey.scope": %w`, err)}
		}
	}
	if _, ok := akc.mutation.Secret(); !ok {
		return &ValidationError{Name: "secret", err: errors.New(`ent: missing required field "APIKey.secret"`)}
	}
	if v, ok := akc.mutation.Secret(); ok {
		if err := apikey.SecretValidator(v); err != nil {
			return &ValidationError{Name: "secret", err: fmt.Errorf(`ent: validator failed for field "APIKey.secret": %w`, err)}
		}
	}
	if _, ok := akc.mutation.IsActive(); !ok {
		return &ValidationError{Name: "is_active", err: errors.New(`ent: missing required field "APIKey.is_active"`)}
	}
	if _, ok := akc.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "APIKey.created_at"`)}
	}
	return nil
}

func (akc *APIKeyCreate) sqlSave(ctx context.Context) (*APIKey, error) {
	if err := akc.check(); err != nil {
		return nil, err
	}
	_node, _spec := akc.createSpec()
	if err := sqlgraph.CreateNode(ctx, akc.driver, _spec); err != nil {
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
	akc.mutation.id = &_node.ID
	akc.mutation.done = true
	return _node, nil
}

func (akc *APIKeyCreate) createSpec() (*APIKey, *sqlgraph.CreateSpec) {
	var (
		_node = &APIKey{config: akc.config}
		_spec = sqlgraph.NewCreateSpec(apikey.Table, sqlgraph.NewFieldSpec(apikey.FieldID, field.TypeUUID))
	)
	if id, ok := akc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := akc.mutation.Name(); ok {
		_spec.SetField(apikey.FieldName, field.TypeString, value)
		_node.Name = value
	}
	if value, ok := akc.mutation.Scope(); ok {
		_spec.SetField(apikey.FieldScope, field.TypeEnum, value)
		_node.Scope = value
	}
	if value, ok := akc.mutation.Secret(); ok {
		_spec.SetField(apikey.FieldSecret, field.TypeString, value)
		_node.Secret = value
	}
	if value, ok := akc.mutation.IsActive(); ok {
		_spec.SetField(apikey.FieldIsActive, field.TypeBool, value)
		_node.IsActive = value
	}
	if value, ok := akc.mutation.CreatedAt(); ok {
		_spec.SetField(apikey.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if nodes := akc.mutation.OwnerIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   apikey.OwnerTable,
			Columns: []string{apikey.OwnerColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.user_api_keys = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := akc.mutation.ProviderProfileIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   apikey.ProviderProfileTable,
			Columns: []string{apikey.ProviderProfileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerprofile.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	if nodes := akc.mutation.PaymentOrdersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   apikey.PaymentOrdersTable,
			Columns: []string{apikey.PaymentOrdersColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(paymentorder.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// APIKeyCreateBulk is the builder for creating many APIKey entities in bulk.
type APIKeyCreateBulk struct {
	config
	builders []*APIKeyCreate
}

// Save creates the APIKey entities in the database.
func (akcb *APIKeyCreateBulk) Save(ctx context.Context) ([]*APIKey, error) {
	specs := make([]*sqlgraph.CreateSpec, len(akcb.builders))
	nodes := make([]*APIKey, len(akcb.builders))
	mutators := make([]Mutator, len(akcb.builders))
	for i := range akcb.builders {
		func(i int, root context.Context) {
			builder := akcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*APIKeyMutation)
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
					_, err = mutators[i+1].Mutate(root, akcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, akcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, akcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (akcb *APIKeyCreateBulk) SaveX(ctx context.Context) []*APIKey {
	v, err := akcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (akcb *APIKeyCreateBulk) Exec(ctx context.Context) error {
	_, err := akcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (akcb *APIKeyCreateBulk) ExecX(ctx context.Context) {
	if err := akcb.Exec(ctx); err != nil {
		panic(err)
	}
}
