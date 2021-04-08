package proxyservice

import (
	"context"
	"sync"
	"testing"

	"go.uber.org/zap"

	"github.com/stretchr/testify/assert"
	"github.com/zilliztech/milvus-distributed/internal/log"
)

func TestBaseTaskQueue_Enqueue(t *testing.T) {
	queue := newBaseTaskQueue()

	num := 10
	var wg sync.WaitGroup

	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			tsk := newMockTask(context.Background())
			err := queue.Enqueue(tsk)
			assert.Equal(t, nil, err)
		}()
	}

	wg.Wait()
}

func TestBaseTaskQueue_FrontTask(t *testing.T) {
	queue := newBaseTaskQueue()

	tsk := queue.FrontTask()
	assert.Equal(t, nil, tsk)

	frontTask := newMockTask(context.Background())
	err := queue.Enqueue(frontTask)
	assert.Equal(t, nil, err)
	tsk = queue.FrontTask()
	assert.NotEqual(t, nil, tsk)
	assert.Equal(t, frontTask.ID(), tsk.ID())
	assert.Equal(t, frontTask.Name(), tsk.Name())

	num := 10

	for i := 0; i < num; i++ {
		tsk := newMockTask(context.Background())
		err := queue.Enqueue(tsk)
		assert.Equal(t, nil, err)

		tskF := queue.FrontTask()
		assert.NotEqual(t, nil, tskF)
		assert.Equal(t, frontTask.ID(), tskF.ID())
		assert.Equal(t, frontTask.Name(), tskF.Name())
	}
}

func TestBaseTaskQueue_PopTask(t *testing.T) {
	queue := newBaseTaskQueue()

	tsk := queue.PopTask()
	assert.Equal(t, nil, tsk)

	num := 10

	for i := 0; i < num; i++ {
		tsk := newMockTask(context.Background())
		err := queue.Enqueue(tsk)
		assert.Equal(t, nil, err)

		tskP := queue.PopTask()
		assert.NotEqual(t, nil, tskP)
	}

	tsk = queue.PopTask()
	assert.Equal(t, nil, tsk)
}

func TestBaseTaskQueue_Chan(t *testing.T) {
	queue := newBaseTaskQueue()

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		for {
			select {
			case <-ctx.Done():
				log.Debug("TestBaseTaskQueue_Chan exit")
				return
			case i := <-queue.Chan():
				log.Debug("TestBaseTaskQueue_Chan", zap.Any("receive", i))
			}
		}
	}()

	num := 10
	var wg sync.WaitGroup
	for i := 0; i < num; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			tsk := newMockTask(context.Background())
			err := queue.Enqueue(tsk)
			assert.Equal(t, nil, err)
		}()
	}

	wg.Wait()

	cancel()
}

func TestBaseTaskQueue_Empty(t *testing.T) {
	queue := newBaseTaskQueue()
	assert.Equal(t, true, queue.Empty())

	num := 10
	for i := 0; i < num; i++ {
		tsk := newMockTask(context.Background())
		err := queue.Enqueue(tsk)
		assert.Equal(t, nil, err)

		assert.Equal(t, false, queue.Empty())
	}

	for !queue.Empty() {
		assert.Equal(t, false, queue.Empty())
		queue.PopTask()
	}

	assert.Equal(t, true, queue.Empty())
}

func TestBaseTaskQueue_Full(t *testing.T) {
	queue := newBaseTaskQueue()

	for !queue.Full() {
		assert.Equal(t, false, queue.Full())

		tsk := newMockTask(context.Background())
		err := queue.Enqueue(tsk)
		assert.Equal(t, nil, err)
	}

	assert.Equal(t, true, queue.Full())
}