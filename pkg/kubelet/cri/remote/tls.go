package remote

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// RemoteTLSOptions remote cri tls options
type RemoteTLSOptions struct {
	Enabled   bool
	TLSRootCA string
	TLSCert   string
	TLSKey    string
}

const (
	DefaulCriRemoteTLSRootCA = "/etc/kubernetes/pki/ca.crt"
	DefaulCriRemoteTLSCert   = "/etc/kubernetes/pki/kubelet-cri-client.crt"
	DefaulCriRemoteTLSKey    = "/etc/kubernetes/pki/kubelet-cri-client.key"

	defaultServerName = "containerd"
)

func defaultRemoteTLSOptions() *RemoteTLSOptions {
	return &RemoteTLSOptions{
		Enabled:   true,
		TLSRootCA: DefaulCriRemoteTLSRootCA,
		TLSCert:   DefaulCriRemoteTLSCert,
		TLSKey:    DefaulCriRemoteTLSKey,
	}
}

func (r RemoteTLSOptions) GrpcDialOption() grpc.DialOption {
	if !r.Enabled {
		return grpc.WithInsecure()
	}

	certificate, err := tls.LoadX509KeyPair(r.TLSCert, r.TLSKey)
	if err != nil {
		panic("Load client certification failed: " + err.Error())
	}

	ca, err := ioutil.ReadFile(r.TLSRootCA)
	if err != nil {
		panic("can't read ca file: " + err.Error())
	}

	capool := x509.NewCertPool()
	capool.AppendCertsFromPEM(ca)

	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{certificate},
		RootCAs:      capool,
		ServerName:   defaultServerName,
	}

	cred := credentials.NewTLS(tlsConfig)

	return grpc.WithTransportCredentials(cred)
}
