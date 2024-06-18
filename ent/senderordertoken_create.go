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
	"github.com/paycrest/protocol/ent/senderordertoken"
	"github.com/paycrest/protocol/ent/senderprofile"
	"github.com/shopspring/decimal"
)

// SenderOrderTokenCreate is the builder for creating a SenderOrderToken entity.
type SenderOrderTokenCreate struct {
	config
	mutation *SenderOrderTokenMutation
	hooks    []Hook
	conflict []sql.ConflictOption
}

// SetSymbol sets the "symbol" field.
func (sotc *SenderOrderTokenCreate) SetSymbol(s string) *SenderOrderTokenCreate {
	sotc.mutation.SetSymbol(s)
	return sotc
}

// SetFeePerTokenUnit sets the "fee_per_token_unit" field.
func (sotc *SenderOrderTokenCreate) SetFeePerTokenUnit(d decimal.Decimal) *SenderOrderTokenCreate {
	sotc.mutation.SetFeePerTokenUnit(d)
	return sotc
}

// SetAddresses sets the "addresses" field.
func (sotc *SenderOrderTokenCreate) SetAddresses(sddaaaa []struct {
	IsDisabled    bool   "json:\"isDisabled\""
	FeeAddress    string "json:\"feeAddress\""
	RefundAddress string "json:\"refundAddress\""
	Network       string "json:\"network\""
}) *SenderOrderTokenCreate {
	sotc.mutation.SetAddresses(sddaaaa)
	return sotc
}

// SetSenderID sets the "sender" edge to the SenderProfile entity by ID.
func (sotc *SenderOrderTokenCreate) SetSenderID(id uuid.UUID) *SenderOrderTokenCreate {
	sotc.mutation.SetSenderID(id)
	return sotc
}

// SetNillableSenderID sets the "sender" edge to the SenderProfile entity by ID if the given value is not nil.
func (sotc *SenderOrderTokenCreate) SetNillableSenderID(id *uuid.UUID) *SenderOrderTokenCreate {
	if id != nil {
		sotc = sotc.SetSenderID(*id)
	}
	return sotc
}

// SetSender sets the "sender" edge to the SenderProfile entity.
func (sotc *SenderOrderTokenCreate) SetSender(s *SenderProfile) *SenderOrderTokenCreate {
	return sotc.SetSenderID(s.ID)
}

// Mutation returns the SenderOrderTokenMutation object of the builder.
func (sotc *SenderOrderTokenCreate) Mutation() *SenderOrderTokenMutation {
	return sotc.mutation
}

