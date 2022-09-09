package mongolock

import (
	"context"
	"testing"
	"time"

	"github.com/kamva/gutil"
	"github.com/kamva/hexa"
	"github.com/kamva/hexa/hexatranslator"
	"github.com/kamva/hexa/hlog"
	"github.com/kamva/mgm/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var cli *mongo.Client
var collection *mongo.Collection

func newCtx(ctx context.Context) context.Context {
	return hexa.NewContext(ctx, hexa.ContextParams{
		CorrelationId:  "_cid",
		Locale:         "def",
		User:           hexa.NewGuest(),
		BaseLogger:     hlog.NewPrinterDriver(hlog.DebugLevel),
		BaseTranslator: hexatranslator.NewEmptyDriver(),
	})
}

func setupDefConnection() {
	gutil.PanicErr(
		mgm.SetDefaultConfig(nil, "locks"),
	)
	var err error
	cli, err = mongo.NewClient(options.Client().ApplyURI("mongodb://root:12345@localhost:27017"))
	gutil.PanicErr(err)

	gutil.PanicErr(cli.Connect(context.Background()))

	collection = cli.Database("lock_labs").Collection(CollectionName)
}
func disconnect() {
	if err := cli.Disconnect(context.Background()); err != nil {
		panic(err)
	}
}

func resetCollection() {
	_, err := collection.DeleteMany(mgm.Ctx(), bson.M{})
	gutil.PanicErr(err)
}

func TestNewDlmDatabaseIndex(t *testing.T) {
	setupDefConnection()
	defer disconnect()

	resetCollection()
	gutil.PanicErr(collection.Drop(mgm.Ctx()))
	dlm, err := NewDlm(DlmOptions{
		Collection:      collection,
		WaitingInterval: time.Millisecond * 200,
		DefaultTTL:      time.Second * 60,
		DefaultOwner:    "lab",
	})

	require.Nil(t, err)
	require.NotNil(t, dlm)

	// calling createIndex() with the same name but different options than an existing
	// index will throw an error MongoError: Index with name: {indexName} already
	// exists with different options, so one time we create the true index and assert
	// error be nil, next time we create wrong options and assert error is IndexExists.
	_, err = collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{bson.E{Key: "expiry", Value: 1}},
		Options: &options.IndexOptions{
			ExpireAfterSeconds: gutil.NewInt32(0),
			Name:               gutil.NewString("expired_locks"),
		},
	})

	require.NotNil(t, err)

	_, err = collection.Indexes().CreateOne(context.Background(), mongo.IndexModel{
		Keys: bson.D{bson.E{Key: "expiry", Value: 1}},
		Options: &options.IndexOptions{
			Name: gutil.NewString("expired_locks"),
		},
	})

	require.Nil(t, err)
}

func TestDlm_NewMutex(t *testing.T) {
	setupDefConnection()
	defer disconnect()

	resetCollection()

	ttl := time.Second * 60
	dlm, err := NewDlm(DlmOptions{
		Collection:      collection,
		WaitingInterval: time.Millisecond * 200,
		DefaultTTL:      ttl,
		DefaultOwner:    "lab",
	})

	require.Nil(t, err)
	m := dlm.NewMutex("abc")
	var mObj = m.(*mutex)
	require.NotNil(t, m)

	now := time.Now().Add(ttl)
	time.Sleep(time.Millisecond * 200)

	assert.Nil(t, m.Lock(context.Background()))
	assert.Equal(t, mObj.ID, "abc")
	assert.Equal(t, mObj.Owner, "lab")
	assert.Equal(t, mObj.ttl, ttl)
	assert.True(t, now.Before(mObj.Expiry)) // now+ttl < lock_time + ttl.
}

