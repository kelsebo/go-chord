package types

import "sync"

const THE_M = 8

//BEGIN Types
type ID_t                       uint64
type IP_t                       string
type Port_t                     uint16
type Size_t                     uint32
type Organizer_t                string
type DHTKey_t                   string
type DHTValue_t                 string
type Storagemap_t               map[DHTKey_t]DHTValue_t
//END Types


//BEGIN Structs
type Slave_t struct {
  KillFunc                      func()
  LeaveFunc                     func()
}

type Host_t struct {
  ID                            ID_t
  IP                            IP_t
  Port                          Port_t
  HttpPort                      Port_t
}

type Chord_t struct {
  Node                          *Host_t
  Predecessor                   *Host_t
  Successor                     *Host_t
  Finger                        []*Host_t
  Successorlist                 *Hostlist_t
  Organizer                     Organizer_t
  Owner                         interface{}
  mutex                         sync.Mutex
  Alive                         bool
}

type Hostelem_t  struct {
  Prev                          *Hostelem_t
  Host                          *Host_t
  Next                          *Hostelem_t
}

type Hostlist_t struct {
  First                         *Hostelem_t
  Last                          *Hostelem_t
  Size                          uint
  Maxsize                       uint
}

type Report_t struct {
  Predecessor, Node, Successor  ID_t
  Size                          Size_t
}

type Sighandler_t struct {
  ExitCallback                  func()
}

type Storage_t struct {
  Chord                         *Chord_t
  size                          uint64
  Storagemap                    Storagemap_t
  Listener                      interface {}
  mutex                         sync.Mutex
  Leaving                       bool
  Joining                       bool
}

type DHTElem_t struct {
  ID                            ID_t
  Key                           DHTKey_t
  Value                         DHTValue_t
}
//END Structs
