package scheduler

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/hashicorp/terraform-ls/internal/document"
	"github.com/hashicorp/terraform-ls/internal/job"
	"github.com/hashicorp/terraform-ls/internal/state"
)

func TestScheduler_closedOnly(t *testing.T) {
	ss, err := state.NewStateStore()
	if err != nil {
		t.Fatal(err)
	}
	ss.SetLogger(testLogger())

	tmpDir := t.TempDir()

	ctx := context.Background()

	s := NewScheduler(&closedDirJobs{js: ss.JobStore}, 2)
	s.SetLogger(testLogger())
	s.Start(ctx)
	t.Cleanup(func() {
		s.Stop()
	})

	var jobsExecuted int64 = 0
	jobsToExecute := 50

	ids := make([]job.ID, 0)
	for i := 0; i < jobsToExecute; i++ {
		i := i
		dirPath := filepath.Join(tmpDir, fmt.Sprintf("folder-%d", i))

		newId, err := ss.JobStore.EnqueueJob(job.Job{
			Func: func(c context.Context) error {
				atomic.AddInt64(&jobsExecuted, 1)
				return nil
			},
			Dir:  document.DirHandleFromPath(dirPath),
			Type: "test-type",
		})
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, newId)
	}

	err = ss.JobStore.WaitForJobs(ctx, ids...)
	if err != nil {
		t.Fatal(err)
	}

	if jobsExecuted != int64(jobsToExecute) {
		t.Fatalf("expected %d jobs to execute, got: %d", jobsToExecute, jobsExecuted)
	}
}

func TestScheduler_closedAndOpen(t *testing.T) {
	ss, err := state.NewStateStore()
	if err != nil {
		t.Fatal(err)
	}
	ss.SetLogger(testLogger())

	tmpDir := t.TempDir()

	var wg sync.WaitGroup

	var closedJobsExecuted int64 = 0
	closedJobsToExecute := 50
	closedIds := make([]job.ID, 0)
	wg.Add(1)
	go func(t *testing.T) {
		defer wg.Done()
		for i := 0; i < closedJobsToExecute; i++ {
			i := i
			dirPath := filepath.Join(tmpDir, fmt.Sprintf("folder-x-%d", i))

			newId, err := ss.JobStore.EnqueueJob(job.Job{
				Func: func(c context.Context) error {
					atomic.AddInt64(&closedJobsExecuted, 1)
					return nil
				},
				Dir:  document.DirHandleFromPath(dirPath),
				Type: "test-type",
			})
			if err != nil {
				t.Error(err)
			}
			closedIds = append(closedIds, newId)
		}
	}(t)

	openJobsToExecute := 50
	var openJobsExecuted int64 = 0
	openIds := make([]job.ID, 0)
	wg.Add(1)
	go func(t *testing.T) {
		defer wg.Done()
		for i := 0; i < openJobsToExecute; i++ {
			i := i
			dirPath := filepath.Join(tmpDir, fmt.Sprintf("folder-y-%d", i))

			newId, err := ss.JobStore.EnqueueJob(job.Job{
				Func: func(c context.Context) error {
					atomic.AddInt64(&openJobsExecuted, 1)
					return nil
				},
				Dir:  document.DirHandleFromPath(dirPath),
				Type: "test-type",
			})
			if err != nil {
				t.Error(err)
			}

			openIds = append(openIds, newId)
		}
	}(t)

	wg.Add(1)
	// we intentionally open the documents in a separate routine,
	// possibly after some of the relevant jobs have been queued (as closed)
	// to better reflect what may happen in reality
	go func(t *testing.T) {
		defer wg.Done()
		for i := 0; i < openJobsToExecute; i++ {
			dirPath := filepath.Join(tmpDir, fmt.Sprintf("folder-y-%d", i))
			dh := document.HandleFromPath(filepath.Join(dirPath, "test.tf"))
			err := ss.DocumentStore.OpenDocument(dh, "", 0, []byte{})
			if err != nil {
				t.Error(err)
			}
		}
	}(t)

	ctx := context.Background()
	if deadline, ok := t.Deadline(); ok {
		var cancelFunc context.CancelFunc
		ctx, cancelFunc = context.WithDeadline(ctx, deadline)
		t.Cleanup(cancelFunc)
	}

	cs := NewScheduler(&closedDirJobs{js: ss.JobStore}, 1)
	cs.SetLogger(testLogger())
	cs.Start(ctx)
	t.Cleanup(func() {
		cs.Stop()
	})

	os := NewScheduler(&openDirJobs{js: ss.JobStore}, 1)
	os.SetLogger(testLogger())
	os.Start(ctx)
	t.Cleanup(func() {
		os.Stop()
	})

	// wait for all scheduling and document opening to finish
	wg.Wait()
	t.Log("finished all scheduling and doc opening")

	allIds := make([]job.ID, 0)
	allIds = append(allIds, closedIds...)
	allIds = append(allIds, openIds...)

	t.Logf("waiting for %d jobs", len(allIds))
	err = ss.JobStore.WaitForJobs(ctx, allIds...)
	if err != nil {
		t.Fatal(err)
	}

	if closedJobsExecuted != int64(closedJobsToExecute) {
		t.Fatalf("expected %d closed jobs to execute, got: %d", closedJobsToExecute, closedJobsExecuted)
	}

	if openJobsExecuted != int64(openJobsToExecute) {
		t.Fatalf("expected %d open jobs to execute, got: %d", openJobsToExecute, openJobsExecuted)
	}
}