// Save creates the SenderOrderToken in the database.
func (sotc *SenderOrderTokenCreate) Save(ctx context.Context) (*SenderOrderToken, error) {
	return withHooks(ctx, sotc.sqlSave, sotc.mutation, sotc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (sotc *SenderOrderTokenCreate) SaveX(ctx context.Context) *SenderOrderToken {
	v, err := sotc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sotc *SenderOrderTokenCreate) Exec(ctx context.Context) error {
	_, err := sotc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sotc *SenderOrderTokenCreate) ExecX(ctx context.Context) {
	if err := sotc.Exec(ctx); err != nil {
		panic(err)
	}
}

// check runs all checks and user-defined validators on the builder.
func (sotc *SenderOrderTokenCreate) check() error {
	if _, ok := sotc.mutation.Symbol(); !ok {
		return &ValidationError{Name: "symbol", err: errors.New(`ent: missing required field "SenderOrderToken.symbol"`)}
	}
	if _, ok := sotc.mutation.FeePerTokenUnit(); !ok {
		return &ValidationError{Name: "fee_per_token_unit", err: errors.New(`ent: missing required field "SenderOrderToken.fee_per_token_unit"`)}
	}
	if _, ok := sotc.mutation.Addresses(); !ok {
		return &ValidationError{Name: "addresses", err: errors.New(`ent: missing required field "SenderOrderToken.addresses"`)}
	}
	return nil
}

func (sotc *SenderOrderTokenCreate) sqlSave(ctx context.Context) (*SenderOrderToken, error) {
	if err := sotc.check(); err != nil {
		return nil, err
	}
	_node, _spec := sotc.createSpec()
	if err := sqlgraph.CreateNode(ctx, sotc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	sotc.mutation.id = &_node.ID
	sotc.mutation.done = true
	return _node, nil
}

func (sotc *SenderOrderTokenCreate) createSpec() (*SenderOrderToken, *sqlgraph.CreateSpec) {
	var (
		_node = &SenderOrderToken{config: sotc.config}
		_spec = sqlgraph.NewCreateSpec(senderordertoken.Table, sqlgraph.NewFieldSpec(senderordertoken.FieldID, field.TypeInt))
	)
	_spec.OnConflict = sotc.conflict
	if value, ok := sotc.mutation.Symbol(); ok {
		_spec.SetField(senderordertoken.FieldSymbol, field.TypeString, value)
		_node.Symbol = value
	}
	if value, ok := sotc.mutation.FeePerTokenUnit(); ok {
		_spec.SetField(senderordertoken.FieldFeePerTokenUnit, field.TypeFloat64, value)
		_node.FeePerTokenUnit = value
	}
	if value, ok := sotc.mutation.Addresses(); ok {
		_spec.SetField(senderordertoken.FieldAddresses, field.TypeJSON, value)
		_node.Addresses = value
	}
	if nodes := sotc.mutation.SenderIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2O,
			Inverse: true,
			Table:   senderordertoken.SenderTable,
			Columns: []string{senderordertoken.SenderColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(senderprofile.FieldID, field.TypeUUID),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_node.sender_profile_order_tokens = &nodes[0]
		_spec.Edges = append(_spec.Edges, edge)
	}
	return _node, _spec
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.SenderOrderToken.Create().
//		SetSymbol(v).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SenderOrderTokenUpsert) {
//			SetSymbol(v+v).
//		}).
//		Exec(ctx)
func (sotc *SenderOrderTokenCreate) OnConflict(opts ...sql.ConflictOption) *SenderOrderTokenUpsertOne {
	sotc.conflict = opts
	return &SenderOrderTokenUpsertOne{
		create: sotc,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.SenderOrderToken.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (sotc *SenderOrderTokenCreate) OnConflictColumns(columns ...string) *SenderOrderTokenUpsertOne {
	sotc.conflict = append(sotc.conflict, sql.ConflictColumns(columns...))
	return &SenderOrderTokenUpsertOne{
		create: sotc,
	}
}

type (
	// SenderOrderTokenUpsertOne is the builder for "upsert"-ing
	//  one SenderOrderToken node.
	SenderOrderTokenUpsertOne struct {
		create *SenderOrderTokenCreate
	}

	// SenderOrderTokenUpsert is the "OnConflict" setter.
	SenderOrderTokenUpsert struct {
		*sql.UpdateSet
	}
)

// SetSymbol sets the "symbol" field.
func (u *SenderOrderTokenUpsert) SetSymbol(v string) *SenderOrderTokenUpsert {
	u.Set(senderordertoken.FieldSymbol, v)
	return u
}

// UpdateSymbol sets the "symbol" field to the value that was provided on create.
func (u *SenderOrderTokenUpsert) UpdateSymbol() *SenderOrderTokenUpsert {
	u.SetExcluded(senderordertoken.FieldSymbol)
	return u
}

// SetFeePerTokenUnit sets the "fee_per_token_unit" field.
func (u *SenderOrderTokenUpsert) SetFeePerTokenUnit(v decimal.Decimal) *SenderOrderTokenUpsert {
	u.Set(senderordertoken.FieldFeePerTokenUnit, v)
	return u
}

// UpdateFeePerTokenUnit sets the "fee_per_token_unit" field to the value that was provided on create.
func (u *SenderOrderTokenUpsert) UpdateFeePerTokenUnit() *SenderOrderTokenUpsert {
	u.SetExcluded(senderordertoken.FieldFeePerTokenUnit)
	return u
}

// AddFeePerTokenUnit adds v to the "fee_per_token_unit" field.
func (u *SenderOrderTokenUpsert) AddFeePerTokenUnit(v decimal.Decimal) *SenderOrderTokenUpsert {
	u.Add(senderordertoken.FieldFeePerTokenUnit, v)
	return u
}

// SetAddresses sets the "addresses" field.
func (u *SenderOrderTokenUpsert) SetAddresses(v []struct {
	IsDisabled    bool   "json:\"isDisabled\""
	FeeAddress    string "json:\"feeAddress\""
	RefundAddress string "json:\"refundAddress\""
	Network       string "json:\"network\""
}) *SenderOrderTokenUpsert {
	u.Set(senderordertoken.FieldAddresses, v)
	return u
}

// UpdateAddresses sets the "addresses" field to the value that was provided on create.
func (u *SenderOrderTokenUpsert) UpdateAddresses() *SenderOrderTokenUpsert {
	u.SetExcluded(senderordertoken.FieldAddresses)
	return u
}

// UpdateNewValues updates the mutable fields using the new values that were set on create.
// Using this option is equivalent to using:
//
//	client.SenderOrderToken.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *SenderOrderTokenUpsertOne) UpdateNewValues() *SenderOrderTokenUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.SenderOrderToken.Create().
//	    OnConflict(sql.ResolveWithIgnore()).
//	    Exec(ctx)
func (u *SenderOrderTokenUpsertOne) Ignore() *SenderOrderTokenUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SenderOrderTokenUpsertOne) DoNothing() *SenderOrderTokenUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SenderOrderTokenCreate.OnConflict
// documentation for more info.
func (u *SenderOrderTokenUpsertOne) Update(set func(*SenderOrderTokenUpsert)) *SenderOrderTokenUpsertOne {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SenderOrderTokenUpsert{UpdateSet: update})
	}))
	return u
}

// SetSymbol sets the "symbol" field.
func (u *SenderOrderTokenUpsertOne) SetSymbol(v string) *SenderOrderTokenUpsertOne {
	return u.Update(func(s *SenderOrderTokenUpsert) {
		s.SetSymbol(v)
	})
}

// UpdateSymbol sets the "symbol" field to the value that was provided on create.
func (u *SenderOrderTokenUpsertOne) UpdateSymbol() *SenderOrderTokenUpsertOne {
	return u.Update(func(s *SenderOrderTokenUpsert) {
		s.UpdateSymbol()
	})
}

// SetFeePerTokenUnit sets the "fee_per_token_unit" field.
func (u *SenderOrderTokenUpsertOne) SetFeePerTokenUnit(v decimal.Decimal) *SenderOrderTokenUpsertOne {
	return u.Update(func(s *SenderOrderTokenUpsert) {
		s.SetFeePerTokenUnit(v)
	})
}

// AddFeePerTokenUnit adds v to the "fee_per_token_unit" field.
func (u *SenderOrderTokenUpsertOne) AddFeePerTokenUnit(v decimal.Decimal) *SenderOrderTokenUpsertOne {
	return u.Update(func(s *SenderOrderTokenUpsert) {
		s.AddFeePerTokenUnit(v)
	})
}

// UpdateFeePerTokenUnit sets the "fee_per_token_unit" field to the value that was provided on create.
func (u *SenderOrderTokenUpsertOne) UpdateFeePerTokenUnit() *SenderOrderTokenUpsertOne {
	return u.Update(func(s *SenderOrderTokenUpsert) {
		s.UpdateFeePerTokenUnit()
	})
}

// SetAddresses sets the "addresses" field.
func (u *SenderOrderTokenUpsertOne) SetAddresses(v []struct {
	IsDisabled    bool   "json:\"isDisabled\""
	FeeAddress    string "json:\"feeAddress\""
	RefundAddress string "json:\"refundAddress\""
	Network       string "json:\"network\""
}) *SenderOrderTokenUpsertOne {
	return u.Update(func(s *SenderOrderTokenUpsert) {
		s.SetAddresses(v)
	})
}

// UpdateAddresses sets the "addresses" field to the value that was provided on create.
func (u *SenderOrderTokenUpsertOne) UpdateAddresses() *SenderOrderTokenUpsertOne {
	return u.Update(func(s *SenderOrderTokenUpsert) {
		s.UpdateAddresses()
	})
}

// Exec executes the query.
func (u *SenderOrderTokenUpsertOne) Exec(ctx context.Context) error {
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SenderOrderTokenCreate.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SenderOrderTokenUpsertOne) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}

