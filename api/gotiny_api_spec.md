


# URL shortener API
  

## Informations

### Version

0.0.1

### License

[Apache 2.0](http://www.apache.org/licenses/LICENSE-2.0.html)

### Contact

Alex Astapchuk alex.astapchuk@gmail.com https://alex.astapchuk.com

## Content negotiation

### URI Schemes
  * http

### Consumes
  * application/json

### Produces
  * application/json

## All endpoints

###  operations

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| DELETE | /api/url/{id} | [delete API URL ID](#delete-api-url-id) |  |
| GET | /api/url/{id} | [get API URL ID](#get-api-url-id) |  |
| POST | /api/url | [post API URL](#post-api-url) | Creates new short url |
  


## Paths

### <span id="delete-api-url-id"></span> delete API URL ID (*DeleteAPIURLID*)

```
DELETE /api/url/{id}
```

Deletes the given short URL

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | string | `string` |  | ✓ |  | id of short URL to delete |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#delete-api-url-id-200) | OK | successful operation |  | [schema](#delete-api-url-id-200-schema) |
| [404](#delete-api-url-id-404) | Not Found | Object not found |  | [schema](#delete-api-url-id-404-schema) |

#### Responses


##### <span id="delete-api-url-id-200"></span> 200 - successful operation
Status: OK

###### <span id="delete-api-url-id-200-schema"></span> Schema
   
  

[DeleteURLResponse](#delete-url-response)

##### <span id="delete-api-url-id-404"></span> 404 - Object not found
Status: Not Found

###### <span id="delete-api-url-id-404-schema"></span> Schema

### <span id="get-api-url-id"></span> get API URL ID (*GetAPIURLID*)

```
GET /api/url/{id}
```

Returns long URL for the given short one.

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| id | `path` | string | `string` |  | ✓ |  | id of short URL to return |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#get-api-url-id-200) | OK | successful operation |  | [schema](#get-api-url-id-200-schema) |

#### Responses


##### <span id="get-api-url-id-200"></span> 200 - successful operation
Status: OK

###### <span id="get-api-url-id-200-schema"></span> Schema
   
  

[GetURLRequest](#get-url-request)

### <span id="post-api-url"></span> Creates new short url (*PostAPIURL*)

```
POST /api/url
```

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| body | `body` | [GetURLRequest](#get-url-request) | `models.GetURLRequest` | | ✓ | | Long URL to be shortened and stored |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#post-api-url-200) | OK | successful operation |  | [schema](#post-api-url-200-schema) |

#### Responses


##### <span id="post-api-url-200"></span> 200 - successful operation
Status: OK

###### <span id="post-api-url-200-schema"></span> Schema
   
  

[][GetURLResponse](#get-url-response)

## Models

### <span id="create-url-request"></span> CreateUrlRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| longUrl | string| `string` |  | |  |  |



### <span id="create-url-response"></span> CreateUrlResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| status | [Status](#status)| `Status` |  | |  |  |
| urlInfo | [URLInfo](#url-info)| `URLInfo` |  | |  |  |



### <span id="delete-url-response"></span> DeleteUrlResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| shortId | string| `string` |  | |  |  |
| status | [Status](#status)| `Status` |  | |  |  |



### <span id="get-url-request"></span> GetUrlRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| shortId | string| `string` |  | |  |  |



### <span id="get-url-response"></span> GetUrlResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| source | string| `string` |  | |  |  |
| status | [Status](#status)| `Status` |  | |  |  |
| urlInfo | [URLInfo](#url-info)| `URLInfo` |  | |  |  |



### <span id="status"></span> Status


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| errorMsg | string| `string` |  | |  |  |
| success | boolean| `bool` |  | |  |  |



### <span id="url-info"></span> UrlInfo


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| created | string| `string` |  | |  |  |
| expires | string| `string` |  | |  |  |
| longUrl | string| `string` |  | |  |  |
| shortId | string| `string` |  | |  |  |


