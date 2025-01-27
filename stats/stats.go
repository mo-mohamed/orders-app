package stats

import (
	"context"
	"math/rand"
	"time"

	"github.com/orders-app/logger"
	"github.com/orders-app/models"
)

const WorkerCount = 3

type statsService struct {
	result    Result
	processed <-chan models.Order
	done      <-chan struct{}
	pStats    chan models.Statistics
}

type StatsService interface {
	GetStats(ctx context.Context) <-chan models.Statistics
}

func New(processed <-chan models.Order, done <-chan struct{}) StatsService {
	s := statsService{
		result:    &result{},
		processed: processed,
		done:      done,
		pStats:    make(chan models.Statistics, WorkerCount),
	}

	for i := 0; i < WorkerCount; i++ {
		go s.processStats()
	}

	go s.reconcile()
	return &s
}

// processStats is the overall processing method that listens to incoming orders
func (s *statsService) processStats() {
	logger.Log.Info("Stats processing started!")
	for {
		select {
		case order := <-s.processed:
			pstats := s.processOrder(order)
			s.pStats <- pstats
		case <-s.done:
			logger.Log.Warn("Stats processing stopped!")
			return
		}
	}
}

// reconcile is a helper method which saves stats object
// back into the statisticsService
func (s *statsService) reconcile() {
	logger.Log.Info("Reconcile started!")
	for {
		select {
		case p := <-s.pStats:
			s.result.Combine(p)
		case <-s.done:
			logger.Log.Warn("Reconcile stopped!")
			return
		}
	}
}

// processOrder is a helper method that incorporates the current order in the stats service
func (s *statsService) processOrder(order models.Order) models.Statistics {
	// simulate processing as a costly operation
	randomSleep()
	// completed orders increment add to the revenue
	if order.Status == string(models.OrderStatus_Completed) {
		return models.Statistics{
			CompletedOrders: 1,
			Revenue:         order.Total,
		}
	}

	// reversed orders remove from the revenue
	if order.Status == string(models.OrderStatus_Reversed) {
		return models.Statistics{
			ReversedOrders: 1,
			Revenue:        -order.Total,
		}
	}
	// otherwise the order is rejected
	return models.Statistics{
		RejectedOrders: 1,
	}
}

// GetStats returns the latest order stats
func (s *statsService) GetStats(ctx context.Context) <-chan models.Statistics {
	stats := make(chan models.Statistics)
	go func() {
		randomSleep()
		select {
		case stats <- s.result.Get():
			logger.Log.Info("Stats fetched successfully")
			return
		case <-ctx.Done():
			logger.Log.Info("Context deadline exceeded")
			return
		}
	}()
	return stats
}

func randomSleep() {
	rand.Seed(time.Now().UnixNano())
	time.Sleep(time.Duration(rand.Intn(800)) * time.Millisecond)
}
