package vmess

import (
	"context"
	"errors"
	"net"
	"net/http"
	"strings"

	"github.com/xxf098/lite-proxy/config"
	"github.com/xxf098/lite-proxy/outbound"

	C "github.com/xxf098/lite-proxy/constant"
)

type Dialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

type DialerFunc func(ctx context.Context, network, address string) (net.Conn, error)

func (fn DialerFunc) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return fn(ctx, network, address)
}

type error_dialer struct {
	err error
}

func (e *error_dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	return nil, e.err
}

type dialer struct {
	d outbound.Dialer
}

func (d *dialer) DialContext(ctx context.Context, network, address string) (net.Conn, error) {
	i := strings.LastIndex(address, ":")
	if i < 0 {
		return nil, errors.New("address error: " + address)
	}

	n := C.TCP
	if network == "udp" {
		n = C.UDP
	}
	meta := &C.Metadata{
		NetWork:  n,
		Type:     C.SOCKS,
		SrcPort:  "",
		DstPort:  address[i+1:],
		AddrType: 3,
		Host:     address[:i],
	}
	return d.d.DialContext(ctx, meta)
}

// link MUST start with vmess:// or ss:// or ssr:// or trogan://
func New(link string) Dialer {
	d, err := config.Link2Dialer(link)
	if err != nil {
		return &error_dialer{err: err}
	}

	return &dialer{d: d}
}

func NewClient(link string) *http.Client {
	dialer := New(link)
	return &http.Client{
		Transport: &http.Transport{
			DialContext: dialer.DialContext,
		},
	}
}
