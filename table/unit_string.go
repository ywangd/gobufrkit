// Code generated by "stringer -type=Unit"; DO NOT EDIT.

package table

import "fmt"

const _Unit_name = "NUMERICSTRINGCODEFLAGCOMMON_CODEBIN"

var _Unit_index = [...]uint8{0, 7, 13, 17, 21, 32, 35}

func (i Unit) String() string {
	if i < 0 || i >= Unit(len(_Unit_index)-1) {
		return fmt.Sprintf("Unit(%d)", i)
	}
	return _Unit_name[_Unit_index[i]:_Unit_index[i+1]]
}
