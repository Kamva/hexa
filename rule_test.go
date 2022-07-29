package hexa

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func successRule() Rule {
	return func() error { return nil }
}

func failedRule(err error) Rule {
	return func() error {
		return err
	}
}

func TestVerifyRules(t *testing.T) {
	err1 := errors.New("err1")
	err2 := errors.New("err2")

	assert.Nil(t, Validate())
	assert.Nil(t, Validate(successRule()))
	assert.Equal(t, err1, Validate(failedRule(err1), successRule(), failedRule(err2)))
	assert.Equal(t, err1, Validate(successRule(), failedRule(err1), failedRule(err2)))
	assert.Equal(t, err2, Validate(failedRule(err2), failedRule(err1)))
}
