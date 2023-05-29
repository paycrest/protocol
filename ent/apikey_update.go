// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/paycrest/paycrest-protocol/ent/apikey"
	"github.com/paycrest/paycrest-protocol/ent/predicate"
	"github.com/paycrest/paycrest-protocol/ent/user"
)

// APIKeyUpdate is the builder for updating APIKey entities.
type APIKeyUpdate struct {
	config
	hooks    []Hook
	mutation *APIKeyMutation
}

// Where appends a list predicates to the APIKeyUpdate builder.
func (aku *APIKeyUpdate) Where(ps ...predicate.APIKey) *APIKeyUpdate {
	aku.mutation.Where(ps...)
	return aku
}

// SetName sets the "name" field.
func (aku *APIKeyUpdate) SetName(s string) *APIKeyUpdate {
	aku.mutation.SetName(s)
	return aku
}

// SetScope sets the "scope" field.
func (aku *APIKeyUpdate) SetScope(a apikey.Scope) *APIKeyUpdate {
	aku.mutation.SetScope(a)
	return aku
}

// SetPair sets the "pair" field.
func (aku *APIKeyUpdate) SetPair(s string) *APIKeyUpdate {
	aku.mutation.SetPair(s)
	return aku
}

// SetIsActive sets the "is_active" field.
func (aku *APIKeyUpdate) SetIsActive(b bool) *APIKeyUpdate {
	aku.mutation.SetIsActive(b)
	return aku
}

// SetNillableIsActive sets the "is_active" field if the given value is not nil.
func (aku *APIKeyUpdate) SetNillableIsActive(b *bool) *APIKeyUpdate {
	if b != nil {
		aku.SetIsActive(*b)
	}
	return aku
}

