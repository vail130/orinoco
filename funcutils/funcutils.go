package funcutils

import (
	"sync"
)

type AnonymousFunction func()

func ExecuteWithWaitGroup(wg *sync.WaitGroup, f AnonymousFunction) {
	defer wg.Done()
	f()
}
