package main

import (
	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net"
	"github.com/eawsy/aws-lambda-go-net/service/lambda/runtime/net/apigatewayproxy"
	"github.com/gin-gonic/gin"
	"net/http"
)

// Handle is the exported handler called by AWS Lambda.
var Handle apigatewayproxy.Handler

var router *gin.Engine

func init() {
	ln := net.Listen()

	Handle = apigatewayproxy.New(ln, nil).Handle

	// Any Go framework complying with the Go http.Handler interface can be used.
	// This includes, but is not limited to, Vanilla Go, Gin, Echo, Gorrila, Goa, etc.
	router = gin.Default()

	initializeRoutes()

	go http.Serve(ln, router)
}

func main() {
	// for local testing
	http.ListenAndServe(":9080", router)
}
