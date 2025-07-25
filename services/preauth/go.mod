module ivpn.net/auth/services/preauth

go 1.24.5

require (
	github.com/redis/go-redis/v9 v9.11.0
	google.golang.org/grpc v1.73.0
	ivpn.net/auth/services/proto v0.0.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace ivpn.net/auth/services/proto => ../../proto
