package dht

import (
  "types"
  "types/host"
  "types/chord"
  "types/organizer"
  "netutils"
  "net/http"
  "net"
  "time"
)


func EnvControl (s *types.Storage_t) {
  if ! s.Joining {
    return
  }

  for {
    if s.Chord.Predecessor != nil && ! s.Chord.Predecessor.Equal (s.Chord.Node) {
      if s.Merge () {
        break
      }
   }

    time.Sleep (500 * time.Millisecond)
  }

}

func register (org *types.Organizer_t, port *int, fingers uint,
                maxs uint, alive *bool)(*types.Chord_t, *net.TCPListener){
  var self, joinhost *types.Host_t
  var listener *net.TCPListener
  var tries = 10

  for joinhost == nil {
    if tries <= 0 {
      *alive = false
    }
    tries --

    h, l := host.SetupHost (port)
    if h == nil || l == nil {
      panic (nil)
    }

    joinhost = org.Register (h)
    if joinhost == nil {
      l.Close()
      time.Sleep (1 * time.Second)
    } else {
      self = h
      listener = l
    }
    if ! *alive {
      return nil, nil
    }
  }
  c := chord.New (self, *org, fingers, maxs)
  if c == nil {
    panic (nil)
  }

  c.Successor = joinhost


  ok := netutils.RPCRegisterName ("Chord", c)
  if ok {
    go netutils.RPCServe (listener)
  }
  for ! c.Node.Alive () {
    time.Sleep (1 * time.Second)
  }
  println ("Serving RPC on: " + c.Node.Address ())

  return c, listener
}

func New (alive *bool, port *int,
          orga *string, orgp *int,
          httpport *int, fingers, maxsuccessors uint) *types.Storage_t {

  org := organizer.New(*orga, *orgp)
  if org == nil {
    panic (nil)
  }

  c, listener := register (org, port, fingers, maxsuccessors, alive)
  if c == nil {
    println ("Failed to setup chord")
    return nil
  }

  httplistener := host.SetupHttpListener (c.Node, httpport)

  dht := new (types.Storage_t)
  dht.Chord = c
  dht.Storagemap = make (types.Storagemap_t)
  dht.Chord.Owner = dht
  dht.Joining = true

  if httplistener != nil {
    dht.Listener = httplistener
  }

  netutils.RPCRegisterName("DHT", dht)

  joinhost := c.Successor
  if joinhost.Equal (c.Node) {
    dht.Joining = false
  }
  err := c.Join (joinhost, nil)
  if err != nil {
    org.ReportDead (c.Node.ID)
    return nil
  }
  println ("Joined chord ring on: " + joinhost.Address ())

  slave := new (types.Slave_t)
  slave.KillFunc = func () {
    println ("Killed by organizer")
    c.Organizer.ReportDead (c.Node.ID)
    listener.Close()
    if httplistener != nil {
      httplistener.Close()
    }
    *alive = false
  }
  slave.LeaveFunc = func () {
    println ("Received order to leave from organizer")
    c.Owner.(*types.Storage_t).Leave (nil, nil)
    *alive = false
  }
  ok := netutils.RPCRegisterName("Slave", slave)

  if ok {
    go chord.RunStabilizer (c)
    go chord.RunPredecessorCheck (c)
    go chord.RunReporter (c)
    go chord.RunFixFingers (c)
    go EnvControl (dht)

    if httplistener != nil {
      http.HandleFunc ("/", dht.Reqhandler)
      go http.Serve (httplistener, nil)
    }
  }
  c.Organizer.HostChange (c.Node)

  return dht
}

