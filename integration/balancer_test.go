package integration

import (
	"fmt"
	"net/http"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
)

const baseAddress = "http://balancer:8090"

var client = http.Client{
	Timeout: 3 * time.Second,
}

func TestBalancer(t *testing.T) {
	resp1, err1 := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
	if err1 != nil {
		t.Error(err1)
	}
	resp2, err2 := client.Get(fmt.Sprintf("%s/api/v1/another-data", baseAddress))
	if err2 != nil {
		t.Error(err2)
	}
	resp3, err3 := client.Get(fmt.Sprintf("%s/api/v1/", baseAddress))
	if err3 != nil {
		t.Error(err3)
	}
	resp4, err4 := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
	if err4 != nil {
		t.Error(err4)
	}
	assert.NotEqual(resp1.Header.Get("lb-from"), resp2.Header.Get("lb-from"))
	assert.NotEqual(resp1.Header.Get("lb-from"), resp3.Header.Get("lb-from"))
	assert.NotEqual(resp2.Header.Get("lb-from"), resp3.Header.Get("lb-from"))
	assert.Equal(resp1.Header.Get("lb-from"), resp4.Header.Get("lb-from"))
}

func BenchmarkBalancer(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := client.Get(fmt.Sprintf("%s/api/v1/some-data", baseAddress))
		if err != nil {
			b.Error(err)
		}
	}
}