// Exec executes the UPSERT query and returns the inserted/updated ID.
func (u *SenderOrderTokenUpsertOne) ID(ctx context.Context) (id int, err error) {
	node, err := u.create.Save(ctx)
	if err != nil {
		return id, err
	}
	return node.ID, nil
}

// IDX is like ID, but panics if an error occurs.
func (u *SenderOrderTokenUpsertOne) IDX(ctx context.Context) int {
	id, err := u.ID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// SenderOrderTokenCreateBulk is the builder for creating many SenderOrderToken entities in bulk.
type SenderOrderTokenCreateBulk struct {
	config
	builders []*SenderOrderTokenCreate
	conflict []sql.ConflictOption
}

// Save creates the SenderOrderToken entities in the database.
func (sotcb *SenderOrderTokenCreateBulk) Save(ctx context.Context) ([]*SenderOrderToken, error) {
	specs := make([]*sqlgraph.CreateSpec, len(sotcb.builders))
	nodes := make([]*SenderOrderToken, len(sotcb.builders))
	mutators := make([]Mutator, len(sotcb.builders))
	for i := range sotcb.builders {
		func(i int, root context.Context) {
			builder := sotcb.builders[i]
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*SenderOrderTokenMutation)
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
					_, err = mutators[i+1].Mutate(root, sotcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					spec.OnConflict = sotcb.conflict
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, sotcb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, sotcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (sotcb *SenderOrderTokenCreateBulk) SaveX(ctx context.Context) []*SenderOrderToken {
	v, err := sotcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (sotcb *SenderOrderTokenCreateBulk) Exec(ctx context.Context) error {
	_, err := sotcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (sotcb *SenderOrderTokenCreateBulk) ExecX(ctx context.Context) {
	if err := sotcb.Exec(ctx); err != nil {
		panic(err)
	}
}

// OnConflict allows configuring the `ON CONFLICT` / `ON DUPLICATE KEY` clause
// of the `INSERT` statement. For example:
//
//	client.SenderOrderToken.CreateBulk(builders...).
//		OnConflict(
//			// Update the row with the new values
//			// the was proposed for insertion.
//			sql.ResolveWithNewValues(),
//		).
//		// Override some of the fields with custom
//		// update values.
//		Update(func(u *ent.SenderOrderTokenUpsert) {
//			SetSymbol(v+v).
//		}).
//		Exec(ctx)
func (sotcb *SenderOrderTokenCreateBulk) OnConflict(opts ...sql.ConflictOption) *SenderOrderTokenUpsertBulk {
	sotcb.conflict = opts
	return &SenderOrderTokenUpsertBulk{
		create: sotcb,
	}
}

// OnConflictColumns calls `OnConflict` and configures the columns
// as conflict target. Using this option is equivalent to using:
//
//	client.SenderOrderToken.Create().
//		OnConflict(sql.ConflictColumns(columns...)).
//		Exec(ctx)
func (sotcb *SenderOrderTokenCreateBulk) OnConflictColumns(columns ...string) *SenderOrderTokenUpsertBulk {
	sotcb.conflict = append(sotcb.conflict, sql.ConflictColumns(columns...))
	return &SenderOrderTokenUpsertBulk{
		create: sotcb,
	}
}

// SenderOrderTokenUpsertBulk is the builder for "upsert"-ing
// a bulk of SenderOrderToken nodes.
type SenderOrderTokenUpsertBulk struct {
	create *SenderOrderTokenCreateBulk
}

// UpdateNewValues updates the mutable fields using the new values that
// were set on create. Using this option is equivalent to using:
//
//	client.SenderOrderToken.Create().
//		OnConflict(
//			sql.ResolveWithNewValues(),
//		).
//		Exec(ctx)
func (u *SenderOrderTokenUpsertBulk) UpdateNewValues() *SenderOrderTokenUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithNewValues())
	return u
}

// Ignore sets each column to itself in case of conflict.
// Using this option is equivalent to using:
//
//	client.SenderOrderToken.Create().
//		OnConflict(sql.ResolveWithIgnore()).
//		Exec(ctx)
func (u *SenderOrderTokenUpsertBulk) Ignore() *SenderOrderTokenUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWithIgnore())
	return u
}

// DoNothing configures the conflict_action to `DO NOTHING`.
// Supported only by SQLite and PostgreSQL.
func (u *SenderOrderTokenUpsertBulk) DoNothing() *SenderOrderTokenUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.DoNothing())
	return u
}

