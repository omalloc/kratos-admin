package task

import (
	"context"
	"sync"
)

var (
	lock        sync.Mutex
	registryMap = make(map[string]Task)
)

type Task interface {
	Name() string
	Do(ctx context.Context, args []byte) error
}

func Register(task Task) {
	lock.Lock()
	defer lock.Unlock()

	registryMap[task.Name()] = task
}

func Deregister(task Task) {
	lock.Lock()
	defer lock.Unlock()

	delete(registryMap, task.Name())
}
