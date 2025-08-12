module ivpn.net/auth/services/generator

go 1.24.5

require (
	github.com/google/uuid v1.6.0
	github.com/jasonlvhit/gocron v0.0.1
	golang.org/x/time v0.12.0
	google.golang.org/grpc v1.73.0
	gorm.io/driver/mysql v1.6.0
	gorm.io/gorm v1.30.0
	ivpn.net/auth/services/proto v0.0.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/go-sql-driver/mysql v1.8.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	golang.org/x/net v0.38.0 // indirect
	golang.org/x/sys v0.31.0 // indirect
	golang.org/x/text v0.23.0 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20250324211829-b45e905df463 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
)

replace ivpn.net/auth/services/proto => ../../proto
