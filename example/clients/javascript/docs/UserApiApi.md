# SwaggerJsClient.UserApiApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getUserInfo**](UserApiApi.md#getUserInfo) | **GET** /api/user/{id} | 获取用户信息
[**login**](UserApiApi.md#login) | **POST** /api/user/login | 登录
[**register**](UserApiApi.md#register) | **POST** /api/user/register | 注册
[**searchUser**](UserApiApi.md#searchUser) | **GET** /api/user/search | 用户搜索


<a name="getUserInfo"></a>
# **getUserInfo**
> UserInfoReply getUserInfo(id, body)

获取用户信息

### Example
```javascript
var SwaggerJsClient = require('Swagger-js-client');

var apiInstance = new SwaggerJsClient.UserApiApi();

var id = "id_example"; // String | 

var body = new SwaggerJsClient.UserInfoReq(); // UserInfoReq | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.getUserInfo(id, body, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **String**|  | 
 **body** | [**UserInfoReq**](UserInfoReq.md)|  | 

### Return type

[**UserInfoReply**](UserInfoReply.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="login"></a>
# **login**
> Object login(body)

登录

### Example
```javascript
var SwaggerJsClient = require('Swagger-js-client');

var apiInstance = new SwaggerJsClient.UserApiApi();

var body = new SwaggerJsClient.LoginReq(); // LoginReq | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.login(body, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**LoginReq**](LoginReq.md)|  | 

### Return type

**Object**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="register"></a>
# **register**
> Object register(body)

注册

注册一个用户

### Example
```javascript
var SwaggerJsClient = require('Swagger-js-client');

var apiInstance = new SwaggerJsClient.UserApiApi();

var body = new SwaggerJsClient.RegisterReq(); // RegisterReq | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.register(body, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**RegisterReq**](RegisterReq.md)|  | 

### Return type

**Object**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

<a name="searchUser"></a>
# **searchUser**
> UserSearchReply searchUser(body)

用户搜索

### Example
```javascript
var SwaggerJsClient = require('Swagger-js-client');

var apiInstance = new SwaggerJsClient.UserApiApi();

var body = new SwaggerJsClient.UserSearchReq(); // UserSearchReq | 


var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.searchUser(body, callback);
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**UserSearchReq**](UserSearchReq.md)|  | 

### Return type

[**UserSearchReply**](UserSearchReply.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

