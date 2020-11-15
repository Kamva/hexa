package hexa

import "github.com/kamva/tracer"

type (
	// Rule is a rule signature.
	Rule func() error

	// BoolRule is a rule signature that specifies the rule is ok or not.
	BoolRule func() (bool, error)

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
			return tracer.Trace(err)
		}
	}
	return nil
}

// CheckRules check rules and returns first broken rule's error and whether its ok or not.
// otherwise returns true,nil.
func (rc RuleChecker) CheckBoolRules(rules ...BoolRule) (bool, error) {
	for _, r := range rules {
		if ok, err := r(); err != nil || !ok {
			return ok, tracer.Trace(err)
		}
	}
	return true, nil
}
