module github.com/LanceLRQ/deer-executor/v2

go 1.16

require (
	github.com/LanceLRQ/deer-common v0.0.7
	github.com/kr/pretty v0.1.0 // indirect
	github.com/pkg/errors v0.9.1
	github.com/satori/go.uuid v1.2.0
	github.com/urfave/cli/v2 v2.2.0
	gopkg.in/check.v1 v1.0.0-20190902080502-41f04d3bba15 // indirect
)

//replace github.com/LanceLRQ/deer-common => ./pkg/deer-common
