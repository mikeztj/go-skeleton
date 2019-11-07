package skeleton

import (
	"context"
	"fmt"
	"math"
	"strconv"
	"sync"

	skeletonEntity "github.com/vilbert/go-skeleton/internal/entity/skeleton"
)

// Data ...
type Data interface {
	GetAllSkeletons(ctx context.Context) ([]skeletonEntity.Skeleton, error)
	GetAllSkeletonsWithPage(ctx context.Context, offset int, length int) ([]skeletonEntity.Skeleton, error)
	GetSkeletonsCount(ctx context.Context) (int, error)
}

// Service ...
type Service struct {
	data Data
}

// New ...
func New(skeletonData Data) Service {
	return Service{
		data: skeletonData,
	}
}

// GetAllSkeletons ...
func (s Service) GetAllSkeletons(ctx context.Context, page int, length int) ([]skeletonEntity.Skeleton, interface{}, error) {
	var (
		wg       sync.WaitGroup
		metadata map[string]string
		offset   int
		result   []skeletonEntity.Skeleton
		count    int
		err      error
	)

	if length != 0 {
		offset = (page - 1) * length
		// Add 2 goroutines to waitgroup queue
		wg.Add(2)

		// Use goroutine for paralel processing
		go func() {
			defer wg.Done()
			result, err = s.data.GetAllSkeletonsWithPage(ctx, offset, length)
		}()

		// Use goroutine for paralel processing
		go func() {
			defer wg.Done()
			count, err = s.data.GetSkeletonsCount(ctx)
		}()
	} else {
		result, err = s.data.GetAllSkeletons(ctx)
	}

	// Wait for goroutines to finish processing
	wg.Wait()

	if page == 0 && length == 0 {
		metadata = map[string]string{
			"current_page": "1",
			"max_page":     "1",
		}
	} else {
		metadata = map[string]string{
			"current_page": strconv.Itoa(page),
			"max_page":     fmt.Sprintf("%.0f", math.Ceil(float64(count)/float64(length))),
		}
	}

	return result, metadata, err
}
