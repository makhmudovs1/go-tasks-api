package logging

import (
	"context"
	"log"
	"sync"
	"time"
)

type LogEvent struct {
	Time    time.Time
	Action  string
	Details string
}

func StartLogger(ctx context.Context, wg *sync.WaitGroup, events <-chan LogEvent) {
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				return
			case e, ok := <-events:
				if !ok {
					return
				}
				log.Printf("[%s] %s — %s\n", e.Time.Format(time.RFC3339), e.Action, e.Details)
			}
		}
	}()
}
