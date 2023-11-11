// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/paycrest/protocol/ent/webhookretryattempt"
)

// WebhookRetryAttemptCreate is the builder for creating a WebhookRetryAttempt entity.
type WebhookRetryAttemptCreate struct {
	config
	mutation *WebhookRetryAttemptMutation
	hooks    []Hook
}

// SetCreatedAt sets the "created_at" field.
func (wrac *WebhookRetryAttemptCreate) SetCreatedAt(t time.Time) *WebhookRetryAttemptCreate {
	wrac.mutation.SetCreatedAt(t)
	return wrac
}

// SetNillableCreatedAt sets the "created_at" field if the given value is not nil.
func (wrac *WebhookRetryAttemptCreate) SetNillableCreatedAt(t *time.Time) *WebhookRetryAttemptCreate {
	if t != nil {
		wrac.SetCreatedAt(*t)
	}
	return wrac
}

// SetUpdatedAt sets the "updated_at" field.
func (wrac *WebhookRetryAttemptCreate) SetUpdatedAt(t time.Time) *WebhookRetryAttemptCreate {
	wrac.mutation.SetUpdatedAt(t)
	return wrac
}

// SetNillableUpdatedAt sets the "updated_at" field if the given value is not nil.
func (wrac *WebhookRetryAttemptCreate) SetNillableUpdatedAt(t *time.Time) *WebhookRetryAttemptCreate {
	if t != nil {
		wrac.SetUpdatedAt(*t)
	}
	return wrac
}

// SetAttemptNumber sets the "attempt_number" field.
func (wrac *WebhookRetryAttemptCreate) SetAttemptNumber(i int) *WebhookRetryAttemptCreate {
	wrac.mutation.SetAttemptNumber(i)
	return wrac
}

// SetNextRetryTime sets the "next_retry_time" field.
func (wrac *WebhookRetryAttemptCreate) SetNextRetryTime(t time.Time) *WebhookRetryAttemptCreate {
	wrac.mutation.SetNextRetryTime(t)
	return wrac
}

// SetNillableNextRetryTime sets the "next_retry_time" field if the given value is not nil.
func (wrac *WebhookRetryAttemptCreate) SetNillableNextRetryTime(t *time.Time) *WebhookRetryAttemptCreate {
	if t != nil {
		wrac.SetNextRetryTime(*t)
	}
	return wrac
}

// SetPayload sets the "payload" field.
func (wrac *WebhookRetryAttemptCreate) SetPayload(m map[string]interface{}) *WebhookRetryAttemptCreate {
	wrac.mutation.SetPayload(m)
	return wrac
}

// SetSignature sets the "signature" field.
func (wrac *WebhookRetryAttemptCreate) SetSignature(s string) *WebhookRetryAttemptCreate {
	wrac.mutation.SetSignature(s)
	return wrac
}

// SetNillableSignature sets the "signature" field if the given value is not nil.
func (wrac *WebhookRetryAttemptCreate) SetNillableSignature(s *string) *WebhookRetryAttemptCreate {
	if s != nil {
		wrac.SetSignature(*s)
	}
	return wrac
}

// SetWebhookURL sets the "webhook_url" field.
func (wrac *WebhookRetryAttemptCreate) SetWebhookURL(s string) *WebhookRetryAttemptCreate {
	wrac.mutation.SetWebhookURL(s)
	return wrac
}

