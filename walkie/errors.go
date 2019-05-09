package walkie

import (
	"fmt"
	"strings"
)

type UnknownFileError struct {
	path string
}

func (e UnknownFileError) Error() string {
	return fmt.Sprintf("File '%s' unknow", e.path)
}

type FileCompareError struct {
	path   string
	origin File
	other  File
}

func (e FileCompareError) Error() string {

	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("File '%s' differs origin!=comp :", e.path))

	if e.origin.Name != e.other.Name {
		sb.WriteString(fmt.Sprintf(" Name: '%s'!='%s'", e.origin.Name, e.other.Name))
	}

	if !e.origin.Mtime.Equal(e.other.Mtime) {
		sb.WriteString(fmt.Sprintf(" Mtime: '%v'!='%v'", e.origin.Mtime, e.other.Mtime))
	}

	if e.origin.Size != e.other.Size {
		sb.WriteString(fmt.Sprintf(" Size: '%v'!='%v'", e.origin.Size, e.other.Size))
	}

	if e.origin.SHA256 != e.other.SHA256 {
		sb.WriteString(fmt.Sprintf(" SHA256: '%v'!='%v'", e.origin.SHA256, e.other.SHA256))
	}

	return sb.String()
}
