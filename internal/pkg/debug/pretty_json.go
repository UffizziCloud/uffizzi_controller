package debug

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PrettyJson(data interface{}) string {
	const (
		empty = ""
		tab   = " "
	)

	buffer := new(bytes.Buffer)

	encoder := json.NewEncoder(buffer)
	encoder.SetIndent(empty, tab)

	err := encoder.Encode(data)
	if err != nil {
		return fmt.Sprintf("%s", err)
	}

	return buffer.String()
}
