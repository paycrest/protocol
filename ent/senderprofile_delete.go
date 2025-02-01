// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/paycrest/aggregator/ent/predicate"
	"github.com/paycrest/aggregator/ent/senderprofile"
)

// SenderProfileDelete is the builder for deleting a SenderProfile entity.
type SenderProfileDelete struct {
	config
	hooks    []Hook
	mutation *SenderProfileMutation
}

// Where appends a list predicates to the SenderProfileDelete builder.
func (spd *SenderProfileDelete) Where(ps ...predicate.SenderProfile) *SenderProfileDelete {
	spd.mutation.Where(ps...)
	return spd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (spd *SenderProfileDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, spd.sqlExec, spd.mutation, spd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (spd *SenderProfileDelete) ExecX(ctx context.Context) int {
	n, err := spd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (spd *SenderProfileDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(senderprofile.Table, sqlgraph.NewFieldSpec(senderprofile.FieldID, field.TypeUUID))
	if ps := spd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, spd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	spd.mutation.done = true
	return affected, err
}

// SenderProfileDeleteOne is the builder for deleting a single SenderProfile entity.
type SenderProfileDeleteOne struct {
	spd *SenderProfileDelete
}

// Where appends a list predicates to the SenderProfileDelete builder.
func (spdo *SenderProfileDeleteOne) Where(ps ...predicate.SenderProfile) *SenderProfileDeleteOne {
	spdo.spd.mutation.Where(ps...)
	return spdo
}

// Exec executes the deletion query.
func (spdo *SenderProfileDeleteOne) Exec(ctx context.Context) error {
	n, err := spdo.spd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{senderprofile.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (spdo *SenderProfileDeleteOne) ExecX(ctx context.Context) {
	if err := spdo.Exec(ctx); err != nil {
		panic(err)
	}
}
