# SwaggerJsClient.GreetApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ping**](GreetApi.md#ping) | **GET** /user/ping | 


<a name="ping"></a>
# **ping**
> Object ping()



### Example
```javascript
var SwaggerJsClient = require('Swagger-js-client');

var apiInstance = new SwaggerJsClient.GreetApi();

var callback = function(error, data, response) {
  if (error) {
    console.error(error);
  } else {
    console.log('API called successfully. Returned data: ' + data);
  }
};
apiInstance.ping(callback);
```

### Parameters
This endpoint does not need any parameter.

### Return type

**Object**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

