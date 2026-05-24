package mgmadapter

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const sampleHex = "507f1f77bcf86cd799439011"

func TestHexaID_FromString(t *testing.T) {
	id := EmptyID()
	require.NoError(t, id.From(sampleHex))
	assert.Equal(t, sampleHex, id.String())

	oid, ok := id.Val().(primitive.ObjectID)
	require.True(t, ok)
	assert.Equal(t, sampleHex, oid.Hex())
}

func TestHexaID_FromObjectID(t *testing.T) {
	oid := primitive.NewObjectID()
	id := EmptyID()
	require.NoError(t, id.From(oid))
	assert.Equal(t, oid, id.Val())
}

func TestHexaID_FromInvalid(t *testing.T) {
	assert.Error(t, EmptyID().From("not-an-object-id"))
	assert.Error(t, EmptyID().From(nil))
	assert.Error(t, EmptyID().From(42)) // unsupported type
}

func TestHexaID_Validate(t *testing.T) {
	id := EmptyID()
	assert.NoError(t, id.Validate(sampleHex))
	assert.Error(t, id.Validate("bad"))
}

func TestHexaID_IsEqual(t *testing.T) {
	a := ID(sampleHex)
	b := ID(sampleHex)
	c := ID("507f1f77bcf86cd799439012")

	assert.True(t, a.IsEqual(b))
	assert.False(t, a.IsEqual(c))
	assert.False(t, a.IsEqual(nil))
}

func TestID_MustFromPanicsOnInvalid(t *testing.T) {
	assert.Panics(t, func() { ID("bad") })
	assert.NotNil(t, ID(sampleHex))
}

func TestIDField(t *testing.T) {
	oid := primitive.NewObjectID()

	f := &IDField{}
	f.SetID(oid)
	assert.Equal(t, oid, f.GetID())

	// PrepareID accepts a hex string and an ObjectID, and rejects bad hex.
	got, err := f.PrepareID(oid.Hex())
	require.NoError(t, err)
	assert.Equal(t, oid, got)

	got, err = f.PrepareID(oid)
	require.NoError(t, err)
	assert.Equal(t, oid, got)

	_, err = f.PrepareID("bad")
	assert.Error(t, err)
}

func TestIDField_SyncingSetsFreshID(t *testing.T) {
	f := &IDField{}
	require.True(t, f.ID.IsZero())
	require.NoError(t, f.Syncing())
	assert.False(t, f.ID.IsZero())

	// Syncing keeps an existing ID.
	existing := f.ID
	require.NoError(t, f.Syncing())
	assert.Equal(t, existing, f.ID)
}
