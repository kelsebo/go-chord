package types

import (
  "netutils"
)


func (o *Organizer_t) Register (host *Host_t) *Host_t {
  if host == nil {
    panic (nil)
  }
  var reply Host_t
  err := netutils.RPCSingleCall(string(*o), "Node.Register", *host, &reply)
  if err != nil {
    return nil
  }
  return &reply
}

func (o *Organizer_t) HostChange (host *Host_t) bool {
  if host == nil {
    panic (nil)
  }
  var ok bool
  netutils.RPCSingleCall (string (*o), "Node.HostChange", *host, &ok)
  return ok
}

func (o *Organizer_t) ReRegister (host *Host_t) {
  if host == nil {
    panic (nil)
  }
  netutils.RPCSingleCall (string (*o), "Node.ReRegister", *host, nil)
}

func (o *Organizer_t) Report (c *Chord_t) {
  if c.Predecessor == nil {
    return
  }
  rep := new (Report_t)
  rep.Node = c.Node.ID
  if c.Predecessor != nil {
    rep.Predecessor = c.Predecessor.ID
  }
  if c.Successor != nil {
    rep.Successor = c.Successor.ID
  }
  rep.Size = Size_t(c.Owner.(*Storage_t).Size(nil, nil))
  err := netutils.RPCSingleCall(string(*o), "Node.Report", rep, nil)
  if err != nil && err.Error() == "Not a registered node!" {
    if c.Predecessor != nil {
      o.ReRegister (c.Node)
    }
  }
}

func (o *Organizer_t) ReportDead (nodeID ID_t) {
  netutils.RPCSingleCall(string(*o), "Node.ReportDead", nodeID, nil)
}
