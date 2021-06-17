package config

import (
	"encoding/json"
	"errors"
	bootconsul "github.com/al8n/micro-boot/consul"
	bootflag "github.com/al8n/micro-boot/flag"
	bootredis "github.com/al8n/micro-boot/goredis"
	bootgrpc "github.com/al8n/micro-boot/grpc"
	boothttp "github.com/al8n/micro-boot/http"
	bootmongo "github.com/al8n/micro-boot/mongo"
	bootprom "github.com/al8n/micro-boot/prometheus"
	"github.com/al8n/kit-auth/authentication/common"
	"github.com/imdario/mergo"
	"github.com/sony/sonyflake"
	"gopkg.in/yaml.v3"
	"io/ioutil"
	"path/filepath"
	"sync"
)

const (
	HS256 = 256
	HS384 = 384
	HS512 = 512
)

var (
	config    *Config
	MachineID uint64
	once sync.Once
)

func GetConfig() *Config {
	once.Do(func() {
		config = &Config{}
	})
	return config
}


type Config struct {
	Host  string `json:"host" yaml:"host"`

	// RPC
	RPC bootgrpc.GRPC `json:"rpc" yaml:"rpc"`

	// HTTP
	HTTP boothttp.HTTP `json:"http" yaml:"http"`

	// HTTPS
	HTTPS boothttp.HTTPS `json:"https" yaml:"https"`

	// Authentication
	Service Register `json:"service" yaml:"service"`

	// Mongo
	Mongo bootmongo.ClientOptions `json:"mongo" yaml:"mongo"`

	// Consul
	Consul bootconsul.Config       `json:"consul" yaml:"consul"`

	// Prometheus
	Prom  bootprom.Config          `json:"prometheus" yaml:"prometheus"`

	// Redis
	Redis bootredis.ClientOptions      `json:"redis" yaml:"redis"`
}

func (c *Config) Initialize(name string) (err error)  {
	var (
		extension string
		file []byte
		newCfg Config
		sf *sonyflake.Sonyflake
	)

	sf = sonyflake.NewSonyflake(sonyflake.Settings{})
	MachineID, err = sf.NextID()
	if err != nil {
		return err
	}

	file, err = ioutil.ReadFile(name)
	if err != nil {
		return err
	}

	extension = filepath.Ext(name)
	switch extension {
	case ".yaml", ".yml":
		err = yaml.Unmarshal(file, &newCfg)
		if err != nil {
			return err
		}
	case ".json":
		err = json.Unmarshal(file, &newCfg)
		if err != nil {
			return err
		}
	default:
		return errors.New("unsupported config file type")
	}

	err = mergo.Merge(config, &newCfg)
	if err != nil {
		return err
	}

	if !config.HTTP.Runnable && !config.HTTPS.Runnable && !config.RPC.Runnable {
		return common.ErrorNoServicesConfig
	}

	err = config.Parse()
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) BindFlags(fs *bootflag.FlagSet)  {
	fs.StringVar(&c.Host, "host", "service host", "the service host")
	c.HTTP.BindFlags(fs)
	c.HTTPS.BindFlags(fs)
	c.RPC.BindFlags(fs)
	c.Service.BindFlags(fs)
	c.Mongo.BindFlags(fs)
	c.Consul.BindFlags(fs)
	c.Prom.BindFlags(fs)
}

func (c *Config) Parse() (err error) {
	err = c.Mongo.Parse()
	if err != nil {
		return err
	}

	err = c.Service.Parse()
	if err != nil {
		return err
	}
	return nil
}

