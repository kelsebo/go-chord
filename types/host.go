package types

import (
  "fmt"
  "netutils"
  "net/rpc"
)


//BEGIN String functions
func (h *Host_t) String () string {
	return fmt.Sprintf ("%x:%s:%d:%d", h.ID, h.IP, h.Port, h.HttpPort)
}
func (h *Host_t) Address() string {
	return fmt.Sprintf ("%s:%d", h.IP, h.Port)
}
func (h *Host_t) HttpAddress () string {
	return fmt.Sprintf ("%s:%d", h.IP, h.HttpPort)
}
//END String functions

//BEGIN Functions
func (h *Host_t) Hash () ID_t {
  return hashfunc (h.Address ())
}
func (h *Host_t) Equal (o *Host_t) bool {
  return h.ID == o.ID
}
func (h *Host_t) SingleCall (serviceMethod string, args interface{}, reply interface{}) error {
  return netutils.RPCSingleCall (h.Address (), serviceMethod, args, reply)
}
func (h *Host_t) Alive () bool {
  return netutils.RPCAlive (h.Address())
}
func (h *Host_t) Connect () *rpc.Client {
  return netutils.RPCConnect (h.Address())
}

//END Functions

