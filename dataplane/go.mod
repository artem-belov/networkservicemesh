module github.com/networkservicemesh/networkservicemesh/dataplane

require (
	github.com/codahale/hdrhistogram v0.0.0-20161010025455-3a0bb77429bd // indirect
	github.com/gogo/protobuf v1.2.2-0.20190723190241-65acae22fc9d
	github.com/golang/protobuf v1.3.2
	github.com/ligato/cn-infra v2.0.0+incompatible
	github.com/ligato/vpp-agent v2.1.1+incompatible
	github.com/networkservicemesh/networkservicemesh v0.1.0
	github.com/networkservicemesh/networkservicemesh/controlplane v0.1.0
	github.com/networkservicemesh/networkservicemesh/controlplane/api v0.1.0
	github.com/networkservicemesh/networkservicemesh/dataplane/api v0.1.0
	github.com/onsi/gomega v1.5.1-0.20190520121345-efe19c39ca10
	github.com/opentracing/opentracing-go v1.1.0
	github.com/rs/xid v1.2.1
	github.com/sirupsen/logrus v1.4.2
	github.com/uber-go/atomic v1.4.0 // indirect
	github.com/vishvananda/netlink v1.0.0
	github.com/vishvananda/netns v0.0.0-20190625233234-7109fa855b0f
	go.uber.org/atomic v1.4.0 // indirect
	google.golang.org/appengine v1.4.0 // indirect
	google.golang.org/grpc v1.23.0
)

replace (
	github.com/networkservicemesh/networkservicemesh => ../
	github.com/networkservicemesh/networkservicemesh/controlplane => ../controlplane
	github.com/networkservicemesh/networkservicemesh/controlplane/api => ../controlplane/api
	github.com/networkservicemesh/networkservicemesh/dataplane => ./
	github.com/networkservicemesh/networkservicemesh/dataplane/api => ./api
	github.com/networkservicemesh/networkservicemesh/k8s/api => ../k8s/api
	github.com/networkservicemesh/networkservicemesh/sdk => ../sdk
	github.com/networkservicemesh/networkservicemesh/side-cars => ../side-cars
)
