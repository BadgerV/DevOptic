package monitor

import (
    "context"
    "log"
    "os"
    "strconv"
    "time"
)

type Scheduler struct {
    cancel context.CancelFunc
    done   chan struct{}
}

func (s *Service) StartScheduler(ctx context.Context) (*Scheduler, error) {
    timer, err := strconv.Atoi(os.Getenv("CHECK_TIMER"))
    if err != nil {
        return nil, err
    }

    endpoints, err := s.LoadAndSyncEndpoints("endpoints.json")
    if err != nil {
        return nil, err
    }

    ctx, cancel := context.WithCancel(ctx)
    ticker := time.NewTicker(time.Duration(timer) * time.Second)
    done := make(chan struct{})

    go func() {
        defer close(done)
        defer ticker.Stop()
        defer log.Println("Scheduler stopped")

        for {
            select {
            case <-ctx.Done():
                return
            case <-ticker.C:
                for _, ep := range endpoints {
                    select {
                    case <-ctx.Done():
                        return
                    default:
                        go func(e Endpoint) {
                            reqCtx, cancelReq := context.WithTimeout(ctx, 5*time.Second)
                            defer cancelReq()
                            if err := s.CheckEndpoint(reqCtx, e.ID, e.URL, e.APIMethod, e.ServerName, e.ServiceName); err != nil {
                                log.Println("Error Checking", e.URL, ":", err)
                            }
                        }(ep)
                    }
                }
            }
        }
    }()

    return &Scheduler{cancel: cancel, done: done}, nil
}

func (sch *Scheduler) Stop() {
    if sch.cancel != nil {
        sch.cancel()
        <-sch.done
    }
}