func BenchmarkScheduler_EnqueueAndWaitForJob_closedOnly(b *testing.B) {
	ss, err := state.NewStateStore()
	if err != nil {
		b.Fatal(err)
	}

	tmpDir := b.TempDir()
	ctx := context.Background()

	s := NewScheduler(&closedDirJobs{js: ss.JobStore}, 1)
	s.Start(ctx)
	b.Cleanup(func() {
		s.Stop()
	})

	ids := make(job.IDs, 0)
	for i := 0; i < b.N; i++ {
		i := i
		dirPath := filepath.Join(tmpDir, fmt.Sprintf("folder-%d", i))

		newId, err := ss.JobStore.EnqueueJob(job.Job{
			Func: func(c context.Context) error {
				return nil
			},
			Dir:  document.DirHandleFromPath(dirPath),
			Type: "test-type",
		})
		if err != nil {
			b.Fatal(err)
		}
		ids = append(ids, newId)
	}

	err = ss.JobStore.WaitForJobs(ctx, ids...)
	if err != nil {
		b.Fatal(err)
	}
}

func TestScheduler_defer(t *testing.T) {
	ss, err := state.NewStateStore()
	if err != nil {
		t.Fatal(err)
	}
	ss.SetLogger(testLogger())

	tmpDir := t.TempDir()

	ctx := context.Background()

	s := NewScheduler(&closedDirJobs{js: ss.JobStore}, 2)
	s.SetLogger(testLogger())
	s.Start(ctx)
	t.Cleanup(func() {
		s.Stop()
	})

	var jobsExecuted, deferredJobsExecuted int64 = 0, 0
	jobsToExecute := 50

	ids := make(job.IDs, 0)
	for i := 0; i < jobsToExecute; i++ {
		i := i
		dirPath := filepath.Join(tmpDir, fmt.Sprintf("folder-%d", i))

		newId, err := ss.JobStore.EnqueueJob(job.Job{
			Func: func(c context.Context) error {
				atomic.AddInt64(&jobsExecuted, 1)
				return nil
			},
			Dir:  document.DirHandleFromPath(dirPath),
			Type: "test-type",
			Defer: func(ctx context.Context, jobErr error) (ids job.IDs) {
				je, err := job.JobStoreFromContext(ctx)
				if err != nil {
					log.Fatal(err)
					return nil
				}

				id1, err := je.EnqueueJob(job.Job{
					Dir:  document.DirHandleFromPath(dirPath),
					Type: "test-1",
					Func: func(c context.Context) error {
						atomic.AddInt64(&deferredJobsExecuted, 1)
						return nil
					},
				})
				if err != nil {
					log.Fatal(err)
					return nil
				}
				ids = append(ids, id1)

				id2, err := je.EnqueueJob(job.Job{
					Dir:  document.DirHandleFromPath(dirPath),
					Type: "test-2",
					Func: func(c context.Context) error {
						atomic.AddInt64(&deferredJobsExecuted, 1)
						return nil
					},
				})
				if err != nil {
					log.Fatal(err)
					return nil
				}
				ids = append(ids, id2)

				return
			},
		})
		if err != nil {
			t.Fatal(err)
		}
		ids = append(ids, newId)
	}

	err = ss.JobStore.WaitForJobs(ctx, ids...)
	if err != nil {
		t.Fatal(err)
	}

	if jobsExecuted != int64(jobsToExecute) {
		t.Fatalf("expected %d jobs to execute, got: %d", jobsToExecute, jobsExecuted)
	}

	expectedDeferredJobs := int64(jobsToExecute * 2)
	if deferredJobsExecuted != expectedDeferredJobs {
		t.Fatalf("expected %d deferred jobs to execute, got: %d", expectedDeferredJobs, deferredJobsExecuted)
	}
}

func testLogger() *log.Logger {
	if testing.Verbose() {
		return log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile)
	}

	return log.New(ioutil.Discard, "", 0)
}

type closedDirJobs struct {
	js *state.JobStore
}

func (js *closedDirJobs) EnqueueJob(newJob job.Job) (job.ID, error) {
	return js.js.EnqueueJob(newJob)
}

func (js *closedDirJobs) AwaitNextJob(ctx context.Context) (job.ID, job.Job, error) {
	return js.js.AwaitNextJob(ctx, false)
}

func (js *closedDirJobs) FinishJob(id job.ID, jobErr error, deferredJobIds ...job.ID) error {
	return js.js.FinishJob(id, jobErr, deferredJobIds...)
}

func (js *closedDirJobs) WaitForJobs(ctx context.Context, jobIds ...job.ID) error {
	return js.js.WaitForJobs(ctx, jobIds...)
}

type openDirJobs struct {
	js *state.JobStore
}

func (js *openDirJobs) EnqueueJob(newJob job.Job) (job.ID, error) {
	return js.js.EnqueueJob(newJob)
}

func (js *openDirJobs) AwaitNextJob(ctx context.Context) (job.ID, job.Job, error) {
	return js.js.AwaitNextJob(ctx, true)
}

func (js *openDirJobs) FinishJob(id job.ID, jobErr error, deferredJobIds ...job.ID) error {
	return js.js.FinishJob(id, jobErr, deferredJobIds...)
}

func (js *openDirJobs) WaitForJobs(ctx context.Context, jobIds ...job.ID) error {
	return js.js.WaitForJobs(ctx, jobIds...)
}
