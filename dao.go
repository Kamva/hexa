package hexa

import (
	"encoding/json"
	"github.com/kamva/gutil"
	"github.com/kamva/tracer"
)

//--------------------------------
// Data Access Object package
//--------------------------------

type (
	// DAO is the Data Access Object interface.
	DAO interface {
		Map() map[string]interface{}
		Json() ([]byte, error)
		MustJson() []byte
	}

	dao struct {
		obj interface{}
	}
)

func (d dao) Map() map[string]interface{} {
	if m, ok := d.obj.(map[string]interface{}); ok {
		return m
	}
	return gutil.StructToMap(d)
}

func (d dao) Json() ([]byte, error) {
	return json.Marshal(d.Map())
}

func (d dao) MustJson() []byte {
	b, err := d.Json()
	gutil.PanicErr(tracer.Trace(err))
	return b
}

// NewDAO returns new instance of the DAO.
func NewDAO(obj interface{}) DAO {
	return &dao{obj}
}

// Assertion
var _ DAO = &dao{}
