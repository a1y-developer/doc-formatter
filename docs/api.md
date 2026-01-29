


# AI Doc Formatter API
API for AI Doc Formatter
  

## Informations

### Version

1.0

### Contact

  

## Content negotiation

### URI Schemes
  * http

### Consumes
  * application/json
  * multipart/form-data

### Produces
  * application/json

## All endpoints

###  auth

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| POST | /api/v1/auth/login | [login](#login) | Login |
| POST | /api/v1/auth/signup | [signup](#signup) | Signup |
  


###  storage

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| GET | /api/v1/storage/files | [list files by user Id](#list-files-by-user-id) | List files by user ID |
| POST | /api/v1/storage/upload | [upload file](#upload-file) | Upload file |
  


## Paths

### <span id="list-files-by-user-id"></span> List files by user ID (*ListFilesByUserId*)

```
GET /api/v1/storage/files
```

Retrieve a list of files uploaded by a specific user

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| user_id | `query` | string | `string` |  | ✓ |  | User ID (UUID) |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#list-files-by-user-id-200) | OK | Success |  | [schema](#list-files-by-user-id-200-schema) |
| [400](#list-files-by-user-id-400) | Bad Request | Bad Request |  | [schema](#list-files-by-user-id-400-schema) |
| [401](#list-files-by-user-id-401) | Unauthorized | Unauthorized |  | [schema](#list-files-by-user-id-401-schema) |
| [404](#list-files-by-user-id-404) | Not Found | Not Found |  | [schema](#list-files-by-user-id-404-schema) |
| [429](#list-files-by-user-id-429) | Too Many Requests | Too Many Requests |  | [schema](#list-files-by-user-id-429-schema) |
| [500](#list-files-by-user-id-500) | Internal Server Error | Internal Server Error |  | [schema](#list-files-by-user-id-500-schema) |

#### Responses


##### <span id="list-files-by-user-id-200"></span> 200 - Success
Status: OK

###### <span id="list-files-by-user-id-200-schema"></span> Schema
   
  

[ListFilesByUserIDOKBody](#list-files-by-user-id-o-k-body)

##### <span id="list-files-by-user-id-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="list-files-by-user-id-400-schema"></span> Schema
   
  

any

##### <span id="list-files-by-user-id-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="list-files-by-user-id-401-schema"></span> Schema
   
  

any

##### <span id="list-files-by-user-id-404"></span> 404 - Not Found
Status: Not Found

###### <span id="list-files-by-user-id-404-schema"></span> Schema
   
  

any

##### <span id="list-files-by-user-id-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="list-files-by-user-id-429-schema"></span> Schema
   
  

any

##### <span id="list-files-by-user-id-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="list-files-by-user-id-500-schema"></span> Schema
   
  

any

###### Inlined models

**<span id="list-files-by-user-id-o-k-body"></span> ListFilesByUserIDOKBody**


  


* composed type [HandlerResponse](#handler-response)
* inlined member (*listFilesByUserIdOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [][ResponseDocument](#response-document)| `[]*models.ResponseDocument` |  | |  |  |



### <span id="login"></span> Login (*Login*)

```
POST /api/v1/auth/login
```

Login user and return JWT token

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| body | `body` | [RequestLoginRequest](#request-login-request) | `models.RequestLoginRequest` | | ✓ | | Login payload |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [200](#login-200) | OK | Success |  | [schema](#login-200-schema) |
| [400](#login-400) | Bad Request | Bad Request |  | [schema](#login-400-schema) |
| [401](#login-401) | Unauthorized | Unauthorized |  | [schema](#login-401-schema) |
| [404](#login-404) | Not Found | Not Found |  | [schema](#login-404-schema) |
| [429](#login-429) | Too Many Requests | Too Many Requests |  | [schema](#login-429-schema) |
| [500](#login-500) | Internal Server Error | Internal Server Error |  | [schema](#login-500-schema) |

#### Responses


##### <span id="login-200"></span> 200 - Success
Status: OK

###### <span id="login-200-schema"></span> Schema
   
  

[LoginOKBody](#login-o-k-body)

##### <span id="login-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="login-400-schema"></span> Schema
   
  

any

##### <span id="login-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="login-401-schema"></span> Schema
   
  

any

##### <span id="login-404"></span> 404 - Not Found
Status: Not Found

###### <span id="login-404-schema"></span> Schema
   
  

any

##### <span id="login-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="login-429-schema"></span> Schema
   
  

any

##### <span id="login-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="login-500-schema"></span> Schema
   
  

any

###### Inlined models

**<span id="login-o-k-body"></span> LoginOKBody**


  


* composed type [HandlerResponse](#handler-response)
* inlined member (*loginOKBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [ResponseLoginResponse](#response-login-response)| `models.ResponseLoginResponse` |  | |  |  |



### <span id="signup"></span> Signup (*Signup*)

```
POST /api/v1/auth/signup
```

Create a new user account

#### Consumes
  * application/json

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| body | `body` | [RequestSignupRequest](#request-signup-request) | `models.RequestSignupRequest` | | ✓ | | Signup payload |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [201](#signup-201) | Created | Success |  | [schema](#signup-201-schema) |
| [400](#signup-400) | Bad Request | Bad Request |  | [schema](#signup-400-schema) |
| [401](#signup-401) | Unauthorized | Unauthorized |  | [schema](#signup-401-schema) |
| [404](#signup-404) | Not Found | Not Found |  | [schema](#signup-404-schema) |
| [429](#signup-429) | Too Many Requests | Too Many Requests |  | [schema](#signup-429-schema) |
| [500](#signup-500) | Internal Server Error | Internal Server Error |  | [schema](#signup-500-schema) |

#### Responses


##### <span id="signup-201"></span> 201 - Success
Status: Created

###### <span id="signup-201-schema"></span> Schema
   
  

[SignupCreatedBody](#signup-created-body)

##### <span id="signup-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="signup-400-schema"></span> Schema
   
  

any

##### <span id="signup-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="signup-401-schema"></span> Schema
   
  

any

##### <span id="signup-404"></span> 404 - Not Found
Status: Not Found

###### <span id="signup-404-schema"></span> Schema
   
  

any

##### <span id="signup-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="signup-429-schema"></span> Schema
   
  

any

##### <span id="signup-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="signup-500-schema"></span> Schema
   
  

any

###### Inlined models

**<span id="signup-created-body"></span> SignupCreatedBody**


  


* composed type [HandlerResponse](#handler-response)
* inlined member (*signupCreatedBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [ResponseSignUpResponse](#response-sign-up-response)| `models.ResponseSignUpResponse` |  | |  |  |



### <span id="upload-file"></span> Upload file (*UploadFile*)

```
POST /api/v1/storage/upload
```

Upload a file for a user

#### Consumes
  * multipart/form-data

#### Produces
  * application/json

#### Parameters

| Name | Source | Type | Go type | Separator | Required | Default | Description |
|------|--------|------|---------|-----------| :------: |---------|-------------|
| file | `formData` | file | `io.ReadCloser` |  | ✓ |  | File to upload |
| user_id | `formData` | string | `string` |  | ✓ |  | User ID (UUID) |

#### All responses
| Code | Status | Description | Has headers | Schema |
|------|--------|-------------|:-----------:|--------|
| [201](#upload-file-201) | Created | Success |  | [schema](#upload-file-201-schema) |
| [400](#upload-file-400) | Bad Request | Bad Request |  | [schema](#upload-file-400-schema) |
| [401](#upload-file-401) | Unauthorized | Unauthorized |  | [schema](#upload-file-401-schema) |
| [404](#upload-file-404) | Not Found | Not Found |  | [schema](#upload-file-404-schema) |
| [429](#upload-file-429) | Too Many Requests | Too Many Requests |  | [schema](#upload-file-429-schema) |
| [500](#upload-file-500) | Internal Server Error | Internal Server Error |  | [schema](#upload-file-500-schema) |

#### Responses


##### <span id="upload-file-201"></span> 201 - Success
Status: Created

###### <span id="upload-file-201-schema"></span> Schema
   
  

[UploadFileCreatedBody](#upload-file-created-body)

##### <span id="upload-file-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="upload-file-400-schema"></span> Schema
   
  

any

##### <span id="upload-file-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="upload-file-401-schema"></span> Schema
   
  

any

##### <span id="upload-file-404"></span> 404 - Not Found
Status: Not Found

###### <span id="upload-file-404-schema"></span> Schema
   
  

any

##### <span id="upload-file-429"></span> 429 - Too Many Requests
Status: Too Many Requests

###### <span id="upload-file-429-schema"></span> Schema
   
  

any

##### <span id="upload-file-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="upload-file-500-schema"></span> Schema
   
  

any

###### Inlined models

**<span id="upload-file-created-body"></span> UploadFileCreatedBody**


  


* composed type [HandlerResponse](#handler-response)
* inlined member (*uploadFileCreatedBodyAO1*)



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| data | [ResponseDocument](#response-document)| `models.ResponseDocument` |  | |  |  |



## Models

### <span id="handler-response"></span> handler.Response


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| costTime | string| `string` |  | | Time taken for the request. |  |
| data | [any](#any)| `any` |  | | Data payload. |  |
| endTime | string| `string` |  | | Request end time. |  |
| message | string| `string` |  | | Descriptive message. |  |
| startTime | string| `string` |  | | Request start time. |  |
| success | boolean| `bool` |  | | Indicates success status. |  |
| traceID | string| `string` |  | | Trace identifier. |  |



### <span id="request-login-request"></span> request.LoginRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| email | string| `string` | ✓ | |  |  |
| password | string| `string` | ✓ | |  |  |



### <span id="request-signup-request"></span> request.SignupRequest


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| email | string| `string` | ✓ | |  |  |
| password | string| `string` | ✓ | |  |  |



### <span id="response-document"></span> response.Document


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| created_at | integer| `int64` |  | |  |  |
| file_name | string| `string` |  | |  |  |
| file_size | integer| `int64` |  | |  |  |
| id | string| `string` |  | |  |  |
| updated_at | integer| `int64` |  | |  |  |
| user_id | string| `string` |  | |  |  |



### <span id="response-login-response"></span> response.LoginResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| access_token | string| `string` |  | |  |  |
| expiry_unix | integer| `int64` |  | |  |  |



### <span id="response-sign-up-response"></span> response.SignUpResponse


  



**Properties**

| Name | Type | Go type | Required | Default | Description | Example |
|------|------|---------|:--------:| ------- |-------------|---------|
| user_id | string| `string` |  | |  |  |


