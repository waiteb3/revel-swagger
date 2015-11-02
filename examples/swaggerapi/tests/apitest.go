package tests

import "github.com/revel/revel/testing"

type ApiTest struct {
	testing.TestSuite
}

func (t *ApiTest) Before() {
	println("Set up")
}

func (t *ApiTest) TestThatIndexPageWorks() {
	t.Get("/")
	t.AssertOk()
	t.AssertContentType("text/html; charset=utf-8")
}

func (t *ApiTest) TestSwaggerEndpoint() {
	t.Get("/api/endpoint")
	t.AssertOk()
	t.AssertContentType("application/json; charset=utf-8")
	t.AssertContains("Swagger routing!")
}

func (t *ApiTest) After() {
	println("Tear down")
}
