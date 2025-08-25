package services

import (
	"athena/infrastructure/api"
	"context"
	"fmt"
	"log"
	"sync"
	"time"
)

type IScheduleService interface {
	StartApplication(ctx context.Context, minutes int)
	StopApplication()
}

type ScheduleService struct {
	trackerService ITrackerService
	api            api.ITrackerAPIClient
	cancelFunc     context.CancelFunc
	wg             sync.WaitGroup
	mu             sync.Mutex
	isRunning      bool
}

var (
	timerInstance *ScheduleService
	timerOnce     sync.Once
)

func (t *ScheduleService) StartApplication(ctx context.Context, minutes int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.isRunning {
		log.Println("Schedule já está em execução")
		return
	}

	ctx, t.cancelFunc = context.WithCancel(ctx)
	t.isRunning = true

	t.wg.Add(1)
	go func() {
		defer t.wg.Done()
		t.run(ctx, minutes)
	}()
}

func (t *ScheduleService) run(ctx context.Context, minutes int) {
	ticker := time.NewTicker(time.Duration(minutes) * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			t.executeTask(ctx)
		case <-ctx.Done():
			log.Println("Schedule parando...")
			return
		}
	}
}

func (t *ScheduleService) executeTask(_ context.Context) {
	resp, err := t.api.ListaVeiculos()
	if err != nil {
		log.Printf("Erro ao receber resposta da API: %v", err)
		return
	}

	if resp != nil {
		fmt.Printf("Resposta recebida: %d veículos\n", len(resp.Dados))
	}

	if err := t.trackerService.SaveTrackerData(); err != nil {
		log.Printf("Erro ao salvar dados: %v", err)
	}
}

func (t *ScheduleService) StopApplication() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.cancelFunc != nil {
		t.cancelFunc()
	}

	t.wg.Wait()
	t.isRunning = false
	log.Println("TimerService parado")
}

func NewTimerService(
	trackerService ITrackerService,
	api api.ITrackerAPIClient,
) IScheduleService {
	timerOnce.Do(func() {
		timerInstance = &ScheduleService{
			trackerService: trackerService,
			api:            api,
			isRunning:      false,
		}
	})
	return timerInstance
}
