# \DefaultAPI

All URIs are relative to *http://localhost:7770/v1*

Method | HTTP request | Description
------------- | ------------- | -------------
[**UrlPost**](DefaultAPI.md#UrlPost) | **Post** /url | Create a short URL
[**UrlShortUrlDelete**](DefaultAPI.md#UrlShortUrlDelete) | **Delete** /url/{short_url} | Delete the short URL
[**UrlShortUrlGet**](DefaultAPI.md#UrlShortUrlGet) | **Get** /url/{short_url} | Redirect to long URL or get long URL info



## UrlPost

> UrlPost201Response UrlPost(ctx).UrlPostRequest(urlPostRequest).Execute()

Create a short URL

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/Leopold1975/url_shortener/internal/shortener/server/restserver/openapi"
)

func main() {
	urlPostRequest := *openapiclient.NewUrlPostRequest("LongUrl_example") // UrlPostRequest | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.UrlPost(context.Background()).UrlPostRequest(urlPostRequest).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.UrlPost``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UrlPost`: UrlPost201Response
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.UrlPost`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiUrlPostRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **urlPostRequest** | [**UrlPostRequest**](UrlPostRequest.md) |  | 

### Return type

[**UrlPost201Response**](UrlPost201Response.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UrlShortUrlDelete

> UrlShortUrlDelete(ctx, shortUrl).Execute()

Delete the short URL

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/Leopold1975/url_shortener/internal/shortener/server/restserver/openapi"
)

func main() {
	shortUrl := "shortUrl_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.DefaultAPI.UrlShortUrlDelete(context.Background(), shortUrl).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.UrlShortUrlDelete``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**shortUrl** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUrlShortUrlDeleteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------


### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: Not defined

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## UrlShortUrlGet

> URL UrlShortUrlGet(ctx, shortUrl).Info(info).Execute()

Redirect to long URL or get long URL info

### Example

```go
package main

import (
	"context"
	"fmt"
	"os"
	openapiclient "github.com/Leopold1975/url_shortener/internal/shortener/server/restserver/openapi"
)

func main() {
	shortUrl := "shortUrl_example" // string | 
	info := true // bool |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.DefaultAPI.UrlShortUrlGet(context.Background(), shortUrl).Info(info).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `DefaultAPI.UrlShortUrlGet``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `UrlShortUrlGet`: URL
	fmt.Fprintf(os.Stdout, "Response from `DefaultAPI.UrlShortUrlGet`: %v\n", resp)
}
```

### Path Parameters


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
**ctx** | **context.Context** | context for authentication, logging, cancellation, deadlines, tracing, etc.
**shortUrl** | **string** |  | 

### Other Parameters

Other parameters are passed through a pointer to a apiUrlShortUrlGetRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------

 **info** | **bool** |  | 

### Return type

[**URL**](URL.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

