package mongolock

import (
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

// Mongo will add this method to its repo, until that time, use this helper.
// [see this topic](https://developer.mongodb.com/community/forums/t/isdup-function-missing-golang/2926)
func isDup(err error) bool {
	var e mongo.WriteException
	if errors.As(err, &e) {
		for _, we := range e.WriteErrors {
			if we.Code == 11000 {
				return true
			}
		}
	}
	return false
}
