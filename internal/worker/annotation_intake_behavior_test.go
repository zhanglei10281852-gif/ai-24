package worker

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/zhanglei10281852-gif/ai/internal/domain"
	"github.com/zhanglei10281852-gif/ai/internal/repository"
)

func TestRunOnceCompletesPlannedRunPayload(t *testing.T) {
	worker, store, ctx, now := workerFixture(t)
	payload, err := json.Marshal(domain.InferenceRun{ID: "run_planned_payload", Reference: "PLANNED-PAYLOAD"})
	if err != nil {
		t.Fatal(err)
	}
	job := domain.OutboxJob{
		ID: "job_planned_payload", Kind: "inference_run_planned", AggregateID: "run_planned_payload",
		Payload: payload, Status: domain.JobPending, MaxAttempts: 5,
		AvailableAt: now, CreatedAt: now, UpdatedAt: now,
	}
	if err := store.WithTx(ctx, func(tx repository.Tx) error { return tx.InsertJob(ctx, job) }); err != nil {
		t.Fatal(err)
	}
	if err := worker.RunOnce(ctx); err != nil {
		t.Fatal(err)
	}
	if err := store.WithTx(ctx, func(tx repository.Tx) error {
		jobs, err := tx.ClaimJobs(ctx, now.Add(10*time.Second), 10)
		if err != nil {
			return err
		}
		if len(jobs) != 0 {
			t.Fatalf("completed planned event was reclaimed: %+v", jobs)
		}
		return nil
	}); err != nil {
		t.Fatal(err)
	}
}
