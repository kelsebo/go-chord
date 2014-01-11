package types


func (s *Hostlist_t) String () string {
  var hosts = ""
  for h := s.First; h != nil; h = h.Next {
    hosts += "[" + h.Host.String() + "] "
    hosts += " "
  }
  return hosts
}


func (s *Hostlist_t) AddFirst (h *Host_t) {
  for s.Size >= s.Maxsize {
    s.RemoveLast ()
  }

  elem := new (Hostelem_t)

  elem.Host = h
  elem.Next = s.First

  s.First = elem
  s.Size++

}


func (s *Hostlist_t) Remove (h *Host_t) {
  elem := s.First
  if elem == nil {
    s.Size = 0
    return
  }

  if elem.Host.Equal (h) {
    s.First = elem.Next
    s.Size--
    return
  }

  prev := elem
  elem = elem.Next

  for elem != nil {
    if elem.Host.Equal (h) {
      prev.Next = elem.Next
      s.Size--
      return
    }
    prev = elem
    elem = elem.Next
  }
}

func (s *Hostlist_t) RemoveFirst () {
  if s.First == nil {
    return
  }

  s.First = s.First.Next
  s.Size--
}

func (s *Hostlist_t) RemoveLast () {
  if s.First == nil || s.First.Next == nil {
    s.First = nil
    s.Size = 0
    return
  }

  var elem *Hostelem_t
  elem = s.First
  for elem.Next.Next != nil {
    elem = elem.Next
  }

  elem.Next = nil
  s.Size--
  return

  //var prev, elem *Hostelem_t
  //prev = s.First
/*  if prev.Next == nil {
    s.First = nil
    s.Size--
    return
  }
  elem = prev.Next
*/
/*  for elem{
    if elem.Next == nil {
      prev.Next = nil
      s.Size--
      return
    }
    prev = elem
    elem = elem.Next
  }*/
}

