package node

import (
  "fmt"
  "math"
  "math/rand"
  "types"
  "errors"
  "github.com/fjukstad/gographer"
  "time"
  "sync"
)

var mutex sync.Mutex

//BEGIN TYPES
type nodemap_t map[types.ID_t]*Node_t
type Node_t struct {
  graph        *gographer.Graph
  Nodemap     *nodemap_t
  Node        *types.Host_t
  Size        types.Size_t
  Predecessor types.ID_t
  Successor   types.ID_t
  Name        string
  Active      bool
}
//END TYPES

//BEGIN FUNCTIONS
func (n *Node_t) HostChange (h types.Host_t, _ok *bool) error {
  nodemap := *n.Nodemap
  if nodemap == nil {
    panic (nil)
  }
  node, ok := nodemap[h.ID]
  if !ok {
    return errors.New ("Not a registered node!")
  }
  node.Node = &h
  *_ok = true
  return nil
}

func (n *Node_t) Report (c types.Report_t, _ign *string) error {
  nodemap := *n.Nodemap
  if nodemap == nil {
    panic (nil)
  }
  node, ok := nodemap[c.Node]
  if !ok {
    return errors.New("Not a registered node!")
  }

  node.Active = true

  if c.Size != node.Size {
    node.Size = c.Size
    n.Name = fmt.Sprintf ("%x  : %d", node.Node.ID, node.Size)
    n.graph.RenameNode (int(node.Node.ID), n.Name)
  }
  if c.Predecessor != node.Predecessor {
   // n.graph.RemoveEdge (int(node.Node.ID), int(node.Predecessor), int(node.Node.ID))
    node.Predecessor = c.Predecessor
   // n.graph.AddEdge (int(node.Node.ID), int(node.Predecessor), int(node.Node.ID), 0)
  }
  if c.Successor != node.Successor {
    n.graph.RemoveEdge (int(node.Node.ID), int(node.Successor), int(node.Node.ID))
    node.Successor = c.Successor
    n.graph.AddEdge (int(node.Node.ID), int(node.Successor), int(node.Node.ID), 0)
  }
  //nodemap[node.Node.ID] = node
  return nil
}

func (n *Node_t) ReportDead (nodeID types.ID_t, _ign *string) error {
  nodemap := *n.Nodemap
  if nodemap == nil {
    panic (nil)
  }
  node, ok := nodemap[nodeID]
  if !ok {
    return nil
  }
  node.Active = false
  n.graph.RemoveEdge (int(node.Node.ID), int(node.Successor), int(node.Node.ID))
  n.graph.RemoveNode (int(node.Node.ID))
  delete (nodemap, nodeID)
  return nil
}

func (n *Node_t) addNode (h *types.Host_t) {
  nodemap := *n.Nodemap
  newnode := new(Node_t)
  newnode.Node = h
  newnode.Successor = newnode.Node.ID
  newnode.Predecessor = newnode.Node.ID
  newnode.Name = fmt.Sprintf ("%x  : %d", newnode.Node.ID, newnode.Size)
  nodemap [h.ID] = newnode
  n.graph.AddNode (int (newnode.Node.ID), newnode.Name, 0, 1)
}

func (n *Node_t) Register (h *types.Host_t, reply *types.Host_t) error {
  nodemap := *n.Nodemap
  fmt.Print ("Register request: '" + h.String() + "'\t... ")
  if nodemap == nil {
    panic (nil)
  }

  _, ok := nodemap[h.ID]
  if ok {
    fmt.Println("Got register request from someone I already know")
    return errors.New ("P O")
  }

  n.getJoinHost (h, reply)
  n.addNode(h)
  fmt.Print ("Telling new host to join on\t'" + reply.String () + "' ... ")
  fmt.Println ("OK")

  return nil
}

func (n *Node_t) ReRegister (h *types.Host_t, ign *string) error {
  nodemap := *n.Nodemap
  if nodemap == nil {
    panic (nil)
  }

  _, ok := nodemap[h.ID]
  if ok {
    return errors.New ("P O")
  }
  n.addNode (h)
  return nil
}

func (n *Node_t) getRandomNode (maxattempts int) *Node_t {
  maxattempts--
  nodemap := *n.Nodemap
  if len (nodemap) == 0 {
    return nil
  }
  rand.Seed (time.Now().UTC().UnixNano())
  hostnum := math.Mod (float64 (rand.Int ()), float64 (len (nodemap)))
  var key types.ID_t
  for key, _ = range nodemap {
    hostnum--
    if hostnum <= 0{
      break
    }
  }
  node, ok := nodemap[key]
  if !ok {
    return nil
  }
  if !node.Active && len (nodemap) > 1 && maxattempts > 0{
    return n.getRandomNode (maxattempts)
  }
  return node

}

func (n *Node_t) getJoinHost (h, reply *types.Host_t) {

  tmpnode := n.getRandomNode (2)
  if tmpnode == nil {
    *reply = *h
    return
  }

  if ! tmpnode.Node.Alive () {
    n.ReportDead (tmpnode.Node.ID, nil)
    n.getJoinHost (h, reply)
  }
  *reply = *tmpnode.Node
}

func (n *Node_t) ControlPresenceOfRandomNode () {
  if len((*n.Nodemap)) == 0 {
    return
  }
  node := n.getRandomNode (2)
  if node == nil {
    return
  }
  if ! node.Node.Alive() {
    n.ReportDead (node.Node.ID, nil)
    return
  }
}

func (n *Node_t) String () string {
  nodemap := *n.Nodemap
  if nodemap == nil || len(nodemap) == 0 {
    return "Empty Nodemap"
  }

  var nodes = ""
  for _, v := range nodemap {
    nodes += fmt.Sprintf("%d\n", v.Node.ID)
  }
  return nodes
}

func New () *Node_t {
  n := new(Node_t)
  nodemap := make (nodemap_t)
  n.Nodemap = &nodemap
  graph := gographer.New()
  n.graph = graph
  return n
}

func (n *Node_t) GetDHTHosts () string {
  var hosts = ""
  nodemap := *n.Nodemap
  for _, v := range nodemap {
    if v.Node.HttpPort > 0 {
      hosts += fmt.Sprintf("%s:%d ", v.Node.IP, v.Node.HttpPort)
    }
  }
  return hosts
}
//END FUNCTIONS

