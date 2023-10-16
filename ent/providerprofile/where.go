// Code generated by ent, DO NOT EDIT.

package providerprofile

import (
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"github.com/paycrest/paycrest-protocol/ent/predicate"
)

// ID filters vertices based on their ID field.
func ID(id string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldID, id))
}

// IDEQ applies the EQ predicate on the ID field.
func IDEQ(id string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldID, id))
}

// IDNEQ applies the NEQ predicate on the ID field.
func IDNEQ(id string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNEQ(FieldID, id))
}

// IDIn applies the In predicate on the ID field.
func IDIn(ids ...string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldIn(FieldID, ids...))
}

// IDNotIn applies the NotIn predicate on the ID field.
func IDNotIn(ids ...string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNotIn(FieldID, ids...))
}

// IDGT applies the GT predicate on the ID field.
func IDGT(id string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldGT(FieldID, id))
}

// IDGTE applies the GTE predicate on the ID field.
func IDGTE(id string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldGTE(FieldID, id))
}

// IDLT applies the LT predicate on the ID field.
func IDLT(id string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldLT(FieldID, id))
}

// IDLTE applies the LTE predicate on the ID field.
func IDLTE(id string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldLTE(FieldID, id))
}

// IDEqualFold applies the EqualFold predicate on the ID field.
func IDEqualFold(id string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEqualFold(FieldID, id))
}

// IDContainsFold applies the ContainsFold predicate on the ID field.
func IDContainsFold(id string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldContainsFold(FieldID, id))
}

// CreatedAt applies equality check predicate on the "created_at" field. It's identical to CreatedAtEQ.
func CreatedAt(v time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldCreatedAt, v))
}

// UpdatedAt applies equality check predicate on the "updated_at" field. It's identical to UpdatedAtEQ.
func UpdatedAt(v time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldUpdatedAt, v))
}

// TradingName applies equality check predicate on the "trading_name" field. It's identical to TradingNameEQ.
func TradingName(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldTradingName, v))
}

// Country applies equality check predicate on the "country" field. It's identical to CountryEQ.
func Country(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldCountry, v))
}

// HostIdentifier applies equality check predicate on the "host_identifier" field. It's identical to HostIdentifierEQ.
func HostIdentifier(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldHostIdentifier, v))
}

// IsPartner applies equality check predicate on the "is_partner" field. It's identical to IsPartnerEQ.
func IsPartner(v bool) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldIsPartner, v))
}

// CreatedAtEQ applies the EQ predicate on the "created_at" field.
func CreatedAtEQ(v time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldCreatedAt, v))
}

// CreatedAtNEQ applies the NEQ predicate on the "created_at" field.
func CreatedAtNEQ(v time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNEQ(FieldCreatedAt, v))
}

// CreatedAtIn applies the In predicate on the "created_at" field.
func CreatedAtIn(vs ...time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldIn(FieldCreatedAt, vs...))
}

// CreatedAtNotIn applies the NotIn predicate on the "created_at" field.
func CreatedAtNotIn(vs ...time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNotIn(FieldCreatedAt, vs...))
}

// CreatedAtGT applies the GT predicate on the "created_at" field.
func CreatedAtGT(v time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldGT(FieldCreatedAt, v))
}

// CreatedAtGTE applies the GTE predicate on the "created_at" field.
func CreatedAtGTE(v time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldGTE(FieldCreatedAt, v))
}

// CreatedAtLT applies the LT predicate on the "created_at" field.
func CreatedAtLT(v time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldLT(FieldCreatedAt, v))
}

// CreatedAtLTE applies the LTE predicate on the "created_at" field.
func CreatedAtLTE(v time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldLTE(FieldCreatedAt, v))
}

// UpdatedAtEQ applies the EQ predicate on the "updated_at" field.
func UpdatedAtEQ(v time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldUpdatedAt, v))
}

// UpdatedAtNEQ applies the NEQ predicate on the "updated_at" field.
func UpdatedAtNEQ(v time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNEQ(FieldUpdatedAt, v))
}

