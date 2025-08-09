# go-fatty-context
Implementing an efficient and concurrency safe way to read write and read from context - resolving fatty contexts in Go.

## Information
this is how a value context in golang stored
```go
type valueCtx struct {
	Context
	key, val any
}
```

since it stores context in a recursive fashion, with a context carrying a child-context with n-1 keys, it grows linearly with each value added. 
writing a new value is O(1) op, since it would be simply creating a new ctx, and wrapping the existing context inside the new. but lookups are costly and often linearly growing with input size (O(n)).

## Test
1. Run multiple tests to compare speedup at 100, 1_000, and 1_000_000 context values (why not 1_000_000? is there a reason on earth to keep those many vars in context - none perhaps <).
2. Test for any consistency issues/race-around with multiple threads reading/writing on a single context.

## Results

Lib               | 100     | 1000   | 1_000_000
------------------|---------|--------|------------
`go/context`      |17.833µs |5.917µs | 5.843375ms
`go-fatty-context`|416ns    |2.333µs | 4.226875ms

