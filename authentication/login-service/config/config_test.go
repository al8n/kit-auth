package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"testing"
)

const (
	mockSecretKey = "R541QRVVTFFGZ2APJDHFSBBF9DMO6XU9PMBQ3C4CPDH4EII86PC9U5DVTFELM3VNK0OWYLRIDM7ROXGCGF84KVCQPNQK71BJC5PEAL4K7CU8XW8AMVKQL0X33HGOF49FLDA8DR2HEDEG4PMZ2RCK3WVI3LCMT7SM5WSEUW7C1R56NHDOHGN8LR7RG0J02KN178PLVVPM5SI84LZ371VDX24ER7SQWNRXWLMKYY5AJCS0YQ91HKB8CT13S1PG7R89IDPHYCW7CPMZGYQSCWAO9J9VJK5CF1C7MEEBH2SL7CGBNZHSCTRPUSI9U0S0N8IZG3I36QWTJ7KRTMECGOVC6WAVVPU9OW7BTR1XYJ3Y43RDG7Q38E831DK9HOS9X8ZEF1LVDNLT1JUJLN3PGRJN5EHLLMGNC2BYAXE7A9QDPHK6U9KMYRSBFBKOEAA6UDTPVMX41ZWOVW5JI2B0GZG2A51IO7OS0I8SW7RCAO8H01TJRR72M6AUAMPCLFU1ZZRW"
)

var (
	config1 = fmt.Sprintf(`
authentication:
  secret: "%s"
  method: HS256
  exp-at: 30d
http:
  port: 8080
  read-timeout: 15s
  write-timeout: 15s
https:
  port: 8043
  read-timeout: 15s
  write-timeout: 15s
  pem: ""
  key: ""
rpc:
  port: 50051
mongo:
  mode: tcp
  user: uni
  password: secret
  hosts: [{host: localhost, port: 27017}, {host: 192.168.0.7, port: 27017}]
  options:
    auth: true
    authSoource: admin
  db: uni
  collection: users
  timeout: 30s
`, mockSecretKey)
	config2 = fmt.Sprintf(`
authentication:
  secret: "%s"
  method: HS256
  exp-at: 30d
http:
  port: 8080
  read-timeout: 15s
  write-timeout: 15s
https:
  port: 8043
  read-timeout: 15s
  write-timeout: 15s
  pem: ""
  key: ""
rpc:
  port: 50051
mongo:
  mode: atlas
  user: uni
  password: secret
  hosts: [{host: localhost, port: 27017}, {host: 192.168.0.7, port: 27017}]
  options:
    auth: true
    authSoource: admin
  db: uni
  collection: users
  timeout: 30s
`, mockSecretKey)
	config3 = fmt.Sprintf(`
authentication:
  secret: "%s"
  method: HS256
  exp-at: 30d
http:
  port: 8080
  read-timeout: 15s
  write-timeout: 15s
https:
  port: 8043
  read-timeout: 15s
  write-timeout: 15s
  pem: ""
  key: ""
rpc:
  port: 50051
mongo:
  mode: tcp
  user: uni
  password: secret
  hosts: [{host: localhost, port: 27017}, {host: 192.168.0.7, port: 27017}]
  db: uni
  collection: users
  timeout: 30s
`, mockSecretKey)
	config4 = fmt.Sprintf(`
authentication:
  secret: "%s"
  method: HS256
  exp-at: 30d
http:
  port: 8080
  read-timeout: 15s
  write-timeout: 15s
https:
  port: 8043
  read-timeout: 15s
  write-timeout: 15s
  pem: ""
  key: ""
rpc:
  port: 50051
mongo:
  mode: atlas
  user: uni
  password: secret
  hosts: [{host: localhost, port: 27017}, {host: 192.168.0.7, port: 27017}]
  db: uni
  collection: users
  timeout: 30s
`, mockSecretKey)
	configs = []string{config1, config2, config3, config4}
)

func TestInit(t *testing.T) {
	for _, val := range configs {
		var conf = &Config{}
		err := yaml.Unmarshal([]byte(val), conf)
		if err != nil {
			t.Error(err)
		}
	}
}
