# MTAuthServer
Micro JWT Authentication Server against Active Directory / LDAP written in GOLANG

The authentication server provide a interface to authenticate against active directory or ldap and gives back all group memberships.

You should use this microservice only with HTTPS!

Have fun :-)

## Authenticate

### CURL

```shell
curl -X 'POST' \
  'https://localhost:8080/api/v1/auth' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "password": "password",
  "username": "username"
}'
```

### Response

```json
{
  "groups": [
    "p_user"
  ],
  "status": "ok",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTA2OTQ0NDQsImdyb3VwcyI6IltcInBfdXNlclwiXSIsInVzZXJuYW1lIjoibXRob21hIn0.qY-BiVf_R-PIrlFRqjTdPjxtFvR_wMRPg49T9UGN0sU",
  "username": "username"
}
```

## Check if the token is valid

### CURL

```shell
curl -X 'GET' \
  'http://localhost:8080/api/v1/verify/TESTAPP' \
  -H 'accept: application/json' \
  -H 'Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTA2OTQ0NDQsImdyb3VwcyI6IltcInBfdXNlclwiXSIsInVzZXJuYW1lIjoibXRob21hIn0.qY-BiVf_R-PIrlFRqjTdPjxtFvR_wMRPg49T9UGN0sU'
```

### Response


```json
{
  "app_name": "TESTAPP",
  "groups": [
    "p_user"
  ],
  "status": "valid",
  "username": "username"
}
```

## Use Introspect

### CURL

```shell
curl -X 'POST' \
  'http://localhost:8080/api/v1/introspect' \
  -H 'accept: application/json' \
  -H 'Content-Type: application/json' \
  -d '{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTA2OTQ0NDQsImdyb3VwcyI6IltcInBfdXNlclwiXSIsInVzZXJuYW1lIjoibXRob21hIn0.qY-BiVf_R-PIrlFRqjTdPjxtFvR_wMRPg49T9UGN0sU"
}'
```

### Response

```json
{
  "groups": [
    "p_user"
  ],
  "status": "valid",
  "username": "username"
}
```

