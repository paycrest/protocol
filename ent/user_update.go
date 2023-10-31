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
	"github.com/paycrest/protocol/ent/predicate"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/senderprofile"
	"github.com/paycrest/protocol/ent/user"
	"github.com/paycrest/protocol/ent/verificationtoken"
)

// UserUpdate is the builder for updating User entities.
type UserUpdate struct {
	config
	hooks    []Hook
	mutation *UserMutation
}

// Where appends a list predicates to the UserUpdate builder.
func (uu *UserUpdate) Where(ps ...predicate.User) *UserUpdate {
	uu.mutation.Where(ps...)
	return uu
}

// SetUpdatedAt sets the "updated_at" field.
func (uu *UserUpdate) SetUpdatedAt(t time.Time) *UserUpdate {
	uu.mutation.SetUpdatedAt(t)
	return uu
}

// SetFirstName sets the "first_name" field.
func (uu *UserUpdate) SetFirstName(s string) *UserUpdate {
	uu.mutation.SetFirstName(s)
	return uu
}

// SetLastName sets the "last_name" field.
func (uu *UserUpdate) SetLastName(s string) *UserUpdate {
	uu.mutation.SetLastName(s)
	return uu
}

// SetEmail sets the "email" field.
func (uu *UserUpdate) SetEmail(s string) *UserUpdate {
	uu.mutation.SetEmail(s)
	return uu
}

// SetPassword sets the "password" field.
func (uu *UserUpdate) SetPassword(s string) *UserUpdate {
	uu.mutation.SetPassword(s)
	return uu
}

// SetScope sets the "scope" field.
func (uu *UserUpdate) SetScope(u user.Scope) *UserUpdate {
	uu.mutation.SetScope(u)
	return uu
}

// SetIsVerified sets the "is_verified" field.
func (uu *UserUpdate) SetIsVerified(b bool) *UserUpdate {
	uu.mutation.SetIsVerified(b)
	return uu
}

// SetNillableIsVerified sets the "is_verified" field if the given value is not nil.
func (uu *UserUpdate) SetNillableIsVerified(b *bool) *UserUpdate {
	if b != nil {
		uu.SetIsVerified(*b)
	}
	return uu
}

// SetSenderProfileID sets the "sender_profile" edge to the SenderProfile entity by ID.
func (uu *UserUpdate) SetSenderProfileID(id uuid.UUID) *UserUpdate {
	uu.mutation.SetSenderProfileID(id)
	return uu
}

// SetNillableSenderProfileID sets the "sender_profile" edge to the SenderProfile entity by ID if the given value is not nil.
func (uu *UserUpdate) SetNillableSenderProfileID(id *uuid.UUID) *UserUpdate {
	if id != nil {
		uu = uu.SetSenderProfileID(*id)
	}
	return uu
}

// SetSenderProfile sets the "sender_profile" edge to the SenderProfile entity.
func (uu *UserUpdate) SetSenderProfile(s *SenderProfile) *UserUpdate {
	return uu.SetSenderProfileID(s.ID)
}

// SetProviderProfileID sets the "provider_profile" edge to the ProviderProfile entity by ID.
func (uu *UserUpdate) SetProviderProfileID(id string) *UserUpdate {
	uu.mutation.SetProviderProfileID(id)
	return uu
}

// SetNillableProviderProfileID sets the "provider_profile" edge to the ProviderProfile entity by ID if the given value is not nil.
func (uu *UserUpdate) SetNillableProviderProfileID(id *string) *UserUpdate {
	if id != nil {
		uu = uu.SetProviderProfileID(*id)
	}
	return uu
}

// SetProviderProfile sets the "provider_profile" edge to the ProviderProfile entity.
func (uu *UserUpdate) SetProviderProfile(p *ProviderProfile) *UserUpdate {
	return uu.SetProviderProfileID(p.ID)
}

// AddVerificationTokenIDs adds the "verification_token" edge to the VerificationToken entity by IDs.
func (uu *UserUpdate) AddVerificationTokenIDs(ids ...uuid.UUID) *UserUpdate {
	uu.mutation.AddVerificationTokenIDs(ids...)
	return uu
}

