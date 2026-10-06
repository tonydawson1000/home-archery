package application

import (
	"fmt"
	"sync"
	"time"
)

type FixedClock struct {
	T time.Time
}

func (c FixedClock) Now() time.Time { return c.T }

type SeqIDs struct {
	mu     sync.Mutex
	prefix string
	n      int
}

func NewSeqIDs(prefix string) *SeqIDs {
	return &SeqIDs{prefix: prefix}
}

func (s *SeqIDs) New() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.n++
	return fmt.Sprintf("%s-%d", s.prefix, s.n)
}
