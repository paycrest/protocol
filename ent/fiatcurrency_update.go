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
	"github.com/paycrest/aggregator/ent/fiatcurrency"
	"github.com/paycrest/aggregator/ent/institution"
	"github.com/paycrest/aggregator/ent/predicate"
	"github.com/paycrest/aggregator/ent/providerordertoken"
	"github.com/paycrest/aggregator/ent/providerprofile"
	"github.com/paycrest/aggregator/ent/provisionbucket"
	"github.com/shopspring/decimal"
)

// FiatCurrencyUpdate is the builder for updating FiatCurrency entities.
type FiatCurrencyUpdate struct {
	config
	hooks    []Hook
	mutation *FiatCurrencyMutation
}

// Where appends a list predicates to the FiatCurrencyUpdate builder.
func (fcu *FiatCurrencyUpdate) Where(ps ...predicate.FiatCurrency) *FiatCurrencyUpdate {
	fcu.mutation.Where(ps...)
	return fcu
}

// SetUpdatedAt sets the "updated_at" field.
func (fcu *FiatCurrencyUpdate) SetUpdatedAt(t time.Time) *FiatCurrencyUpdate {
	fcu.mutation.SetUpdatedAt(t)
	return fcu
}

// SetCode sets the "code" field.
func (fcu *FiatCurrencyUpdate) SetCode(s string) *FiatCurrencyUpdate {
	fcu.mutation.SetCode(s)
	return fcu
}

// SetNillableCode sets the "code" field if the given value is not nil.
func (fcu *FiatCurrencyUpdate) SetNillableCode(s *string) *FiatCurrencyUpdate {
	if s != nil {
		fcu.SetCode(*s)
	}
	return fcu
}

// SetShortName sets the "short_name" field.
func (fcu *FiatCurrencyUpdate) SetShortName(s string) *FiatCurrencyUpdate {
	fcu.mutation.SetShortName(s)
	return fcu
}

// SetNillableShortName sets the "short_name" field if the given value is not nil.
func (fcu *FiatCurrencyUpdate) SetNillableShortName(s *string) *FiatCurrencyUpdate {
	if s != nil {
		fcu.SetShortName(*s)
	}
	return fcu
}

// SetDecimals sets the "decimals" field.
func (fcu *FiatCurrencyUpdate) SetDecimals(i int) *FiatCurrencyUpdate {
	fcu.mutation.ResetDecimals()
	fcu.mutation.SetDecimals(i)
	return fcu
}

// SetNillableDecimals sets the "decimals" field if the given value is not nil.
func (fcu *FiatCurrencyUpdate) SetNillableDecimals(i *int) *FiatCurrencyUpdate {
	if i != nil {
		fcu.SetDecimals(*i)
	}
	return fcu
}

// AddDecimals adds i to the "decimals" field.
func (fcu *FiatCurrencyUpdate) AddDecimals(i int) *FiatCurrencyUpdate {
	fcu.mutation.AddDecimals(i)
	return fcu
}

// SetSymbol sets the "symbol" field.
func (fcu *FiatCurrencyUpdate) SetSymbol(s string) *FiatCurrencyUpdate {
	fcu.mutation.SetSymbol(s)
	return fcu
}

// SetNillableSymbol sets the "symbol" field if the given value is not nil.
func (fcu *FiatCurrencyUpdate) SetNillableSymbol(s *string) *FiatCurrencyUpdate {
	if s != nil {
		fcu.SetSymbol(*s)
	}
	return fcu
}

// SetName sets the "name" field.
func (fcu *FiatCurrencyUpdate) SetName(s string) *FiatCurrencyUpdate {
	fcu.mutation.SetName(s)
	return fcu
}

// SetNillableName sets the "name" field if the given value is not nil.
func (fcu *FiatCurrencyUpdate) SetNillableName(s *string) *FiatCurrencyUpdate {
	if s != nil {
		fcu.SetName(*s)
	}
	return fcu
}

// SetMarketRate sets the "market_rate" field.
func (fcu *FiatCurrencyUpdate) SetMarketRate(d decimal.Decimal) *FiatCurrencyUpdate {
	fcu.mutation.ResetMarketRate()
	fcu.mutation.SetMarketRate(d)
	return fcu
}

