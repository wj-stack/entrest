// Code generated by ent, DO NOT EDIT.

package privacy

import (
	"context"

	"github.com/lrstanley/entrest/_examples/kitchensink/database/ent"

	"entgo.io/ent/privacy"
)

var (
	// Allow may be returned by rules to indicate that the policy
	// evaluation should terminate with allow decision.
	Allow = privacy.Allow

	// Deny may be returned by rules to indicate that the policy
	// evaluation should terminate with deny decision.
	Deny = privacy.Deny

	// Skip may be returned by rules to indicate that the policy
	// evaluation should continue to the next rule.
	Skip = privacy.Skip
)

// Allowf returns a formatted wrapped Allow decision.
func Allowf(format string, a ...any) error {
	return privacy.Allowf(format, a...)
}

// Denyf returns a formatted wrapped Deny decision.
func Denyf(format string, a ...any) error {
	return privacy.Denyf(format, a...)
}

// Skipf returns a formatted wrapped Skip decision.
func Skipf(format string, a ...any) error {
	return privacy.Skipf(format, a...)
}

// DecisionContext creates a new context from the given parent context with
// a policy decision attach to it.
func DecisionContext(parent context.Context, decision error) context.Context {
	return privacy.DecisionContext(parent, decision)
}

// DecisionFromContext retrieves the policy decision from the context.
func DecisionFromContext(ctx context.Context) (error, bool) {
	return privacy.DecisionFromContext(ctx)
}

type (
	// Policy groups query and mutation policies.
	Policy = privacy.Policy

	// QueryRule defines the interface deciding whether a
	// query is allowed and optionally modify it.
	QueryRule = privacy.QueryRule
	// QueryPolicy combines multiple query rules into a single policy.
	QueryPolicy = privacy.QueryPolicy

	// MutationRule defines the interface which decides whether a
	// mutation is allowed and optionally modifies it.
	MutationRule = privacy.MutationRule
	// MutationPolicy combines multiple mutation rules into a single policy.
	MutationPolicy = privacy.MutationPolicy
	// MutationRuleFunc type is an adapter which allows the use of
	// ordinary functions as mutation rules.
	MutationRuleFunc = privacy.MutationRuleFunc

	// QueryMutationRule is an interface which groups query and mutation rules.
	QueryMutationRule = privacy.QueryMutationRule
)

// QueryRuleFunc type is an adapter to allow the use of
// ordinary functions as query rules.
type QueryRuleFunc func(context.Context, ent.Query) error

// Eval returns f(ctx, q).
func (f QueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	return f(ctx, q)
}

// AlwaysAllowRule returns a rule that returns an allow decision.
func AlwaysAllowRule() QueryMutationRule {
	return privacy.AlwaysAllowRule()
}

// AlwaysDenyRule returns a rule that returns a deny decision.
func AlwaysDenyRule() QueryMutationRule {
	return privacy.AlwaysDenyRule()
}

// ContextQueryMutationRule creates a query/mutation rule from a context eval func.
func ContextQueryMutationRule(eval func(context.Context) error) QueryMutationRule {
	return privacy.ContextQueryMutationRule(eval)
}

// OnMutationOperation evaluates the given rule only on a given mutation operation.
func OnMutationOperation(rule MutationRule, op ent.Op) MutationRule {
	return privacy.OnMutationOperation(rule, op)
}

// DenyMutationOperationRule returns a rule denying specified mutation operation.
func DenyMutationOperationRule(op ent.Op) MutationRule {
	rule := MutationRuleFunc(func(_ context.Context, m ent.Mutation) error {
		return Denyf("ent/privacy: operation %s is not allowed", m.Op())
	})
	return OnMutationOperation(rule, op)
}

// The CategoryQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type CategoryQueryRuleFunc func(context.Context, *ent.CategoryQuery) error

// EvalQuery return f(ctx, q).
func (f CategoryQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.CategoryQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.CategoryQuery", q)
}

// The CategoryMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type CategoryMutationRuleFunc func(context.Context, *ent.CategoryMutation) error

// EvalMutation calls f(ctx, m).
func (f CategoryMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.CategoryMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.CategoryMutation", m)
}

// The FollowsQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type FollowsQueryRuleFunc func(context.Context, *ent.FollowsQuery) error

// EvalQuery return f(ctx, q).
func (f FollowsQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.FollowsQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.FollowsQuery", q)
}

// The FollowsMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type FollowsMutationRuleFunc func(context.Context, *ent.FollowsMutation) error

// EvalMutation calls f(ctx, m).
func (f FollowsMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.FollowsMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.FollowsMutation", m)
}

// The FriendshipQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type FriendshipQueryRuleFunc func(context.Context, *ent.FriendshipQuery) error

// EvalQuery return f(ctx, q).
func (f FriendshipQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.FriendshipQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.FriendshipQuery", q)
}

// The FriendshipMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type FriendshipMutationRuleFunc func(context.Context, *ent.FriendshipMutation) error

// EvalMutation calls f(ctx, m).
func (f FriendshipMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.FriendshipMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.FriendshipMutation", m)
}

// The PetQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type PetQueryRuleFunc func(context.Context, *ent.PetQuery) error

// EvalQuery return f(ctx, q).
func (f PetQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.PetQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.PetQuery", q)
}

// The PetMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type PetMutationRuleFunc func(context.Context, *ent.PetMutation) error

// EvalMutation calls f(ctx, m).
func (f PetMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.PetMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.PetMutation", m)
}

// The SettingsQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type SettingsQueryRuleFunc func(context.Context, *ent.SettingsQuery) error

// EvalQuery return f(ctx, q).
func (f SettingsQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.SettingsQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.SettingsQuery", q)
}

// The SettingsMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type SettingsMutationRuleFunc func(context.Context, *ent.SettingsMutation) error

// EvalMutation calls f(ctx, m).
func (f SettingsMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.SettingsMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.SettingsMutation", m)
}

// The SkippedQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type SkippedQueryRuleFunc func(context.Context, *ent.SkippedQuery) error

// EvalQuery return f(ctx, q).
func (f SkippedQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.SkippedQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.SkippedQuery", q)
}

// The SkippedMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type SkippedMutationRuleFunc func(context.Context, *ent.SkippedMutation) error

// EvalMutation calls f(ctx, m).
func (f SkippedMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.SkippedMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.SkippedMutation", m)
}

// The UserQueryRuleFunc type is an adapter to allow the use of ordinary
// functions as a query rule.
type UserQueryRuleFunc func(context.Context, *ent.UserQuery) error

// EvalQuery return f(ctx, q).
func (f UserQueryRuleFunc) EvalQuery(ctx context.Context, q ent.Query) error {
	if q, ok := q.(*ent.UserQuery); ok {
		return f(ctx, q)
	}
	return Denyf("ent/privacy: unexpected query type %T, expect *ent.UserQuery", q)
}

// The UserMutationRuleFunc type is an adapter to allow the use of ordinary
// functions as a mutation rule.
type UserMutationRuleFunc func(context.Context, *ent.UserMutation) error

// EvalMutation calls f(ctx, m).
func (f UserMutationRuleFunc) EvalMutation(ctx context.Context, m ent.Mutation) error {
	if m, ok := m.(*ent.UserMutation); ok {
		return f(ctx, m)
	}
	return Denyf("ent/privacy: unexpected mutation type %T, expect *ent.UserMutation", m)
}