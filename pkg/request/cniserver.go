package request

import (
	"context"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"net"
	"net/http"
)

type CniServerClient struct {
	*gorequest.SuperAgent
}

type PodRequest struct {
	PodName      string `json:"pod_name"`
	PodNamespace string `json:"pod_namespace"`
	ContainerID  string `json:"container_id"`
	NetNs        string `json:"net_ns"`
}

type PodResponse struct {
	IpAddress  string `json:"address"`
	MacAddress string `json:"mac_address"`
	CIDR       string `json:"cidr"`
	Gateway    string `json:"gateway"`
	Mtu        int    `json:"mtu"`
}

func NewCniServerClient(socketAddress string) CniServerClient {
	request := gorequest.New()
	request.Transport = &http.Transport{DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
		return net.Dial("unix", socketAddress)
	}}
	return CniServerClient{request}
}

func (csc CniServerClient) Add(podRequest PodRequest) (*PodResponse, error) {
	resp := PodResponse{}
	res, body, errors := csc.Post("http://dummy/api/v1/add").Send(podRequest).EndStruct(&resp)
	if len(errors) != 0 {
		return nil, errors[0]
	}
	if res.StatusCode != 200 {
		return nil, fmt.Errorf("request ip return %d %s", res.StatusCode, body)
	}
	return &resp, nil
}

func (csc CniServerClient) Del(podRequest PodRequest) error {
	res, body, errors := csc.Post("http://dummy/api/v1/del").Send(podRequest).End()
	if len(errors) != 0 {
		return errors[0]
	}
	if res.StatusCode != 204 {
		return fmt.Errorf("delete ip return %d %s", res.StatusCode, body)
	}
	return nil
}