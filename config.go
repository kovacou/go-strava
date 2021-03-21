// Copyright © 2020 Alexandre KOVAC <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package strava

import "github.com/kovacou/go-env"

var (
	// cfgEnviron contains the loaded configuration from environment.
	cfgEnviron Config
)

// init loads the global configuration.
func init() {
	_ = env.Unmarshal(&cfgEnviron)
}

// Config is the configuration for the client to request the Strava API.
type Config struct {
	Host         string `json:"host" env:"STRAVA_HOST"`
	ClientID     string `json:"client_id" env:"STRAVA_ID"`
	ClientSecret string `json:"client_secret" env:"STRAVA_SECRET"`
	RedirectURI  string `json:"redirect_uri" env:"STRAVA_REDIRECT_URI"`
	Timeout      uint16 `json:"timeout" env:"STRAVA_TIMEOUT"`
	Scope        string `json:"scope" env:"STRAVA_SCOPE"`
}
