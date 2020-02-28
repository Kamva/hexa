package jobs

import (
	"errors"
	"github.com/Kamva/gutil"
	"github.com/Kamva/kitty"
	"github.com/Kamva/tracer"
	"github.com/contribsys/faktory/client"
	worker "github.com/contribsys/faktory_worker_go"
	"time"
)

type (
	// faktoryJobs is implementation of kitty Jobs using faktory
	faktoryJobs struct {
		p *client.Pool
	}

	// faktoryWorker is implementation of kitty worker using faktory.
	faktoryWorker struct {
		w  *worker.Manager
		uf kitty.UserFinder
		l  kitty.Logger
		t  kitty.Translator
	}
)

func (j *faktoryJobs) prepare(c kitty.Context, job *kitty.Job) *client.Job {

	if job.Queue == "" {
		job.Queue = "default"
	}

	ctxMap := c.ToMap()
	return &client.Job{
		Jid:       client.RandomJid(),
		Type:      job.Name,
		Queue:     job.Queue,
		Args:      []interface{}{ctxMap, gutil.StructToMap(job.Payload)},
		CreatedAt: time.Now().UTC().Format(time.RFC3339Nano),
		Retry:     job.Retry,
		// We don't using this custom data in any middleware, but just put it here :)
		Custom: ctxMap,
	}
}

func (j *faktoryJobs) Push(ctx kitty.Context, job *kitty.Job) error {
	if job == nil || job.Name == "" {
		return tracer.Trace(errors.New("job is not valid (enter job name please)"))
	}

	return j.p.With(func(conn *client.Client) error {
		return conn.Push(j.prepare(ctx, job))
	})
}

func (w *faktoryWorker) handler(h kitty.JobHandlerFunc) worker.Perform {
	return func(ctx worker.Context, args ...interface{}) error {

		var payload kitty.Payload
		ctxMap := args[0].(map[string]interface{})
		err := gutil.MapToStruct(args[1].(map[string]interface{}), &payload)

		if err != nil {
			return tracer.Trace(err)
		}

		kCtx, err := kitty.CtxFromMap(ctxMap, w.uf, w.l, w.t)

		if err != nil {
			return tracer.Trace(err)
		}
		return h(kCtx, payload)
	}
}

func (w *faktoryWorker) Register(name string, h kitty.JobHandlerFunc) error {
	w.w.Register(name, w.handler(h))
	return nil
}

func (w *faktoryWorker) Concurrency(c int) error {
	w.w.Concurrency = c
	return nil
}

func (w *faktoryWorker) Process(queues ...string) error {
	w.w.ProcessStrictPriorityQueues(queues...)
	// Run method does not return.
	w.w.Run()

	return nil
}

// NewFaktoryJobsDriver returns new instance of Jobs driver for the faktory
func NewFaktoryJobsDriver(p *client.Pool) kitty.Jobs {
	return &faktoryJobs{p}
}

// NewFaktoryWorkerDriver returns new instance of kitty Worker driver for the faktory
func NewFaktoryWorkerDriver(w *worker.Manager, uf kitty.UserFinder, l kitty.Logger, t kitty.Translator) kitty.Worker {
	return &faktoryWorker{
		w:  w,
		uf: uf,
		l:  l,
		t:  t,
	}
}

var _ kitty.Jobs = &faktoryJobs{}
var _ kitty.Worker = &faktoryWorker{}
