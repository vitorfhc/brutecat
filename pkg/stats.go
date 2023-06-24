package brutecat

import (
	"fmt"
	"sync"
	"time"
)

type RunStats struct {
	totalCredentials   uint64
	validCredentials   uint64
	invalidCredentials uint64
	errors             uint64
	start              *time.Time
	end                *time.Time

	mut *sync.Mutex
}

func NewStats(totalCredentials uint64) *RunStats {
	mut := &sync.Mutex{}

	return &RunStats{
		totalCredentials:   totalCredentials,
		validCredentials:   0,
		invalidCredentials: 0,
		errors:             0,
		start:              nil,
		end:                nil,
		mut:                mut,
	}
}

func (s *RunStats) Start() {
	s.mut.Lock()
	defer s.mut.Unlock()

	now := time.Now()
	s.start = &now
}

func (s *RunStats) End() {
	s.mut.Lock()
	defer s.mut.Unlock()

	now := time.Now()
	s.end = &now
}

func (s *RunStats) GetStart() *time.Time {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.start
}

func (s *RunStats) GetEnd() *time.Time {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.end
}

func (s *RunStats) AddValidCredentials() {
	s.mut.Lock()
	defer s.mut.Unlock()

	s.validCredentials++
}

func (s *RunStats) AddInvalidCredentials() {
	s.mut.Lock()
	defer s.mut.Unlock()

	s.invalidCredentials++
}

func (s *RunStats) AddError() {
	s.mut.Lock()
	defer s.mut.Unlock()

	s.errors++
}

func (s *RunStats) GetTotalCredentials() uint64 {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.totalCredentials
}

func (s *RunStats) GetValidCredentials() uint64 {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.validCredentials
}

func (s *RunStats) GetInvalidCredentials() uint64 {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.invalidCredentials
}

func (s *RunStats) GetErrors() uint64 {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.errors
}

func (s *RunStats) GetElapsedTime() time.Duration {
	if s.GetStart() == nil {
		return 0
	}

	end := s.GetEnd()

	var until time.Time
	if end == nil {
		until = time.Now()
	} else {
		until = *end
	}

	return until.Sub(*s.GetStart())
}

func (s *RunStats) GetAverageCredentialsPerSecond() float64 {
	start := s.GetStart()
	if start == nil {
		return 0
	}

	elapsedTime := s.GetElapsedTime()
	return float64(s.GetTriedCredentials()) / elapsedTime.Seconds()
}

func (s *RunStats) GetTriedCredentials() uint64 {
	return s.GetValidCredentials() + s.GetInvalidCredentials()
}

func (s *RunStats) GetRemainingTime() time.Duration {
	start := s.GetStart()
	if start == nil {
		return 0
	}

	averageCredentialsPerSecond := float64(s.GetAverageCredentialsPerSecond())
	remainingCredentials := float64(s.GetTotalCredentials() - s.GetTriedCredentials())

	if averageCredentialsPerSecond == 0 {
		return 0
	}

	return time.Duration(remainingCredentials/averageCredentialsPerSecond) * time.Second
}

func (s *RunStats) String() string {
	tried := s.GetTriedCredentials()
	total := s.GetTotalCredentials()
	valid := s.GetValidCredentials()
	elapsedTime := s.GetElapsedTime()
	averageCredentialsPerSecond := s.GetAverageCredentialsPerSecond()
	remainingTime := s.GetRemainingTime()

	elapsedTime = elapsedTime.Round(time.Second)

	msg := "Tried %d/%d in %s (%.1f/s) - Valid: %d - Remaining time: %s"
	fmtted := fmt.Sprintf(msg, tried, total, elapsedTime, averageCredentialsPerSecond, valid, remainingTime)

	return fmtted
}
