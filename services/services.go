package services

import "time"

type ServiceStatus struct {
	Service   *ServiceConfig
	Status    string
	Pid       int
	StartTime time.Time
	Ports     []string
}

type ServiceOrGroup interface {
	GetName() string
	Build() error
	Start() error
	Stop() error
	Status() ([]ServiceStatus, error)
	IsSudo() bool
	Watch() ([]ServiceWatch, error)
}