func TestDlm_NewMutexWithTTL(t *testing.T) {
	setupDefConnection()
	defer disconnect()

	resetCollection()

	ttl := time.Second * 30
	dlm, err := NewDlm(DlmOptions{
		Collection:      collection,
		WaitingInterval: time.Millisecond * 200,
		DefaultTTL:      ttl,
		DefaultOwner:    "lab",
	})

	require.Nil(t, err)
	mttl := time.Second * 60
	m := dlm.NewMutexWithTTL("abc", mttl)
	var mObj = m.(*mutex)
	require.NotNil(t, m)

	now := time.Now().Add(ttl)
	time.Sleep(time.Millisecond * 200)

	assert.Nil(t, m.Lock(context.Background()))
	assert.Equal(t, mObj.ID, "abc")
	assert.Equal(t, mObj.Owner, "lab")
	assert.Equal(t, mObj.ttl, mttl)
	assert.True(t, now.Before(mObj.Expiry)) // now+ttl < lock_time + ttl.
}

func TestDlm_NewMutexValuesWithMutexOptions(t *testing.T) {
	setupDefConnection()
	defer disconnect()

	resetCollection()

	ttl := time.Second * 60
	dlm, err := NewDlm(DlmOptions{
		Collection:      collection,
		WaitingInterval: time.Millisecond * 200,
		DefaultTTL:      ttl,
		DefaultOwner:    "lab",
	})

	require.Nil(t, err)
	m := dlm.NewMutexWithOptions(hexa.MutexOptions{
		Key: "abc",
		TTL: time.Second,
	})
	mObj := m.(*mutex)
	require.NotNil(t, m)
	assert.Equal(t, mObj.ID, "abc")
	assert.Equal(t, mObj.Owner, "lab")
	assert.Equal(t, mObj.ttl, time.Second)

	m = dlm.NewMutexWithOptions(hexa.MutexOptions{
		Key:   "abcd",
		Owner: "123",
		TTL:   time.Second * 2,
	})
	mObj = m.(*mutex)
	require.NotNil(t, m)
	assert.Equal(t, mObj.ID, "abcd")
	assert.Equal(t, mObj.Owner, "123")
	assert.Equal(t, mObj.ttl, time.Second*2)
}

func TestDlm_NewMutexRefreshAndMultipleCall(t *testing.T) {
	setupDefConnection()
	defer disconnect()
	resetCollection()

	ttl := time.Second * 60
	dlm, err := NewDlm(DlmOptions{
		Collection:      collection,
		WaitingInterval: time.Millisecond * 200,
		DefaultTTL:      ttl,
		DefaultOwner:    "lab",
	})

	require.Nil(t, err)
	m := dlm.NewMutex("abc")
	var mObj = m.(*mutex)
	require.NotNil(t, m)

	oldExpiry := mObj.Expiry
	time.Sleep(time.Second)
	// re-lock must update expiry:
	ctx := context.Background()
	assert.Nil(t, m.Lock(ctx))
	assert.NotEqual(t, mObj.Expiry, oldExpiry)
	assert.True(t, mObj.Expiry.After(oldExpiry))

	assert.Nil(t, m.Unlock(ctx))
	assert.Nil(t, m.Lock(ctx))
	assert.Nil(t, m.Lock(ctx))
	assert.Nil(t, m.TryLock(ctx))
	assert.Nil(t, m.Lock(ctx))

	// all data must be remained untouched.
	assert.Equal(t, mObj.ID, "abc")
	assert.Equal(t, mObj.Owner, "lab")
	assert.Equal(t, mObj.ttl, ttl)
}

func TestMutexDataInDB(t *testing.T) {
	setupDefConnection()
	defer disconnect()

	resetCollection()

	dlm, err := NewDlm(DlmOptions{
		Collection:      collection,
		WaitingInterval: time.Millisecond * 200,
		DefaultTTL:      time.Second * 60,
		DefaultOwner:    "lab",
	})

	require.Nil(t, err)
	m := dlm.NewMutex("abc")
	var mObj = m.(*mutex)

	ctx := context.Background()
	require.NotNil(t, m)
	assert.Nil(t, m.Lock(ctx))
	count, err := collection.CountDocuments(ctx, bson.M{})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), count)
	res := collection.FindOne(ctx, bson.M{})
	require.Nil(t, res.Err())

	var doc mutex
	require.Nil(t, res.Decode(&doc))
	assert.Equal(t, mObj.ID, doc.ID)
	assert.Equal(t, mObj.Owner, doc.Owner)
	assert.Equal(t, mObj.Expiry.Unix(), doc.Expiry.Unix())
}

