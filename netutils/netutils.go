package netutils

import (
  "fmt"
  "net"
)

/**
 * Will return the first valid ip that is not prefixed with '127'
 * from interface nr ifacenum (starting at 0)
 */
func GetInetIP (ifacenum int) net.IP {
  ifaces, err := net.Interfaces()
  if err != nil {
    return nil
  }
  for _, i := range ifaces {
    if ifacenum > 0 {
      ifacenum--
      continue
    }
    addrs, err := i.Addrs()
    if err != nil {
      return nil
    }
    for _, a := range addrs {
      ipnet, ok := a.(*net.IPNet)
      if !ok {
        continue
      }
      v4 := ipnet.IP.To4()
      if v4 == nil || v4[0] == 127 { //|| v4[0] == 10 {
        continue
      }
      return v4
    }
  }
  return nil
}

/**
 * Get TCPAddr based on IP object and a port number
 * If nil is received as port, a random available port is chosen
 */
func GetTCPAddr (ip net.IP, port *int) *net.TCPAddr {
  var _port int
  if port == nil {
    _port = 0
  } else {
    _port = *port
  }
  addr, err := net.ResolveTCPAddr ("tcp", fmt.Sprintf("%s:%d", ip.String(), _port))
  if err != nil {
    return nil
  }
  return addr
}

/**
 * Get a TCP listener for this node
 * bind: serve only on first valid ip-address that is not localhost
 * port: specific OR 0 for random available
 */
func GetInetTCPListener (bind bool, port *int) *net.TCPListener {
  var addr *net.TCPAddr
  if bind {
    ip := GetInetIP (0)
    if ip == nil {
      return nil
    }
    addr = GetTCPAddr (ip, port)
    if addr != nil {
      return nil
    }
  } else {
    var err error
    addr, err = net.ResolveTCPAddr ("tcp4", fmt.Sprintf("0:%d", *port))
    if err != nil {
      return nil
    }
  }
  listener, err := net.ListenTCP ("tcp4", addr)
  if err != nil {
    return nil
  }
  return listener
}
