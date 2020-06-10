package main

import (
	"github.com/stretchr/testify/assert"
	"strconv"
	"strings"
	"testing"
)

type TestCase struct {
	name           string
	url            string
	serverHealth   string
	expectedServer string
}

var (
	serversPool2 = []string{
		"server1:8080",
		"server2:8080",
		"server3:8080",
	}

	cases = []TestCase{
		{"All servers working", "test/path/to/target/file/two",
			"1 1 1", serversPool2[0]},
		{"All servers working", "test/path/new",
			"1 1 1", serversPool2[1]},
		{"All servers working", "test/path/to/target/file",
			"1 1 1", serversPool2[2]},
		{"Server not working", "test/path/to/target/file/two",
			"0 1 1", serversPool2[1]},
		{"Server not working", "test/path/new",
			"1 0 1", serversPool2[0]},
		{"Server not working", "test/path/to/target/file",
			"1 1 0", serversPool2[0]},
		{"All servers down", "test/path/to/target/file",
			"0 0 0", ""},
	}
)

func TestBalancer(t *testing.T) {

	//serversPool = serversPool2
	safeServer = SafeServer{v: make([]Server, len(serversPool))}

	for _, tcase := range cases {
		t.Run(tcase.name, func(t *testing.T) {
			serverHealth := strings.Split(tcase.serverHealth, " ")
			for index, health := range serverHealth {
				healthBoolean, _ := strconv.ParseBool(health)
				safeServer.v[index] = Server{
					IsHealthy: healthBoolean}

			}

			assert.Equal(t, tcase.expectedServer, chooseHealthyServer(tcase.url))

		})
	}

}
