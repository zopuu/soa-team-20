module github.com/zopuu/soa-team-20/Backend/services/followers_service

go 1.24.6

require (
	github.com/neo4j/neo4j-go-driver/v5 v5.28.3
	github.com/zopuu/soa-team-20/common/obs v0.0.0
	google.golang.org/grpc v1.75.0
	google.golang.org/protobuf v1.36.8
)

replace github.com/zopuu/soa-team-20/common/obs => ../../common/obs

require (
	github.com/google/uuid v1.6.0 // indirect
	go.uber.org/multierr v1.10.0 // indirect
	go.uber.org/zap v1.27.0 // indirect
	golang.org/x/net v0.41.0 // indirect
	golang.org/x/sys v0.33.0 // indirect
	golang.org/x/text v0.26.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250707201910-8d1bb00bc6a7 // indirect
)
