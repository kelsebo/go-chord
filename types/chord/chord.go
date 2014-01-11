package chord

import (
  "types"
  "types/hostlist"
  "time"
  "runtime"
)

/**
 * Create a new chord instance with successor = host
 */
func New (h *types.Host_t, organizer types.Organizer_t, fingers, maxsuccessors uint) *types.Chord_t {
  if h == nil {
    return nil
  }
  c := new (types.Chord_t)
  c.Predecessor = nil
  c.Node = h
  c.Alive = true

  c.Finger = make ([]*types.Host_t, fingers)
  for i := uint (0); i < fingers; i++ {
    c.Finger[i] = c.Node
  }

  c.Successor = h
  c.Successorlist = hostlist.New (maxsuccessors)
  for i := uint (0); i < maxsuccessors; i++ {
    c.Successorlist.AddFirst (c.Node)
  }
  c.Organizer = organizer
  return c
}

func RunReporter (c *types.Chord_t) {
  for {
    if !c.Alive {
      runtime.Goexit ()
    }
    c.Organizer.Report(c)
    time.Sleep (2500 * time.Millisecond)
  }
}

func RunStabilizer (c *types.Chord_t) {
  for {
    if !c.Alive {
      runtime.Goexit ()
    }
    var change bool
    c.Stabilize (nil, &change)
    if ! change {
      time.Sleep (1200 * time.Millisecond)
    } else {
      time.Sleep (50 * time.Millisecond)
    }
  }
}

func RunPredecessorCheck (c *types.Chord_t) {
  for {
    if !c.Alive {
      runtime.Goexit ()
    }
    c.CheckPredecessor (nil, nil)
    time.Sleep (1500 * time.Millisecond)
  }
}

func RunFixFingers (c *types.Chord_t) {
  runtime.Goexit() //XXX XXX XXX
  var next = 0
  for {
    if !c.Alive {
      runtime.Goexit ()
    }
    c.FixFingers (&next, nil)
    time.Sleep (1200 * time.Millisecond)
  }
}

//END Functions
