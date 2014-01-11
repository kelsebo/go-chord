package main

import (
  "errors"
  "fmt"
  "strconv"
  "types"
  "console"
)

func startconsole (exitfunc func()) {
  console.RegisterCommand ("kill", "Kill node. 'ALL', '[0-n]', 'id'", killnode)
  console.RegisterCommand ("leave", "Make node leave. 'ALL', '[0-n]', 'id'", leavenode)
  console.RegisterCommand ("list", "List nodes. 'dht', 'rpc', 'size'", listnodes)
  console.RegisterCommand ("size", "Show the total storage size for the ring", sumsize)
  console.RunConsole (exitfunc)
}

func leavenode (args []string) error {
  if len (args) < 1 {
    return errors.New ("Argument required")
  }
  nodemap := *n.Nodemap
  if args[0] == "ALL" {
    i := len (nodemap)
    for _, v := range nodemap {
      v.Node.SingleCall ("Slave.Leave", "", nil)
    }
    return errors.New (fmt.Sprintf ("Sent leave message to %d nodes", i))
  }

  if args[0] == "node" {
    if len (args) < 2 {
      return errors.New ("Need ID of at least one node")
    }
    for i := 1; i < len(args); i++ {
      nodeid, err := strconv.ParseUint(args[i], 16, 0)
      if err == nil {
        v, ok := nodemap[types.ID_t(nodeid)]
        if ok {
          v.Node.SingleCall ("Slave.Leave", "", nil)
          fmt.Println("Sent leave to node'" + args[i] + "'")
        }
      }
    }
    return nil
  }


  num, err := strconv.ParseUint(args[0], 10, 0)
  if err != nil {
    num, err = strconv.ParseUint(args[0], 16, 0)
  }
  if err == nil {
    if num < uint64(len (nodemap)) {
      for _, v := range nodemap {
        num--
        v.Node.SingleCall ("Slave.Leave", "", nil)
        if num <= 0 {
          return errors.New (fmt.Sprintf ("Sent leave message to %s nodes", args[0]))
        }
      }
    }
  }
  return nil
}

func killnode (args []string) error {
  if len (args) < 1 {
    return errors.New ("Argument required")
  }
  nodemap := *n.Nodemap
  if args[0] == "ALL" {
    i := len (nodemap)
    for _, v := range nodemap {
      v.Node.SingleCall ("Slave.Kill", "", nil)
    }
    return errors.New (fmt.Sprintf ("Sent killpill to %d nodes", i))
  }

  if args[0] == "node" {
    if len (args) < 2 {
      return errors.New ("Need ID of at least one node")
    }
    for i := 1; i < len(args); i++ {
      nodeid, err := strconv.ParseUint(args[i], 16, 0)
      if err == nil {
        v, ok := nodemap[types.ID_t(nodeid)]
        if ok {
          v.Node.SingleCall ("Slave.Kill", "", nil)
          fmt.Println("Killed node '" + args[i] + "'")
        }
      }
    }
    return nil
  }


  num, err := strconv.ParseUint(args[0], 10, 0)
  if err != nil {
    num, err = strconv.ParseUint(args[0], 16, 0)
  }
  if err == nil {
    if num < uint64(len (nodemap)) {
      for _, v := range nodemap {
        num--
        v.Node.SingleCall ("Slave.Kill", "", nil)
        if num <= 0 {
          return errors.New (fmt.Sprintf ("Sent killpill to %s nodes", args[0]))
        }
      }
    }
  }
  return nil
}

func listnodes (args []string) error {
  if len (args) == 1 {
    switch args[0] {
    case "dht" :
      printdhthosts ()
      return nil
    case "rpc" :
      printrpchosts ()
      return nil
    case "size" :
      printsize ()
      return nil
    default :
      return errors.New ("Unknown list command")

    }
  }
  nodemap := *n.Nodemap
  fmt.Printf ("      ========================== %3d Registered nodes =============================\n", len(nodemap))
  fmt.Println("            ID    |         IP         | PORT  |  HTTPPORT |   SIZE   |  ACTIVE ")
  fmt.Println("                  |                    |       |           |          |")
  for k, v := range nodemap {
    fmt.Printf ("        %8x  |    %s  | %5d |   %5d   |   %5d  |  ", k, v.Node.IP, v.Node.Port, v.Node.HttpPort, v.Size)
    fmt.Println (v.Active)
  }
  fmt.Println("     _____________________________________________________________________________")
  return nil
}

func columns (maplen int) int {
  linefeed := 0
  if maplen >= 10 {
    linefeed += 2
  }
  if maplen >= 20 {
    linefeed ++
  }
  if maplen >= 30 {
    linefeed ++
  }
  if maplen >= 40 {
    linefeed ++
  }
  return linefeed
}

func printdhthosts () {
  nodemap := *n.Nodemap
  linefeed := columns (len(nodemap))
  newline := linefeed
  for _, v := range nodemap {
    if v.Node.HttpPort > 0 {
      if newline == 0 {
        fmt.Printf ("%s:%d\n", v.Node.IP, v.Node.HttpPort)
        newline = linefeed
      } else {
        fmt.Printf ("%s:%d\t", v.Node.IP, v.Node.HttpPort)
        newline--
      }
    }
  }
  fmt.Println ()
}

func printrpchosts () {
  nodemap := *n.Nodemap
  linefeed := columns (len(nodemap))
  newline := linefeed
  for _, v := range nodemap {
    if newline == 0 {
      fmt.Printf ("%s:%d\n", v.Node.IP, v.Node.Port)
      newline = linefeed
    } else {
      fmt.Printf ("%s:%d\t", v.Node.IP, v.Node.Port)
      newline--
    }
  }
  fmt.Println ()
}

func printsize () {
  nodemap := *n.Nodemap
  linefeed := columns (len(nodemap))
  newline := linefeed
  for _, v := range nodemap {
    if newline == 0 {
      fmt.Printf("%x:%d\n", v.Node.ID, v.Size)
      newline = linefeed
    } else {
      fmt.Printf("%x:%d\t", v.Node.ID, v.Size)
      newline--
    }
  }
}

func sumsize (args []string) error {
  nodemap := *n.Nodemap
  sum := types.Size_t(0)
  for _, v := range nodemap {
    sum += v.Size
  }
  return errors.New (fmt.Sprintf("Total ring size: %d", sum))
}

