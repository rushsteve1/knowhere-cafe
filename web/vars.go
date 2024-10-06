// Enable the rpc, expvar and pprof endpoints
// This file (among others) also sets some expvars

package web

import (
	"expvar"
	_ "net/http/pprof"
	"net/rpc"
	"os"
	"time"
)

func init() {
	rpc.HandleHTTP()

	expvar.Publish("environ", expvar.Func(func() any {
		return os.Environ()
	}))

	expvar.Publish("time", expvar.Func(func() any {
		return time.Now()
	}))

	expvar.Publish("cwd", expvar.Func(func() any {
		cwd, _ := os.Getwd()
		return cwd
	}))
}
