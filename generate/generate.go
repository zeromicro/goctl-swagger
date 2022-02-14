package generate

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/zeromicro/go-zero/tools/goctl/plugin"
)

func Do(filename string, host string, basePath string, in *plugin.Plugin) error {
	swagger, err := applyGenerate(in, host, basePath)
	if err != nil {
		fmt.Println(err)
	}
	var formatted bytes.Buffer
	enc := json.NewEncoder(&formatted)
	enc.SetIndent("", "  ")

	if err := enc.Encode(swagger); err != nil {
		fmt.Println(err)
	}

	output := in.Dir + "/" + filename

	err = ioutil.WriteFile(output, formatted.Bytes(), 0666)
	if err != nil {
		fmt.Println(err)
	}
	return err
}
