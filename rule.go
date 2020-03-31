package hexa

type (
	// Rule is a rule signature.
	Rule func() error

	// RuleChecker get a list of rules and check if a rule is broken,
	// returns that rule error message.
	// You can use rule checker in your aggregates,
	//entities, valueObjects or services,...
	RuleChecker struct{}
)

// CheckRules check rules and returns first broken rule's error.
// otherwise returns nil.
func (rc RuleChecker) CheckRules(rules ...Rule) error {
	for _, r := range rules {
		if err := r(); err != nil {
			return err
		}
	}
	return nil
}
