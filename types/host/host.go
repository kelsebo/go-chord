package host

import (
  "types"
  "net"
  "netutils"
  "fmt"
)

func New(addr string, port string) (*types.Host_t, error) {
  h := new(types.Host_t)
  tcpaddr, err := net.ResolveTCPAddr("tcp4", addr+":"+port)
  if err != nil {
    return nil, err
  }

  h.IP = types.IP_t(tcpaddr.IP.String())
  h.Port = types.Port_t(tcpaddr.Port)
  h.ID = h.Hash()

  return h, nil
}


func SetupHost (port *int) (*types.Host_t, *net.TCPListener) {
  listener := netutils.GetInetTCPListener (false, port)
  if listener == nil {
    return nil, nil
  }
  addr, err := net.ResolveTCPAddr ("tcp4", listener.Addr().String())
  if err != nil {
    return nil, nil
  }
  ip := netutils.GetInetIP(0)
  if ip == nil {
    return nil, nil
  }
  host, err := New(ip.String(), fmt.Sprintf("%d", addr.Port))
  if err != nil {
    return nil, nil
  }
  return host, listener
}

func SetupHttpListener (host *types.Host_t, port *int) *net.TCPListener {
  if types.Port_t(*port) == host.Port {
    host.HttpPort = types.Port_t(*port)
    return nil
  }

  listener := netutils.GetInetTCPListener (false, port)
  if listener == nil {
    return nil
  }
  addr, err := net.ResolveTCPAddr ("tcp4", listener.Addr().String())
  if err != nil {
    return nil
  }
  host.HttpPort = types.Port_t(addr.Port)
  return listener
}
//END Functions
