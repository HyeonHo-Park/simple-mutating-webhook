package main

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"net/http"
	"net/http/httptest"
	"testing"
)

type payload struct {
	path   string
	method string
}

type response struct {
	httpResponse *httptest.ResponseRecorder
}

type request struct {
	httpRequest *http.Request
}

type httpTest struct {
	router *gin.Engine
	test   *testing.T
	payload
	response
	request
}

func NewHttpTest(router *gin.Engine, test *testing.T) *httpTest {
	return &httpTest{router: router, test: test}
}

func (ctx *httpTest) Get(url string) {
	ctx.payload.path = url
	ctx.payload.method = http.MethodGet
}

func (ctx *httpTest) Post(url string) {
	ctx.payload.path = url
	ctx.payload.method = http.MethodPost
}

func (ctx *httpTest) Delete(url string) {
	ctx.payload.path = url
	ctx.payload.method = http.MethodDelete
}

func (ctx *httpTest) Put(url string) {
	ctx.payload.path = url
	ctx.payload.method = http.MethodPut
}

func (ctx *httpTest) Patch(url string) {
	ctx.payload.path = url
	ctx.payload.method = http.MethodPatch
}

func (ctx *httpTest) Head(url string) {
	ctx.payload.path = url
	ctx.payload.method = http.MethodHead
}

func (ctx *httpTest) Options(url string) {
	ctx.payload.path = url
	ctx.payload.method = http.MethodOptions
}

func (ctx *httpTest) Set(key, value string) {
	ctx.request.httpRequest.Header.Set(key, value)
}

func (ctx *httpTest) Send(payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		ctx.test.Error(err.Error())
		return
	}

	req, err := http.NewRequest(ctx.payload.method, ctx.payload.path, bytes.NewBuffer(response))
	if err != nil {
		ctx.test.Error(err.Error())
		return
	}

	req.Header.Add("Access-Control-Allow-Origin", "*")
	req.Header.Add("Access-Control-Allow-Headers", "*")
	req.Header.Add("Access-Control-Expose-Headers", "*")
	req.Header.Add("User-Agent", "go-supertest/0.0.1")

	ctx.request.httpRequest = req
	ctx.response.httpResponse = httptest.NewRecorder()
}

func (ctx *httpTest) End(handleFunc func(req *http.Request, rr *httptest.ResponseRecorder)) {
	ctx.router.ServeHTTP(ctx.response.httpResponse, ctx.request.httpRequest)
	handleFunc(ctx.request.httpRequest, ctx.response.httpResponse)
}
