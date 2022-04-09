package mgmadapter

import (
	"context"
	"errors"

	"github.com/kamva/gutil"
	"github.com/kamva/mgm/v3"
	"github.com/kamva/tracer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Store struct {
}

func (r Store) ReplaceErr(err error, notfoundErr error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return notfoundErr
	}

	return tracer.Trace(err)
}

func (r Store) CreateIndexIfNotExist(coll *mgm.Collection, name string, keys ...interface{}) error {
	return r.CreateIndexWithOptionsIfNotExist(coll, &options.IndexOptions{Name: gutil.NewString(name)}, keys...)
}

func (r Store) CreateUniqueIndexIfNotExist(coll *mgm.Collection, name string, keys ...interface{}) error {
	o := options.IndexOptions{
		Name:   gutil.NewString(name),
		Unique: gutil.NewBool(true),
	}
	return r.CreateIndexWithOptionsIfNotExist(coll, &o, keys...)
}

// CreateIndexWithOptionsIfNotExist creates index if it doesn't exist. fields value could be either string or bson.E.
// One exception is if the first field value's type is bson.D, we will use it and ignore all other fields values.
// otherwise if you send another type as field it will panic.
func (r *Store) CreateIndexWithOptionsIfNotExist(coll *mgm.Collection, o *options.IndexOptions, fields ...interface{}) error {
	var keys bson.D
	if len(fields) == 1 {
		if d, ok := fields[0].(bson.D); ok {
			keys = d
		}
	}

	if keys == nil { // If it not initialized yet
		keys = make(bson.D, len(fields))
		for i, k := range fields {
			if key, ok := k.(string); ok {
				keys[i] = bson.E{Key: key, Value: 1}
				continue
			}

			keys[i] = k.(bson.E)
		}
	}

	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    keys,
		Options: o,
	})
	return tracer.Trace(err)
}
