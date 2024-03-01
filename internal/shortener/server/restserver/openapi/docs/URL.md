# URL

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**LongUrl** | Pointer to **string** |  | [optional] 
**ShortUrl** | Pointer to **string** |  | [optional] 
**CreatedAt** | Pointer to **time.Time** |  | [optional] 
**Clicks** | Pointer to **int64** |  | [optional] 

## Methods

### NewURL

`func NewURL() *URL`

NewURL instantiates a new URL object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewURLWithDefaults

`func NewURLWithDefaults() *URL`

NewURLWithDefaults instantiates a new URL object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLongUrl

`func (o *URL) GetLongUrl() string`

GetLongUrl returns the LongUrl field if non-nil, zero value otherwise.

### GetLongUrlOk

`func (o *URL) GetLongUrlOk() (*string, bool)`

GetLongUrlOk returns a tuple with the LongUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLongUrl

`func (o *URL) SetLongUrl(v string)`

SetLongUrl sets LongUrl field to given value.

### HasLongUrl

`func (o *URL) HasLongUrl() bool`

HasLongUrl returns a boolean if a field has been set.

### GetShortUrl

`func (o *URL) GetShortUrl() string`

GetShortUrl returns the ShortUrl field if non-nil, zero value otherwise.

### GetShortUrlOk

`func (o *URL) GetShortUrlOk() (*string, bool)`

GetShortUrlOk returns a tuple with the ShortUrl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetShortUrl

`func (o *URL) SetShortUrl(v string)`

SetShortUrl sets ShortUrl field to given value.

### HasShortUrl

`func (o *URL) HasShortUrl() bool`

HasShortUrl returns a boolean if a field has been set.

### GetCreatedAt

`func (o *URL) GetCreatedAt() time.Time`

GetCreatedAt returns the CreatedAt field if non-nil, zero value otherwise.

### GetCreatedAtOk

`func (o *URL) GetCreatedAtOk() (*time.Time, bool)`

GetCreatedAtOk returns a tuple with the CreatedAt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCreatedAt

`func (o *URL) SetCreatedAt(v time.Time)`

SetCreatedAt sets CreatedAt field to given value.

### HasCreatedAt

`func (o *URL) HasCreatedAt() bool`

HasCreatedAt returns a boolean if a field has been set.

### GetClicks

`func (o *URL) GetClicks() int64`

GetClicks returns the Clicks field if non-nil, zero value otherwise.

### GetClicksOk

`func (o *URL) GetClicksOk() (*int64, bool)`

GetClicksOk returns a tuple with the Clicks field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClicks

`func (o *URL) SetClicks(v int64)`

SetClicks sets Clicks field to given value.

### HasClicks

`func (o *URL) HasClicks() bool`

HasClicks returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


