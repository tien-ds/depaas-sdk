module signer

go 1.14

require golang.org/x/crypto v0.0.0-20200728195943-123391ffb6de // indirect

replace github.com/ethereum/go-ethereum v1.8.20 => github.com/loomnetwork/go-ethereum v1.8.17-0.20191122084538-6128fa1a8c76

require (
	gitee.com/aifuturewell/gojni v0.0.0-20210507105514-3201d9b6ae5d
	github.com/cespare/cp v1.1.1 // indirect
	github.com/ethereum/go-ethereum v1.8.20
	github.com/fjl/memsize v0.0.0-20190710130421-bcb5799ab5e5 // indirect
	github.com/gogo/protobuf v1.3.1
	github.com/huin/goupnp v1.0.0 // indirect
	github.com/jackpal/go-nat-pmp v1.0.2 // indirect
	github.com/karalabe/hid v1.0.0 // indirect
	github.com/loomnetwork/go-loom v0.0.0
	github.com/miguelmota/go-solidity-sha3 v0.1.0
	github.com/prometheus/prometheus v1.8.2 // indirect
	github.com/rjeczalik/notify v0.9.2 // indirect
	github.com/rs/cors v1.7.0 // indirect
	github.com/sirupsen/logrus v1.6.0
	gopkg.in/natefinch/npipe.v2 v2.0.0-20160621034901-c1b8fa8bdcce // indirect
	gopkg.in/urfave/cli.v1 v1.20.0 // indirect
)

replace github.com/loomnetwork/go-loom v0.0.0 => ../go-loom

replace github.com/phonkee/go-pubsub v0.0.0 => github.com/loomnetwork/go-pubsub v0.0.0-20180626134536-2d1454660ed1
