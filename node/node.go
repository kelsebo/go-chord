package main

import (
  "fmt"
  "flag"
//  "console"
  "os"
  "types"
  "types/dht"
  "types/sighandler"
  "time"
  "runtime"
)

//BEGIN FLAGS
var orga = flag.String("orga", "localhost", "Organizer address")
var orgp = flag.Int("orgp", 1234, "Organizer port")
var port = flag.Int("port", 0, "Port to listen for RPC connections")
var httpport = flag.Int("http", 0, "Port to listen for HTTP connections")
//END FLAGS

//BEGIN VARIABLES
var storage *types.Storage_t
var alive = true
//var chordnode *types.Chord_t
//END VARIABLES

//BEGIN TYPES
//END TYPES



//BEGIN FUNCTIONS
func exitfunc () {
  alive = false
}

func sigexitfunc () {
  fmt.Println ("Killed by signal")
  if storage != nil {
    storage.Chord.Organizer.ReportDead (storage.Chord.Node.ID)
  }
  os.Exit(1)
}

func main () {
  flag.Parse()

  sig := sighandler.New()
  sig.RegisterExitHandler (sigexitfunc)

/*
  console.RegisterCommand ("info", "Show info about node", nodeinfo)
  console.RegisterCommand ("leave", "Gracefully leave the ring", leave)
  go console.RunConsole(exitfunc)
*/
  storage = nil
  runtime.GOMAXPROCS(2)

  storage = dht.New(&alive, port, orga, orgp, httpport, 3 ,0)

  if alive && storage != nil {
    for alive {
      time.Sleep(2 * time.Second)
    }
  }
  if storage != nil {
    storage.Chord.Organizer.ReportDead (storage.Chord.Node.ID)
  }
}
//END FUNCTIONS

