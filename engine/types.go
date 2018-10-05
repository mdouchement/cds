package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/fatih/structs"

	"github.com/ovh/cds/engine/api"
	"github.com/ovh/cds/engine/api/observability"
	"github.com/ovh/cds/engine/elasticsearch"
	"github.com/ovh/cds/engine/hatchery/kubernetes"
	"github.com/ovh/cds/engine/hatchery/local"
	"github.com/ovh/cds/engine/hatchery/marathon"
	"github.com/ovh/cds/engine/hatchery/openstack"
	"github.com/ovh/cds/engine/hatchery/swarm"
	"github.com/ovh/cds/engine/hatchery/vsphere"
	"github.com/ovh/cds/engine/hooks"
	"github.com/ovh/cds/engine/migrateservice"
	"github.com/ovh/cds/engine/repositories"
	"github.com/ovh/cds/engine/vcs"
)

// Configuration contains CDS Configuration and toml description
type Configuration struct {
	Log struct {
		Level   string `toml:"level" default:"warning" comment:"Log Level: debug, info, warning, notice, critical" json:"level"`
		Graylog struct {
			Host       string `toml:"host" comment:"Example: thot.ovh.com" json:"host"`
			Port       int    `toml:"port" comment:"Example: 12202" json:"port"`
			Protocol   string `toml:"protocol" default:"tcp" comment:"tcp or udp" json:"protocol"`
			ExtraKey   string `toml:"extraKey" comment:"Example: X-OVH-TOKEN. You can use many keys: aaa,bbb" json:"extraKey"`
			ExtraValue string `toml:"extraValue" comment:"value for extraKey field. For many keys: valueaaa,valuebbb" json:"extraValue"`
		} `toml:"graylog"`
	} `toml:"log" comment:"#####################\n CDS Logs Settings \n####################"`
	Debug           *DebugConfiguration           `toml:"debug" comment:"#####################\n Debug with gops \n####################" json:"debug"`
	API             *api.Configuration            `toml:"api" comment:"#####################\n API Configuration \n####################" json:"api"`
	Hatchery        *HatcheryConfiguration        `toml:"hatchery" json:"hatchery"`
	Hooks           *hooks.Configuration          `toml:"hooks" comment:"######################\n CDS Hooks Settings \n######################" json:"hooks"`
	VCS             *vcs.Configuration            `toml:"vcs" comment:"######################\n CDS VCS Settings \n######################" json:"vcs"`
	Repositories    *repositories.Configuration   `toml:"repositories" comment:"######################\n CDS Repositories Settings \n######################" json:"repositories"`
	ElasticSearch   *elasticsearch.Configuration  `toml:"elasticsearch" comment:"######################\n CDS ElasticSearch Settings \n This is use for CDS timeline and is optional\n######################" json:"elasticsearch"`
	DatabaseMigrate *migrateservice.Configuration `toml:"databaseMigrate" comment:"######################\n CDS DB Migrate Service Settings \n######################" json:"databaseMigrate"`
	Tracing         *observability.Configuration  `toml:"tracing" comment:"###########################\n CDS Tracing Settings \n##########################" json:"tracing"`
}

// DebugConfiguration contains debug configuration
type DebugConfiguration struct {
	Enable         bool   `toml:"enable" default:"false" comment:"allow debugging with gops" json:"enable"`
	RemoteDebugURL string `toml:"remoteDebugURL" comment:"start a gops agent on specified URL. Ex: localhost:9999" json:"remoteDebugURL"`
}

// HatcheryConfiguration contains subsection of Hatchery configuration
type HatcheryConfiguration struct {
	Local      *local.HatcheryConfiguration      `toml:"local" comment:"Hatchery Local." json:"local"`
	Kubernetes *kubernetes.HatcheryConfiguration `toml:"kubernetes" comment:"Hatchery Kubernetes." json:"kubernetes"`
	Marathon   *marathon.HatcheryConfiguration   `toml:"marathon" comment:"Hatchery Marathon." json:"marathon"`
	Openstack  *openstack.HatcheryConfiguration  `toml:"openstack" comment:"Hatchery OpenStack. Doc: https://ovh.github.io/cds/advanced/advanced.hatcheries.openstack/" json:"openstack"`
	Swarm      *swarm.HatcheryConfiguration      `toml:"swarm" comment:"Hatchery Swarm. Doc: https://ovh.github.io/cds/advanced/advanced.hatcheries.swarm/" json:"swarm"`
	VSphere    *vsphere.HatcheryConfiguration    `toml:"vsphere" comment:"Hatchery VShpere. Doc: https://ovh.github.io/cds/advanced/advanced.hatcheries.vsphere/" json:"vshpere"`
}

// AsEnvVariables returns the object attributes as env variables. It used for configuration structs
func AsEnvVariables(o interface{}, prefix string, skipCommented bool) map[string]string {
	r := map[string]string{}
	prefix = strings.ToUpper(prefix)
	delim := "_"
	if prefix == "" {
		delim = ""
	}
	fields := structs.Fields(o)
	for _, f := range fields {
		if skipCommented {
			if commented, _ := strconv.ParseBool(f.Tag("commented")); commented {
				continue
			}
		}
		if structs.IsStruct(f.Value()) {
			rf := AsEnvVariables(f.Value(), prefix+delim+f.Name(), skipCommented)
			for k, v := range rf {
				r[k] = v
			}
		} else {
			r[prefix+"_"+strings.ToUpper(f.Name())] = fmt.Sprintf("%v", f.Value())
		}
	}
	return r
}
