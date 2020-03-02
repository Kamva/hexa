package mgmadapter

import (
	"errors"
	"github.com/Kamva/tracer"
	"go.mongodb.org/mongo-driver/mongo"
)

type Repository struct {
}

func (a Repository) ReplaceErr(err error, notfoundErr error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return notfoundErr
	}

	return tracer.Trace(err)
}
