package interceptor

import "github.com/go-cinderella/cinderella-engine/engine"

type AbstractCommandInterceptor struct {
	Next engine.Interceptor
}
