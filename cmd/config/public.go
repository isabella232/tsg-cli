//
//  Copyright (c) 2018, Joyent, Inc. All rights reserved.
//
//  This Source Code Form is subject to the terms of the Mozilla Public
//  License, v. 2.0. If a copy of the MPL was not distributed with this
//  file, You can obtain one at http://mozilla.org/MPL/2.0/.
//

package config

import (
	"encoding/pem"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/joyent/triton-go"
	"github.com/joyent/triton-go/authentication"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
	"github.com/joyent/tsg-cli/cmd/internal/config"
)

type TritonClientConfig struct {
	Config *triton.ClientConfig
}

func New() (*TritonClientConfig, error) {
	viper.AutomaticEnv()

	var signer authentication.Signer
	var err error

	keyMaterial := GetTritonKeyMaterial()
	if keyMaterial == "" {
		signer, err = authentication.NewSSHAgentSigner(authentication.SSHAgentSignerInput{
			KeyID:       GetTritonKeyID(),
			AccountName: GetTritonAccount(),
		})
		if err != nil {
			log.Fatal().Str("func", "initConfig").Msg(err.Error())
			panic(err)
			return nil, err
		}
	} else {
		var keyBytes []byte
		if _, err = os.Stat(keyMaterial); err == nil {
			keyBytes, err = ioutil.ReadFile(keyMaterial)
			if err != nil {
				return nil, fmt.Errorf("error reading key material from %s: %s",
					keyMaterial, err)
			}
			block, _ := pem.Decode(keyBytes)
			if block == nil {
				return nil, fmt.Errorf(
					"failed to read key material '%s': no key found", keyMaterial)
			}

			if block.Headers["Proc-Type"] == "4,ENCRYPTED" {
				return nil, fmt.Errorf(
					"failed to read key '%s': password protected keys are\n"+
						"not currently supported. Please decrypt the key prior to use.", keyMaterial)
			}

		} else {
			keyBytes = []byte(keyMaterial)
		}

		signer, err = authentication.NewPrivateKeySigner(authentication.PrivateKeySignerInput{
			KeyID:              GetTritonKeyID(),
			PrivateKeyMaterial: keyBytes,
			AccountName:        GetTritonAccount(),
		})
		if err != nil {
			return nil, errors.Wrap(err, "Error Creating SSH Private Key Signer")
		}
	}

	config := &triton.ClientConfig{
		TritonURL:   GetTritonUrl(),
		AccountName: GetTritonAccount(),
		Signers:     []authentication.Signer{signer},
	}

	return &TritonClientConfig{
		Config: config,
	}, nil
}

func GetTritonUrl() string {
	return viper.GetString(config.KeyUrl)
}

func GetTritonKeyMaterial() string {
	return viper.GetString(config.KeySshKeyMaterial)
}

func GetTritonAccount() string {
	return viper.GetString(config.KeyAccount)
}

func GetTritonKeyID() string {
	return viper.GetString(config.KeySshKeyID)
}

func GetPkgID() string {
	return viper.GetString(config.KeyPackageId)
}

func GetImgID() string {
	return viper.GetString(config.KeyImageId)
}

func GetExpectedMachineCount() int {
	return viper.GetInt(config.KeyInstanceCount)
}

func GetTsgName() string {
	return viper.GetString(config.KeyTsgGroupName)
}

func GetMachineName() string {
	return viper.GetString(config.KeyInstanceName)
}

func GetMachineNamePrefix() string {
	return viper.GetString(config.KeyInstanceNamePrefix)
}

func GetMachineBrand() string {
	return viper.GetString(config.KeyInstanceBrand)
}

func GetMachineFirewall() bool {
	return viper.GetBool(config.KeyInstanceFirewall)
}

func GetMachineNetworks() []string {
	if viper.IsSet(config.KeyInstanceNetwork) {
		var networks []string
		cfg := viper.GetStringSlice(config.KeyInstanceNetwork)
		for _, i := range cfg {
			networks = append(networks, i)
		}

		return networks
	}
	return nil
}

func GetMachineAffinityRules() []string {
	if viper.IsSet(config.KeyInstanceAffinityRule) {
		var rules []string
		cfg := viper.GetStringSlice(config.KeyInstanceAffinityRule)
		for _, i := range cfg {
			rules = append(rules, i)
		}

		return rules
	}
	return nil
}

func GetMachineTags() map[string]string {
	if viper.IsSet(config.KeyInstanceTag) {
		tags := make(map[string]string, 0)
		cfg := viper.GetStringSlice(config.KeyInstanceTag)
		for _, i := range cfg {
			m := strings.Split(i, "=")
			tags[m[0]] = m[1]
		}

		return tags
	}

	return nil
}

func GetSearchTags() map[string]interface{} {
	if viper.IsSet(config.KeyInstanceTag) {
		tags := make(map[string]interface{}, 0)
		cfg := viper.GetStringSlice(config.KeyInstanceTag)
		for _, i := range cfg {
			m := strings.Split(i, "=")
			tags[m[0]] = m[1]
		}

		return tags
	}

	return nil
}

func GetMachineMetadata() map[string]string {
	if viper.IsSet(config.KeyInstanceMetadata) {
		metadata := make(map[string]string, 0)
		cfg := viper.GetStringSlice(config.KeyInstanceMetadata)
		for _, i := range cfg {
			m := strings.Split(i, "=")
			metadata[m[0]] = m[1]
		}

		return metadata
	}

	return nil
}

func GetMachineUserdata() string {
	return viper.GetString(config.KeyInstanceUserdata)
}
