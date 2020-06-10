package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/dkalytovskyi/go-lab-3/httptools"
	"github.com/dkalytovskyi/go-lab-3/signal"
	"hash/fnv"
	"io"
	"log"
	"net/http"
	"sync"
	"time"
)

type Server struct {
	IsHealthy bool
}

type SafeServer struct {
	v   []Server
	mux sync.Mutex
}

var (
	port       = flag.Int("port", 8090, "load balancer port")
	timeoutSec = flag.Int("timeout-sec", 3, "request timeout time in seconds")
	https      = flag.Bool("https", false, "whether backends support HTTPs")

	traceEnabled = flag.Bool("trace", false, "whether to include tracing information into responses")
	safeServer   = SafeServer{v: make([]Server, len(serversPool))}
)

var (
	timeout     = time.Duration(*timeoutSec) * time.Second
	serversPool = []string{
		"server1:8080",
		"server2:8080",
		"server3:8080",
	}
)

func scheme() string {
	if *https {
		return "https"
	}
	return "http"
}

func health(dst string) bool {
	ctx, _ := context.WithTimeout(context.Background(), timeout)
	req, _ := http.NewRequestWithContext(ctx, "GET",
		fmt.Sprintf("%s://%s/health", scheme(), dst), nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return false
	}
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}

func forward(dst string, rw http.ResponseWriter, r *http.Request) error {
	if len(dst) < 1 {
		return fmt.Errorf("no healthy servers found")

	}
	ctx, _ := context.WithTimeout(r.Context(), timeout)
	fwdRequest := r.Clone(ctx)
	fwdRequest.RequestURI = ""
	fwdRequest.URL.Host = dst
	fwdRequest.URL.Scheme = scheme()
	fwdRequest.Host = dst

	resp, err := http.DefaultClient.Do(fwdRequest)
	if err == nil {
		for k, values := range resp.Header {
			for _, value := range values {
				rw.Header().Add(k, value)
			}
		}
		if *traceEnabled {
			rw.Header().Set("lb-from", dst)
		}
		log.Println("fwd", resp.StatusCode, resp.Request.URL)
		rw.WriteHeader(resp.StatusCode)
		defer resp.Body.Close()
		_, err := io.Copy(rw, resp.Body)
		if err != nil {
			log.Printf("Failed to write response: %s", err)
		}
		return nil
	} else {
		log.Printf("Failed to get response from %s: %s", dst, err)
		rw.WriteHeader(http.StatusServiceUnavailable)
		return err
	}
}

func determineServerByURL(URL string, ) int {

	return int(hash(URL)) % (len(serversPool))
}

func chooseHealthyServer(URL string) string {
	if serverInt := determineServerByURL(URL); safeServer.v[serverInt].IsHealthy == true {
		log.Println(serverInt)
		return serversPool[serverInt]
	} else {
		for index, server := range safeServer.v {
			if server.IsHealthy {
				return serversPool[index]
			}

		}

		return ""

	}

}

func hash(s string) uint32 {
	h := fnv.New32a()
	h.Write([]byte(s))
	return h.Sum32()
}

func main() {
	flag.Parse()

	for i, server := range serversPool {
		server:= server
		i := i

		safeServer.mux.Lock()
		safeServer.v[i] = Server{false}
		safeServer.mux.Unlock()

		go func() {
			for range time.Tick(10 * time.Second) {
				log.Println(server, health(server))
			}
		}()

		go func() {
			for range time.Tick(1 * time.Second) {
				safeServer.mux.Lock()
				safeServer.v[i].IsHealthy = health(server)
				safeServer.mux.Unlock()
			}
		}()

	}

	frontend := httptools.CreateServer(*port, http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		forward(chooseHealthyServer(r.URL.Path), rw, r)

	}))

	log.Println("Starting load balancer...")
	log.Printf("Tracing support enabled: %t", *traceEnabled)
	frontend.Start()
	signal.WaitForTerminationSignal()
}
