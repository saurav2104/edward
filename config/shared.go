package config

import (
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/pkg/errors"
	"github.com/yext/edward/services"
)

var groupMap map[string]*services.ServiceGroupConfig
var serviceMap map[string]*services.ServiceConfig

func GetServiceMap() map[string]*services.ServiceConfig {
	return serviceMap
}

func LoadSharedConfig(configPath string, edwardVersion string, logger *log.Logger) error {
	InitEmptyConfig()
	if configPath != "" {
		r, err := os.Open(configPath)
		if err != nil {
			return errors.WithStack(err)
		}
		cfg, err := LoadConfigWithDir(r, filepath.Dir(configPath), edwardVersion, logger)
		if err != nil {
			workingDir, _ := os.Getwd()
			configRel, _ := filepath.Rel(workingDir, configPath)
			return errors.WithMessage(err, configRel)
		}

		serviceMap = cfg.ServiceMap
		groupMap = cfg.GroupMap
		return nil
	}

	return errors.New("No config file found")

}

func GetServicesOrGroups(names []string) ([]services.ServiceOrGroup, error) {
	var outSG []services.ServiceOrGroup
	for _, name := range names {
		sg, err := GetServiceOrGroup(name)
		if err != nil {
			return nil, err
		}
		outSG = append(outSG, sg)
	}
	return outSG, nil
}

func GetServiceOrGroup(name string) (services.ServiceOrGroup, error) {
	if group, ok := groupMap[name]; ok {
		return group, nil
	}
	if service, ok := serviceMap[name]; ok {
		return service, nil
	}
	return nil, errors.New("Service or group not found")
}

func GetAllServiceNames() []string {
	var serviceNames []string
	for name := range serviceMap {
		serviceNames = append(serviceNames, name)
	}
	return serviceNames
}

func GetAllGroupNames() []string {
	var groupNames []string
	for name := range groupMap {
		groupNames = append(groupNames, name)
	}
	return groupNames
}

func GetAllServicesSorted() []services.ServiceOrGroup {
	var as []services.ServiceOrGroup
	for _, service := range serviceMap {
		as = append(as, service)
	}
	sort.Sort(serviceOrGroupByName(as))
	return as
}

type serviceOrGroupByName []services.ServiceOrGroup

func (s serviceOrGroupByName) Len() int {
	return len(s)
}
func (s serviceOrGroupByName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s serviceOrGroupByName) Less(i, j int) bool {
	return s[i].GetName() < s[j].GetName()
}

func InitEmptyConfig() {
	groupMap = make(map[string]*services.ServiceGroupConfig)
	serviceMap = make(map[string]*services.ServiceConfig)
}
