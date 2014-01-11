package types

import (
  "errors"
  "fmt"
  "net/http"
  "io/ioutil"
  "net"
)

func (s *Storage_t) Get (key *DHTKey_t, val *DHTValue_t) error {
  v, ok := s.Storagemap[*key]
  if ! ok {
    if s.Leaving || s.Joining {
      if ! s.Chord.Successor.Equal (s.Chord.Node) {
        return s.Chord.Successor.SingleCall ("DHT.Get", key, val)
      }
    }
    *val = ""
    return errors.New ("404")
  }

  *val = v
  println ("GET: ", *key, " : ", v)
  return nil
}

func (s *Storage_t) Put (e *DHTElem_t, _ign *string) error {
  if s.Leaving {
    if ! s.Chord.Successor.Equal (s.Chord.Node) {
      return s.Chord.Successor.SingleCall ("DHT.Put", e, nil)
    }
    return errors.New ("404")
  }
  s.Storagemap[e.Key] = e.Value
  println ("PUT: ", e.Key, " : ", e.Value)
//  s.size += uint64 (len (e.Value))
  s.size++
  return nil
}

func (s *Storage_t) Size (ign *string, _ign *string) uint64 {
  return s.size
}

func (s *Storage_t) Reqhandler (w http.ResponseWriter, r *http.Request) {
  switch r.Method {
	case "GET":
		s.getHandler(&w, r)
	case "PUT":
		s.putHandler(&w, r)
  default:
    fmt.Fprintf (w, "Invalid method ", r.Method)
  }
}

func getReqDHTKeyVal(r *http.Request) *DHTElem_t {
	q := r.URL.Query()

	qkey, ok := q["key"]
	if !ok {
		return nil
	}

	key := qkey[0]

	dkv := DHTElem_t {
		ID:      hashfunc (key),
		Key:     DHTKey_t(key),
	}

	if r.Method == "PUT" {
		qval, ok := q["value"]
		if ok {
			dkv.Value = DHTValue_t(qval[0])
		} else {
			b, _ := ioutil.ReadAll(r.Body)
			dkv.Value = DHTValue_t(b)
		}
	}

	return &dkv
}

func (s *Storage_t) getHandler(w *http.ResponseWriter, r *http.Request) {
	dkv := getReqDHTKeyVal(r)
  if dkv == nil {
    fmt.Fprintf (*w, "500")
    return
  }

  var keyhost Host_t
  err := s.Chord.FindSuccessor (dkv.ID, &keyhost)
  if err != nil {
    fmt.Fprintf (*w, "501")
    return
  }

  if keyhost.Equal (s.Chord.Node) {
    s.Get (&dkv.Key, &dkv.Value)
  } else {
    err = keyhost.SingleCall ("DHT.Get", &dkv.Key, &dkv.Value)
    if err != nil {
      fmt.Fprintf (*w, err.Error ())
      return
    }
  }
  fmt.Fprintf(*w, string(dkv.Value))
}

func (s *Storage_t) putHandler(w *http.ResponseWriter, r *http.Request) {
	dkv := getReqDHTKeyVal(r)
  if dkv == nil {
    fmt.Fprintf (*w, "400")
    return
  }

  var keyhost Host_t
  err := s.Chord.FindSuccessor (dkv.ID, &keyhost)
  if err != nil {
    fmt.Fprintf (*w, "501")
    return
  }

  if keyhost.Equal (s.Chord.Node) {
    s.Put (dkv, nil)
  } else {
    err = keyhost.SingleCall ("DHT.Put", dkv, nil)
    if err != nil {
      fmt.Fprintf (*w, err.Error())
      return
    }
  }
  fmt.Fprintf (*w, "200 OK")
}

func (s *Storage_t) Leave (ign *string, _ign *string) error {
  var err error
  s.Leaving = true
  if s.Listener != nil {
    s.Listener.(*net.TCPListener).Close()
  }
  if ! s.Chord.Successor.Equal (s.Chord.Node) {
    err = s.Chord.Successor.SingleCall ("DHT.MoveDataToSuccessor", s.Storagemap, nil)
  }
  if err != nil {
    return err
  }
  s.Storagemap = nil
  s.size = 0
  s.Chord.Leave (nil, nil)
  s.Chord.Organizer.ReportDead (s.Chord.Node.ID)
  return nil
}

func (s *Storage_t) MoveDataToSuccessor (data Storagemap_t, ign *string) error {
  if s.Leaving {
    return s.Chord.Successor.SingleCall ("DHT.MoveDataToSuccessor", data, nil)
  }
  if s.Storagemap == nil {
    panic (nil)
  }
  if data != nil && len (data) > 0 {
    for k, v := range data {
      s.Storagemap [k] = v
    }
  }
  s.size = uint64 (len (s.Storagemap))
  return nil
}


func (s *Storage_t) DeleteAfterMerge (callerID ID_t, ign *string) error {
  if s.Storagemap == nil {
    panic (nil)
  }
  for k, _ := range s.Storagemap {
    if ! hashfunc (string (k)).InIntervalIncludeUpper (callerID, s.Chord.Node.ID) {
      delete (s.Storagemap, k)
    }
  }
  s.size = uint64 (len (s.Storagemap))
  return nil
}

func (s *Storage_t) GetData (callerID ID_t, data *Storagemap_t) error {
  if s.Joining || s.Leaving {
    return errors.New ("Wait")
  }
  if s.Storagemap == nil {
    panic (nil)
  }
  _data := make (Storagemap_t)

  for k, v := range s.Storagemap {
    if ! hashfunc (string (k)).InIntervalIncludeUpper (callerID, s.Chord.Node.ID) {
      _data[k] = v
    }
  }
  *data = _data
  return nil

}

func (s *Storage_t) Merge () bool{
  var pre Host_t
  err := s.Chord.Predecessor.SingleCall ("Chord.FindSuccessor", s.Chord.Node.ID, &pre)
  if err != nil || ! pre.Equal (s.Chord.Node) {
    return false
  }

  var suc Host_t
  err = s.Chord.Successor.SingleCall ("Chord.GetPredecessor", "", &suc)
  if err != nil || ! suc.Equal (s.Chord.Node) {
    return false
  }

  var data Storagemap_t
  err = s.Chord.Successor.SingleCall ("DHT.GetData", s.Chord.Node.ID, &data)
  if err != nil {
    return false
  }
  if len (data) > 0 {
    for k, v := range data {
      s.Storagemap[k] = v
    }
    s.Chord.Successor.SingleCall ("DHT.DeleteAfterMerge", s.Chord.Node.ID, nil)
  }
  s.size = uint64 (len (s.Storagemap))
  s.Joining = false
  return true
}

