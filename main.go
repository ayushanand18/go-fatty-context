package main

import "context"

func getContextImplementationsAsCtx() []context.Context {
	implementations := []context.Context{
		context.Background(),
		ThinBackgroundContext(),
	}

	return implementations
}

func runBasicTestsOnCtx(ctx context.Context) {

}

func main() {
	for _, ctx := range getContextImplementationsAsCtx() {
		runBasicTestsOnCtx(ctx)
	}
}
