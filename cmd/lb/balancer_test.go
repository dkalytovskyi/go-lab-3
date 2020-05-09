package main

import (
	"testing"
	"github.com/stretchr/testify/assert"
)
func TestBalancer(t *testing.T) {
	assert := assert.New(t)
	var serversPool = []string {
		"server1:8080",
		"server2:8080",
		"server3:8080",
		"server4:8080",
		"server5:8080",
		"server6:8080",
	}
	var URLs = []string {
		"test/path/to/target/file",
		"test/path/to/target/file/one",
		"test/path/to/target/file/two",
		"another/test/path/animals",
		"another/test/path/cars",
		"another/test/path/movies",
	}
	var expectedHash = []uint32 {
		4193102960,
		2044097447,
		466081617,
		884529959,
		191875093,
		4268647763,
	}

	for i, _ := range serversPool {
		assert.Equal(expectedHash[i], hash(URLs[i]))
		assert.Equal(int(expectedHash[i])%len(serversPool), determineServerByURL(URLs[i], serversPool))
	}
}