// SetNillableMarketRate sets the "market_rate" field if the given value is not nil.
func (fcu *FiatCurrencyUpdate) SetNillableMarketRate(d *decimal.Decimal) *FiatCurrencyUpdate {
	if d != nil {
		fcu.SetMarketRate(*d)
	}
	return fcu
}

// AddMarketRate adds d to the "market_rate" field.
func (fcu *FiatCurrencyUpdate) AddMarketRate(d decimal.Decimal) *FiatCurrencyUpdate {
	fcu.mutation.AddMarketRate(d)
	return fcu
}

// SetIsEnabled sets the "is_enabled" field.
func (fcu *FiatCurrencyUpdate) SetIsEnabled(b bool) *FiatCurrencyUpdate {
	fcu.mutation.SetIsEnabled(b)
	return fcu
}

// SetNillableIsEnabled sets the "is_enabled" field if the given value is not nil.
func (fcu *FiatCurrencyUpdate) SetNillableIsEnabled(b *bool) *FiatCurrencyUpdate {
	if b != nil {
		fcu.SetIsEnabled(*b)
	}
	return fcu
}

// AddProviderIDs adds the "providers" edge to the ProviderProfile entity by IDs.
func (fcu *FiatCurrencyUpdate) AddProviderIDs(ids ...string) *FiatCurrencyUpdate {
	fcu.mutation.AddProviderIDs(ids...)
	return fcu
}

// AddProviders adds the "providers" edges to the ProviderProfile entity.
func (fcu *FiatCurrencyUpdate) AddProviders(p ...*ProviderProfile) *FiatCurrencyUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fcu.AddProviderIDs(ids...)
}

// AddProvisionBucketIDs adds the "provision_buckets" edge to the ProvisionBucket entity by IDs.
func (fcu *FiatCurrencyUpdate) AddProvisionBucketIDs(ids ...int) *FiatCurrencyUpdate {
	fcu.mutation.AddProvisionBucketIDs(ids...)
	return fcu
}

// AddProvisionBuckets adds the "provision_buckets" edges to the ProvisionBucket entity.
func (fcu *FiatCurrencyUpdate) AddProvisionBuckets(p ...*ProvisionBucket) *FiatCurrencyUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fcu.AddProvisionBucketIDs(ids...)
}

// AddInstitutionIDs adds the "institutions" edge to the Institution entity by IDs.
func (fcu *FiatCurrencyUpdate) AddInstitutionIDs(ids ...int) *FiatCurrencyUpdate {
	fcu.mutation.AddInstitutionIDs(ids...)
	return fcu
}

// AddInstitutions adds the "institutions" edges to the Institution entity.
func (fcu *FiatCurrencyUpdate) AddInstitutions(i ...*Institution) *FiatCurrencyUpdate {
	ids := make([]int, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return fcu.AddInstitutionIDs(ids...)
}

// AddProviderSettingIDs adds the "provider_settings" edge to the ProviderOrderToken entity by IDs.
func (fcu *FiatCurrencyUpdate) AddProviderSettingIDs(ids ...int) *FiatCurrencyUpdate {
	fcu.mutation.AddProviderSettingIDs(ids...)
	return fcu
}

// AddProviderSettings adds the "provider_settings" edges to the ProviderOrderToken entity.
func (fcu *FiatCurrencyUpdate) AddProviderSettings(p ...*ProviderOrderToken) *FiatCurrencyUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fcu.AddProviderSettingIDs(ids...)
}

// Mutation returns the FiatCurrencyMutation object of the builder.
func (fcu *FiatCurrencyUpdate) Mutation() *FiatCurrencyMutation {
	return fcu.mutation
}

// ClearProviders clears all "providers" edges to the ProviderProfile entity.
func (fcu *FiatCurrencyUpdate) ClearProviders() *FiatCurrencyUpdate {
	fcu.mutation.ClearProviders()
	return fcu
}

// RemoveProviderIDs removes the "providers" edge to ProviderProfile entities by IDs.
func (fcu *FiatCurrencyUpdate) RemoveProviderIDs(ids ...string) *FiatCurrencyUpdate {
	fcu.mutation.RemoveProviderIDs(ids...)
	return fcu
}

// RemoveProviders removes "providers" edges to ProviderProfile entities.
func (fcu *FiatCurrencyUpdate) RemoveProviders(p ...*ProviderProfile) *FiatCurrencyUpdate {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fcu.RemoveProviderIDs(ids...)
}

