package main

import (
  "fmt"
  "flag"
  "netutils"
  "organizer/node"
  "net/http"
  "math/rand"
  "time"
  "types/host"
  "types/sighandler"
)

//BEGIN FLAGS
var port = flag.Int("port", 1234, "Port to listen for RPC connections")
var killnodesonexit = flag.Bool("killnodes", true, "Send kill to all registered nodes on exit")
//END FLAGS

//BEGIN VARIABLES
var n *node.Node_t
//END VARIABLES

//BEGIN TYPES
//END TYPES

//BEGIN FUNCTIONS
func handleNodeReq (w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf (w, n.GetDHTHosts())
}

func graphServe () {
  rand.Seed(time.Now().UTC().UnixNano())
  panic (http.ListenAndServe (":8080", http.FileServer (http.Dir ("./root_serve_dir"))))
}

func main() {
  var ALIVE = true
  exitfunc := func () {
    ALIVE = false
  }
  flag.Parse()

  sighandler := sighandler.New()

  n = node.New()
  netutils.RPCRegisterName("Node", n)
  host, listener := host.SetupHost(port)
   if host == nil || listener == nil {
    panic(nil)
  }
  fmt.Println("Serving from: " + host.Address())

  sighandler.RegisterExitHandler (exitfunc)

  http.HandleFunc ("/Nodes", handleNodeReq)

  go netutils.RPCServe (listener)
  go graphServe ()
  go startconsole (exitfunc)
  go func() {
    for {
      n.ControlPresenceOfRandomNode()
      time.Sleep (1 * time.Second)
    }
  }()
  for ALIVE {
    time.Sleep(100 * time.Millisecond)
  }
  if *killnodesonexit {
    err := killnode ([]string {"ALL"});
    if err != nil {
      fmt.Println(err.Error())
    }
  }
}
//END FUNCTIONS