// UpdatedAtIn applies the In predicate on the "updated_at" field.
func UpdatedAtIn(vs ...time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldIn(FieldUpdatedAt, vs...))
}

// UpdatedAtNotIn applies the NotIn predicate on the "updated_at" field.
func UpdatedAtNotIn(vs ...time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNotIn(FieldUpdatedAt, vs...))
}

// UpdatedAtGT applies the GT predicate on the "updated_at" field.
func UpdatedAtGT(v time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldGT(FieldUpdatedAt, v))
}

// UpdatedAtGTE applies the GTE predicate on the "updated_at" field.
func UpdatedAtGTE(v time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldGTE(FieldUpdatedAt, v))
}

// UpdatedAtLT applies the LT predicate on the "updated_at" field.
func UpdatedAtLT(v time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldLT(FieldUpdatedAt, v))
}

// UpdatedAtLTE applies the LTE predicate on the "updated_at" field.
func UpdatedAtLTE(v time.Time) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldLTE(FieldUpdatedAt, v))
}

// TradingNameEQ applies the EQ predicate on the "trading_name" field.
func TradingNameEQ(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldTradingName, v))
}

// TradingNameNEQ applies the NEQ predicate on the "trading_name" field.
func TradingNameNEQ(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNEQ(FieldTradingName, v))
}

// TradingNameIn applies the In predicate on the "trading_name" field.
func TradingNameIn(vs ...string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldIn(FieldTradingName, vs...))
}

// TradingNameNotIn applies the NotIn predicate on the "trading_name" field.
func TradingNameNotIn(vs ...string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNotIn(FieldTradingName, vs...))
}

// TradingNameGT applies the GT predicate on the "trading_name" field.
func TradingNameGT(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldGT(FieldTradingName, v))
}

// TradingNameGTE applies the GTE predicate on the "trading_name" field.
func TradingNameGTE(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldGTE(FieldTradingName, v))
}

// TradingNameLT applies the LT predicate on the "trading_name" field.
func TradingNameLT(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldLT(FieldTradingName, v))
}

// TradingNameLTE applies the LTE predicate on the "trading_name" field.
func TradingNameLTE(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldLTE(FieldTradingName, v))
}

// TradingNameContains applies the Contains predicate on the "trading_name" field.
func TradingNameContains(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldContains(FieldTradingName, v))
}

// TradingNameHasPrefix applies the HasPrefix predicate on the "trading_name" field.
func TradingNameHasPrefix(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldHasPrefix(FieldTradingName, v))
}

// TradingNameHasSuffix applies the HasSuffix predicate on the "trading_name" field.
func TradingNameHasSuffix(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldHasSuffix(FieldTradingName, v))
}

// TradingNameEqualFold applies the EqualFold predicate on the "trading_name" field.
func TradingNameEqualFold(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEqualFold(FieldTradingName, v))
}

// TradingNameContainsFold applies the ContainsFold predicate on the "trading_name" field.
func TradingNameContainsFold(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldContainsFold(FieldTradingName, v))
}

// CountryEQ applies the EQ predicate on the "country" field.
func CountryEQ(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldCountry, v))
}

// CountryNEQ applies the NEQ predicate on the "country" field.
func CountryNEQ(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNEQ(FieldCountry, v))
}

// CountryIn applies the In predicate on the "country" field.
func CountryIn(vs ...string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldIn(FieldCountry, vs...))
}

// CountryNotIn applies the NotIn predicate on the "country" field.
func CountryNotIn(vs ...string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNotIn(FieldCountry, vs...))
}

// CountryGT applies the GT predicate on the "country" field.
func CountryGT(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldGT(FieldCountry, v))
}

// CountryGTE applies the GTE predicate on the "country" field.
func CountryGTE(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldGTE(FieldCountry, v))
}

// CountryLT applies the LT predicate on the "country" field.
func CountryLT(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldLT(FieldCountry, v))
}

// CountryLTE applies the LTE predicate on the "country" field.
func CountryLTE(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldLTE(FieldCountry, v))
}

