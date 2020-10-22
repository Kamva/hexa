package mgmadapter

import (
	"context"
	"errors"
	"github.com/Kamva/mgm/v3"
	"github.com/kamva/gutil"
	"github.com/kamva/tracer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Repository struct {
}

func (r Repository) ReplaceErr(err error, notfoundErr error) error {
	if errors.Is(err, mongo.ErrNoDocuments) {
		return notfoundErr
	}

	return tracer.Trace(err)
}

func (r Repository) CreateIndexIfNotExists(coll *mgm.Collection, name string, fields ...string) error {
	return r.CreateIndexWithOptionsIfNotExists(coll, &options.IndexOptions{Name: gutil.NewString(name)}, fields...)
}

func (r Repository) CreateUniqueIndexIfNotExists(coll *mgm.Collection, name string, fields ...string) error {
	o := options.IndexOptions{
		Name:   gutil.NewString(name),
		Unique: gutil.NewBool(true),
	}
	return r.CreateIndexWithOptionsIfNotExists(coll, &o, fields...)
}

func (r *Repository) CreateIndexWithOptionsIfNotExists(coll *mgm.Collection, o *options.IndexOptions, fields ...string) error {
	keys := make(bson.D, len(fields))
	for i, k := range fields {
		keys[i] = bson.E{Key: k, Value: 1}
	}
	_, err := coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    keys,
		Options: o,
	})
	return tracer.Trace(err)
}
