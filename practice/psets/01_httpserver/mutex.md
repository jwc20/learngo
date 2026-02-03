# Mutexes

## Refactor

We are almost there! Lets take some effort to prevent concurrency errors like these

```bash
fatal error: concurrent map read and map write
```

By adding mutexes, we enforce concurrency safety especially for the counter in our RecordWin function.

```go
package main

import "sync"

type InMemoryPlayerStore struct {
 mu    sync.RWMutex
 store map[string]int
}

func NewInMemoryPlayerStore() *InMemoryPlayerStore {
 return &InMemoryPlayerStore{store: map[string]int{}}
}

func (i *InMemoryPlayerStore) RecordWin(name string) {
 i.mu.Lock()
 defer i.mu.Unlock()
 i.store[name]++
}

func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
 i.mu.RLock()
 defer i.mu.RUnlock()
 return i.store[name]
}
```

### Why RWMutex Instead of Mutex?

Using RWMutex is more performant when you have many reads and few writes. Multiple goroutines can read scores simultaneously, but writes get exclusive access. A regular Mutex would block all access even for concurrent reads.
