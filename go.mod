module github.com/esakat/trimg

go 1.13

// https://github.com/golang/go/issues/28489#issuecomment-524088969
replace github.com/Sirupsen/logrus v1.4.2 => github.com/sirupsen/logrus v1.4.2

require (
	github.com/Azure/go-ansiterm v0.0.0-20170929234023-d6e3b3328b78 // indirect
	github.com/Microsoft/go-winio v0.4.14 // indirect
	github.com/Sirupsen/logrus v1.4.2 // indirect
	github.com/VividCortex/ewma v1.1.1 // indirect
	github.com/aws/aws-sdk-go v1.25.19
	github.com/docker/distribution v2.7.1+incompatible
	github.com/docker/docker v1.13.1
	github.com/docker/go-connections v0.4.0 // indirect
	github.com/docker/go-units v0.4.0 // indirect
	github.com/kr/pretty v0.1.0 // indirect
	github.com/mattn/go-isatty v0.0.11 // indirect
	github.com/opencontainers/go-digest v1.0.0-rc1 // indirect
	github.com/pkg/errors v0.8.1
	github.com/spf13/cobra v0.0.5
	github.com/stretchr/testify v1.3.0 // indirect
	github.com/vbauerster/mpb v3.4.0+incompatible
	golang.org/x/net v0.0.0-20190613194153-d28f0bde5980 // indirect
	gopkg.in/check.v1 v1.0.0-20180628173108-788fd7840127 // indirect
	gopkg.in/yaml.v2 v2.2.2
)
