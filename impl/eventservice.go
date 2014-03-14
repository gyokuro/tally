package impl

import (
	_ "errors"
	"github.com/golang/glog"
	_ "github.com/gyokuro/tally"
	"github.com/gyokuro/tally/proto"
)

type mockEventService struct {
}

func (s *mockEventService) Put(events []Tally.Event) error {
	glog.Info("Received events", events)
	return nil
}

func NewMockEventService() *mockEventService {
	return &mockEventService{}
}