// CountryContains applies the Contains predicate on the "country" field.
func CountryContains(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldContains(FieldCountry, v))
}

// CountryHasPrefix applies the HasPrefix predicate on the "country" field.
func CountryHasPrefix(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldHasPrefix(FieldCountry, v))
}

// CountryHasSuffix applies the HasSuffix predicate on the "country" field.
func CountryHasSuffix(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldHasSuffix(FieldCountry, v))
}

// CountryEqualFold applies the EqualFold predicate on the "country" field.
func CountryEqualFold(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEqualFold(FieldCountry, v))
}

// CountryContainsFold applies the ContainsFold predicate on the "country" field.
func CountryContainsFold(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldContainsFold(FieldCountry, v))
}

// HostIdentifierEQ applies the EQ predicate on the "host_identifier" field.
func HostIdentifierEQ(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldHostIdentifier, v))
}

// HostIdentifierNEQ applies the NEQ predicate on the "host_identifier" field.
func HostIdentifierNEQ(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNEQ(FieldHostIdentifier, v))
}

// HostIdentifierIn applies the In predicate on the "host_identifier" field.
func HostIdentifierIn(vs ...string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldIn(FieldHostIdentifier, vs...))
}

// HostIdentifierNotIn applies the NotIn predicate on the "host_identifier" field.
func HostIdentifierNotIn(vs ...string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNotIn(FieldHostIdentifier, vs...))
}

// HostIdentifierGT applies the GT predicate on the "host_identifier" field.
func HostIdentifierGT(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldGT(FieldHostIdentifier, v))
}

// HostIdentifierGTE applies the GTE predicate on the "host_identifier" field.
func HostIdentifierGTE(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldGTE(FieldHostIdentifier, v))
}

// HostIdentifierLT applies the LT predicate on the "host_identifier" field.
func HostIdentifierLT(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldLT(FieldHostIdentifier, v))
}

// HostIdentifierLTE applies the LTE predicate on the "host_identifier" field.
func HostIdentifierLTE(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldLTE(FieldHostIdentifier, v))
}

// HostIdentifierContains applies the Contains predicate on the "host_identifier" field.
func HostIdentifierContains(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldContains(FieldHostIdentifier, v))
}

// HostIdentifierHasPrefix applies the HasPrefix predicate on the "host_identifier" field.
func HostIdentifierHasPrefix(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldHasPrefix(FieldHostIdentifier, v))
}

// HostIdentifierHasSuffix applies the HasSuffix predicate on the "host_identifier" field.
func HostIdentifierHasSuffix(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldHasSuffix(FieldHostIdentifier, v))
}

// HostIdentifierIsNil applies the IsNil predicate on the "host_identifier" field.
func HostIdentifierIsNil() predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldIsNull(FieldHostIdentifier))
}

// HostIdentifierNotNil applies the NotNil predicate on the "host_identifier" field.
func HostIdentifierNotNil() predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNotNull(FieldHostIdentifier))
}

// HostIdentifierEqualFold applies the EqualFold predicate on the "host_identifier" field.
func HostIdentifierEqualFold(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEqualFold(FieldHostIdentifier, v))
}

// HostIdentifierContainsFold applies the ContainsFold predicate on the "host_identifier" field.
func HostIdentifierContainsFold(v string) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldContainsFold(FieldHostIdentifier, v))
}

// ProvisionModeEQ applies the EQ predicate on the "provision_mode" field.
func ProvisionModeEQ(v ProvisionMode) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldProvisionMode, v))
}

// ProvisionModeNEQ applies the NEQ predicate on the "provision_mode" field.
func ProvisionModeNEQ(v ProvisionMode) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNEQ(FieldProvisionMode, v))
}

// ProvisionModeIn applies the In predicate on the "provision_mode" field.
func ProvisionModeIn(vs ...ProvisionMode) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldIn(FieldProvisionMode, vs...))
}

// ProvisionModeNotIn applies the NotIn predicate on the "provision_mode" field.
func ProvisionModeNotIn(vs ...ProvisionMode) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNotIn(FieldProvisionMode, vs...))
}

