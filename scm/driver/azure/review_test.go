package azure

import (
	"context"
	"testing"

	"github.com/jenkins-x/go-scm/scm"
)

func TestReviewFind(t *testing.T) {
	reviewService := &reviewService{&wrapper{
		Project: "test-project",
	}}
	_, _, err := reviewService.Find(context.Background(), "", 0, 0)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestReviewList(t *testing.T) {
	reviewService := &reviewService{&wrapper{
		Project: "test-project",
	}}
	_, _, err := reviewService.List(context.Background(), "", 0, scm.ListOptions{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestReviewCreate(t *testing.T) {
	reviewService := &reviewService{&wrapper{
		Project: "test-project",
	}}
	_, _, err := reviewService.Create(context.Background(), "", 0, &scm.ReviewInput{})
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}

func TestReviewDelete(t *testing.T) {
	reviewService := &reviewService{&wrapper{
		Project: "test-project",
	}}
	_, err := reviewService.Delete(context.Background(), "", 0, 0)
	if err != scm.ErrNotSupported {
		t.Errorf("Expect Not Supported error")
	}
}
