package keepalive

import (
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

type NewClient struct {
	Maxconn        int
	MaxDialTimeOut int
	RepTimeOut     int
	client         *http.Client
}

func (cfg *NewClient) GetClient() *http.Client {
	if cfg.client == nil {
		cfg.client = &http.Client{
			Transport: &http.Transport{
				ResponseHeaderTimeout: time.Millisecond * time.Duration(cfg.RepTimeOut),
				MaxIdleConnsPerHost:   cfg.Maxconn,
				Dial: func(netw, addr string) (net.Conn, error) {
					conn, err := net.DialTimeout(netw, addr, time.Millisecond*time.Duration(cfg.MaxDialTimeOut))
					if err == nil {
						if tcpConn, ok := conn.(*net.TCPConn); ok {
							tcpConn.SetLinger(0)
						}
					}

					return conn, err
				},
			},
		}
	}

	return cfg.client
}

func (cfg *NewClient) Request(req *http.Request) ([]byte, error) {
	resp, err := cfg.GetClient().Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