// ClearProvisionBuckets clears all "provision_buckets" edges to the ProvisionBucket entity.
func (fcu *FiatCurrencyUpdate) ClearProvisionBuckets() *FiatCurrencyUpdate {
	fcu.mutation.ClearProvisionBuckets()
	return fcu
}

// RemoveProvisionBucketIDs removes the "provision_buckets" edge to ProvisionBucket entities by IDs.
func (fcu *FiatCurrencyUpdate) RemoveProvisionBucketIDs(ids ...int) *FiatCurrencyUpdate {
	fcu.mutation.RemoveProvisionBucketIDs(ids...)
	return fcu
}

// RemoveProvisionBuckets removes "provision_buckets" edges to ProvisionBucket entities.
func (fcu *FiatCurrencyUpdate) RemoveProvisionBuckets(p ...*ProvisionBucket) *FiatCurrencyUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fcu.RemoveProvisionBucketIDs(ids...)
}

// ClearInstitutions clears all "institutions" edges to the Institution entity.
func (fcu *FiatCurrencyUpdate) ClearInstitutions() *FiatCurrencyUpdate {
	fcu.mutation.ClearInstitutions()
	return fcu
}

// RemoveInstitutionIDs removes the "institutions" edge to Institution entities by IDs.
func (fcu *FiatCurrencyUpdate) RemoveInstitutionIDs(ids ...int) *FiatCurrencyUpdate {
	fcu.mutation.RemoveInstitutionIDs(ids...)
	return fcu
}

