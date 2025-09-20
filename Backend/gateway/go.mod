module github.com/zopuu/soa-team-20/Backend/gateway

go 1.24.6

require (
	github.com/didip/tollbooth/v7 v7.0.2
	github.com/didip/tollbooth_chi v0.0.0-20250112173903-88de5e56a7cc
	github.com/go-chi/chi/v5 v5.0.12
	github.com/go-chi/cors v1.2.1
	github.com/golang-jwt/jwt/v5 v5.2.1
	github.com/zopuu/soa-team-20/Backend/services/followers_service v0.0.0-00010101000000-000000000000
	github.com/zopuu/soa-team-20/Backend/services/shopping_service v0.0.0-00010101000000-000000000000
	github.com/zopuu/soa-team-20/Backend/services/tour v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.75.0
)

require (
	github.com/go-pkgz/expirable-cache/v3 v3.0.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
	google.golang.org/protobuf v1.36.8 // indirect
)

replace github.com/zopuu/soa-team-20/Backend/services/followers_service => ../services/followers_service
replace github.com/zopuu/soa-team-20/Backend/services/shopping_service => ../services/shopping_service

replace github.com/zopuu/soa-team-20/Backend/services/tour => ../services/tour
