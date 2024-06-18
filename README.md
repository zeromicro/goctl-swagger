# goctl-swagger

### 1. 编译goctl-swagger插件

```
GOPROXY=https://goproxy.cn/,direct go install github.com/zeromicro/goctl-swagger@latest
```

### 2. 配置环境

将$GOPATH/bin中的goctl-swagger添加到环境变量

### 3. 使用姿势

* 创建api文件

    ```go
    info(
     title: "type title here"
     desc: "type desc here"
     author: "type author here"
     email: "type email here"
     version: "type version here"
    )
    
    
    type (
     RegisterReq {
      Username string `json:"username"`
      Password string `json:"password"`
      Mobile string `json:"mobile"`
     }
     
     LoginReq {
      Username string `json:"username"`
      Password string `json:"password"`
     }
     
     UserInfoReq {
      Id string `path:"id"`
     }
     
     UserInfoReply {
      Name string `json:"name"`
      Age int `json:"age"`
      Birthday string `json:"birthday"`
      Description string `json:"description"`
      Tag []string `json:"tag"`
     }
     
     UserSearchReq {
      KeyWord string `form:"keyWord"`
     }
  
     UploadReq {
      Type int    `form:"type"`
      Key  string `form:"key,optional"`
     }
     
     UploadResp {
      Name string `json:"name"`
      Age int `json:"age"`
      Birthday string `json:"birthday"`
      Description string `json:"description"`
      Tag []string `json:"tag"`
     }
    )
    
    @server(
      // 这里是为了给service添加swagger的security，用于使用Authorize添加到接口的header
      // 如果添加了jwt，这个描述可以省略
      security: true
    )
    service user-api {
     @doc(
      summary: "注册"
     )
     @handler register
     post /api/user/register (RegisterReq)
     
     @doc(
      summary: "登录"
     )
     @handler login
     post /api/user/login (LoginReq)
     
     @doc(
      summary: "获取用户信息"
     )
     @handler getUserInfo
     get /api/user/:id (UserInfoReq) returns (UserInfoReply)
     
     @doc(
      summary: "用户搜索"
     )
     @handler searchUser
     get /api/user/search (UserSearchReq) returns (UserInfoReply)
  
     @doc(
      summary: "上传文件，body修改为formData，并添加一个file参数"
      inject_formdata_param: "file"
     )
     @handler UploadHandler
     post /file/upload (UploadReq) returns (UploadResp)
    }
    ```

* 生成swagger.json 文件

    ```shell script
    goctl api plugin -plugin goctl-swagger="swagger -filename user.json" -api user.api -dir .
    ```

* 指定Host，basePath，schemes [api-host-and-base-path](https://swagger.io/docs/specification/2-0/api-host-and-base-path/)

    ```shell script
    goctl api plugin -plugin goctl-swagger="swagger -filename user.json -host 127.0.0.2 -basepath /api -schemes https,wss" -api user.api -dir .
    ```

* swagger ui 查看生成的文档

    ```shell script
     docker run --rm -p 8083:8080 -e SWAGGER_JSON=/foo/user.json -v $PWD:/foo swaggerapi/swagger-ui
   ```

* Swagger Codegen 生成客户端调用代码(go,javascript,php)

  ```shell script
  for l in go javascript php; do
    docker run --rm -v "$(pwd):/go-work" swaggerapi/swagger-codegen-cli generate \
      -i "/go-work/rest.swagger.json" \
      -l "$l" \
      -o "/go-work/clients/$l"
  done
   ```
