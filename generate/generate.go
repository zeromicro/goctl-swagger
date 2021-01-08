package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
)

func Do(in Plugin) error {

	swagger, err := applyGenerate(in)
	if err != nil {
		fmt.Println(err)
	}
	var formatted bytes.Buffer
	enc := json.NewEncoder(&formatted)
	enc.SetIndent("", "  ")

	if err := enc.Encode(swagger); err != nil {

	}

	output := in.Dir + "/rest.swagger.json"

	err = ioutil.WriteFile(output, formatted.Bytes(), 0666)

	return err
}
