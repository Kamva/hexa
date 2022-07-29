package hexa

import "context"

// WrapRule wraps the Rule and returns a RuleWithContext.
func WrapRule(r Rule) RuleWithContext {
	return func(_ context.Context) error {
		return r()
	}
}

// Rule is a rule signature.
type Rule func() error
type RuleWithContext func(ctx context.Context) error

// Validate validates the provided rules and returns first broken
// rule's error. otherwise it returns nil.
func Validate(rules ...Rule) error {
	for _, r := range rules {
		if err := r(); err != nil {
			return err
		}
	}
	return nil
}

// ValidateWithContext validates rules with a context.
func ValidateWithContext(ctx context.Context, rules ...RuleWithContext) error {
	for _, r := range rules {
		if err := r(ctx); err != nil {
			return err
		}

		// We can check if context is done.
		//if err := ctx.Err(); err != nil {
		//	return nil
		//}
	}
	return nil
}
