package mongolock

import (
	"context"
	"errors"
	"time"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/db/mgmadapter"
	"github.com/kamva/tracer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// CollectionName is default collection name.
const CollectionName = "locks"

type DlmOptions struct {
	Collection *mongo.Collection
	// If a lock already held by another mutex, we need to interval
	// to check it again.
	WaitingInterval time.Duration
	// Default ttl value for a lock. e.g, 2s.
	DefaultTTL time.Duration
	// usually your machine name. e.g., k8s-pod-199831uf.
	DefaultOwner string
}

// dlm implements the Hexa DLM.
type dlm struct {
	hexa.Health

	coll *mongo.Collection
	// owner is default lock owner.
	owner string
	// ttl is default lock ttl.
	ttl time.Duration
	// interval is waiting time before try to lock again if
	// lock already held by another mutex.
	interval time.Duration
}

func NewDlm(o DlmOptions) (hexa.DLM, error) {
	dlm := &dlm{
		Health: mgmadapter.NewDBHealth("distributed_locks", o.Collection.Database().Client()),

		coll:     o.Collection,
		ttl:      o.DefaultTTL,
		owner:    o.DefaultOwner,
		interval: o.WaitingInterval,
	}

	return dlm, tracer.Trace(dlm.createIndexesIfNotExists())
}

func (m *dlm) createIndexesIfNotExists() error {
	// Please note this index doesn't have any effect on the mutex behavior,
	// its just for cleanup.
	_, err := m.coll.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys:    bson.D{bson.E{Key: "expiry", Value: 1}},
		Options: &options.IndexOptions{Name: gutil.NewString("expired_locks")},
	})
	return tracer.Trace(err)
}

func (m *dlm) NewMutex(Key string) hexa.Mutex {
	return m.NewMutexWithOptions(hexa.MutexOptions{
		Key:   Key,
		Owner: m.owner,
		TTL:   m.ttl,
	})
}

func (m *dlm) NewMutexWithTTL(Key string, ttl time.Duration) hexa.Mutex {
	return m.NewMutexWithOptions(hexa.MutexOptions{Key: Key, TTL: ttl})
}

func (m *dlm) NewMutexWithOptions(o hexa.MutexOptions) hexa.Mutex {
	if o.Owner == "" {
		o.Owner = m.owner
	}

	return &mutex{
		coll:     m.coll,
		ttl:      o.TTL,
		interval: m.interval,

		ID:    o.Key,
		Owner: o.Owner,
	}
}

// mutex implements hexa Mutex distributed lock using MongoDB.
// it uses this query to update locks:
// coll.Query({_id:"my_key",$or:[{expiry:{$lt:now}},{owner:"me"}]},{new_data},{upsert: true})
// this query will get update my key or if key is expired, otherwise try to insert new key which
// either create the new key or get duplicate key error which means key held by another mutex.
type mutex struct {
	coll *mongo.Collection
	ttl  time.Duration

	// interval is waiting time before try to lock again if
	// lock already held by another mutex.
	interval time.Duration

	ID    string `json:"key" bson:"_id"`
	Owner string `json:"owner" bson:"owner"`
	// Expiry begins when we lock the mutex.
	Expiry time.Time `json:"expiry" bson:"expiry"`
}

// Lock try to lock and if lock is held by another mutex, it wait and
// try it again.
func (m *mutex) Lock(c hexa.Context) error {
	for {
		err := m.TryLock(c)

		if !errors.Is(err, hexa.ErrLockAlreadyAcquired) {
			return tracer.Trace(err)
		}

		select {
		case <-c.Done():
			return c.Err()
		case <-time.After(m.interval):
		}
	}
}

func (m *mutex) TryLock(c hexa.Context) error {
	m.Expiry = time.Now().Add(m.ttl)

	filter := bson.M{"_id": m.ID, "$or": []bson.M{
		{"owner": m.Owner},
		{"expiry": bson.M{"$lt": time.Now()}},
	}}

	_, err := m.coll.UpdateOne(c, filter, bson.M{"$set": &m}, &options.UpdateOptions{
		Upsert: gutil.NewBool(true),
	})

	if isDup(err) {
		err = hexa.ErrLockAlreadyAcquired
	}

	return tracer.Trace(err)
}

func (m *mutex) Unlock(c hexa.Context) error {
	_, err := m.coll.DeleteOne(c, bson.M{"_id": m.ID, "owner": m.Owner})
	return tracer.Trace(err)
}

var _ hexa.DLM = &dlm{}
var _ hexa.Mutex = &mutex{}
