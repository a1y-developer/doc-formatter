


# AI Doc Formatter API Gateway
API Gateway for AI Doc Formatter
  

## Informations

### Version

1.0

### Contact

  

## Content negotiation

### URI Schemes
  * http

### Consumes
  * application/json

### Produces
  * application/json

## All endpoints

###  auth

| Method  | URI     | Name   | Summary |
|---------|---------|--------|---------|
| POST | /api/v1/auth/login | [post API v1 auth login](#post-api-v1-auth-login) | Login |
| POST | /api/v1/auth/signup | [post API v1 auth signup](#post-api-v1-auth-signup) | Signup |
  


## Paths

### <span id="post-api-v1-auth-login"></span> Login (*PostAPIV1AuthLogin*)

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
| [200](#post-api-v1-auth-login-200) | OK | OK |  | [schema](#post-api-v1-auth-login-200-schema) |
| [400](#post-api-v1-auth-login-400) | Bad Request | Bad Request |  | [schema](#post-api-v1-auth-login-400-schema) |
| [401](#post-api-v1-auth-login-401) | Unauthorized | Unauthorized |  | [schema](#post-api-v1-auth-login-401-schema) |
| [500](#post-api-v1-auth-login-500) | Internal Server Error | Internal Server Error |  | [schema](#post-api-v1-auth-login-500-schema) |

#### Responses


##### <span id="post-api-v1-auth-login-200"></span> 200 - OK
Status: OK

###### <span id="post-api-v1-auth-login-200-schema"></span> Schema
   
  

[ResponseLoginResponse](#response-login-response)

##### <span id="post-api-v1-auth-login-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-api-v1-auth-login-400-schema"></span> Schema
   
  

map of string

##### <span id="post-api-v1-auth-login-401"></span> 401 - Unauthorized
Status: Unauthorized

###### <span id="post-api-v1-auth-login-401-schema"></span> Schema
   
  

map of string

##### <span id="post-api-v1-auth-login-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="post-api-v1-auth-login-500-schema"></span> Schema
   
  

map of string

### <span id="post-api-v1-auth-signup"></span> Signup (*PostAPIV1AuthSignup*)

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
| [201](#post-api-v1-auth-signup-201) | Created | Created |  | [schema](#post-api-v1-auth-signup-201-schema) |
| [400](#post-api-v1-auth-signup-400) | Bad Request | Bad Request |  | [schema](#post-api-v1-auth-signup-400-schema) |
| [500](#post-api-v1-auth-signup-500) | Internal Server Error | Internal Server Error |  | [schema](#post-api-v1-auth-signup-500-schema) |

#### Responses


##### <span id="post-api-v1-auth-signup-201"></span> 201 - Created
Status: Created

###### <span id="post-api-v1-auth-signup-201-schema"></span> Schema
   
  

[ResponseSignUpResponse](#response-sign-up-response)

##### <span id="post-api-v1-auth-signup-400"></span> 400 - Bad Request
Status: Bad Request

###### <span id="post-api-v1-auth-signup-400-schema"></span> Schema
   
  

map of string

##### <span id="post-api-v1-auth-signup-500"></span> 500 - Internal Server Error
Status: Internal Server Error

###### <span id="post-api-v1-auth-signup-500-schema"></span> Schema
   
  

map of string

## Models

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


