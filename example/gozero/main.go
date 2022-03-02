package main

import (
	"fmt"
	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
	"github.com/zeromicro/go-zero/rest/httpx"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const config = `{"Name": "gozero", "Host": "0.0.0.0", "Port": 8888}`

func main() {
	var c rest.RestConf
	var env string = "dev"
	err := conf.LoadConfigFromYamlBytes([]byte(config), &c)
	if err != nil {
		panic(err)
	}
	server, err := rest.NewServer(c)
	if err != nil {
		panic(err)
	}
	defer server.Stop()

	swaggerFile, err := os.Open("example/gozero/swagger.json")
	if err != nil {
		log.Println(err)
	}
	defer swaggerFile.Close()
	SwaggerByte, err := ioutil.ReadAll(swaggerFile)
	if err != nil {
		log.Println(err)
	}

	server.AddRoutes([]rest.Route{
		{
			Method:  http.MethodGet,
			Path:    "/swagger",
			Handler: Doc("/swagger", env),
		},
		{
			Method: http.MethodGet,
			Path:   "/swagger-json",
			Handler: func(writer http.ResponseWriter, request *http.Request) {
				writer.Header().Set("Content-Type", "application/json; charset=utf-8")
				_, err := writer.Write(SwaggerByte)
				if err != nil {
					httpx.Error(writer, err)
				}
			},
		},
	})

	fmt.Printf("Starting server at http://%s:%d...\n", c.Host, c.Port)
	server.Start()
}
