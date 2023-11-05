package utils

import (
	"fmt"

	"github.com/hashicorp/consul/api"
)

type Registry struct {
	host   string
	port   int
	client *api.Client
}

type RegistryClient interface {
	Register(ipAddr string, port int, name string, tags []string, id string) error
	Deregister(serviceId string) error
}

func NewRegistryClient(host string, port int) RegistryClient {
	cfg := api.DefaultConfig()
	cfg.Address = fmt.Sprintf("%s:%d", host, port)

	client, err := api.NewClient(cfg)
	if err != nil {
		panic(err)
	}
	return &Registry{
		host:   host,
		port:   port,
		client: client,
	}
}

func (r *Registry) Register(ipAddr string, port int, name string, tags []string, id string) error {
	// 生成注册对象
	registration := new(api.AgentServiceRegistration)
	registration.Address = ipAddr
	registration.Port = port
	registration.Name = name
	registration.Tags = tags
	registration.ID = id
	// 配置 gRPC 服务健康检查
	registration.Check = &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", ipAddr, port),
		Timeout:                        "5s",
		Interval:                       "5s",
		DeregisterCriticalServiceAfter: "15s",
	}

	// 注册服务
	return r.client.Agent().ServiceRegister(registration)
}

func (r *Registry) Deregister(serviceID string) error {
	return r.client.Agent().ServiceDeregister(serviceID)
}
