package hostlist

import "types"

func New (maxsize uint) *types.Hostlist_t {
  list := new (types.Hostlist_t)
  list.Maxsize = maxsize
  return list
}
