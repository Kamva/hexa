package hexa

// Rule is a rule signature.
type Rule func() error

// VerifyRules verifies the provided rules and returns first broken
// rule's error. otherwise it returns nil.
func VerifyRules(rules ...Rule) error {
	for _, r := range rules {
		if err := r(); err != nil {
			return err
		}
	}
	return nil
}
