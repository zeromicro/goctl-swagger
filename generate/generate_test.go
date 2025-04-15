package generate

import (
	"testing"

	"github.com/zeromicro/go-zero/tools/goctl/api/parser"
	"github.com/zeromicro/go-zero/tools/goctl/plugin"
)

func TestDo(t *testing.T) {

	type args struct {
		apiFile  string
		filename string
		host     string
		basePath string
		schemes  string
		in       *plugin.Plugin
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "",
			args: args{
				apiFile:  "./upload.api",
				filename: "upload.json",
				host:     "127.0.0.1:8890",
				basePath: "/",
				schemes:  "http",
				in: &plugin.Plugin{
					Api:         nil,
					ApiFilePath: "./upload.api",
					Style:       "",
					Dir:         "./",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		// 解析api文件
		api, err := parser.Parse(tt.args.apiFile)
		if err != nil {
			t.Errorf("Parse() error = %v", err)
		}
		tt.args.in.Api = api

		t.Run(tt.name, func(t *testing.T) {
			if err := Do(tt.args.filename, tt.args.host, tt.args.basePath, tt.args.schemes, tt.args.in); (err != nil) != tt.wantErr {
				t.Errorf("Do() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
