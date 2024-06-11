// Code generated by ent, DO NOT EDIT.

package hook

import (
	"context"
	"fmt"

	"github.com/paycrest/protocol/ent"
)

// The APIKeyFunc type is an adapter to allow the use of ordinary
// function as APIKey mutator.
type APIKeyFunc func(context.Context, *ent.APIKeyMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f APIKeyFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.APIKeyMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.APIKeyMutation", m)
}

// The FiatCurrencyFunc type is an adapter to allow the use of ordinary
// function as FiatCurrency mutator.
type FiatCurrencyFunc func(context.Context, *ent.FiatCurrencyMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f FiatCurrencyFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.FiatCurrencyMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.FiatCurrencyMutation", m)
}

// The InstitutionFunc type is an adapter to allow the use of ordinary
// function as Institution mutator.
type InstitutionFunc func(context.Context, *ent.InstitutionMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f InstitutionFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.InstitutionMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.InstitutionMutation", m)
}

// The LockOrderFulfillmentFunc type is an adapter to allow the use of ordinary
// function as LockOrderFulfillment mutator.
type LockOrderFulfillmentFunc func(context.Context, *ent.LockOrderFulfillmentMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f LockOrderFulfillmentFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.LockOrderFulfillmentMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.LockOrderFulfillmentMutation", m)
}

// The LockPaymentOrderFunc type is an adapter to allow the use of ordinary
// function as LockPaymentOrder mutator.
type LockPaymentOrderFunc func(context.Context, *ent.LockPaymentOrderMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f LockPaymentOrderFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.LockPaymentOrderMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.LockPaymentOrderMutation", m)
}

// The NetworkFunc type is an adapter to allow the use of ordinary
// function as Network mutator.
type NetworkFunc func(context.Context, *ent.NetworkMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f NetworkFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.NetworkMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.NetworkMutation", m)
}

// The PaymentOrderFunc type is an adapter to allow the use of ordinary
// function as PaymentOrder mutator.
type PaymentOrderFunc func(context.Context, *ent.PaymentOrderMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f PaymentOrderFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.PaymentOrderMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.PaymentOrderMutation", m)
}

// The PaymentOrderRecipientFunc type is an adapter to allow the use of ordinary
// function as PaymentOrderRecipient mutator.
type PaymentOrderRecipientFunc func(context.Context, *ent.PaymentOrderRecipientMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f PaymentOrderRecipientFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.PaymentOrderRecipientMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.PaymentOrderRecipientMutation", m)
}

// The ProviderOrderTokenFunc type is an adapter to allow the use of ordinary
// function as ProviderOrderToken mutator.
type ProviderOrderTokenFunc func(context.Context, *ent.ProviderOrderTokenMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ProviderOrderTokenFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.ProviderOrderTokenMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ProviderOrderTokenMutation", m)
}

// The ProviderProfileFunc type is an adapter to allow the use of ordinary
// function as ProviderProfile mutator.
type ProviderProfileFunc func(context.Context, *ent.ProviderProfileMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ProviderProfileFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.ProviderProfileMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ProviderProfileMutation", m)
}

// The ProviderRatingFunc type is an adapter to allow the use of ordinary
// function as ProviderRating mutator.
type ProviderRatingFunc func(context.Context, *ent.ProviderRatingMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ProviderRatingFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.ProviderRatingMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ProviderRatingMutation", m)
}

// The ProvisionBucketFunc type is an adapter to allow the use of ordinary
// function as ProvisionBucket mutator.
type ProvisionBucketFunc func(context.Context, *ent.ProvisionBucketMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ProvisionBucketFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.ProvisionBucketMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ProvisionBucketMutation", m)
}

// The ReceiveAddressFunc type is an adapter to allow the use of ordinary
// function as ReceiveAddress mutator.
type ReceiveAddressFunc func(context.Context, *ent.ReceiveAddressMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f ReceiveAddressFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.ReceiveAddressMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.ReceiveAddressMutation", m)
}

// The SenderProfileFunc type is an adapter to allow the use of ordinary
// function as SenderProfile mutator.
type SenderProfileFunc func(context.Context, *ent.SenderProfileMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f SenderProfileFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.SenderProfileMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.SenderProfileMutation", m)
}

// The TokenFunc type is an adapter to allow the use of ordinary
// function as Token mutator.
type TokenFunc func(context.Context, *ent.TokenMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f TokenFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.TokenMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.TokenMutation", m)
}

// The UserFunc type is an adapter to allow the use of ordinary
// function as User mutator.
type UserFunc func(context.Context, *ent.UserMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f UserFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.UserMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.UserMutation", m)
}

// The VerificationTokenFunc type is an adapter to allow the use of ordinary
// function as VerificationToken mutator.
type VerificationTokenFunc func(context.Context, *ent.VerificationTokenMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f VerificationTokenFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.VerificationTokenMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.VerificationTokenMutation", m)
}

