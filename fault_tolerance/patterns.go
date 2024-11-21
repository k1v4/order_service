package fault_tolerance

import (
	"fmt"
	"math"
	"time"
)

func Retry(operation func() error, maxRetries int, baseDelay int) error {
	n := 1
	var err error

	for n <= maxRetries {
		err = operation()
		if err == nil {
			return err
		} else {
			pause := baseDelay * int(math.Pow(2, float64(n)))
			time.Sleep(time.Duration(pause) * time.Millisecond)
			n++
		}
	}

	return err
}

func Timeout(operation func() error, timeout int) error {
	done := make(chan error)

	go func() { done <- operation() }()

	select {
	case err := <-done:
		return err
	case <-time.After(time.Duration(timeout) * time.Millisecond):
		return fmt.Errorf("timeout after %d milliseconds", timeout)
	}
}

func ProcessWithDLQ(messages []string, operation func(msg string) error) []string {
	var q []string

	for _, msg := range messages {
		if err := operation(msg); err != nil {
			q = append(q, msg)
		}
	}

	return q
}
