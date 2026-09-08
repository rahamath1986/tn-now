package cron

import (
	"context"
	"testing"
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
