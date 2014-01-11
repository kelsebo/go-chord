package types

import (
  "fmt"
)

//BEGIN String functions
func (i *ID_t) String () string {
	return fmt.Sprintf ("%d", i)
}
//END String functions

//BEGIN Functions
func (id ID_t) InIntervalIncludeUpper(from, to ID_t) bool {
	if from < to {
		if id > from && id <= to {
			return true
		}
	} else {
		if id > from || id <= to {
			return true
		}
	}
	return false
}
func (id ID_t) InInterval(from, to ID_t) bool {
	if from < to {
		if id > from && id < to {
			return true
		}
	} else {
		if id > from || id < to {
			return true
		}
	}
  return false
}
func (id ID_t) InIntervalIncludeLower(from, to ID_t) bool {
	if from < to {
		if id >= from && id < to {
			return true
		}
	} else {
		if id >= from || id < to {
			return true
		}
	}
	return false
}
func (id ID_t) InIntervalIncludeBoth(from, to ID_t) bool {
	if from < to {
		if id >= from && id <= to {
			return true
		}
	} else {
		if id >= from || id <= to {
			return true
		}
	}
	return false
}

//END Functions

