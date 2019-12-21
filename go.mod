module github.com/jingweno/upterm

go 1.13

require (
	github.com/anmitsu/go-shlex v0.0.0-20161002113705-648efa622239 // indirect
	github.com/creack/pty v1.1.9
	github.com/flynn/go-shlex v0.0.0-20150515145356-3f9db97f8568 // indirect
	github.com/gliderlabs/ssh v0.2.2
	github.com/go-openapi/errors v0.19.2
	github.com/go-openapi/runtime v0.19.5
	github.com/go-openapi/strfmt v0.19.3
	github.com/go-openapi/swag v0.19.5
	github.com/golang/protobuf v1.3.2
	github.com/google/go-cmp v0.3.1
	github.com/google/shlex v0.0.0-20181106134648-c34317bd91bf
	github.com/grpc-ecosystem/grpc-gateway v1.12.1
	github.com/oklog/run v1.0.0
	github.com/pborman/ansi v0.0.0-20160920233902-86f499584b0a
	github.com/rs/xid v1.2.1
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.5
	github.com/spf13/pflag v1.0.5 // indirect
	golang.org/x/crypto v0.0.0-20191029031824-8986dd9e96cf
	golang.org/x/net v0.0.0-20191209160850-c0dbc17a3553 // indirect
	golang.org/x/sys v0.0.0-20191210023423-ac6580df4449 // indirect
	google.golang.org/genproto v0.0.0-20191206224255-0243a4be9c8f
	google.golang.org/grpc v1.25.1
	gopkg.in/yaml.v2 v2.2.7 // indirect
)

replace (
	github.com/gliderlabs/ssh => github.com/jingweno/ssh v0.2.3-0.20191221201824-4cd54473e34e
	golang.org/x/crypto => github.com/jingweno/upterm.crypto v0.0.0-20191221200714-6f9467940236
)
