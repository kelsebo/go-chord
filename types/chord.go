package types

import (
  "fmt"
  "errors"
  "runtime"
)

func (c *Chord_t) GetSuccessors (num uint,  reply *Host_t) error {
  elem := c.Successorlist.First

  if num > c.Successorlist.Size || num > c.Successorlist.Maxsize {
    return errors.New ("Index out of bound")
  }
  for i := uint(0); i < num && elem.Next != nil; i++ {
    elem = elem.Next
  }
  if elem == nil {
    return errors.New ("NIL ELEM")
  }
  *reply = *elem.Host

  return nil
}

/**
 * If returned error is nil, then the successor is the callers new successor,
 * otherwise, iterate further
 */
func (c *Chord_t) FindSuccessor (callerID ID_t, successor *Host_t) error {
  mysuccessor := c.Successor
  if mysuccessor == nil {
    panic (nil)
  }
  if callerID.InIntervalIncludeUpper (c.Node.ID, mysuccessor.ID) {
    *successor = *c.Successor
    return nil
  }
  return mysuccessor.SingleCall ("Chord.FindSuccessor", callerID, successor)
}

/**
 * Caller is predecessor if I have no predecessor OR it is in the interval
 */
func (c *Chord_t) Notify (caller *Host_t, _ign *string) error {
  c.mutex.Lock()
  defer c.mutex.Unlock()
  if c.Predecessor == nil || caller.ID.InInterval(c.Predecessor.ID, c.Node.ID) {
    c.Predecessor = caller
    println ("New predecessor:", c.Predecessor.String())
  }
  return nil
}

/**
 * Returns my predecessor, or an error if I do not have one
 */
func (c *Chord_t) GetPredecessor (_ign string, predecessor *Host_t) error {
  c.mutex.Lock()
  defer c.mutex.Unlock()
  if c.Predecessor == nil {
    return errors.New ("I have no predecessor")
  }
  *predecessor = *c.Predecessor
  return nil
}

func (c *Chord_t) CheckPredecessor (ign *string, _ign *string) error {
  if c.Predecessor != nil && ! c.Predecessor.Alive () {
    c.mutex.Lock ()
    c.Predecessor = nil
    c.mutex.Unlock ()
  }
  return nil
}


func (c *Chord_t) Join (h *Host_t, ign *string) error {
  c.Predecessor = nil

  var successor Host_t
  err := h.SingleCall("Chord.FindSuccessor", &c.Node.ID, &successor)
  if err != nil {
    return err
  }
  c.Successor = &successor

  err = successor.SingleCall("Chord.Notify", c.Node, nil)

  return err
}


func (c *Chord_t) Stabilize (ign *string, change *bool) error {
  c.mutex.Lock()
  runtime.Gosched ()
  c.mutex.Unlock()
  var successorsPredecessor Host_t
  var ch = false
  if c.Successor == nil {
    *change = false
    return nil
  }
  conn := c.Successor.Connect ()
  if conn == nil {
    c.Organizer.ReportDead (c.Successor.ID)
    c.mutex.Lock()
    c.Successor = c.Node
    *change = true
    c.mutex.Unlock()
    return nil
  }
  err := conn.Call("Chord.GetPredecessor", "", &successorsPredecessor)
  if err == nil {
    if successorsPredecessor.ID.InInterval(c.Node.ID, c.Successor.ID) {
      c.Successor = &successorsPredecessor
      ch = true
    }
  }
  conn.Call("Chord.Notify", c.Node, nil)
  conn.Close()

  *change = ch
  return nil
}


func (c *Chord_t) ClosestPrecedingNode (id ID_t, reply *Host_t) error {
  for i := 0; i < len (c.Finger); i++ {
    if c.Finger[i].ID.InInterval (c.Node.ID, id) {
      *reply = *c.Finger[i]
      return nil
    }
  }
  return nil
}

func (c *Chord_t) FixFingers (n *int, ign *string) error {
  next := uint(*n + 1)
  if next >= uint(len (c.Finger)) {
    next = 0
  }

  nextid := int(c.Node.ID) + int(2<<(next - 1))
  var finger Host_t

  err := c.FindSuccessor (ID_t(nextid), &finger)
  if err != nil {
    println (err.Error())
    return nil
  }
  c.Finger[next] = &finger
  fmt.Println (c.Finger)
  *n = int (next)
  return nil
}

func (c *Chord_t) NewPredecessor (predecessor Host_t, ign *string) error {
  c.mutex.Lock ()
  if !c.Alive {
    c.mutex.Unlock ()
    c.Successor.SingleCall ("Chord.NewPredecessor", predecessor, nil)
    return nil
  }
  println ("New predecessor!")
  c.Predecessor = &predecessor
  c.mutex.Unlock()
  return nil
}

func (c *Chord_t) NewSuccessor (successor Host_t, ign *string) error {
  c.mutex.Lock ()
  if !c.Alive {
    c.mutex.Unlock ()
    return c.Predecessor.SingleCall ("Chord.NewSuccessor", successor, nil)
  }
  println ("New successor!")
  c.Successor = &successor
  c.mutex.Unlock ()
  c.Successor.SingleCall ("Chord.NewPredecessor", *c.Node, nil)
  c.Organizer.Report (c)

  return nil
}

func (c *Chord_t) Leave (ign *string, _ign *string) error {
  c.Alive = false
  if c.Successor.Equal (c.Node) {
    return nil
  }
  c.Predecessor.SingleCall ("Chord.NewSuccessor", *c.Successor, nil)

  return nil
}

//END Functions

