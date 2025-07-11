package api

import (
	"context"
	"log/slog"

	"golang.org/x/sync/errgroup"
)

type Job struct {
	Execute func()
}

type IJobQueue interface {
	Start(context.Context)
	Add(*Job)
}

type jobQueue struct {
	g    *errgroup.Group
	jobs chan *Job
	name string
}

func NeqJobQueue(name string, g *errgroup.Group) IJobQueue {
	slog.Info("NewJobQueue::", "name", name)
	return &jobQueue{
		jobs: make(chan *Job),
		g:    g,
		name: name,
	}
}

// Start starts a dispatcher.
// This dispatcher will stops when it receive a value from `ctx.Done`.
func (jq *jobQueue) Start(ctx context.Context) {
	slog.Info("JobQueue", jq.name, "::Start")
	jq.g.Go(func() error {
		defer close(jq.jobs)
	Loop:
		for {
			slog.Info("JobQueue", jq.name, "::Wait for Job")
			select {
			case <-ctx.Done():
				slog.Info("JobQueue", jq.name, "::Finish")
				break Loop

			case job := <-jq.jobs:

				slog.Info("JobQueue", jq.name, "::Do Job")
				jq.g.Go(func() error {
					job.Execute()
					return nil
				})

			}
		}
		slog.Info("JobQueue", jq.name, "::Done")
		return nil
	})
}

// Add enqueues a job into the queue.
// If the number of enqueued jobs has already reached to the maximum size,
// this will block until the other job has finish and the queue has space to accept a new job.
func (jq *jobQueue) Add(job *Job) {
	slog.Info("JobQueue", jq.name, "::Add")
	jq.g.Go(func() error {
		jq.jobs <- job
		return nil
	})
	slog.Info("JobQueue", jq.name, "::Added")
}
