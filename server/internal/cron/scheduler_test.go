package cron

import (
	"context"
	"testing"
	"time"
)

func TestSchedulerOperations(t *testing.T) {
	scheduler := NewScheduler(nil)

	// Verify default registered jobs
	jobs := scheduler.GetJobs(context.Background())
	if len(jobs) < 4 {
		t.Errorf("Expected at least 4 default cron jobs, got %d", len(jobs))
	}

	// Test triggering a job
	res, err := scheduler.TriggerJob(context.Background(), "dead_link_checker")
	if err != nil {
		t.Fatalf("Failed to trigger dead_link_checker: %v", err)
	}
	if res.Status != "SUCCESS" {
		t.Errorf("Expected status SUCCESS, got %s", res.Status)
	}

	// Test toggling active state
	active, err := scheduler.ToggleJob(context.Background(), "dead_link_checker")
	if err != nil {
		t.Fatalf("Failed to toggle job: %v", err)
	}
	if active != false {
		t.Errorf("Expected job to be paused (active=false), got %v", active)
	}

	// Test creating a new job
	newJob, err := scheduler.CreateJob(context.Background(), "custom_test_task", "Test Task", "Custom test task", "5m", "CUSTOM")
	if err != nil {
		t.Fatalf("Failed to create job: %v", err)
	}
	if newJob.ID != "custom_test_task" {
		t.Errorf("Expected job ID custom_test_task, got %s", newJob.ID)
	}

	// Verify created job exists in job list
	found := false
	for _, j := range scheduler.GetJobs(context.Background()) {
		if j.ID == "custom_test_task" {
			found = true
			break
		}
	}
	if !found {
		t.Error("Newly created job not found in scheduler job list")
	}
}

func TestCheckAndRunJobsNoDeadlock(t *testing.T) {
	scheduler := NewScheduler(nil)

	ran := false
	past := time.Now().Add(-1 * time.Minute)
	testJob := &CronJob{
		ID:               "test_deadlock_job",
		Name:             "Test Deadlock Job",
		ScheduleInterval: "1m",
		IsActive:         true,
		NextRunAt:        &past,
		Handler: func(ctx context.Context) (string, error) {
			ran = true
			return "executed successfully", nil
		},
	}
	scheduler.RegisterJob(testJob)

	// Ensure NextRunAt is in the past
	testJob.NextRunAt = &past

	// checkAndRunJobs must not prematurely lock IsRunning and block TriggerJob
	scheduler.checkAndRunJobs()

	// Wait briefly for the async goroutine to execute
	time.Sleep(100 * time.Millisecond)

	if !ran {
		t.Error("expected job handler to execute when checkAndRunJobs was invoked")
	}

	// Verify job is NOT permanently stuck in IsRunning
	jobs := scheduler.GetJobs(context.Background())
	for _, j := range jobs {
		if j.ID == "test_deadlock_job" && j.IsRunning {
			t.Error("expected test_deadlock_job to have IsRunning=false after execution")
		}
	}
}

func TestTNLiveNewsCronConfig(t *testing.T) {
	scheduler := NewScheduler(nil)
	jobs := scheduler.GetJobs(context.Background())

	var liveJob *CronJob
	for _, j := range jobs {
		if j.ID == "tn_live_news_cron" {
			liveJob = j
			break
		}
	}

	if liveJob == nil {
		t.Fatal("expected tn_live_news_cron to be registered")
	}
	if !liveJob.IsActive {
		t.Error("expected tn_live_news_cron to be active")
	}
	if len(liveJob.SourceURLs) < 4 {
		t.Errorf("expected at least 4 source URLs for tn_live_news_cron, got %d", len(liveJob.SourceURLs))
	}
}