// SetUserID sets the "user" edge to the User entity by ID.
func (aku *APIKeyUpdate) SetUserID(id uuid.UUID) *APIKeyUpdate {
	aku.mutation.SetUserID(id)
	return aku
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (aku *APIKeyUpdate) SetNillableUserID(id *uuid.UUID) *APIKeyUpdate {
	if id != nil {
		aku = aku.SetUserID(*id)
	}
	return aku
}

// SetUser sets the "user" edge to the User entity.
func (aku *APIKeyUpdate) SetUser(u *User) *APIKeyUpdate {
	return aku.SetUserID(u.ID)
}

// Mutation returns the APIKeyMutation object of the builder.
func (aku *APIKeyUpdate) Mutation() *APIKeyMutation {
	return aku.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (aku *APIKeyUpdate) ClearUser() *APIKeyUpdate {
	aku.mutation.ClearUser()
	return aku
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (aku *APIKeyUpdate) Save(ctx context.Context) (int, error) {
	return withHooks[int, APIKeyMutation](ctx, aku.sqlSave, aku.mutation, aku.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (aku *APIKeyUpdate) SaveX(ctx context.Context) int {
	affected, err := aku.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (aku *APIKeyUpdate) Exec(ctx context.Context) error {
	_, err := aku.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (aku *APIKeyUpdate) ExecX(ctx context.Context) {
	if err := aku.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (aku *APIKeyUpdate) check() error {
	if v, ok := aku.mutation.Scope(); ok {
		if err := apikey.ScopeValidator(v); err != nil {
			return &ValidationError{Name: "scope", err: fmt.Errorf(`ent: validator failed for field "APIKey.scope": %w`, err)}
		}
	}
	if v, ok := aku.mutation.Pair(); ok {
		if err := apikey.PairValidator(v); err != nil {
			return &ValidationError{Name: "pair", err: fmt.Errorf(`ent: validator failed for field "APIKey.pair": %w`, err)}
		}
	}
	return nil
}

func (aku *APIKeyUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := aku.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(apikey.Table, apikey.Columns, sqlgraph.NewFieldSpec(apikey.FieldID, field.TypeInt))
	if ps := aku.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := aku.mutation.Name(); ok {
		_spec.SetField(apikey.FieldName, field.TypeString, value)
	}
	if value, ok := aku.mutation.Scope(); ok {
		_spec.SetField(apikey.FieldScope, field.TypeEnum, value)
	}
	if value, ok := aku.mutation.Pair(); ok {
		_spec.SetField(apikey.FieldPair, field.TypeString, value)
	}
	if value, ok := aku.mutation.IsActive(); ok {
		_spec.SetField(apikey.FieldIsActive, field.TypeBool, value)
	}
	if aku.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   apikey.UserTable,
			Columns: []string{apikey.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := aku.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   apikey.UserTable,
			Columns: []string{apikey.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, aku.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{apikey.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	aku.mutation.done = true
	return n, nil
}

// APIKeyUpdateOne is the builder for updating a single APIKey entity.
type APIKeyUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *APIKeyMutation
}

// SetName sets the "name" field.
func (akuo *APIKeyUpdateOne) SetName(s string) *APIKeyUpdateOne {
	akuo.mutation.SetName(s)
	return akuo
}

// SetScope sets the "scope" field.
func (akuo *APIKeyUpdateOne) SetScope(a apikey.Scope) *APIKeyUpdateOne {
	akuo.mutation.SetScope(a)
	return akuo
}

// SetPair sets the "pair" field.
func (akuo *APIKeyUpdateOne) SetPair(s string) *APIKeyUpdateOne {
	akuo.mutation.SetPair(s)
	return akuo
}

// SetIsActive sets the "is_active" field.
func (akuo *APIKeyUpdateOne) SetIsActive(b bool) *APIKeyUpdateOne {
	akuo.mutation.SetIsActive(b)
	return akuo
}

// SetNillableIsActive sets the "is_active" field if the given value is not nil.
func (akuo *APIKeyUpdateOne) SetNillableIsActive(b *bool) *APIKeyUpdateOne {
	if b != nil {
		akuo.SetIsActive(*b)
	}
	return akuo
}

// SetUserID sets the "user" edge to the User entity by ID.
func (akuo *APIKeyUpdateOne) SetUserID(id uuid.UUID) *APIKeyUpdateOne {
	akuo.mutation.SetUserID(id)
	return akuo
}

// SetNillableUserID sets the "user" edge to the User entity by ID if the given value is not nil.
func (akuo *APIKeyUpdateOne) SetNillableUserID(id *uuid.UUID) *APIKeyUpdateOne {
	if id != nil {
		akuo = akuo.SetUserID(*id)
	}
	return akuo
}

// SetUser sets the "user" edge to the User entity.
func (akuo *APIKeyUpdateOne) SetUser(u *User) *APIKeyUpdateOne {
	return akuo.SetUserID(u.ID)
}

// Mutation returns the APIKeyMutation object of the builder.
func (akuo *APIKeyUpdateOne) Mutation() *APIKeyMutation {
	return akuo.mutation
}

// ClearUser clears the "user" edge to the User entity.
func (akuo *APIKeyUpdateOne) ClearUser() *APIKeyUpdateOne {
	akuo.mutation.ClearUser()
	return akuo
}

// Where appends a list predicates to the APIKeyUpdate builder.
func (akuo *APIKeyUpdateOne) Where(ps ...predicate.APIKey) *APIKeyUpdateOne {
	akuo.mutation.Where(ps...)
	return akuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (akuo *APIKeyUpdateOne) Select(field string, fields ...string) *APIKeyUpdateOne {
	akuo.fields = append([]string{field}, fields...)
	return akuo
}

// Save executes the query and returns the updated APIKey entity.
func (akuo *APIKeyUpdateOne) Save(ctx context.Context) (*APIKey, error) {
	return withHooks[*APIKey, APIKeyMutation](ctx, akuo.sqlSave, akuo.mutation, akuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (akuo *APIKeyUpdateOne) SaveX(ctx context.Context) *APIKey {
	node, err := akuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (akuo *APIKeyUpdateOne) Exec(ctx context.Context) error {
	_, err := akuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (akuo *APIKeyUpdateOne) ExecX(ctx context.Context) {
	if err := akuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (akuo *APIKeyUpdateOne) check() error {
	if v, ok := akuo.mutation.Scope(); ok {
		if err := apikey.ScopeValidator(v); err != nil {
			return &ValidationError{Name: "scope", err: fmt.Errorf(`ent: validator failed for field "APIKey.scope": %w`, err)}
		}
	}
	if v, ok := akuo.mutation.Pair(); ok {
		if err := apikey.PairValidator(v); err != nil {
			return &ValidationError{Name: "pair", err: fmt.Errorf(`ent: validator failed for field "APIKey.pair": %w`, err)}
		}
	}
	return nil
}

func (akuo *APIKeyUpdateOne) sqlSave(ctx context.Context) (_node *APIKey, err error) {
	if err := akuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(apikey.Table, apikey.Columns, sqlgraph.NewFieldSpec(apikey.FieldID, field.TypeInt))
	id, ok := akuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "APIKey.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := akuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, apikey.FieldID)
		for _, f := range fields {
			if !apikey.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != apikey.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := akuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := akuo.mutation.Name(); ok {
		_spec.SetField(apikey.FieldName, field.TypeString, value)
	}
	if value, ok := akuo.mutation.Scope(); ok {
		_spec.SetField(apikey.FieldScope, field.TypeEnum, value)
	}
	if value, ok := akuo.mutation.Pair(); ok {
		_spec.SetField(apikey.FieldPair, field.TypeString, value)
	}
	if value, ok := akuo.mutation.IsActive(); ok {
		_spec.SetField(apikey.FieldIsActive, field.TypeBool, value)
	}
	if akuo.mutation.UserCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   apikey.UserTable,
			Columns: []string{apikey.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := akuo.mutation.UserIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   apikey.UserTable,
			Columns: []string{apikey.UserColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &APIKey{config: akuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, akuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{apikey.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	akuo.mutation.done = true
	return _node, nil
}