// IsPartnerEQ applies the EQ predicate on the "is_partner" field.
func IsPartnerEQ(v bool) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldEQ(FieldIsPartner, v))
}

// IsPartnerNEQ applies the NEQ predicate on the "is_partner" field.
func IsPartnerNEQ(v bool) predicate.ProviderProfile {
	return predicate.ProviderProfile(sql.FieldNEQ(FieldIsPartner, v))
}

// HasUser applies the HasEdge predicate on the "user" edge.
func HasUser() predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, UserTable, UserColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasUserWith applies the HasEdge predicate on the "user" edge with a given conditions (other predicates).
func HasUserWith(preds ...predicate.User) predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		step := newUserStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasCurrency applies the HasEdge predicate on the "currency" edge.
func HasCurrency() predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, true, CurrencyTable, CurrencyColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasCurrencyWith applies the HasEdge predicate on the "currency" edge with a given conditions (other predicates).
func HasCurrencyWith(preds ...predicate.FiatCurrency) predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		step := newCurrencyStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProvisionBuckets applies the HasEdge predicate on the "provision_buckets" edge.
func HasProvisionBuckets() predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.M2M, true, ProvisionBucketsTable, ProvisionBucketsPrimaryKey...),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProvisionBucketsWith applies the HasEdge predicate on the "provision_buckets" edge with a given conditions (other predicates).
func HasProvisionBucketsWith(preds ...predicate.ProvisionBucket) predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		step := newProvisionBucketsStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasOrderTokens applies the HasEdge predicate on the "order_tokens" edge.
func HasOrderTokens() predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, OrderTokensTable, OrderTokensColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasOrderTokensWith applies the HasEdge predicate on the "order_tokens" edge with a given conditions (other predicates).
func HasOrderTokensWith(preds ...predicate.ProviderOrderToken) predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		step := newOrderTokensStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasAvailability applies the HasEdge predicate on the "availability" edge.
func HasAvailability() predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, AvailabilityTable, AvailabilityColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasAvailabilityWith applies the HasEdge predicate on the "availability" edge with a given conditions (other predicates).
func HasAvailabilityWith(preds ...predicate.ProviderAvailability) predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		step := newAvailabilityStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasProviderRating applies the HasEdge predicate on the "provider_rating" edge.
func HasProviderRating() predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2O, false, ProviderRatingTable, ProviderRatingColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasProviderRatingWith applies the HasEdge predicate on the "provider_rating" edge with a given conditions (other predicates).
func HasProviderRatingWith(preds ...predicate.ProviderRating) predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		step := newProviderRatingStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// HasAssignedOrders applies the HasEdge predicate on the "assigned_orders" edge.
func HasAssignedOrders() predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		step := sqlgraph.NewStep(
			sqlgraph.From(Table, FieldID),
			sqlgraph.Edge(sqlgraph.O2M, false, AssignedOrdersTable, AssignedOrdersColumn),
		)
		sqlgraph.HasNeighbors(s, step)
	})
}

// HasAssignedOrdersWith applies the HasEdge predicate on the "assigned_orders" edge with a given conditions (other predicates).
func HasAssignedOrdersWith(preds ...predicate.LockPaymentOrder) predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		step := newAssignedOrdersStep()
		sqlgraph.HasNeighborsWith(s, step, func(s *sql.Selector) {
			for _, p := range preds {
				p(s)
			}
		})
	})
}

// And groups predicates with the AND operator between them.
func And(predicates ...predicate.ProviderProfile) predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for _, p := range predicates {
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Or groups predicates with the OR operator between them.
func Or(predicates ...predicate.ProviderProfile) predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		s1 := s.Clone().SetP(nil)
		for i, p := range predicates {
			if i > 0 {
				s1.Or()
			}
			p(s1)
		}
		s.Where(s1.P())
	})
}

// Not applies the not operator on the given predicate.
func Not(p predicate.ProviderProfile) predicate.ProviderProfile {
	return predicate.ProviderProfile(func(s *sql.Selector) {
		p(s.Not())
	})
}