// Update allows overriding fields `UPDATE` values. See the SenderOrderTokenCreateBulk.OnConflict
// documentation for more info.
func (u *SenderOrderTokenUpsertBulk) Update(set func(*SenderOrderTokenUpsert)) *SenderOrderTokenUpsertBulk {
	u.create.conflict = append(u.create.conflict, sql.ResolveWith(func(update *sql.UpdateSet) {
		set(&SenderOrderTokenUpsert{UpdateSet: update})
	}))
	return u
}

// SetSymbol sets the "symbol" field.
func (u *SenderOrderTokenUpsertBulk) SetSymbol(v string) *SenderOrderTokenUpsertBulk {
	return u.Update(func(s *SenderOrderTokenUpsert) {
		s.SetSymbol(v)
	})
}

// UpdateSymbol sets the "symbol" field to the value that was provided on create.
func (u *SenderOrderTokenUpsertBulk) UpdateSymbol() *SenderOrderTokenUpsertBulk {
	return u.Update(func(s *SenderOrderTokenUpsert) {
		s.UpdateSymbol()
	})
}

// SetFeePerTokenUnit sets the "fee_per_token_unit" field.
func (u *SenderOrderTokenUpsertBulk) SetFeePerTokenUnit(v decimal.Decimal) *SenderOrderTokenUpsertBulk {
	return u.Update(func(s *SenderOrderTokenUpsert) {
		s.SetFeePerTokenUnit(v)
	})
}

