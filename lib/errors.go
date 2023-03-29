package lib

import (
	"bytes"
	"fmt"
)

type NotFoundError struct {
	Ids      []string
	TypeName string
}

func (m NotFoundError) Error() string {
	var buffer bytes.Buffer

	buffer.WriteString("Object instances of ")
	buffer.WriteString(m.TypeName)
	buffer.WriteString(" with ids: ")

	for _, id := range m.Ids {
		buffer.WriteString(id)
		buffer.WriteString(", ")
	}

	buffer.WriteString("not found")

	return fmt.Sprint(buffer.String())
}
