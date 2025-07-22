# go-fatty-context
Implementing an efficient and concurrency safe way to read write and read from context - resolving fatty contexts in Go.

## Test
1. Run multiple tests to compare speedup at 100, 1_000, and 1_000_000 context values (why not 1_000_000? is there a reason on earth to keep those many vars in context - none perhaps <).
2. Test for any consistency issues/race-around with multiple threads reading/writing on a single context.

thanks