func TestMutex_Lock(t *testing.T) {
	setupDefConnection()
	defer disconnect()

	resetCollection()

	ttl := time.Second * 2
	dlm, err := NewDlm(DlmOptions{
		Collection:      collection,
		WaitingInterval: time.Millisecond * 200,
		DefaultTTL:      ttl,
		DefaultOwner:    "machine-1",
	})

	require.Nil(t, err)
	m1 := dlm.NewMutex("abc")
	m2 := dlm.NewMutexWithOptions(hexa.MutexOptions{
		Key:   "abc",
		Owner: "machine-2",
		TTL:   ttl,
	})

	assert.Nil(t, m1.Lock(context.Background()))

	start := time.Now()
	// ctx will expire in ttl + one second later, so lock should be acquired before ttl+a second.
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(ttl+time.Second))
	defer cancel()
	require.Nil(t, m2.Lock(newCtx(ctx)))

	// we should acquire lock after ttl time.
	assert.True(t, time.Since(start) > ttl)

	//--------------------------------
	// Repeat it for the m1 again
	//--------------------------------

	start = time.Now()
	// ctx will expire in ttl + one second later, so lock should be acquired before ttl+a second.
	ctx, cancel = context.WithDeadline(context.Background(), time.Now().Add(ttl+time.Second))
	defer cancel()
	require.Nil(t, m1.Lock(newCtx(ctx)))

	// we should acquire lock after ttl time.
	assert.True(t, time.Since(start) > ttl)
}

func TestMutex_LockAndReLock(t *testing.T) {
	setupDefConnection()
	defer disconnect()

	resetCollection()

	ttl := time.Second * 2
	dlm, err := NewDlm(DlmOptions{
		Collection:      collection,
		WaitingInterval: time.Millisecond * 200,
		DefaultTTL:      ttl,
		DefaultOwner:    "machine-1",
	})

	require.Nil(t, err)
	m1 := dlm.NewMutex("abc")
	m2 := dlm.NewMutexWithOptions(hexa.MutexOptions{
		Key:   "abc",
		Owner: "machine-2",
		TTL:   ttl,
	})
	ctx := context.Background()
	assert.Nil(t, m1.Lock(ctx))

	// re-lock the m1 mutex two times, each time after one second.
	releaseAfter := ttl + time.Second*2
	go func() {
		// waiting for one second and re-lock it again, so the lock should release after ttl + 2
		time.Sleep(time.Second)
		assert.Nil(t, m1.Lock(ctx))
		time.Sleep(time.Second)
		assert.Nil(t, m1.Lock(ctx))
	}()

	start := time.Now()
	// ctx will expire in "releaseAfter + one second" later, so lock should be acquired before ctx expiry.
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(releaseAfter+time.Second))
	defer cancel()
	require.Nil(t, m2.Lock(newCtx(ctx)))

	// we should acquire lock after releaseAfter time elapsed.
	assert.True(t, time.Since(start) > releaseAfter)
}

func TestMutex_UnlockBeforeExpiration(t *testing.T) {
	setupDefConnection()
	defer disconnect()
	resetCollection()

	ttl := time.Second * 3
	dlm, err := NewDlm(DlmOptions{
		Collection:      collection,
		WaitingInterval: time.Millisecond * 200,
		DefaultTTL:      ttl,
		DefaultOwner:    "machine-1",
	})

	require.Nil(t, err)
	m1 := dlm.NewMutex("abc")
	m2 := dlm.NewMutexWithOptions(hexa.MutexOptions{
		Key:   "abc",
		Owner: "machine-2",
		TTL:   ttl,
	})

	ctx := context.Background()
	assert.Nil(t, m1.Lock(ctx))

	// unlock after one second.
	releaseAfter := time.Second
	go func() {
		// waiting for one second and unlock it.
		time.Sleep(time.Second)
		assert.Nil(t, m1.Unlock(ctx))
	}()

	start := time.Now()
	// ctx will expire in "releaseAfter + one second" later, so lock should be acquired before ctx expiry.
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(releaseAfter+time.Second))
	defer cancel()
	require.Nil(t, m2.Lock(newCtx(ctx)))

	// we should acquire lock after releaseAfter time elapsed.
	assert.True(t, time.Since(start) > releaseAfter)
}
