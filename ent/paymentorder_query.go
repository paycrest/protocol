// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"database/sql/driver"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"github.com/paycrest/aggregator/ent/linkedaddress"
	"github.com/paycrest/aggregator/ent/paymentorder"
	"github.com/paycrest/aggregator/ent/paymentorderrecipient"
	"github.com/paycrest/aggregator/ent/predicate"
	"github.com/paycrest/aggregator/ent/receiveaddress"
	"github.com/paycrest/aggregator/ent/senderprofile"
	"github.com/paycrest/aggregator/ent/token"
	"github.com/paycrest/aggregator/ent/transactionlog"
)

// PaymentOrderQuery is the builder for querying PaymentOrder entities.
type PaymentOrderQuery struct {
	config
	ctx                *QueryContext
	order              []paymentorder.OrderOption
	inters             []Interceptor
	predicates         []predicate.PaymentOrder
	withSenderProfile  *SenderProfileQuery
	withToken          *TokenQuery
	withLinkedAddress  *LinkedAddressQuery
	withReceiveAddress *ReceiveAddressQuery
	withRecipient      *PaymentOrderRecipientQuery
	withTransactions   *TransactionLogQuery
	withFKs            bool
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the PaymentOrderQuery builder.
func (poq *PaymentOrderQuery) Where(ps ...predicate.PaymentOrder) *PaymentOrderQuery {
	poq.predicates = append(poq.predicates, ps...)
	return poq
}

// Limit the number of records to be returned by this query.
func (poq *PaymentOrderQuery) Limit(limit int) *PaymentOrderQuery {
	poq.ctx.Limit = &limit
	return poq
}

// Offset to start from.
func (poq *PaymentOrderQuery) Offset(offset int) *PaymentOrderQuery {
	poq.ctx.Offset = &offset
	return poq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (poq *PaymentOrderQuery) Unique(unique bool) *PaymentOrderQuery {
	poq.ctx.Unique = &unique
	return poq
}

// Order specifies how the records should be ordered.
func (poq *PaymentOrderQuery) Order(o ...paymentorder.OrderOption) *PaymentOrderQuery {
	poq.order = append(poq.order, o...)
	return poq
}

// QuerySenderProfile chains the current query on the "sender_profile" edge.
func (poq *PaymentOrderQuery) QuerySenderProfile() *SenderProfileQuery {
	query := (&SenderProfileClient{config: poq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := poq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := poq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(paymentorder.Table, paymentorder.FieldID, selector),
			sqlgraph.To(senderprofile.Table, senderprofile.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, paymentorder.SenderProfileTable, paymentorder.SenderProfileColumn),
		)
		fromU = sqlgraph.SetNeighbors(poq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryToken chains the current query on the "token" edge.
func (poq *PaymentOrderQuery) QueryToken() *TokenQuery {
	query := (&TokenClient{config: poq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := poq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := poq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(paymentorder.Table, paymentorder.FieldID, selector),
			sqlgraph.To(token.Table, token.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, paymentorder.TokenTable, paymentorder.TokenColumn),
		)
		fromU = sqlgraph.SetNeighbors(poq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryLinkedAddress chains the current query on the "linked_address" edge.
func (poq *PaymentOrderQuery) QueryLinkedAddress() *LinkedAddressQuery {
	query := (&LinkedAddressClient{config: poq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := poq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := poq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(paymentorder.Table, paymentorder.FieldID, selector),
			sqlgraph.To(linkedaddress.Table, linkedaddress.FieldID),
			sqlgraph.Edge(sqlgraph.M2O, true, paymentorder.LinkedAddressTable, paymentorder.LinkedAddressColumn),
		)
		fromU = sqlgraph.SetNeighbors(poq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryReceiveAddress chains the current query on the "receive_address" edge.
func (poq *PaymentOrderQuery) QueryReceiveAddress() *ReceiveAddressQuery {
	query := (&ReceiveAddressClient{config: poq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := poq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := poq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(paymentorder.Table, paymentorder.FieldID, selector),
			sqlgraph.To(receiveaddress.Table, receiveaddress.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, paymentorder.ReceiveAddressTable, paymentorder.ReceiveAddressColumn),
		)
		fromU = sqlgraph.SetNeighbors(poq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryRecipient chains the current query on the "recipient" edge.
func (poq *PaymentOrderQuery) QueryRecipient() *PaymentOrderRecipientQuery {
	query := (&PaymentOrderRecipientClient{config: poq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := poq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := poq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(paymentorder.Table, paymentorder.FieldID, selector),
			sqlgraph.To(paymentorderrecipient.Table, paymentorderrecipient.FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, paymentorder.RecipientTable, paymentorder.RecipientColumn),
		)
		fromU = sqlgraph.SetNeighbors(poq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// QueryTransactions chains the current query on the "transactions" edge.
func (poq *PaymentOrderQuery) QueryTransactions() *TransactionLogQuery {
	query := (&TransactionLogClient{config: poq.config}).Query()
	query.path = func(ctx context.Context) (fromU *sql.Selector, err error) {
		if err := poq.prepareQuery(ctx); err != nil {
			return nil, err
		}
		selector := poq.sqlQuery(ctx)
		if err := selector.Err(); err != nil {
			return nil, err
		}
		step := sqlgraph.NewStep(
			sqlgraph.From(paymentorder.Table, paymentorder.FieldID, selector),
			sqlgraph.To(transactionlog.Table, transactionlog.FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, paymentorder.TransactionsTable, paymentorder.TransactionsColumn),
		)
		fromU = sqlgraph.SetNeighbors(poq.driver.Dialect(), step)
		return fromU, nil
	}
	return query
}

// First returns the first PaymentOrder entity from the query.
// Returns a *NotFoundError when no PaymentOrder was found.
func (poq *PaymentOrderQuery) First(ctx context.Context) (*PaymentOrder, error) {
	nodes, err := poq.Limit(1).All(setContextOp(ctx, poq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{paymentorder.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (poq *PaymentOrderQuery) FirstX(ctx context.Context) *PaymentOrder {
	node, err := poq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first PaymentOrder ID from the query.
// Returns a *NotFoundError when no PaymentOrder ID was found.
func (poq *PaymentOrderQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = poq.Limit(1).IDs(setContextOp(ctx, poq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{paymentorder.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (poq *PaymentOrderQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := poq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single PaymentOrder entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one PaymentOrder entity is found.
// Returns a *NotFoundError when no PaymentOrder entities are found.
func (poq *PaymentOrderQuery) Only(ctx context.Context) (*PaymentOrder, error) {
	nodes, err := poq.Limit(2).All(setContextOp(ctx, poq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{paymentorder.Label}
	default:
		return nil, &NotSingularError{paymentorder.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (poq *PaymentOrderQuery) OnlyX(ctx context.Context) *PaymentOrder {
	node, err := poq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only PaymentOrder ID in the query.
// Returns a *NotSingularError when more than one PaymentOrder ID is found.
// Returns a *NotFoundError when no entities are found.
func (poq *PaymentOrderQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = poq.Limit(2).IDs(setContextOp(ctx, poq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{paymentorder.Label}
	default:
		err = &NotSingularError{paymentorder.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (poq *PaymentOrderQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := poq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of PaymentOrders.
func (poq *PaymentOrderQuery) All(ctx context.Context) ([]*PaymentOrder, error) {
	ctx = setContextOp(ctx, poq.ctx, ent.OpQueryAll)
	if err := poq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*PaymentOrder, *PaymentOrderQuery]()
	return withInterceptors[[]*PaymentOrder](ctx, poq, qr, poq.inters)
}

// AllX is like All, but panics if an error occurs.
func (poq *PaymentOrderQuery) AllX(ctx context.Context) []*PaymentOrder {
	nodes, err := poq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of PaymentOrder IDs.
func (poq *PaymentOrderQuery) IDs(ctx context.Context) (ids []uuid.UUID, err error) {
	if poq.ctx.Unique == nil && poq.path != nil {
		poq.Unique(true)
	}
	ctx = setContextOp(ctx, poq.ctx, ent.OpQueryIDs)
	if err = poq.Select(paymentorder.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (poq *PaymentOrderQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := poq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (poq *PaymentOrderQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, poq.ctx, ent.OpQueryCount)
	if err := poq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, poq, querierCount[*PaymentOrderQuery](), poq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (poq *PaymentOrderQuery) CountX(ctx context.Context) int {
	count, err := poq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (poq *PaymentOrderQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, poq.ctx, ent.OpQueryExist)
	switch _, err := poq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (poq *PaymentOrderQuery) ExistX(ctx context.Context) bool {
	exist, err := poq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the PaymentOrderQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (poq *PaymentOrderQuery) Clone() *PaymentOrderQuery {
	if poq == nil {
		return nil
	}
	return &PaymentOrderQuery{
		config:             poq.config,
		ctx:                poq.ctx.Clone(),
		order:              append([]paymentorder.OrderOption{}, poq.order...),
		inters:             append([]Interceptor{}, poq.inters...),
		predicates:         append([]predicate.PaymentOrder{}, poq.predicates...),
		withSenderProfile:  poq.withSenderProfile.Clone(),
		withToken:          poq.withToken.Clone(),
		withLinkedAddress:  poq.withLinkedAddress.Clone(),
		withReceiveAddress: poq.withReceiveAddress.Clone(),
		withRecipient:      poq.withRecipient.Clone(),
		withTransactions:   poq.withTransactions.Clone(),
		// clone intermediate query.
		sql:  poq.sql.Clone(),
		path: poq.path,
	}
}

// WithSenderProfile tells the query-builder to eager-load the nodes that are connected to
// the "sender_profile" edge. The optional arguments are used to configure the query builder of the edge.
func (poq *PaymentOrderQuery) WithSenderProfile(opts ...func(*SenderProfileQuery)) *PaymentOrderQuery {
	query := (&SenderProfileClient{config: poq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	poq.withSenderProfile = query
	return poq
}

// WithToken tells the query-builder to eager-load the nodes that are connected to
// the "token" edge. The optional arguments are used to configure the query builder of the edge.
func (poq *PaymentOrderQuery) WithToken(opts ...func(*TokenQuery)) *PaymentOrderQuery {
	query := (&TokenClient{config: poq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	poq.withToken = query
	return poq
}

// WithLinkedAddress tells the query-builder to eager-load the nodes that are connected to
// the "linked_address" edge. The optional arguments are used to configure the query builder of the edge.
func (poq *PaymentOrderQuery) WithLinkedAddress(opts ...func(*LinkedAddressQuery)) *PaymentOrderQuery {
	query := (&LinkedAddressClient{config: poq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	poq.withLinkedAddress = query
	return poq
}

// WithReceiveAddress tells the query-builder to eager-load the nodes that are connected to
// the "receive_address" edge. The optional arguments are used to configure the query builder of the edge.
func (poq *PaymentOrderQuery) WithReceiveAddress(opts ...func(*ReceiveAddressQuery)) *PaymentOrderQuery {
	query := (&ReceiveAddressClient{config: poq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	poq.withReceiveAddress = query
	return poq
}

// WithRecipient tells the query-builder to eager-load the nodes that are connected to
// the "recipient" edge. The optional arguments are used to configure the query builder of the edge.
func (poq *PaymentOrderQuery) WithRecipient(opts ...func(*PaymentOrderRecipientQuery)) *PaymentOrderQuery {
	query := (&PaymentOrderRecipientClient{config: poq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	poq.withRecipient = query
	return poq
}

// WithTransactions tells the query-builder to eager-load the nodes that are connected to
// the "transactions" edge. The optional arguments are used to configure the query builder of the edge.
func (poq *PaymentOrderQuery) WithTransactions(opts ...func(*TransactionLogQuery)) *PaymentOrderQuery {
	query := (&TransactionLogClient{config: poq.config}).Query()
	for _, opt := range opts {
		opt(query)
	}
	poq.withTransactions = query
	return poq
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
//	client.PaymentOrder.Query().
//		GroupBy(paymentorder.FieldCreatedAt).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (poq *PaymentOrderQuery) GroupBy(field string, fields ...string) *PaymentOrderGroupBy {
	poq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &PaymentOrderGroupBy{build: poq}
	grbuild.flds = &poq.ctx.Fields
	grbuild.label = paymentorder.Label
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
//	client.PaymentOrder.Query().
//		Select(paymentorder.FieldCreatedAt).
//		Scan(ctx, &v)
func (poq *PaymentOrderQuery) Select(fields ...string) *PaymentOrderSelect {
	poq.ctx.Fields = append(poq.ctx.Fields, fields...)
	sbuild := &PaymentOrderSelect{PaymentOrderQuery: poq}
	sbuild.label = paymentorder.Label
	sbuild.flds, sbuild.scan = &poq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a PaymentOrderSelect configured with the given aggregations.
func (poq *PaymentOrderQuery) Aggregate(fns ...AggregateFunc) *PaymentOrderSelect {
	return poq.Select().Aggregate(fns...)
}

func (poq *PaymentOrderQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range poq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, poq); err != nil {
				return err
			}
		}
	}
	for _, f := range poq.ctx.Fields {
		if !paymentorder.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if poq.path != nil {
		prev, err := poq.path(ctx)
		if err != nil {
			return err
		}
		poq.sql = prev
	}
	return nil
}

func (poq *PaymentOrderQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*PaymentOrder, error) {
	var (
		nodes       = []*PaymentOrder{}
		withFKs     = poq.withFKs
		_spec       = poq.querySpec()
		loadedTypes = [6]bool{
			poq.withSenderProfile != nil,
			poq.withToken != nil,
			poq.withLinkedAddress != nil,
			poq.withReceiveAddress != nil,
			poq.withRecipient != nil,
			poq.withTransactions != nil,
		}
	)
	if poq.withSenderProfile != nil || poq.withToken != nil || poq.withLinkedAddress != nil {
		withFKs = true
	}
	if withFKs {
		_spec.Node.Columns = append(_spec.Node.Columns, paymentorder.ForeignKeys...)
	}
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*PaymentOrder).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &PaymentOrder{config: poq.config}
		nodes = append(nodes, node)
		node.Edges.loadedTypes = loadedTypes
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, poq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	if query := poq.withSenderProfile; query != nil {
		if err := poq.loadSenderProfile(ctx, query, nodes, nil,
			func(n *PaymentOrder, e *SenderProfile) { n.Edges.SenderProfile = e }); err != nil {
			return nil, err
		}
	}
	if query := poq.withToken; query != nil {
		if err := poq.loadToken(ctx, query, nodes, nil,
			func(n *PaymentOrder, e *Token) { n.Edges.Token = e }); err != nil {
			return nil, err
		}
	}
	if query := poq.withLinkedAddress; query != nil {
		if err := poq.loadLinkedAddress(ctx, query, nodes, nil,
			func(n *PaymentOrder, e *LinkedAddress) { n.Edges.LinkedAddress = e }); err != nil {
			return nil, err
		}
	}
	if query := poq.withReceiveAddress; query != nil {
		if err := poq.loadReceiveAddress(ctx, query, nodes, nil,
			func(n *PaymentOrder, e *ReceiveAddress) { n.Edges.ReceiveAddress = e }); err != nil {
			return nil, err
		}
	}
	if query := poq.withRecipient; query != nil {
		if err := poq.loadRecipient(ctx, query, nodes, nil,
			func(n *PaymentOrder, e *PaymentOrderRecipient) { n.Edges.Recipient = e }); err != nil {
			return nil, err
		}
	}
	if query := poq.withTransactions; query != nil {
		if err := poq.loadTransactions(ctx, query, nodes,
			func(n *PaymentOrder) { n.Edges.Transactions = []*TransactionLog{} },
			func(n *PaymentOrder, e *TransactionLog) { n.Edges.Transactions = append(n.Edges.Transactions, e) }); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

func (poq *PaymentOrderQuery) loadSenderProfile(ctx context.Context, query *SenderProfileQuery, nodes []*PaymentOrder, init func(*PaymentOrder), assign func(*PaymentOrder, *SenderProfile)) error {
	ids := make([]uuid.UUID, 0, len(nodes))
	nodeids := make(map[uuid.UUID][]*PaymentOrder)
	for i := range nodes {
		if nodes[i].sender_profile_payment_orders == nil {
			continue
		}
		fk := *nodes[i].sender_profile_payment_orders
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(senderprofile.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "sender_profile_payment_orders" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (poq *PaymentOrderQuery) loadToken(ctx context.Context, query *TokenQuery, nodes []*PaymentOrder, init func(*PaymentOrder), assign func(*PaymentOrder, *Token)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*PaymentOrder)
	for i := range nodes {
		if nodes[i].token_payment_orders == nil {
			continue
		}
		fk := *nodes[i].token_payment_orders
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(token.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "token_payment_orders" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (poq *PaymentOrderQuery) loadLinkedAddress(ctx context.Context, query *LinkedAddressQuery, nodes []*PaymentOrder, init func(*PaymentOrder), assign func(*PaymentOrder, *LinkedAddress)) error {
	ids := make([]int, 0, len(nodes))
	nodeids := make(map[int][]*PaymentOrder)
	for i := range nodes {
		if nodes[i].linked_address_payment_orders == nil {
			continue
		}
		fk := *nodes[i].linked_address_payment_orders
		if _, ok := nodeids[fk]; !ok {
			ids = append(ids, fk)
		}
		nodeids[fk] = append(nodeids[fk], nodes[i])
	}
	if len(ids) == 0 {
		return nil
	}
	query.Where(linkedaddress.IDIn(ids...))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		nodes, ok := nodeids[n.ID]
		if !ok {
			return fmt.Errorf(`unexpected foreign-key "linked_address_payment_orders" returned %v`, n.ID)
		}
		for i := range nodes {
			assign(nodes[i], n)
		}
	}
	return nil
}
func (poq *PaymentOrderQuery) loadReceiveAddress(ctx context.Context, query *ReceiveAddressQuery, nodes []*PaymentOrder, init func(*PaymentOrder), assign func(*PaymentOrder, *ReceiveAddress)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*PaymentOrder)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
	}
	query.withFKs = true
	query.Where(predicate.ReceiveAddress(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(paymentorder.ReceiveAddressColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.payment_order_receive_address
		if fk == nil {
			return fmt.Errorf(`foreign-key "payment_order_receive_address" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "payment_order_receive_address" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (poq *PaymentOrderQuery) loadRecipient(ctx context.Context, query *PaymentOrderRecipientQuery, nodes []*PaymentOrder, init func(*PaymentOrder), assign func(*PaymentOrder, *PaymentOrderRecipient)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*PaymentOrder)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
	}
	query.withFKs = true
	query.Where(predicate.PaymentOrderRecipient(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(paymentorder.RecipientColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.payment_order_recipient
		if fk == nil {
			return fmt.Errorf(`foreign-key "payment_order_recipient" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "payment_order_recipient" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}
func (poq *PaymentOrderQuery) loadTransactions(ctx context.Context, query *TransactionLogQuery, nodes []*PaymentOrder, init func(*PaymentOrder), assign func(*PaymentOrder, *TransactionLog)) error {
	fks := make([]driver.Value, 0, len(nodes))
	nodeids := make(map[uuid.UUID]*PaymentOrder)
	for i := range nodes {
		fks = append(fks, nodes[i].ID)
		nodeids[nodes[i].ID] = nodes[i]
		if init != nil {
			init(nodes[i])
		}
	}
	query.withFKs = true
	query.Where(predicate.TransactionLog(func(s *sql.Selector) {
		s.Where(sql.InValues(s.C(paymentorder.TransactionsColumn), fks...))
	}))
	neighbors, err := query.All(ctx)
	if err != nil {
		return err
	}
	for _, n := range neighbors {
		fk := n.payment_order_transactions
		if fk == nil {
			return fmt.Errorf(`foreign-key "payment_order_transactions" is nil for node %v`, n.ID)
		}
		node, ok := nodeids[*fk]
		if !ok {
			return fmt.Errorf(`unexpected referenced foreign-key "payment_order_transactions" returned %v for node %v`, *fk, n.ID)
		}
		assign(node, n)
	}
	return nil
}

func (poq *PaymentOrderQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := poq.querySpec()
	_spec.Node.Columns = poq.ctx.Fields
	if len(poq.ctx.Fields) > 0 {
		_spec.Unique = poq.ctx.Unique != nil && *poq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, poq.driver, _spec)
}

func (poq *PaymentOrderQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(paymentorder.Table, paymentorder.Columns, sqlgraph.NewFieldSpec(paymentorder.FieldID, field.TypeUUID))
	_spec.From = poq.sql
	if unique := poq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if poq.path != nil {
		_spec.Unique = true
	}
	if fields := poq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, paymentorder.FieldID)
		for i := range fields {
			if fields[i] != paymentorder.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := poq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := poq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := poq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := poq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (poq *PaymentOrderQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(poq.driver.Dialect())
	t1 := builder.Table(paymentorder.Table)
	columns := poq.ctx.Fields
	if len(columns) == 0 {
		columns = paymentorder.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if poq.sql != nil {
		selector = poq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if poq.ctx.Unique != nil && *poq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range poq.predicates {
		p(selector)
	}
	for _, p := range poq.order {
		p(selector)
	}
	if offset := poq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := poq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// PaymentOrderGroupBy is the group-by builder for PaymentOrder entities.
type PaymentOrderGroupBy struct {
	selector
	build *PaymentOrderQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (pogb *PaymentOrderGroupBy) Aggregate(fns ...AggregateFunc) *PaymentOrderGroupBy {
	pogb.fns = append(pogb.fns, fns...)
	return pogb
}

// Scan applies the selector query and scans the result into the given value.
func (pogb *PaymentOrderGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pogb.build.ctx, ent.OpQueryGroupBy)
	if err := pogb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PaymentOrderQuery, *PaymentOrderGroupBy](ctx, pogb.build, pogb, pogb.build.inters, v)
}

func (pogb *PaymentOrderGroupBy) sqlScan(ctx context.Context, root *PaymentOrderQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(pogb.fns))
	for _, fn := range pogb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*pogb.flds)+len(pogb.fns))
		for _, f := range *pogb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*pogb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pogb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// PaymentOrderSelect is the builder for selecting fields of PaymentOrder entities.
type PaymentOrderSelect struct {
	*PaymentOrderQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (pos *PaymentOrderSelect) Aggregate(fns ...AggregateFunc) *PaymentOrderSelect {
	pos.fns = append(pos.fns, fns...)
	return pos
}

// Scan applies the selector query and scans the result into the given value.
func (pos *PaymentOrderSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, pos.ctx, ent.OpQuerySelect)
	if err := pos.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*PaymentOrderQuery, *PaymentOrderSelect](ctx, pos.PaymentOrderQuery, pos, pos.inters, v)
}

func (pos *PaymentOrderSelect) sqlScan(ctx context.Context, root *PaymentOrderQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(pos.fns))
	for _, fn := range pos.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*pos.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := pos.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}