// AddVerificationToken adds the "verification_token" edges to the VerificationToken entity.
func (uu *UserUpdate) AddVerificationToken(v ...*VerificationToken) *UserUpdate {
	ids := make([]uuid.UUID, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return uu.AddVerificationTokenIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uu *UserUpdate) Mutation() *UserMutation {
	return uu.mutation
}

// ClearSenderProfile clears the "sender_profile" edge to the SenderProfile entity.
func (uu *UserUpdate) ClearSenderProfile() *UserUpdate {
	uu.mutation.ClearSenderProfile()
	return uu
}

// ClearProviderProfile clears the "provider_profile" edge to the ProviderProfile entity.
func (uu *UserUpdate) ClearProviderProfile() *UserUpdate {
	uu.mutation.ClearProviderProfile()
	return uu
}

// ClearVerificationToken clears all "verification_token" edges to the VerificationToken entity.
func (uu *UserUpdate) ClearVerificationToken() *UserUpdate {
	uu.mutation.ClearVerificationToken()
	return uu
}

// RemoveVerificationTokenIDs removes the "verification_token" edge to VerificationToken entities by IDs.
func (uu *UserUpdate) RemoveVerificationTokenIDs(ids ...uuid.UUID) *UserUpdate {
	uu.mutation.RemoveVerificationTokenIDs(ids...)
	return uu
}

// RemoveVerificationToken removes "verification_token" edges to VerificationToken entities.
func (uu *UserUpdate) RemoveVerificationToken(v ...*VerificationToken) *UserUpdate {
	ids := make([]uuid.UUID, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return uu.RemoveVerificationTokenIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (uu *UserUpdate) Save(ctx context.Context) (int, error) {
	if err := uu.defaults(); err != nil {
		return 0, err
	}
	return withHooks(ctx, uu.sqlSave, uu.mutation, uu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uu *UserUpdate) SaveX(ctx context.Context) int {
	affected, err := uu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (uu *UserUpdate) Exec(ctx context.Context) error {
	_, err := uu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uu *UserUpdate) ExecX(ctx context.Context) {
	if err := uu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uu *UserUpdate) defaults() error {
	if _, ok := uu.mutation.UpdatedAt(); !ok {
		if user.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized user.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := user.UpdateDefaultUpdatedAt()
		uu.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (uu *UserUpdate) check() error {
	if v, ok := uu.mutation.FirstName(); ok {
		if err := user.FirstNameValidator(v); err != nil {
			return &ValidationError{Name: "first_name", err: fmt.Errorf(`ent: validator failed for field "User.first_name": %w`, err)}
		}
	}
	if v, ok := uu.mutation.LastName(); ok {
		if err := user.LastNameValidator(v); err != nil {
			return &ValidationError{Name: "last_name", err: fmt.Errorf(`ent: validator failed for field "User.last_name": %w`, err)}
		}
	}
	if v, ok := uu.mutation.Scope(); ok {
		if err := user.ScopeValidator(v); err != nil {
			return &ValidationError{Name: "scope", err: fmt.Errorf(`ent: validator failed for field "User.scope": %w`, err)}
		}
	}
	return nil
}

func (uu *UserUpdate) sqlSave(ctx context.Context) (n int, err error) {
	if err := uu.check(); err != nil {
		return n, err
	}
	_spec := sqlgraph.NewUpdateSpec(user.Table, user.Columns, sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID))
	if ps := uu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uu.mutation.UpdatedAt(); ok {
		_spec.SetField(user.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := uu.mutation.FirstName(); ok {
		_spec.SetField(user.FieldFirstName, field.TypeString, value)
	}
	if value, ok := uu.mutation.LastName(); ok {
		_spec.SetField(user.FieldLastName, field.TypeString, value)
	}
	if value, ok := uu.mutation.Email(); ok {
		_spec.SetField(user.FieldEmail, field.TypeString, value)
	}
	if value, ok := uu.mutation.Password(); ok {
		_spec.SetField(user.FieldPassword, field.TypeString, value)
	}
	if value, ok := uu.mutation.Scope(); ok {
		_spec.SetField(user.FieldScope, field.TypeEnum, value)
	}
	if value, ok := uu.mutation.IsVerified(); ok {
		_spec.SetField(user.FieldIsVerified, field.TypeBool, value)
	}
	if uu.mutation.SenderProfileCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   user.SenderProfileTable,
			Columns: []string{user.SenderProfileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(senderprofile.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.SenderProfileIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   user.SenderProfileTable,
			Columns: []string{user.SenderProfileColumn},
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
	if uu.mutation.ProviderProfileCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   user.ProviderProfileTable,
			Columns: []string{user.ProviderProfileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerprofile.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.ProviderProfileIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   user.ProviderProfileTable,
			Columns: []string{user.ProviderProfileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerprofile.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if uu.mutation.VerificationTokenCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.VerificationTokenTable,
			Columns: []string{user.VerificationTokenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(verificationtoken.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.RemovedVerificationTokenIDs(); len(nodes) > 0 && !uu.mutation.VerificationTokenCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.VerificationTokenTable,
			Columns: []string{user.VerificationTokenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(verificationtoken.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uu.mutation.VerificationTokenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.VerificationTokenTable,
			Columns: []string{user.VerificationTokenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(verificationtoken.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, uu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	uu.mutation.done = true
	return n, nil
}

// UserUpdateOne is the builder for updating a single User entity.
type UserUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *UserMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (uuo *UserUpdateOne) SetUpdatedAt(t time.Time) *UserUpdateOne {
	uuo.mutation.SetUpdatedAt(t)
	return uuo
}

// SetFirstName sets the "first_name" field.
func (uuo *UserUpdateOne) SetFirstName(s string) *UserUpdateOne {
	uuo.mutation.SetFirstName(s)
	return uuo
}

// SetLastName sets the "last_name" field.
func (uuo *UserUpdateOne) SetLastName(s string) *UserUpdateOne {
	uuo.mutation.SetLastName(s)
	return uuo
}

// SetEmail sets the "email" field.
func (uuo *UserUpdateOne) SetEmail(s string) *UserUpdateOne {
	uuo.mutation.SetEmail(s)
	return uuo
}

// SetPassword sets the "password" field.
func (uuo *UserUpdateOne) SetPassword(s string) *UserUpdateOne {
	uuo.mutation.SetPassword(s)
	return uuo
}

// SetScope sets the "scope" field.
func (uuo *UserUpdateOne) SetScope(u user.Scope) *UserUpdateOne {
	uuo.mutation.SetScope(u)
	return uuo
}

// SetIsVerified sets the "is_verified" field.
func (uuo *UserUpdateOne) SetIsVerified(b bool) *UserUpdateOne {
	uuo.mutation.SetIsVerified(b)
	return uuo
}

// SetNillableIsVerified sets the "is_verified" field if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableIsVerified(b *bool) *UserUpdateOne {
	if b != nil {
		uuo.SetIsVerified(*b)
	}
	return uuo
}

// SetSenderProfileID sets the "sender_profile" edge to the SenderProfile entity by ID.
func (uuo *UserUpdateOne) SetSenderProfileID(id uuid.UUID) *UserUpdateOne {
	uuo.mutation.SetSenderProfileID(id)
	return uuo
}

// SetNillableSenderProfileID sets the "sender_profile" edge to the SenderProfile entity by ID if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableSenderProfileID(id *uuid.UUID) *UserUpdateOne {
	if id != nil {
		uuo = uuo.SetSenderProfileID(*id)
	}
	return uuo
}

// SetSenderProfile sets the "sender_profile" edge to the SenderProfile entity.
func (uuo *UserUpdateOne) SetSenderProfile(s *SenderProfile) *UserUpdateOne {
	return uuo.SetSenderProfileID(s.ID)
}

// SetProviderProfileID sets the "provider_profile" edge to the ProviderProfile entity by ID.
func (uuo *UserUpdateOne) SetProviderProfileID(id string) *UserUpdateOne {
	uuo.mutation.SetProviderProfileID(id)
	return uuo
}

// SetNillableProviderProfileID sets the "provider_profile" edge to the ProviderProfile entity by ID if the given value is not nil.
func (uuo *UserUpdateOne) SetNillableProviderProfileID(id *string) *UserUpdateOne {
	if id != nil {
		uuo = uuo.SetProviderProfileID(*id)
	}
	return uuo
}

// SetProviderProfile sets the "provider_profile" edge to the ProviderProfile entity.
func (uuo *UserUpdateOne) SetProviderProfile(p *ProviderProfile) *UserUpdateOne {
	return uuo.SetProviderProfileID(p.ID)
}

// AddVerificationTokenIDs adds the "verification_token" edge to the VerificationToken entity by IDs.
func (uuo *UserUpdateOne) AddVerificationTokenIDs(ids ...uuid.UUID) *UserUpdateOne {
	uuo.mutation.AddVerificationTokenIDs(ids...)
	return uuo
}

// AddVerificationToken adds the "verification_token" edges to the VerificationToken entity.
func (uuo *UserUpdateOne) AddVerificationToken(v ...*VerificationToken) *UserUpdateOne {
	ids := make([]uuid.UUID, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return uuo.AddVerificationTokenIDs(ids...)
}

// Mutation returns the UserMutation object of the builder.
func (uuo *UserUpdateOne) Mutation() *UserMutation {
	return uuo.mutation
}

// ClearSenderProfile clears the "sender_profile" edge to the SenderProfile entity.
func (uuo *UserUpdateOne) ClearSenderProfile() *UserUpdateOne {
	uuo.mutation.ClearSenderProfile()
	return uuo
}

// ClearProviderProfile clears the "provider_profile" edge to the ProviderProfile entity.
func (uuo *UserUpdateOne) ClearProviderProfile() *UserUpdateOne {
	uuo.mutation.ClearProviderProfile()
	return uuo
}

// ClearVerificationToken clears all "verification_token" edges to the VerificationToken entity.
func (uuo *UserUpdateOne) ClearVerificationToken() *UserUpdateOne {
	uuo.mutation.ClearVerificationToken()
	return uuo
}

// RemoveVerificationTokenIDs removes the "verification_token" edge to VerificationToken entities by IDs.
func (uuo *UserUpdateOne) RemoveVerificationTokenIDs(ids ...uuid.UUID) *UserUpdateOne {
	uuo.mutation.RemoveVerificationTokenIDs(ids...)
	return uuo
}

// RemoveVerificationToken removes "verification_token" edges to VerificationToken entities.
func (uuo *UserUpdateOne) RemoveVerificationToken(v ...*VerificationToken) *UserUpdateOne {
	ids := make([]uuid.UUID, len(v))
	for i := range v {
		ids[i] = v[i].ID
	}
	return uuo.RemoveVerificationTokenIDs(ids...)
}

// Where appends a list predicates to the UserUpdate builder.
func (uuo *UserUpdateOne) Where(ps ...predicate.User) *UserUpdateOne {
	uuo.mutation.Where(ps...)
	return uuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (uuo *UserUpdateOne) Select(field string, fields ...string) *UserUpdateOne {
	uuo.fields = append([]string{field}, fields...)
	return uuo
}

// Save executes the query and returns the updated User entity.
func (uuo *UserUpdateOne) Save(ctx context.Context) (*User, error) {
	if err := uuo.defaults(); err != nil {
		return nil, err
	}
	return withHooks(ctx, uuo.sqlSave, uuo.mutation, uuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (uuo *UserUpdateOne) SaveX(ctx context.Context) *User {
	node, err := uuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (uuo *UserUpdateOne) Exec(ctx context.Context) error {
	_, err := uuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (uuo *UserUpdateOne) ExecX(ctx context.Context) {
	if err := uuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (uuo *UserUpdateOne) defaults() error {
	if _, ok := uuo.mutation.UpdatedAt(); !ok {
		if user.UpdateDefaultUpdatedAt == nil {
			return fmt.Errorf("ent: uninitialized user.UpdateDefaultUpdatedAt (forgotten import ent/runtime?)")
		}
		v := user.UpdateDefaultUpdatedAt()
		uuo.mutation.SetUpdatedAt(v)
	}
	return nil
}

// check runs all checks and user-defined validators on the builder.
func (uuo *UserUpdateOne) check() error {
	if v, ok := uuo.mutation.FirstName(); ok {
		if err := user.FirstNameValidator(v); err != nil {
			return &ValidationError{Name: "first_name", err: fmt.Errorf(`ent: validator failed for field "User.first_name": %w`, err)}
		}
	}
	if v, ok := uuo.mutation.LastName(); ok {
		if err := user.LastNameValidator(v); err != nil {
			return &ValidationError{Name: "last_name", err: fmt.Errorf(`ent: validator failed for field "User.last_name": %w`, err)}
		}
	}
	if v, ok := uuo.mutation.Scope(); ok {
		if err := user.ScopeValidator(v); err != nil {
			return &ValidationError{Name: "scope", err: fmt.Errorf(`ent: validator failed for field "User.scope": %w`, err)}
		}
	}
	return nil
}

func (uuo *UserUpdateOne) sqlSave(ctx context.Context) (_node *User, err error) {
	if err := uuo.check(); err != nil {
		return _node, err
	}
	_spec := sqlgraph.NewUpdateSpec(user.Table, user.Columns, sqlgraph.NewFieldSpec(user.FieldID, field.TypeUUID))
	id, ok := uuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "User.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := uuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, user.FieldID)
		for _, f := range fields {
			if !user.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != user.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := uuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := uuo.mutation.UpdatedAt(); ok {
		_spec.SetField(user.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := uuo.mutation.FirstName(); ok {
		_spec.SetField(user.FieldFirstName, field.TypeString, value)
	}
	if value, ok := uuo.mutation.LastName(); ok {
		_spec.SetField(user.FieldLastName, field.TypeString, value)
	}
	if value, ok := uuo.mutation.Email(); ok {
		_spec.SetField(user.FieldEmail, field.TypeString, value)
	}
	if value, ok := uuo.mutation.Password(); ok {
		_spec.SetField(user.FieldPassword, field.TypeString, value)
	}
	if value, ok := uuo.mutation.Scope(); ok {
		_spec.SetField(user.FieldScope, field.TypeEnum, value)
	}
	if value, ok := uuo.mutation.IsVerified(); ok {
		_spec.SetField(user.FieldIsVerified, field.TypeBool, value)
	}
	if uuo.mutation.SenderProfileCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   user.SenderProfileTable,
			Columns: []string{user.SenderProfileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(senderprofile.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.SenderProfileIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   user.SenderProfileTable,
			Columns: []string{user.SenderProfileColumn},
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
	if uuo.mutation.ProviderProfileCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   user.ProviderProfileTable,
			Columns: []string{user.ProviderProfileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerprofile.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.ProviderProfileIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2O,
			Inverse: false,
			Table:   user.ProviderProfileTable,
			Columns: []string{user.ProviderProfileColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerprofile.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if uuo.mutation.VerificationTokenCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.VerificationTokenTable,
			Columns: []string{user.VerificationTokenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(verificationtoken.FieldID, field.TypeUUID),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.RemovedVerificationTokenIDs(); len(nodes) > 0 && !uuo.mutation.VerificationTokenCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.VerificationTokenTable,
			Columns: []string{user.VerificationTokenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(verificationtoken.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := uuo.mutation.VerificationTokenIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   user.VerificationTokenTable,
			Columns: []string{user.VerificationTokenColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(verificationtoken.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &User{config: uuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, uuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{user.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	uuo.mutation.done = true
	return _node, nil
}
