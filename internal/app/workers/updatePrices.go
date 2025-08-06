package worker

import (
	"context"
	coingecko "cryptoObserver/internal/app/coingeko"
	"cryptoObserver/internal/app/store/sqlstore"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

type WorkerPool struct {
	client      coingecko.CryptoInterface
	db          sqlstore.StoreInterface
	mu          sync.RWMutex
	currencies  map[string]struct{}
	workers     int
	interval    time.Duration
	ctx         context.Context
	cancel      context.CancelFunc
	wg          sync.WaitGroup
	log         *logrus.Logger
	taskChan    chan string     // Канал для распределения задач
	activeTasks map[string]bool // Трекер активных задач
	taskMu      sync.Mutex      // Защита activeTasks
}

func NewWorkerPool(
	ctx context.Context,
	client coingecko.CryptoInterface,
	db sqlstore.StoreInterface,
	workers int,
	interval time.Duration,
	log *logrus.Logger,
) *WorkerPool {
	poolCtx, cancel := context.WithCancel(ctx)

	return &WorkerPool{
		client:      client,
		db:          db,
		currencies:  make(map[string]struct{}),
		workers:     workers,
		interval:    interval,
		ctx:         poolCtx,
		cancel:      cancel,
		log:         log,
		taskChan:    make(chan string, 100), // Буферизованный канал
		activeTasks: make(map[string]bool),
	}
}

// Добавляем валюту в список отслеживания
func (wp *WorkerPool) AddCurrency(currencyID string) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	if _, exists := wp.currencies[currencyID]; !exists {
		wp.currencies[currencyID] = struct{}{}
		wp.log.Infof("Currency added: %s", currencyID)
	}
}

// Удаляем валюту из списка отслеживания
func (wp *WorkerPool) RemoveCurrency(currencyID string) {
	wp.mu.Lock()
	defer wp.mu.Unlock()

	delete(wp.currencies, currencyID)
	wp.log.Infof("Currency removed: %s", currencyID)
}

// Запускаем воркер-пул
func (wp *WorkerPool) Start() {
	// Загружаем список валют из БД
	currencyList, err := wp.db.Currency().GetCurrencyList()
	if err != nil {
		wp.log.Errorf("Failed to get currency list from DB: %v", err)
	} else if len(currencyList) > 0 {
		wp.mu.Lock()
		for _, id := range currencyList {
			wp.currencies[id] = struct{}{}
		}
		wp.mu.Unlock()
		wp.log.Infof("Loaded %d currencies from database", len(currencyList))
	}

	// Запускаем распределитель задач
	go wp.taskDispatcher()

	// Запускаем воркеров
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.runWorker()
	}
}

// Распределитель задач
func (wp *WorkerPool) taskDispatcher() {
	ticker := time.NewTicker(wp.interval)
	defer ticker.Stop()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case <-ticker.C:
			wp.dispatchTasks()
		}
	}
}

// Распределяем задачи по воркерам
func (wp *WorkerPool) dispatchTasks() {
	wp.mu.RLock()
	defer wp.mu.RUnlock()

	wp.taskMu.Lock()
	defer wp.taskMu.Unlock()

	for currencyID := range wp.currencies {
		if !wp.activeTasks[currencyID] {
			wp.activeTasks[currencyID] = true
			wp.taskChan <- currencyID
		}
	}
}

// Воркер
func (wp *WorkerPool) runWorker() {
	defer wp.wg.Done()

	for {
		select {
		case <-wp.ctx.Done():
			return
		case currencyID := <-wp.taskChan:
			wp.processCurrency(currencyID)
		}
	}
}

// Обработка одной валюты
func (wp *WorkerPool) processCurrency(currencyID string) {
	defer func() {
		wp.taskMu.Lock()
		delete(wp.activeTasks, currencyID)
		wp.taskMu.Unlock()
	}()

	price, err := wp.client.GetCryptoPrice(wp.ctx, currencyID)
	if err != nil {
		wp.log.Errorf("Failed to fetch %s: %v", currencyID, err)
		return
	}

	if err := wp.db.Currency().UpdatePrice(currencyID, price.CurrentPrice, time.Now().Unix()); err != nil {
		wp.log.Errorf("Failed to save %s: %v", currencyID, err)
	}
}

// Остановка воркер-пула
func (wp *WorkerPool) Stop() {
	wp.cancel()
	wp.wg.Wait()
	close(wp.taskChan)
	wp.log.Info("Worker pool stopped")
}
