package circuitbreaker

import (
	"context"
	"net/http"
	"time"

	"github.com/vulcand/oxy/v2/cbreaker"
	ptypes "github.com/traefik/paerser/types"
)

const typeName = "CircuitBreaker"

type CircuitBreakerConfig struct {
	// Expression defines the expression that, once matched, opens the circuit breaker and applies the fallback mechanism instead of calling the services.
	Expression string `json:"expression,omitempty" toml:"expression,omitempty" yaml:"expression,omitempty" export:"true"`
	// CheckPeriod is the interval between successive checks of the circuit breaker condition (when in standby state).
	CheckPeriod ptypes.Duration `json:"checkPeriod,omitempty" toml:"checkPeriod,omitempty" yaml:"checkPeriod,omitempty" export:"true"`
	// FallbackDuration is the duration for which the circuit breaker will wait before trying to recover (from a tripped state).
	FallbackDuration ptypes.Duration `json:"fallbackDuration,omitempty" toml:"fallbackDuration,omitempty" yaml:"fallbackDuration,omitempty" export:"true"`
	// RecoveryDuration is the duration for which the circuit breaker will try to recover (as soon as it is in recovering state).
	RecoveryDuration ptypes.Duration `json:"recoveryDuration,omitempty" toml:"recoveryDuration,omitempty" yaml:"recoveryDuration,omitempty" export:"true"`
	// ResponseCode is the code that the circuit breaker will return while it is in the tripped state.
	ResponseCode int `json:"responseCode,omitempty" toml:"responseCode,omitempty" yaml:"responseCode,omitempty" export:"true"`
}

func (c *CircuitBreakerConfig) SetDefaults() {
	c.CheckPeriod = ptypes.Duration(100 * time.Millisecond)
	c.FallbackDuration = ptypes.Duration(10 * time.Second)
	c.RecoveryDuration = ptypes.Duration(10 * time.Second)
	c.ResponseCode = 503
}

type circuitBreaker struct {
	circuitBreaker *cbreaker.CircuitBreaker
	name           string
}

// New creates a new circuit breaker middleware.
func New(ctx context.Context, next http.Handler, confCircuitBreaker CircuitBreakerConfig, name string) (http.Handler, error) {
	expression := confCircuitBreaker.Expression

	responseCode := confCircuitBreaker.ResponseCode

	cbOpts := []cbreaker.Option{
		cbreaker.Fallback(http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
			rw.WriteHeader(responseCode)
			rw.Write([]byte(http.StatusText(responseCode)))
		})),
	}

	if confCircuitBreaker.CheckPeriod > 0 {
		cbOpts = append(cbOpts, cbreaker.CheckPeriod(time.Duration(confCircuitBreaker.CheckPeriod)))
	}

	if confCircuitBreaker.FallbackDuration > 0 {
		cbOpts = append(cbOpts, cbreaker.FallbackDuration(time.Duration(confCircuitBreaker.FallbackDuration)))
	}

	if confCircuitBreaker.RecoveryDuration > 0 {
		cbOpts = append(cbOpts, cbreaker.RecoveryDuration(time.Duration(confCircuitBreaker.RecoveryDuration)))
	}

	oxyCircuitBreaker, err := cbreaker.New(next, expression, cbOpts...)
	if err != nil {
		return nil, err
	}

	return &circuitBreaker{
		circuitBreaker: oxyCircuitBreaker,
		name:           name,
	}, nil
}

// func (c *circuitBreaker) GetTracingInformation() (string, ext.SpanKindEnum) {
// 	return c.name, tracing.SpanKindNoneEnum
// }

func (c *circuitBreaker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c.circuitBreaker.ServeHTTP(rw, req)
}



// Package plugindemo a demo plugin.
// package plugindemo
// 
// import (
// 	"bytes"
// 	"context"
// 	"fmt"
// 	"net/http"
// 	"text/template"
// )
// 
// 
// // Demo a Demo plugin.
// type Demo struct {
// 	next     http.Handler
// 	headers  map[string]string
// 	name     string
// 	template *template.Template
// }
// 
// // New created a new Demo plugin.
// func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
// 	if len(config.Headers) == 0 {
// 		return nil, fmt.Errorf("headers cannot be empty")
// 	}
// 
// 	return &Demo{
// 		headers:  config.Headers,
// 		next:     next,
// 		name:     name,
// 		template: template.New("demo").Delims("[[", "]]"),
// 	}, nil
// }
// 
// func (a *Demo) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
// 	for key, value := range a.headers {
// 		tmpl, err := a.template.Parse(value)
// 		if err != nil {
// 			http.Error(rw, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 
// 		writer := &bytes.Buffer{}
// 
// 		err = tmpl.Execute(writer, req)
// 		if err != nil {
// 			http.Error(rw, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 
// 		req.Header.Set(key, writer.String())
// 	}
// 
// 	a.next.ServeHTTP(rw, req)
// }
