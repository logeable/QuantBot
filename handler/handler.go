package handler

import (
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"

	"github.com/hprose/hprose-golang/rpc"
	"github.com/phonegapX/QuantBot/config"
	"github.com/phonegapX/QuantBot/constant"
)

type response struct {
	Success bool
	Message string
	Data    interface{}
}

type event struct{}

func (e event) OnSendHeader(ctx *rpc.HTTPContext) {
	ctx.Response.Header().Set("Access-Control-Allow-Headers", "Authorization")
}

// Server ...
func Server() {
	service := rpc.NewHTTPService()
	handler := struct {
		User      user
		Exchange  exchange
		Algorithm algorithm
		Trader    runner
		Log       logger
	}{}
	service.Event = event{}
	service.AddBeforeFilterHandler(func(request []byte, ctx rpc.Context, next rpc.NextFilterHandler) (response []byte, err error) {
		ctx.SetInt64("start", time.Now().UnixNano())
		httpContext := ctx.(*rpc.HTTPContext)
		if httpContext != nil {
			ctx.SetString("username", parseToken(httpContext.Request.Header.Get("Authorization")))
		}
		return next(request, ctx)
	})
	service.AddInvokeHandler(func(name string, args []reflect.Value, ctx rpc.Context, next rpc.NextInvokeHandler) (results []reflect.Value, err error) {
		name = strings.Replace(name, "_", ".", 1)
		results, err = next(name, args, ctx)
		spend := time.Now().UnixNano() - ctx.GetInt64("start")
		spendInfo := time.Duration(spend).Round(time.Millisecond)
		log.Printf("%16s() spend %v", name, spendInfo)
		return
	})
	service.AddAllMethods(handler)
	http.Handle("/api", service)
	http.Handle("/", http.FileServer(http.Dir("web/dist")))
	fmt.Printf("%v  Version %v\n", constant.Banner, constant.Version)
	log.Printf("Running at http://%v\n", config.Config.Server.Addr)
	log.Fatal(http.ListenAndServe(config.Config.Server.Addr, nil))
}
