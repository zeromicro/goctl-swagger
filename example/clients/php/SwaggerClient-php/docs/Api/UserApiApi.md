# Swagger\Client\UserApiApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**getUserInfo**](UserApiApi.md#getUserInfo) | **GET** /api/user/{id} | 获取用户信息
[**login**](UserApiApi.md#login) | **POST** /api/user/login | 登录
[**register**](UserApiApi.md#register) | **POST** /api/user/register | 注册
[**searchUser**](UserApiApi.md#searchUser) | **GET** /api/user/search | 用户搜索


# **getUserInfo**
> \Swagger\Client\Model\UserInfoReply getUserInfo($id, $body)

获取用户信息

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$apiInstance = new Swagger\Client\Api\UserApiApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client()
);
$id = "id_example"; // string | 
$body = new \Swagger\Client\Model\UserInfoReq(); // \Swagger\Client\Model\UserInfoReq | 

try {
    $result = $apiInstance->getUserInfo($id, $body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling UserApiApi->getUserInfo: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **id** | **string**|  |
 **body** | [**\Swagger\Client\Model\UserInfoReq**](../Model/UserInfoReq.md)|  |

### Return type

[**\Swagger\Client\Model\UserInfoReply**](../Model/UserInfoReply.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **login**
> object login($body)

登录

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$apiInstance = new Swagger\Client\Api\UserApiApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client()
);
$body = new \Swagger\Client\Model\LoginReq(); // \Swagger\Client\Model\LoginReq | 

try {
    $result = $apiInstance->login($body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling UserApiApi->login: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**\Swagger\Client\Model\LoginReq**](../Model/LoginReq.md)|  |

### Return type

**object**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **register**
> object register($body)

注册

注册一个用户

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$apiInstance = new Swagger\Client\Api\UserApiApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client()
);
$body = new \Swagger\Client\Model\RegisterReq(); // \Swagger\Client\Model\RegisterReq | 

try {
    $result = $apiInstance->register($body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling UserApiApi->register: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**\Swagger\Client\Model\RegisterReq**](../Model/RegisterReq.md)|  |

### Return type

**object**

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

# **searchUser**
> \Swagger\Client\Model\UserSearchReply searchUser($body)

用户搜索

### Example
```php
<?php
require_once(__DIR__ . '/vendor/autoload.php');

$apiInstance = new Swagger\Client\Api\UserApiApi(
    // If you want use custom http client, pass your client which implements `GuzzleHttp\ClientInterface`.
    // This is optional, `GuzzleHttp\Client` will be used as default.
    new GuzzleHttp\Client()
);
$body = new \Swagger\Client\Model\UserSearchReq(); // \Swagger\Client\Model\UserSearchReq | 

try {
    $result = $apiInstance->searchUser($body);
    print_r($result);
} catch (Exception $e) {
    echo 'Exception when calling UserApiApi->searchUser: ', $e->getMessage(), PHP_EOL;
}
?>
```

### Parameters

Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **body** | [**\Swagger\Client\Model\UserSearchReq**](../Model/UserSearchReq.md)|  |

### Return type

[**\Swagger\Client\Model\UserSearchReply**](../Model/UserSearchReply.md)

### Authorization

No authorization required

### HTTP request headers

 - **Content-Type**: application/json
 - **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../../README.md#documentation-for-api-endpoints) [[Back to Model list]](../../README.md#documentation-for-models) [[Back to README]](../../README.md)