// The WebhookRetryAttemptFunc type is an adapter to allow the use of ordinary
// function as WebhookRetryAttempt mutator.
type WebhookRetryAttemptFunc func(context.Context, *ent.WebhookRetryAttemptMutation) (ent.Value, error)

// Mutate calls f(ctx, m).
func (f WebhookRetryAttemptFunc) Mutate(ctx context.Context, m ent.Mutation) (ent.Value, error) {
	if mv, ok := m.(*ent.WebhookRetryAttemptMutation); ok {
		return f(ctx, mv)
	}
	return nil, fmt.Errorf("unexpected mutation type %T. expect *ent.WebhookRetryAttemptMutation", m)
}

// Condition is a hook condition function.
type Condition func(context.Context, ent.Mutation) bool

// And groups conditions with the AND operator.
func And(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		if !first(ctx, m) || !second(ctx, m) {
			return false
		}
		for _, cond := range rest {
			if !cond(ctx, m) {
				return false
			}
		}
		return true
	}
}

// Or groups conditions with the OR operator.
func Or(first, second Condition, rest ...Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		if first(ctx, m) || second(ctx, m) {
			return true
		}
		for _, cond := range rest {
			if cond(ctx, m) {
				return true
			}
		}
		return false
	}
}

// Not negates a given condition.
func Not(cond Condition) Condition {
	return func(ctx context.Context, m ent.Mutation) bool {
		return !cond(ctx, m)
	}
}

// HasOp is a condition testing mutation operation.
func HasOp(op ent.Op) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		return m.Op().Is(op)
	}
}

// HasAddedFields is a condition validating `.AddedField` on fields.
func HasAddedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if _, exists := m.AddedField(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.AddedField(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasClearedFields is a condition validating `.FieldCleared` on fields.
func HasClearedFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if exists := m.FieldCleared(field); !exists {
			return false
		}
		for _, field := range fields {
			if exists := m.FieldCleared(field); !exists {
				return false
			}
		}
		return true
	}
}

// HasFields is a condition validating `.Field` on fields.
func HasFields(field string, fields ...string) Condition {
	return func(_ context.Context, m ent.Mutation) bool {
		if _, exists := m.Field(field); !exists {
			return false
		}
		for _, field := range fields {
			if _, exists := m.Field(field); !exists {
				return false
			}
		}
		return true
	}
}

// If executes the given hook under condition.
//
//	hook.If(ComputeAverage, And(HasFields(...), HasAddedFields(...)))
func If(hk ent.Hook, cond Condition) ent.Hook {
	return func(next ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(ctx context.Context, m ent.Mutation) (ent.Value, error) {
			if cond(ctx, m) {
				return hk(next).Mutate(ctx, m)
			}
			return next.Mutate(ctx, m)
		})
	}
}

// On executes the given hook only for the given operation.
//
//	hook.On(Log, ent.Delete|ent.Create)
func On(hk ent.Hook, op ent.Op) ent.Hook {
	return If(hk, HasOp(op))
}

// Unless skips the given hook only for the given operation.
//
//	hook.Unless(Log, ent.Update|ent.UpdateOne)
func Unless(hk ent.Hook, op ent.Op) ent.Hook {
	return If(hk, Not(HasOp(op)))
}

// FixedError is a hook returning a fixed error.
func FixedError(err error) ent.Hook {
	return func(ent.Mutator) ent.Mutator {
		return ent.MutateFunc(func(context.Context, ent.Mutation) (ent.Value, error) {
			return nil, err
		})
	}
}

// Reject returns a hook that rejects all operations that match op.
//
//	func (T) Hooks() []ent.Hook {
//		return []ent.Hook{
//			Reject(ent.Delete|ent.Update),
//		}
//	}
func Reject(op ent.Op) ent.Hook {
	hk := FixedError(fmt.Errorf("%s operation is not allowed", op))
	return On(hk, op)
}

// Chain acts as a list of hooks and is effectively immutable.
// Once created, it will always hold the same set of hooks in the same order.
type Chain struct {
	hooks []ent.Hook
}

// NewChain creates a new chain of hooks.
func NewChain(hooks ...ent.Hook) Chain {
	return Chain{append([]ent.Hook(nil), hooks...)}
}

// Hook chains the list of hooks and returns the final hook.
func (c Chain) Hook() ent.Hook {
	return func(mutator ent.Mutator) ent.Mutator {
		for i := len(c.hooks) - 1; i >= 0; i-- {
			mutator = c.hooks[i](mutator)
		}
		return mutator
	}
}

// Append extends a chain, adding the specified hook
// as the last ones in the mutation flow.
func (c Chain) Append(hooks ...ent.Hook) Chain {
	newHooks := make([]ent.Hook, 0, len(c.hooks)+len(hooks))
	newHooks = append(newHooks, c.hooks...)
	newHooks = append(newHooks, hooks...)
	return Chain{newHooks}
}

// Extend extends a chain, adding the specified chain
// as the last ones in the mutation flow.
func (c Chain) Extend(chain Chain) Chain {
	return c.Append(chain.hooks...)
}
