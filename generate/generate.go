package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	plugin2 "github.com/tal-tech/go-zero/tools/goctl/plugin"
	"io/ioutil"
)

func Do(filename string, in *plugin2.Plugin) error {

	swagger, err := applyGenerate(in)
	if err != nil {
		fmt.Println(err)
	}
	var formatted bytes.Buffer
	enc := json.NewEncoder(&formatted)
	enc.SetIndent("", "  ")

	if err := enc.Encode(swagger); err != nil {

	}

	output := in.Dir + "/" + filename

	err = ioutil.WriteFile(output, formatted.Bytes(), 0666)

	return err
}
