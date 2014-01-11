package main

import (
	"flag"
	"fmt"
	"github.com/egraff/inf-3200-1-frontend/frontend"
	"github.com/egraff/inf-3200-1-frontend/frontendtest"
	"io/ioutil"
	"math"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"sync"
	"time"
)

var NUM_TESTS = flag.Int("tests", 100000, "LOL")

//type Nodes []string

//var nodes *[]string

var orgaddr string

var fail = 0

func GetNodes(orgaddr string) []string {
	req, _ := http.NewRequest("GET", "http://"+orgaddr+"/Nodes", nil)
	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		fmt.Println(err)
		return nil
	}
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	//nodes = strings.Fields(string(bs))
	return strings.Fields(string(bs))

}

type DHTFrontendHandler struct {
	storageBackendNodes []string
}

func (this *DHTFrontendHandler) GET(key string) ([]byte, error) {
	for {
		n := rand.Intn(len(this.storageBackendNodes))
		req, _ := http.NewRequest("GET", "http://"+this.storageBackendNodes[int(n)]+"/?key="+key, nil)
	//	fmt.Println(req)
		resp, err := http.DefaultClient.Do(req)
    if resp != nil {
      defer resp.Body.Close()
    }
		if err == nil {
			bs, err := ioutil.ReadAll(resp.Body)
			if err == nil {
				fail = 0
        if n == 0 {
          this.storageBackendNodes = GetNodes (orgaddr)
        }
				return bs, nil
			}
			fail++

			if fail > 5 {
				return nil, err
			} else {
				this.storageBackendNodes = GetNodes (orgaddr)
        fmt.Println (this.storageBackendNodes)

				time.Sleep(200 * time.Millisecond)
				continue
			}
		}
	}
}

func (this *DHTFrontendHandler) PUT(key string, value []byte) error {
	for {
		rand.Seed(time.Now().UTC().UnixNano())
		n := rand.Intn(len(this.storageBackendNodes))
		req, _ := http.NewRequest("PUT", "http://"+this.storageBackendNodes[int(n)]+"/?key="+key+"&value="+string(value) /*bytes.NewBuffer(value)*/, nil)
	//	fmt.Println(req)
		resp, err := http.DefaultClient.Do(req)
		if resp != nil {
			defer resp.Body.Close()
		}
		if err != nil {
			fail++
			if fail > 5 {
				return err
			} else {
				this.storageBackendNodes = GetNodes(orgaddr)
        fmt.Println (this.storageBackendNodes)

				time.Sleep(200 * time.Millisecond)
				continue
			}
		}
    if n == 0 {
      this.storageBackendNodes = GetNodes (orgaddr)
    }
		fail = 0
		return nil
	}
}

func main() {
	var err error
	var runTests bool
	var httpServerPort uint
	var listener net.Listener
	var handler frontend.StorageServerFrontendHandler

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
		//		fmt.Fprintf(os.Stderr, "  compute-1-1 [compute-1-2 ... compute-N-M]\n")
	}

	flag.UintVar(&httpServerPort, "port", 8000, "portnumber(default=8000)")
	flag.BoolVar(&runTests, "runtests", false, "")
	flag.StringVar(&orgaddr, "org", "", "Organizer address")

	flag.Parse()
	nodes := flag.Args()
	/*if len(nodes) == 0 {
		flag.Usage()
		os.Exit(1)
	}
	*/

	for len(nodes) == 0 {
		nodes = GetNodes(orgaddr)
		time.Sleep(5 * time.Second)
	}
	fmt.Println(nodes)

	/************************************************************
	 ** When you have implemented a proper handler, you should **
	 ** comment in the line below. As long as the line is      **
	 ** commented out, the frontend will use a local hashmap   **
	 ** to implement the key-value database, and the tests     **
	 ** pass!
	 ***********************************************************/
	handler = &DHTFrontendHandler{nodes}

	if httpServerPort > math.MaxUint16 {
		fmt.Println("Invalid port %d", httpServerPort)
		flag.Usage()
		os.Exit(2)
	}

	wg := new(sync.WaitGroup)
	done := make(chan bool)
	frontend := frontend.New(handler)
	http.Handle("/", frontend)

	listener, err = net.Listen("tcp", fmt.Sprintf(":%d", httpServerPort))
	if err != nil {
		fmt.Println("Failed to listen on port", httpServerPort, "with error", err.Error())
		os.Exit(3)
	}

	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, os.Interrupt)
		<-c
		listener.Close()
		close(done)
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		http.Serve(listener, nil)
	}()

	if runTests {
		fmt.Println("Running tests...")

		wg.Add(1)
		go func() {
			defer wg.Done()
			result := frontendtest.Run(fmt.Sprintf("http://localhost:%d", httpServerPort), *NUM_TESTS, done)
			if result {
				fmt.Println("Test passed!")
			} else {
				fmt.Println("Test failed!")
			}

			listener.Close()
		}()
	}

	wg.Wait()
	fmt.Println("Bye, bye!")
}
