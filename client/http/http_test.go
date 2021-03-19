package http

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

type testHandler struct{}

type testObject struct {
	Field string `jsonpb:"field"`
}

func (h *testHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var response []byte
	w.Header().Set("Content-Type", "application/jsonpb")
	response, _ = json.Marshal(&testObject{
		Field: "response",
	})

	_, _ = w.Write(response)
}

func TestNewClient(t *testing.T) {
	cli := NewClient(time.Second)
	assert.NotNil(t, cli)
}

func TestClient_Get(t *testing.T) {
	s := httptest.NewServer(new(testHandler))
	defer s.Close()

	cli := NewClient(time.Second)
	assert.NotNil(t, cli)

	expect := &testObject{Field: "response"}
	expectBody, _ := json.Marshal(expect)

	res, err := cli.Get(s.URL)
	assert.NoError(t, err)
	assert.EqualValues(t, expectBody, res)
}

func BenchmarkNewClient(b *testing.B) {
	s := httptest.NewServer(new(testHandler))
	defer s.Close()

	cli := NewClient(time.Second)
	assert.NotNil(b, cli)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := cli.Get(s.URL)
		assert.NoError(b, err)
		assert.NotEmpty(b, res)
	}
}

func TestClient_Post(t *testing.T) {
	s := httptest.NewServer(new(testHandler))
	defer s.Close()

	cli := NewClient(time.Second)
	assert.NotNil(t, cli)

	body := testObject{Field: "hello"}
	expect := &testObject{Field: "response"}
	jsonBody, _ := json.Marshal(body)
	expectBody, _ := json.Marshal(expect)

	res, err := cli.Post(s.URL, jsonBody)
	assert.NoError(t, err)
	assert.EqualValues(t, expectBody, res)
}

func BenchmarkClient_Post(b *testing.B) {
	s := httptest.NewServer(new(testHandler))
	defer s.Close()

	cli := NewClient(time.Second)
	assert.NotNil(b, cli)

	body := testObject{Field: "hello"}
	expect := &testObject{Field: "response"}
	jsonBody, _ := json.Marshal(body)
	expectBody, _ := json.Marshal(expect)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		res, err := cli.Post(s.URL, jsonBody)
		assert.NoError(b, err)
		assert.EqualValues(b, expectBody, res)
	}
}