// SetStatus sets the "status" field.
func (wrac *WebhookRetryAttemptCreate) SetStatus(w webhookretryattempt.Status) *WebhookRetryAttemptCreate {
	wrac.mutation.SetStatus(w)
	return wrac
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (wrac *WebhookRetryAttemptCreate) SetNillableStatus(w *webhookretryattempt.Status) *WebhookRetryAttemptCreate {
	if w != nil {
		wrac.SetStatus(*w)
	}
	return wrac
}

// Mutation returns the WebhookRetryAttemptMutation object of the builder.
func (wrac *WebhookRetryAttemptCreate) Mutation() *WebhookRetryAttemptMutation {
	return wrac.mutation
}

// Save creates the WebhookRetryAttempt in the database.
func (wrac *WebhookRetryAttemptCreate) Save(ctx context.Context) (*WebhookRetryAttempt, error) {
	wrac.defaults()
	return withHooks(ctx, wrac.sqlSave, wrac.mutation, wrac.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (wrac *WebhookRetryAttemptCreate) SaveX(ctx context.Context) *WebhookRetryAttempt {
	v, err := wrac.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (wrac *WebhookRetryAttemptCreate) Exec(ctx context.Context) error {
	_, err := wrac.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wrac *WebhookRetryAttemptCreate) ExecX(ctx context.Context) {
	if err := wrac.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (wrac *WebhookRetryAttemptCreate) defaults() {
	if _, ok := wrac.mutation.CreatedAt(); !ok {
		v := webhookretryattempt.DefaultCreatedAt()
		wrac.mutation.SetCreatedAt(v)
	}
	if _, ok := wrac.mutation.UpdatedAt(); !ok {
		v := webhookretryattempt.DefaultUpdatedAt()
		wrac.mutation.SetUpdatedAt(v)
	}
	if _, ok := wrac.mutation.NextRetryTime(); !ok {
		v := webhookretryattempt.DefaultNextRetryTime()
		wrac.mutation.SetNextRetryTime(v)
	}
	if _, ok := wrac.mutation.Status(); !ok {
		v := webhookretryattempt.DefaultStatus
		wrac.mutation.SetStatus(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (wrac *WebhookRetryAttemptCreate) check() error {
	if _, ok := wrac.mutation.CreatedAt(); !ok {
		return &ValidationError{Name: "created_at", err: errors.New(`ent: missing required field "WebhookRetryAttempt.created_at"`)}
	}
	if _, ok := wrac.mutation.UpdatedAt(); !ok {
		return &ValidationError{Name: "updated_at", err: errors.New(`ent: missing required field "WebhookRetryAttempt.updated_at"`)}
	}
	if _, ok := wrac.mutation.AttemptNumber(); !ok {
		return &ValidationError{Name: "attempt_number", err: errors.New(`ent: missing required field "WebhookRetryAttempt.attempt_number"`)}
	}
	if _, ok := wrac.mutation.NextRetryTime(); !ok {
		return &ValidationError{Name: "next_retry_time", err: errors.New(`ent: missing required field "WebhookRetryAttempt.next_retry_time"`)}
	}
	if _, ok := wrac.mutation.Payload(); !ok {
		return &ValidationError{Name: "payload", err: errors.New(`ent: missing required field "WebhookRetryAttempt.payload"`)}
	}
	if _, ok := wrac.mutation.WebhookURL(); !ok {
		return &ValidationError{Name: "webhook_url", err: errors.New(`ent: missing required field "WebhookRetryAttempt.webhook_url"`)}
	}
	if _, ok := wrac.mutation.Status(); !ok {
		return &ValidationError{Name: "status", err: errors.New(`ent: missing required field "WebhookRetryAttempt.status"`)}
	}
	if v, ok := wrac.mutation.Status(); ok {
		if err := webhookretryattempt.StatusValidator(v); err != nil {
			return &ValidationError{Name: "status", err: fmt.Errorf(`ent: validator failed for field "WebhookRetryAttempt.status": %w`, err)}
		}
	}
	return nil
}

func (wrac *WebhookRetryAttemptCreate) sqlSave(ctx context.Context) (*WebhookRetryAttempt, error) {
	if err := wrac.check(); err != nil {
		return nil, err
	}
	_node, _spec := wrac.createSpec()
	if err := sqlgraph.CreateNode(ctx, wrac.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	id := _spec.ID.Value.(int64)
	_node.ID = int(id)
	wrac.mutation.id = &_node.ID
	wrac.mutation.done = true
	return _node, nil
}

func (wrac *WebhookRetryAttemptCreate) createSpec() (*WebhookRetryAttempt, *sqlgraph.CreateSpec) {
	var (
		_node = &WebhookRetryAttempt{config: wrac.config}
		_spec = sqlgraph.NewCreateSpec(webhookretryattempt.Table, sqlgraph.NewFieldSpec(webhookretryattempt.FieldID, field.TypeInt))
	)
	if value, ok := wrac.mutation.CreatedAt(); ok {
		_spec.SetField(webhookretryattempt.FieldCreatedAt, field.TypeTime, value)
		_node.CreatedAt = value
	}
	if value, ok := wrac.mutation.UpdatedAt(); ok {
		_spec.SetField(webhookretryattempt.FieldUpdatedAt, field.TypeTime, value)
		_node.UpdatedAt = value
	}
	if value, ok := wrac.mutation.AttemptNumber(); ok {
		_spec.SetField(webhookretryattempt.FieldAttemptNumber, field.TypeInt, value)
		_node.AttemptNumber = value
	}
	if value, ok := wrac.mutation.NextRetryTime(); ok {
		_spec.SetField(webhookretryattempt.FieldNextRetryTime, field.TypeTime, value)
		_node.NextRetryTime = value
	}
	if value, ok := wrac.mutation.Payload(); ok {
		_spec.SetField(webhookretryattempt.FieldPayload, field.TypeJSON, value)
		_node.Payload = value
	}
	if value, ok := wrac.mutation.Signature(); ok {
		_spec.SetField(webhookretryattempt.FieldSignature, field.TypeString, value)
		_node.Signature = value
	}
	if value, ok := wrac.mutation.WebhookURL(); ok {
		_spec.SetField(webhookretryattempt.FieldWebhookURL, field.TypeString, value)
		_node.WebhookURL = value
	}
	if value, ok := wrac.mutation.Status(); ok {
		_spec.SetField(webhookretryattempt.FieldStatus, field.TypeEnum, value)
		_node.Status = value
	}
	return _node, _spec
}

// WebhookRetryAttemptCreateBulk is the builder for creating many WebhookRetryAttempt entities in bulk.
type WebhookRetryAttemptCreateBulk struct {
	config
	builders []*WebhookRetryAttemptCreate
}

// Save creates the WebhookRetryAttempt entities in the database.
func (wracb *WebhookRetryAttemptCreateBulk) Save(ctx context.Context) ([]*WebhookRetryAttempt, error) {
	specs := make([]*sqlgraph.CreateSpec, len(wracb.builders))
	nodes := make([]*WebhookRetryAttempt, len(wracb.builders))
	mutators := make([]Mutator, len(wracb.builders))
	for i := range wracb.builders {
		func(i int, root context.Context) {
			builder := wracb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*WebhookRetryAttemptMutation)
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
					_, err = mutators[i+1].Mutate(root, wracb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, wracb.driver, spec); err != nil {
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
		if _, err := mutators[0].Mutate(ctx, wracb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (wracb *WebhookRetryAttemptCreateBulk) SaveX(ctx context.Context) []*WebhookRetryAttempt {
	v, err := wracb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (wracb *WebhookRetryAttemptCreateBulk) Exec(ctx context.Context) error {
	_, err := wracb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (wracb *WebhookRetryAttemptCreateBulk) ExecX(ctx context.Context) {
	if err := wracb.Exec(ctx); err != nil {
		panic(err)
	}
}
