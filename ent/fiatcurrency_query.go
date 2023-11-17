// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/paycrest/protocol/ent/fiatcurrency"
	"github.com/paycrest/protocol/ent/predicate"
	"github.com/paycrest/protocol/ent/providerprofile"
	"github.com/paycrest/protocol/ent/provisionbucket"
)

// FiatCurrencyQuery is the builder for querying FiatCurrency entities.
type FiatCurrencyQuery struct {
	config
	ctx                  *QueryContext
	order                []fiatcurrency.OrderOption
	inters               []Interceptor
	predicates           []predicate.FiatCurrency
	withProviders        *ProviderProfileQuery
	withProvisionBuckets *ProvisionBucketQuery
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the FiatCurrencyQuery builder.
func (fcq *FiatCurrencyQuery) Where(ps ...predicate.FiatCurrency) *FiatCurrencyQuery {
	fcq.predicates = append(fcq.predicates, ps...)
	return fcq
}

// Limit the number of records to be returned by this query.
func (fcq *FiatCurrencyQuery) Limit(limit int) *FiatCurrencyQuery {
	fcq.ctx.Limit = &limit
	return fcq
}

// Offset to start from.
func (fcq *FiatCurrencyQuery) Offset(offset int) *FiatCurrencyQuery {
	fcq.ctx.Offset = &offset
	return fcq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (fcq *FiatCurrencyQuery) Unique(unique bool) *FiatCurrencyQuery {
	fcq.ctx.Unique = &unique
	return fcq
}

// Order specifies how the records should be ordered.
func (fcq *FiatCurrencyQuery) Order(o ...fiatcurrency.OrderOption) *FiatCurrencyQuery {
	fcq.order = append(fcq.order, o...)
	return fcq
}

// QueryProviders chains the current query on the "providers" edge.
func (fcq *FiatCurrencyQuery) QueryProviders() *ProviderProfileQuery {
	query := (&ProviderProfileClient{config: fcq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := fcq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(fiatcurrency.Table, fiatcurrency.FieldID, selector),
			sqlgraph.To(providerprofile.Table, providerprofile.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, fiatcurrency.ProvidersTable, fiatcurrency.ProvidersColumn),
		)
		fromU = sqlgraph.SetNeighbors(fcq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryProvisionBuckets chains the current query on the "provision_buckets" edge.
func (fcq *FiatCurrencyQuery) QueryProvisionBuckets() *ProvisionBucketQuery {
	query := (&ProvisionBucketClient{config: fcq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := fcq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := fcq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(fiatcurrency.Table, fiatcurrency.FieldID, selector),
			sqlgraph.To(provisionbucket.Table, provisionbucket.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, fiatcurrency.ProvisionBucketsTable, fiatcurrency.ProvisionBucketsColumn),
		)
		fromU = sqlgraph.SetNeighbors(fcq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first FiatCurrency entity from the query.
// Returns a *NotFoundError when no FiatCurrency was found.
func (fcq *FiatCurrencyQuery) First(ctx context.Context) (*FiatCurrency, error) {
	nodes, err := fcq.Limit(1).All(setContextOp(ctx, fcq.ctx, "First"))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{fiatcurrency.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (fcq *FiatCurrencyQuery) FirstX(ctx context.Context) *FiatCurrency {
	node, err := fcq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first FiatCurrency ID from the query.
// Returns a *NotFoundError when no FiatCurrency ID was found.
func (fcq *FiatCurrencyQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = fcq.Limit(1).IDs(setContextOp(ctx, fcq.ctx, "FirstID")); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{fiatcurrency.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (fcq *FiatCurrencyQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := fcq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single FiatCurrency entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one FiatCurrency entity is found.
// Returns a *NotFoundError when no FiatCurrency entities are found.
func (fcq *FiatCurrencyQuery) Only(ctx context.Context) (*FiatCurrency, error) {
	nodes, err := fcq.Limit(2).All(setContextOp(ctx, fcq.ctx, "Only"))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{fiatcurrency.Label}
	default:
		return nil, &NotSingularError{fiatcurrency.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (fcq *FiatCurrencyQuery) OnlyX(ctx context.Context) *FiatCurrency {
	node, err := fcq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only FiatCurrency ID in the query.
// Returns a *NotSingularError when more than one FiatCurrency ID is found.
// Returns a *NotFoundError when no entities are found.
func (fcq *FiatCurrencyQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = fcq.Limit(2).IDs(setContextOp(ctx, fcq.ctx, "OnlyID")); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{fiatcurrency.Label}
	default:
		err = &NotSingularError{fiatcurrency.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (fcq *FiatCurrencyQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := fcq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of FiatCurrencies.
func (fcq *FiatCurrencyQuery) All(ctx context.Context) ([]*FiatCurrency, error) {
	ctx = setContextOp(ctx, fcq.ctx, "All")
	if err := fcq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*FiatCurrency, *FiatCurrencyQuery]()
	return withInterceptors[[]*FiatCurrency](ctx, fcq, qr, fcq.inters)
}

// AllX is like All, but panics if an error occurs.
func (fcq *FiatCurrencyQuery) AllX(ctx context.Context) []*FiatCurrency {
	nodes, err := fcq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of FiatCurrency IDs.
func (fcq *FiatCurrencyQuery) IDs(ctx context.Context) (ids []uuid.UUID, err error) {
	if fcq.ctx.Unique == nil && fcq.path != nil {
		fcq.Unique(true)
	}
	ctx = setContextOp(ctx, fcq.ctx, "IDs")
	if err = fcq.Select(fiatcurrency.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (fcq *FiatCurrencyQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := fcq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (fcq *FiatCurrencyQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, fcq.ctx, "Count")
	if err := fcq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, fcq, querierCount[*FiatCurrencyQuery](), fcq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (fcq *FiatCurrencyQuery) CountX(ctx context.Context) int {
	count, err := fcq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (fcq *FiatCurrencyQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, fcq.ctx, "Exist")
	switch _, err := fcq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (fcq *FiatCurrencyQuery) ExistX(ctx context.Context) bool {
	exist, err := fcq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the FiatCurrencyQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (fcq *FiatCurrencyQuery) Clone() *FiatCurrencyQuery {
	if fcq == nil {
		return nil
	}
	return &FiatCurrencyQuery{
		config:               fcq.config,
		ctx:                  fcq.ctx.Clone(),
		order:                append([]fiatcurrency.OrderOption{}, fcq.order...),
		inters:               append([]Interceptor{}, fcq.inters...),
		predicates:           append([]predicate.FiatCurrency{}, fcq.predicates...),
		withProviders:        fcq.withProviders.Clone(),
		withProvisionBuckets: fcq.withProvisionBuckets.Clone(),
		// clone intermediate query.
		sql:  fcq.sql.Clone(),
		path: fcq.path,
	}
}

// WithProviders tells the query-builder to eager-load the nodes that are connected to
// the "providers" edge. The optional arguments are used to configure the query builder of the edge.
func (fcq *FiatCurrencyQuery) WithProviders(opts ...func(*ProviderProfileQuery)) *FiatCurrencyQuery {
	query := (&ProviderProfileClient{config: fcq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	fcq.withProviders = query
	return fcq
}

// WithProvisionBuckets tells the query-builder to eager-load the nodes that are connected to
// the "provision_buckets" edge. The optional arguments are used to configure the query builder of the edge.
func (fcq *FiatCurrencyQuery) WithProvisionBuckets(opts ...func(*ProvisionBucketQuery)) *FiatCurrencyQuery {
	query := (&ProvisionBucketClient{config: fcq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	fcq.withProvisionBuckets = query
	return fcq
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.FiatCurrency.Query().
//		GroupBy(fiatcurrency.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (fcq *FiatCurrencyQuery) GroupBy(field string, fields ...string) *FiatCurrencyGroupBy {
	fcq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &FiatCurrencyGroupBy{build: fcq}
	grbuild.flds = &fcq.ctx.Fields
	grbuild.label = fiatcurrency.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		CreatedAt time.Time `json:"created_at,omitempty"`
//	}
//
//	client.FiatCurrency.Query().
//		Select(fiatcurrency.FieldCreatedAt).
//		Scan(ctx, &v)
func (fcq *FiatCurrencyQuery) Select(fields ...string) *FiatCurrencySelect {
	fcq.ctx.Fields = append(fcq.ctx.Fields, fields...)
	sbuild := &FiatCurrencySelect{FiatCurrencyQuery: fcq}
	sbuild.label = fiatcurrency.Label
	sbuild.flds, sbuild.scan = &fcq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a FiatCurrencySelect configured with the given aggregations.
func (fcq *FiatCurrencyQuery) Aggregate(fns ...AggregateFunc) *FiatCurrencySelect {
	return fcq.Select().Aggregate(fns...)
}

func (fcq *FiatCurrencyQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range fcq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, fcq); err != nil {
				return err
			}
		}
	}
	for _, f := range fcq.ctx.Fields {
		if !fiatcurrency.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if fcq.path != nil {
		prev, err := fcq.path(ctx)
		if err != nil {
			return err
		}
		fcq.sql = prev
	}
	return nil
}

func (fcq *FiatCurrencyQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*FiatCurrency, error) {
	var (
		nodes       = []*FiatCurrency{}
		_spec       = fcq.querySpec()
		loadedTypes = [2]bool{
			fcq.withProviders != nil,
			fcq.withProvisionBuckets != nil,
		}
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*FiatCurrency).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &FiatCurrency{config: fcq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, fcq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := fcq.withProviders; query != nil {
		if err := fcq.loadProviders(ctx, query, nodes,
			func(n *FiatCurrency) { n.Edges.Providers = []*ProviderProfile{} },
			func(n *FiatCurrency, e *ProviderProfile) { n.Edges.Providers = append(n.Edges.Providers, e) }); err != nil {
			return nil, err
		}
	}
	if query := fcq.withProvisionBuckets; query != nil {
		if err := fcq.loadProvisionBuckets(ctx, query, nodes,
			func(n *FiatCurrency) { n.Edges.ProvisionBuckets = []*ProvisionBucket{} },
			func(n *FiatCurrency, e *ProvisionBucket) {
				n.Edges.ProvisionBuckets = append(n.Edges.ProvisionBuckets, e)
			}); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (fcq *FiatCurrencyQuery) loadProviders(ctx context.Context, query *ProviderProfileQuery, nodes []*FiatCurrency, init func(*FiatCurrency), assign func(*FiatCurrency, *ProviderProfile)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*FiatCurrency)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.ProviderProfile(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(fiatcurrency.ProvidersColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.fiat_currency_providers
		if fk == nil {
			return fmt.Errorf(`foreign-key "fiat_currency_providers" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "fiat_currency_providers" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (fcq *FiatCurrencyQuery) loadProvisionBuckets(ctx context.Context, query *ProvisionBucketQuery, nodes []*FiatCurrency, init func(*FiatCurrency), assign func(*FiatCurrency, *ProvisionBucket)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*FiatCurrency)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.ProvisionBucket(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(fiatcurrency.ProvisionBucketsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.fiat_currency_provision_buckets
		if fk == nil {
			return fmt.Errorf(`foreign-key "fiat_currency_provision_buckets" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "fiat_currency_provision_buckets" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (fcq *FiatCurrencyQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := fcq.querySpec()
	_spec.Node.Columns = fcq.ctx.Fields
	if len(fcq.ctx.Fields) > 0 {
		_spec.Unique = fcq.ctx.Unique != nil && *fcq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, fcq.driver, _spec)
}

func (fcq *FiatCurrencyQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(fiatcurrency.Table, fiatcurrency.Columns, sqlgraph.NewFieldSpec(fiatcurrency.FieldID, field.TypeUUID))
	_spec.From = fcq.sql
	if unique := fcq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if fcq.path != nil {
		_spec.Unique = true
	}
	if fields := fcq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, fiatcurrency.FieldID)
		for i := range fields {
			if fields[i] != fiatcurrency.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := fcq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := fcq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := fcq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := fcq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (fcq *FiatCurrencyQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(fcq.driver.Dialect())
	t1 := builder.Table(fiatcurrency.Table)
	columns := fcq.ctx.Fields
	if len(columns) == 0 {
		columns = fiatcurrency.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if fcq.sql != nil {
		selector = fcq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if fcq.ctx.Unique != nil && *fcq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range fcq.predicates {
		p(selector)
	}
	for _, p := range fcq.order {
		p(selector)
	}
	if offset := fcq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := fcq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// FiatCurrencyGroupBy is the group-by builder for FiatCurrency entities.
type FiatCurrencyGroupBy struct {
	selector
	build *FiatCurrencyQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (fcgb *FiatCurrencyGroupBy) Aggregate(fns ...AggregateFunc) *FiatCurrencyGroupBy {
	fcgb.fns = append(fcgb.fns, fns...)
	return fcgb
}

// Scan applies the selector query and scans the result into the given value.
func (fcgb *FiatCurrencyGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, fcgb.build.ctx, "GroupBy")
	if err := fcgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*FiatCurrencyQuery, *FiatCurrencyGroupBy](ctx, fcgb.build, fcgb, fcgb.build.inters, v)
}

func (fcgb *FiatCurrencyGroupBy) sqlScan(ctx context.Context, root *FiatCurrencyQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(fcgb.fns))
	for _, fn := range fcgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*fcgb.flds)+len(fcgb.fns))
		for _, f := range *fcgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*fcgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := fcgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// FiatCurrencySelect is the builder for selecting fields of FiatCurrency entities.
type FiatCurrencySelect struct {
	*FiatCurrencyQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (fcs *FiatCurrencySelect) Aggregate(fns ...AggregateFunc) *FiatCurrencySelect {
	fcs.fns = append(fcs.fns, fns...)
	return fcs
}

// Scan applies the selector query and scans the result into the given value.
func (fcs *FiatCurrencySelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, fcs.ctx, "Select")
	if err := fcs.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*FiatCurrencyQuery, *FiatCurrencySelect](ctx, fcs.FiatCurrencyQuery, fcs, fcs.inters, v)
}

func (fcs *FiatCurrencySelect) sqlScan(ctx context.Context, root *FiatCurrencyQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(fcs.fns))
	for _, fn := range fcs.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*fcs.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := fcs.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
