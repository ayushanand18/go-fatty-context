package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type TestContext struct {
	Name    string
	Context context.Context
}

// get the size of payloads we are going to test
func getContextInputPayloadSize() []int64 {
	return []int64{100, 1000, 1000000}
}

// get all the contexts to run tests on
func getContextImplementationsAsCtx() []TestContext {
	implementations := []TestContext{
		{Context: context.Background(), Name: "go/context"},
		{Context: ThinBackgroundContext(), Name: "go-fatty-context"},
	}

	return implementations
}

func runBasicTestsOnCtx(tctx TestContext) {
	fmt.Printf("Running context benchmarks for: %s\n", tctx.Name)

	for _, size := range getContextInputPayloadSize() {
		benchmarkContextLookup(tctx.Context, size)
	}

	testConcurrentContextAccess(1000)

	fmt.Printf("=======================================\n")
}

// multiple consecutive context lookups on single thread
func benchmarkContextLookup(ctx context.Context, n int64) {
	fmt.Printf("Benchmarking context with %d values...\n", n)

	var firstKey string

	for i := int64(0); i < n; i++ {
		key := fmt.Sprintf("key%d", i)
		if i == 0 {
			firstKey = key
		}
		ctx = context.WithValue(ctx, key, fmt.Sprintf("value%d", i))
	}

	start := time.Now()
	_ = ctx.Value(firstKey)
	elapsed := time.Since(start)

	fmt.Printf("Lookup time for first key: %s\n\n", elapsed)
}

// Concurrent reads after consecutive single-threaded writes
func testConcurrentContextAccess(keyCount int) {
	fmt.Println("Running concurrency test with multiple readers...")

	ctx := context.Background()

	for i := 0; i < keyCount; i++ {
		ctx = context.WithValue(ctx, fmt.Sprintf("key%d", i), fmt.Sprintf("val%d", i))
	}

	var wg sync.WaitGroup
	numReaders := 100
	lookupsPerReader := 1000

	for i := 0; i < numReaders; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < lookupsPerReader; j++ {
				key := fmt.Sprintf("key%d", j%keyCount)
				val := ctx.Value(key)
				if val == nil {
					fmt.Printf("Goroutine %d: missing value for key %s\n", id, key)
				}
			}
		}(i)
	}

	wg.Wait()
	fmt.Printf("All goroutines finished.\n")
}

func main() {
	for _, ctx := range getContextImplementationsAsCtx() {
		runBasicTestsOnCtx(ctx)
	}
}