// AddFeePerTokenUnit adds v to the "fee_per_token_unit" field.
func (u *SenderOrderTokenUpsertBulk) AddFeePerTokenUnit(v decimal.Decimal) *SenderOrderTokenUpsertBulk {
	return u.Update(func(s *SenderOrderTokenUpsert) {
		s.AddFeePerTokenUnit(v)
	})
}

// UpdateFeePerTokenUnit sets the "fee_per_token_unit" field to the value that was provided on create.
func (u *SenderOrderTokenUpsertBulk) UpdateFeePerTokenUnit() *SenderOrderTokenUpsertBulk {
	return u.Update(func(s *SenderOrderTokenUpsert) {
		s.UpdateFeePerTokenUnit()
	})
}

// SetAddresses sets the "addresses" field.
func (u *SenderOrderTokenUpsertBulk) SetAddresses(v []struct {
	IsDisabled    bool   "json:\"isDisabled\""
	FeeAddress    string "json:\"feeAddress\""
	RefundAddress string "json:\"refundAddress\""
	Network       string "json:\"network\""
}) *SenderOrderTokenUpsertBulk {
	return u.Update(func(s *SenderOrderTokenUpsert) {
		s.SetAddresses(v)
	})
}

// UpdateAddresses sets the "addresses" field to the value that was provided on create.
func (u *SenderOrderTokenUpsertBulk) UpdateAddresses() *SenderOrderTokenUpsertBulk {
	return u.Update(func(s *SenderOrderTokenUpsert) {
		s.UpdateAddresses()
	})
}

// Exec executes the query.
func (u *SenderOrderTokenUpsertBulk) Exec(ctx context.Context) error {
	for i, b := range u.create.builders {
		if len(b.conflict) != 0 {
			return fmt.Errorf("ent: OnConflict was set for builder %d. Set it on the SenderOrderTokenCreateBulk instead", i)
		}
	}
	if len(u.create.conflict) == 0 {
		return errors.New("ent: missing options for SenderOrderTokenCreateBulk.OnConflict")
	}
	return u.create.Exec(ctx)
}

// ExecX is like Exec, but panics if an error occurs.
func (u *SenderOrderTokenUpsertBulk) ExecX(ctx context.Context) {
	if err := u.create.Exec(ctx); err != nil {
		panic(err)
	}
}