// RemoveInstitutions removes "institutions" edges to Institution entities.
func (fcu *FiatCurrencyUpdate) RemoveInstitutions(i ...*Institution) *FiatCurrencyUpdate {
	ids := make([]int, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return fcu.RemoveInstitutionIDs(ids...)
}

// ClearProviderSettings clears all "provider_settings" edges to the ProviderOrderToken entity.
func (fcu *FiatCurrencyUpdate) ClearProviderSettings() *FiatCurrencyUpdate {
	fcu.mutation.ClearProviderSettings()
	return fcu
}

// RemoveProviderSettingIDs removes the "provider_settings" edge to ProviderOrderToken entities by IDs.
func (fcu *FiatCurrencyUpdate) RemoveProviderSettingIDs(ids ...int) *FiatCurrencyUpdate {
	fcu.mutation.RemoveProviderSettingIDs(ids...)
	return fcu
}

// RemoveProviderSettings removes "provider_settings" edges to ProviderOrderToken entities.
func (fcu *FiatCurrencyUpdate) RemoveProviderSettings(p ...*ProviderOrderToken) *FiatCurrencyUpdate {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fcu.RemoveProviderSettingIDs(ids...)
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (fcu *FiatCurrencyUpdate) Save(ctx context.Context) (int, error) {
	fcu.defaults()
	return withHooks(ctx, fcu.sqlSave, fcu.mutation, fcu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (fcu *FiatCurrencyUpdate) SaveX(ctx context.Context) int {
	affected, err := fcu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (fcu *FiatCurrencyUpdate) Exec(ctx context.Context) error {
	_, err := fcu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fcu *FiatCurrencyUpdate) ExecX(ctx context.Context) {
	if err := fcu.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fcu *FiatCurrencyUpdate) defaults() {
	if _, ok := fcu.mutation.UpdatedAt(); !ok {
		v := fiatcurrency.UpdateDefaultUpdatedAt()
		fcu.mutation.SetUpdatedAt(v)
	}
}

func (fcu *FiatCurrencyUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(fiatcurrency.Table, fiatcurrency.Columns, sqlgraph.NewFieldSpec(fiatcurrency.FieldID, field.TypeUUID))
	if ps := fcu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fcu.mutation.UpdatedAt(); ok {
		_spec.SetField(fiatcurrency.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := fcu.mutation.Code(); ok {
		_spec.SetField(fiatcurrency.FieldCode, field.TypeString, value)
	}
	if value, ok := fcu.mutation.ShortName(); ok {
		_spec.SetField(fiatcurrency.FieldShortName, field.TypeString, value)
	}
	if value, ok := fcu.mutation.Decimals(); ok {
		_spec.SetField(fiatcurrency.FieldDecimals, field.TypeInt, value)
	}
	if value, ok := fcu.mutation.AddedDecimals(); ok {
		_spec.AddField(fiatcurrency.FieldDecimals, field.TypeInt, value)
	}
	if value, ok := fcu.mutation.Symbol(); ok {
		_spec.SetField(fiatcurrency.FieldSymbol, field.TypeString, value)
	}
	if value, ok := fcu.mutation.Name(); ok {
		_spec.SetField(fiatcurrency.FieldName, field.TypeString, value)
	}
	if value, ok := fcu.mutation.MarketRate(); ok {
		_spec.SetField(fiatcurrency.FieldMarketRate, field.TypeFloat64, value)
	}
	if value, ok := fcu.mutation.AddedMarketRate(); ok {
		_spec.AddField(fiatcurrency.FieldMarketRate, field.TypeFloat64, value)
	}
	if value, ok := fcu.mutation.IsEnabled(); ok {
		_spec.SetField(fiatcurrency.FieldIsEnabled, field.TypeBool, value)
	}
	if fcu.mutation.ProvidersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   fiatcurrency.ProvidersTable,
			Columns: fiatcurrency.ProvidersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerprofile.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcu.mutation.RemovedProvidersIDs(); len(nodes) > 0 && !fcu.mutation.ProvidersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   fiatcurrency.ProvidersTable,
			Columns: fiatcurrency.ProvidersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerprofile.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcu.mutation.ProvidersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   fiatcurrency.ProvidersTable,
			Columns: fiatcurrency.ProvidersPrimaryKey,
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
	if fcu.mutation.ProvisionBucketsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.ProvisionBucketsTable,
			Columns: []string{fiatcurrency.ProvisionBucketsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(provisionbucket.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcu.mutation.RemovedProvisionBucketsIDs(); len(nodes) > 0 && !fcu.mutation.ProvisionBucketsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.ProvisionBucketsTable,
			Columns: []string{fiatcurrency.ProvisionBucketsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(provisionbucket.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcu.mutation.ProvisionBucketsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.ProvisionBucketsTable,
			Columns: []string{fiatcurrency.ProvisionBucketsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(provisionbucket.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fcu.mutation.InstitutionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.InstitutionsTable,
			Columns: []string{fiatcurrency.InstitutionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(institution.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcu.mutation.RemovedInstitutionsIDs(); len(nodes) > 0 && !fcu.mutation.InstitutionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.InstitutionsTable,
			Columns: []string{fiatcurrency.InstitutionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(institution.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcu.mutation.InstitutionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.InstitutionsTable,
			Columns: []string{fiatcurrency.InstitutionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(institution.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fcu.mutation.ProviderSettingsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.ProviderSettingsTable,
			Columns: []string{fiatcurrency.ProviderSettingsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerordertoken.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcu.mutation.RemovedProviderSettingsIDs(); len(nodes) > 0 && !fcu.mutation.ProviderSettingsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.ProviderSettingsTable,
			Columns: []string{fiatcurrency.ProviderSettingsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerordertoken.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcu.mutation.ProviderSettingsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.ProviderSettingsTable,
			Columns: []string{fiatcurrency.ProviderSettingsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerordertoken.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, fcu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{fiatcurrency.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	fcu.mutation.done = true
	return n, nil
}

// FiatCurrencyUpdateOne is the builder for updating a single FiatCurrency entity.
type FiatCurrencyUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *FiatCurrencyMutation
}

// SetUpdatedAt sets the "updated_at" field.
func (fcuo *FiatCurrencyUpdateOne) SetUpdatedAt(t time.Time) *FiatCurrencyUpdateOne {
	fcuo.mutation.SetUpdatedAt(t)
	return fcuo
}

// SetCode sets the "code" field.
func (fcuo *FiatCurrencyUpdateOne) SetCode(s string) *FiatCurrencyUpdateOne {
	fcuo.mutation.SetCode(s)
	return fcuo
}

// SetNillableCode sets the "code" field if the given value is not nil.
func (fcuo *FiatCurrencyUpdateOne) SetNillableCode(s *string) *FiatCurrencyUpdateOne {
	if s != nil {
		fcuo.SetCode(*s)
	}
	return fcuo
}

// SetShortName sets the "short_name" field.
func (fcuo *FiatCurrencyUpdateOne) SetShortName(s string) *FiatCurrencyUpdateOne {
	fcuo.mutation.SetShortName(s)
	return fcuo
}

// SetNillableShortName sets the "short_name" field if the given value is not nil.
func (fcuo *FiatCurrencyUpdateOne) SetNillableShortName(s *string) *FiatCurrencyUpdateOne {
	if s != nil {
		fcuo.SetShortName(*s)
	}
	return fcuo
}

// SetDecimals sets the "decimals" field.
func (fcuo *FiatCurrencyUpdateOne) SetDecimals(i int) *FiatCurrencyUpdateOne {
	fcuo.mutation.ResetDecimals()
	fcuo.mutation.SetDecimals(i)
	return fcuo
}

// SetNillableDecimals sets the "decimals" field if the given value is not nil.
func (fcuo *FiatCurrencyUpdateOne) SetNillableDecimals(i *int) *FiatCurrencyUpdateOne {
	if i != nil {
		fcuo.SetDecimals(*i)
	}
	return fcuo
}

// AddDecimals adds i to the "decimals" field.
func (fcuo *FiatCurrencyUpdateOne) AddDecimals(i int) *FiatCurrencyUpdateOne {
	fcuo.mutation.AddDecimals(i)
	return fcuo
}

// SetSymbol sets the "symbol" field.
func (fcuo *FiatCurrencyUpdateOne) SetSymbol(s string) *FiatCurrencyUpdateOne {
	fcuo.mutation.SetSymbol(s)
	return fcuo
}

// SetNillableSymbol sets the "symbol" field if the given value is not nil.
func (fcuo *FiatCurrencyUpdateOne) SetNillableSymbol(s *string) *FiatCurrencyUpdateOne {
	if s != nil {
		fcuo.SetSymbol(*s)
	}
	return fcuo
}

// SetName sets the "name" field.
func (fcuo *FiatCurrencyUpdateOne) SetName(s string) *FiatCurrencyUpdateOne {
	fcuo.mutation.SetName(s)
	return fcuo
}

// SetNillableName sets the "name" field if the given value is not nil.
func (fcuo *FiatCurrencyUpdateOne) SetNillableName(s *string) *FiatCurrencyUpdateOne {
	if s != nil {
		fcuo.SetName(*s)
	}
	return fcuo
}

// SetMarketRate sets the "market_rate" field.
func (fcuo *FiatCurrencyUpdateOne) SetMarketRate(d decimal.Decimal) *FiatCurrencyUpdateOne {
	fcuo.mutation.ResetMarketRate()
	fcuo.mutation.SetMarketRate(d)
	return fcuo
}

// SetNillableMarketRate sets the "market_rate" field if the given value is not nil.
func (fcuo *FiatCurrencyUpdateOne) SetNillableMarketRate(d *decimal.Decimal) *FiatCurrencyUpdateOne {
	if d != nil {
		fcuo.SetMarketRate(*d)
	}
	return fcuo
}

// AddMarketRate adds d to the "market_rate" field.
func (fcuo *FiatCurrencyUpdateOne) AddMarketRate(d decimal.Decimal) *FiatCurrencyUpdateOne {
	fcuo.mutation.AddMarketRate(d)
	return fcuo
}

// SetIsEnabled sets the "is_enabled" field.
func (fcuo *FiatCurrencyUpdateOne) SetIsEnabled(b bool) *FiatCurrencyUpdateOne {
	fcuo.mutation.SetIsEnabled(b)
	return fcuo
}

// SetNillableIsEnabled sets the "is_enabled" field if the given value is not nil.
func (fcuo *FiatCurrencyUpdateOne) SetNillableIsEnabled(b *bool) *FiatCurrencyUpdateOne {
	if b != nil {
		fcuo.SetIsEnabled(*b)
	}
	return fcuo
}

// AddProviderIDs adds the "providers" edge to the ProviderProfile entity by IDs.
func (fcuo *FiatCurrencyUpdateOne) AddProviderIDs(ids ...string) *FiatCurrencyUpdateOne {
	fcuo.mutation.AddProviderIDs(ids...)
	return fcuo
}

// AddProviders adds the "providers" edges to the ProviderProfile entity.
func (fcuo *FiatCurrencyUpdateOne) AddProviders(p ...*ProviderProfile) *FiatCurrencyUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fcuo.AddProviderIDs(ids...)
}

// AddProvisionBucketIDs adds the "provision_buckets" edge to the ProvisionBucket entity by IDs.
func (fcuo *FiatCurrencyUpdateOne) AddProvisionBucketIDs(ids ...int) *FiatCurrencyUpdateOne {
	fcuo.mutation.AddProvisionBucketIDs(ids...)
	return fcuo
}

// AddProvisionBuckets adds the "provision_buckets" edges to the ProvisionBucket entity.
func (fcuo *FiatCurrencyUpdateOne) AddProvisionBuckets(p ...*ProvisionBucket) *FiatCurrencyUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fcuo.AddProvisionBucketIDs(ids...)
}

// AddInstitutionIDs adds the "institutions" edge to the Institution entity by IDs.
func (fcuo *FiatCurrencyUpdateOne) AddInstitutionIDs(ids ...int) *FiatCurrencyUpdateOne {
	fcuo.mutation.AddInstitutionIDs(ids...)
	return fcuo
}

// AddInstitutions adds the "institutions" edges to the Institution entity.
func (fcuo *FiatCurrencyUpdateOne) AddInstitutions(i ...*Institution) *FiatCurrencyUpdateOne {
	ids := make([]int, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return fcuo.AddInstitutionIDs(ids...)
}

// AddProviderSettingIDs adds the "provider_settings" edge to the ProviderOrderToken entity by IDs.
func (fcuo *FiatCurrencyUpdateOne) AddProviderSettingIDs(ids ...int) *FiatCurrencyUpdateOne {
	fcuo.mutation.AddProviderSettingIDs(ids...)
	return fcuo
}

// AddProviderSettings adds the "provider_settings" edges to the ProviderOrderToken entity.
func (fcuo *FiatCurrencyUpdateOne) AddProviderSettings(p ...*ProviderOrderToken) *FiatCurrencyUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fcuo.AddProviderSettingIDs(ids...)
}

// Mutation returns the FiatCurrencyMutation object of the builder.
func (fcuo *FiatCurrencyUpdateOne) Mutation() *FiatCurrencyMutation {
	return fcuo.mutation
}

// ClearProviders clears all "providers" edges to the ProviderProfile entity.
func (fcuo *FiatCurrencyUpdateOne) ClearProviders() *FiatCurrencyUpdateOne {
	fcuo.mutation.ClearProviders()
	return fcuo
}

// RemoveProviderIDs removes the "providers" edge to ProviderProfile entities by IDs.
func (fcuo *FiatCurrencyUpdateOne) RemoveProviderIDs(ids ...string) *FiatCurrencyUpdateOne {
	fcuo.mutation.RemoveProviderIDs(ids...)
	return fcuo
}

// RemoveProviders removes "providers" edges to ProviderProfile entities.
func (fcuo *FiatCurrencyUpdateOne) RemoveProviders(p ...*ProviderProfile) *FiatCurrencyUpdateOne {
	ids := make([]string, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fcuo.RemoveProviderIDs(ids...)
}

// ClearProvisionBuckets clears all "provision_buckets" edges to the ProvisionBucket entity.
func (fcuo *FiatCurrencyUpdateOne) ClearProvisionBuckets() *FiatCurrencyUpdateOne {
	fcuo.mutation.ClearProvisionBuckets()
	return fcuo
}

// RemoveProvisionBucketIDs removes the "provision_buckets" edge to ProvisionBucket entities by IDs.
func (fcuo *FiatCurrencyUpdateOne) RemoveProvisionBucketIDs(ids ...int) *FiatCurrencyUpdateOne {
	fcuo.mutation.RemoveProvisionBucketIDs(ids...)
	return fcuo
}

// RemoveProvisionBuckets removes "provision_buckets" edges to ProvisionBucket entities.
func (fcuo *FiatCurrencyUpdateOne) RemoveProvisionBuckets(p ...*ProvisionBucket) *FiatCurrencyUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fcuo.RemoveProvisionBucketIDs(ids...)
}

// ClearInstitutions clears all "institutions" edges to the Institution entity.
func (fcuo *FiatCurrencyUpdateOne) ClearInstitutions() *FiatCurrencyUpdateOne {
	fcuo.mutation.ClearInstitutions()
	return fcuo
}

// RemoveInstitutionIDs removes the "institutions" edge to Institution entities by IDs.
func (fcuo *FiatCurrencyUpdateOne) RemoveInstitutionIDs(ids ...int) *FiatCurrencyUpdateOne {
	fcuo.mutation.RemoveInstitutionIDs(ids...)
	return fcuo
}

// RemoveInstitutions removes "institutions" edges to Institution entities.
func (fcuo *FiatCurrencyUpdateOne) RemoveInstitutions(i ...*Institution) *FiatCurrencyUpdateOne {
	ids := make([]int, len(i))
	for j := range i {
		ids[j] = i[j].ID
	}
	return fcuo.RemoveInstitutionIDs(ids...)
}

// ClearProviderSettings clears all "provider_settings" edges to the ProviderOrderToken entity.
func (fcuo *FiatCurrencyUpdateOne) ClearProviderSettings() *FiatCurrencyUpdateOne {
	fcuo.mutation.ClearProviderSettings()
	return fcuo
}

// RemoveProviderSettingIDs removes the "provider_settings" edge to ProviderOrderToken entities by IDs.
func (fcuo *FiatCurrencyUpdateOne) RemoveProviderSettingIDs(ids ...int) *FiatCurrencyUpdateOne {
	fcuo.mutation.RemoveProviderSettingIDs(ids...)
	return fcuo
}

// RemoveProviderSettings removes "provider_settings" edges to ProviderOrderToken entities.
func (fcuo *FiatCurrencyUpdateOne) RemoveProviderSettings(p ...*ProviderOrderToken) *FiatCurrencyUpdateOne {
	ids := make([]int, len(p))
	for i := range p {
		ids[i] = p[i].ID
	}
	return fcuo.RemoveProviderSettingIDs(ids...)
}

// Where appends a list predicates to the FiatCurrencyUpdate builder.
func (fcuo *FiatCurrencyUpdateOne) Where(ps ...predicate.FiatCurrency) *FiatCurrencyUpdateOne {
	fcuo.mutation.Where(ps...)
	return fcuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (fcuo *FiatCurrencyUpdateOne) Select(field string, fields ...string) *FiatCurrencyUpdateOne {
	fcuo.fields = append([]string{field}, fields...)
	return fcuo
}

// Save executes the query and returns the updated FiatCurrency entity.
func (fcuo *FiatCurrencyUpdateOne) Save(ctx context.Context) (*FiatCurrency, error) {
	fcuo.defaults()
	return withHooks(ctx, fcuo.sqlSave, fcuo.mutation, fcuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (fcuo *FiatCurrencyUpdateOne) SaveX(ctx context.Context) *FiatCurrency {
	node, err := fcuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (fcuo *FiatCurrencyUpdateOne) Exec(ctx context.Context) error {
	_, err := fcuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (fcuo *FiatCurrencyUpdateOne) ExecX(ctx context.Context) {
	if err := fcuo.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (fcuo *FiatCurrencyUpdateOne) defaults() {
	if _, ok := fcuo.mutation.UpdatedAt(); !ok {
		v := fiatcurrency.UpdateDefaultUpdatedAt()
		fcuo.mutation.SetUpdatedAt(v)
	}
}

func (fcuo *FiatCurrencyUpdateOne) sqlSave(ctx context.Context) (_node *FiatCurrency, err error) {
	_spec := sqlgraph.NewUpdateSpec(fiatcurrency.Table, fiatcurrency.Columns, sqlgraph.NewFieldSpec(fiatcurrency.FieldID, field.TypeUUID))
	id, ok := fcuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "FiatCurrency.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := fcuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, fiatcurrency.FieldID)
		for _, f := range fields {
			if !fiatcurrency.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != fiatcurrency.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := fcuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := fcuo.mutation.UpdatedAt(); ok {
		_spec.SetField(fiatcurrency.FieldUpdatedAt, field.TypeTime, value)
	}
	if value, ok := fcuo.mutation.Code(); ok {
		_spec.SetField(fiatcurrency.FieldCode, field.TypeString, value)
	}
	if value, ok := fcuo.mutation.ShortName(); ok {
		_spec.SetField(fiatcurrency.FieldShortName, field.TypeString, value)
	}
	if value, ok := fcuo.mutation.Decimals(); ok {
		_spec.SetField(fiatcurrency.FieldDecimals, field.TypeInt, value)
	}
	if value, ok := fcuo.mutation.AddedDecimals(); ok {
		_spec.AddField(fiatcurrency.FieldDecimals, field.TypeInt, value)
	}
	if value, ok := fcuo.mutation.Symbol(); ok {
		_spec.SetField(fiatcurrency.FieldSymbol, field.TypeString, value)
	}
	if value, ok := fcuo.mutation.Name(); ok {
		_spec.SetField(fiatcurrency.FieldName, field.TypeString, value)
	}
	if value, ok := fcuo.mutation.MarketRate(); ok {
		_spec.SetField(fiatcurrency.FieldMarketRate, field.TypeFloat64, value)
	}
	if value, ok := fcuo.mutation.AddedMarketRate(); ok {
		_spec.AddField(fiatcurrency.FieldMarketRate, field.TypeFloat64, value)
	}
	if value, ok := fcuo.mutation.IsEnabled(); ok {
		_spec.SetField(fiatcurrency.FieldIsEnabled, field.TypeBool, value)
	}
	if fcuo.mutation.ProvidersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   fiatcurrency.ProvidersTable,
			Columns: fiatcurrency.ProvidersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerprofile.FieldID, field.TypeString),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcuo.mutation.RemovedProvidersIDs(); len(nodes) > 0 && !fcuo.mutation.ProvidersCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   fiatcurrency.ProvidersTable,
			Columns: fiatcurrency.ProvidersPrimaryKey,
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerprofile.FieldID, field.TypeString),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcuo.mutation.ProvidersIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.M2M,
			Inverse: false,
			Table:   fiatcurrency.ProvidersTable,
			Columns: fiatcurrency.ProvidersPrimaryKey,
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
	if fcuo.mutation.ProvisionBucketsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.ProvisionBucketsTable,
			Columns: []string{fiatcurrency.ProvisionBucketsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(provisionbucket.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcuo.mutation.RemovedProvisionBucketsIDs(); len(nodes) > 0 && !fcuo.mutation.ProvisionBucketsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.ProvisionBucketsTable,
			Columns: []string{fiatcurrency.ProvisionBucketsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(provisionbucket.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcuo.mutation.ProvisionBucketsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.ProvisionBucketsTable,
			Columns: []string{fiatcurrency.ProvisionBucketsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(provisionbucket.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fcuo.mutation.InstitutionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.InstitutionsTable,
			Columns: []string{fiatcurrency.InstitutionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(institution.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcuo.mutation.RemovedInstitutionsIDs(); len(nodes) > 0 && !fcuo.mutation.InstitutionsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.InstitutionsTable,
			Columns: []string{fiatcurrency.InstitutionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(institution.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcuo.mutation.InstitutionsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.InstitutionsTable,
			Columns: []string{fiatcurrency.InstitutionsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(institution.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	if fcuo.mutation.ProviderSettingsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.ProviderSettingsTable,
			Columns: []string{fiatcurrency.ProviderSettingsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerordertoken.FieldID, field.TypeInt),
			},
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcuo.mutation.RemovedProviderSettingsIDs(); len(nodes) > 0 && !fcuo.mutation.ProviderSettingsCleared() {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.ProviderSettingsTable,
			Columns: []string{fiatcurrency.ProviderSettingsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerordertoken.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Clear = append(_spec.Edges.Clear, edge)
	}
	if nodes := fcuo.mutation.ProviderSettingsIDs(); len(nodes) > 0 {
		edge := &sqlgraph.EdgeSpec{
			Rel:     sqlgraph.O2M,
			Inverse: false,
			Table:   fiatcurrency.ProviderSettingsTable,
			Columns: []string{fiatcurrency.ProviderSettingsColumn},
			Bidi:    false,
			Target: &sqlgraph.EdgeTarget{
				IDSpec: sqlgraph.NewFieldSpec(providerordertoken.FieldID, field.TypeInt),
			},
		}
		for _, k := range nodes {
			edge.Target.Nodes = append(edge.Target.Nodes, k)
		}
		_spec.Edges.Add = append(_spec.Edges.Add, edge)
	}
	_node = &FiatCurrency{config: fcuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, fcuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{fiatcurrency.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	fcuo.mutation.done = true
	return _node, nil
}
