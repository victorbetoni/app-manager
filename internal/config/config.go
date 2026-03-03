package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	UseTLS       bool      `json:"use_tls"`
	UseAuth      bool      `json:"use_auth"`
	Port         int       `json:"port"`
	Origins      []string  `json:"origins"`
	AppsJsonPath string    `json:"apps_json_path"`
	TlsCert      KeyConfig `json:"tls_cert"`
	Jwt          JwtConfig `json:"jwt"`
}

type JwtConfig struct {
	CookieKey      string          `json:"cookie_key"`
	CheckUserAgent bool            `json:"check_user_agent"`
	CheckIp        bool            `json:"check_ip"`
	ClaimKeys      ClaimKeysConfig `json:"claim_keys"`
	Keys           KeyConfig       `json:"keys"`
}

type ClaimKeysConfig struct {
	Identifier string `json:"identifier"`
	AdminFlag  string `json:"admin_flag"`
	UserAgent  string `json:"user_agent"`
	Ip         string `json:"ip"`
}

type KeyConfig struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

var cfg Config

func Load() {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	mustReadAndAssign(pwd, "config.json", &cfg)
}

func mustReadAndAssign(pwd, relativeDir string, target interface{}) {
	f, err := os.ReadFile(fmt.Sprintf("%s/%s", pwd, relativeDir))
	if err != nil {
		panic(err)
	}
	if err := json.Unmarshal(f, &target); err != nil {
		panic(err)
	}
}

func GetConfig() *Config {
	return &cfg
}
