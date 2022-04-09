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

func (r Store) CreateIndexIfNotExists(coll *mgm.Collection, name string, keys bson.D) error {
	return r.CreateIndexWithOptionsIfNotExists(coll, &options.IndexOptions{Name: gutil.NewString(name)}, keys)
}

func (r Store) CreateUniqueIndexIfNotExists(coll *mgm.Collection, name string, keys bson.D) error {
	o := options.IndexOptions{
		Name:   gutil.NewString(name),
		Unique: gutil.NewBool(true),
	}
	return r.CreateIndexWithOptionsIfNotExists(coll, &o, keys)
}

func (r *Store) CreateIndexWithOptionsIfNotExists(coll *mgm.Collection, o *options.IndexOptions, keys bson.D) error {
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    keys,
		Options: o,
	})
	return tracer.Trace(err)
}

