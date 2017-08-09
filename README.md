# REST API Skeleton written in golang

This is a very simple version supporting common features used in REST API implementations. 
Can be used as a start for creating more advanced versions.

## Installation
```
go get github.com/cseeger-epages/rest-api-go-skeleton
```
or simply clone the repo
```
git clone https://github.com/cseeger-epages/rest-api-go-skeleton
```

## Configuration
add users by adding
```
[[user]]
username = "<username>"
password = "<password>"
```

to conf/api.conf

change database settings in
conf/api.conf section `[database]`

customize your Database functions in src/Database.go for your needs

## build and run
you can build the binary via
```
./build.sh
```
and the following flags are supported
```
-crt  <certificate file>
-key  <certificate key file>
-c    <config file>
```

## Further Implementation
add your custom Handlers to src/Handler.go and add them to src/Routes.go with the following pattern
```
Route{
  "<some route name>",
  "<request method>",
  "/yourcustomroute",
  "<description used for /help/[cmd]"
  "<handler name>"
}
```

a default Handler template

```
func Handler(w http.ResponseWriter, r*http.Request) {
        // caching stuff is handler specific
        w.Header().Set("Cache-Control", "no-store")
        // add more hanlder specific Headers here or create a wrapper function
  
        // used for filter parameters, where qs (QueryStrings) hold these parameters
        qs := ParseQueryStrings(r)

        // your code here
        msg := HelpMsg{Message: "im a default Handler"}

        // ... the name speaks for itself
        EncodeAndSend(w, r, qs, msg)
}
```

## Supported Features
- path routing using gorilla mux
- versioning
- Database wrapper 
- TLS
- pretty print
- Etag / If-None-Match Clientside caching
- Rate limiting and headers using trottled middleware
- basic auth
- config using TOML format
- error handler
- logging

## Not (yet) implemented
- X-HTTP-Method-Override
- caching serverside (varnish ?)
- Authentication - oauth(2) 

## additional Notes
The Etag implementation lacks a bit of efficency since the data needs to be generated all the time to generate the Etag hash and it only saves some amount of traffic. 
When you are implementing your data think about how you can generate the Etag on using lesser resources.

Also only the basic Auth is implemented yet wich isn't very efficient, maybe some OAuth2 will be implemented later. 
If you send less data via many requests an OAuth implementation will save some bandwidth.

X-HTTP-Method-Override should be implemented when using more than GET or POST methods, e.g when you implement PUT, PATCH or DELETE method to support all kinds of clients.

## Ratelimit Headers
```
X-Ratelimit-Limit - The number of allowed requests in the current period
X-Ratelimit-Remaining - The number of remaining requests in the current period
X-Ratelimit-Reset - The number of seconds left in the current period
```

## generate certificates
```
cd certs
# Key considerations for algorithm "RSA" ≥ 2048-bit
openssl genrsa -out server.key 2048

# Key considerations for algorithm "ECDSA" ≥ secp384r1
# List ECDSA the supported curves (openssl ecparam -list_curves)
openssl ecparam -genkey -name secp384r1 -out server.key
openssl req -new -x509 -sha256 -key server.key -out server.crt -days 3650
```

## some curls
```
curl -k -X GET -H "Authorization: Basic dGVzdHVzZXI6dGVzdHBhc3MK" https://localhost:8443/
curl -k -X GET -H "Authorization: Basic dGVzdHVzZXI6dGVzdHBhc3MK" https://localhost:8443/help
curl -k -X GET -H "Authorization: Basic dGVzdHVzZXI6dGVzdHBhc3MK" https://localhost:8443/projects
curl -k -X GET -H "Authorization: Basic dGVzdHVzZXI6dGVzdHBhc3MK" https://localhost:8443/project/[name|id]

#pretify test
curl -v -H "Authorization: Basic dGVzdHVzZXI6dGVzdHBhc3MK" -k -X GET https://localhost:8443/projects\?prettify
#etag test
curl -v -H "Authorization: Basic dGVzdHVzZXI6dGVzdHBhc3MK" -H "If-None-Match: <some etag>" -k -X GET https://localhost:8443/projects\?prettify
```

## basic auth test stuff
```
testuser:testpass - dGVzdHVzZXI6dGVzdHBhc3MK
username:password - dXNlcm5hbWU6cGFzc3dvcmQK
```
