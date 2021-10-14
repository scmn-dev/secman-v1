// Package insert handles adding a new site to the password store.
package insert

import (
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"

	"github.com/scmn-dev/secman-v1/pkg/pc"
	"github.com/scmn-dev/secman-v1/pkg/pio"
	"golang.org/x/crypto/nacl/box"
)

const (
	// PassPrompt is the string formatter that should be used
	PassPrompt = "Enter password for %s"
)

func Password(name string) {
	var c pio.ConfigFile
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Could not generate site key: %s", err.Error())
	}

	config, err := pio.GetConfigPath()
	if err != nil {
		log.Fatalf("Could not get config file name: %s", err.Error())
	}

	// Read the master public key.
	configContents, err := ioutil.ReadFile(config)

	if err != nil {
		log.Fatalf("Could not get config file contents: %s", err.Error())
	}

	err = json.Unmarshal(configContents, &c)

	if err != nil {
		log.Fatalf("Could not unmarshal config file contents: %s", err.Error())
	}

	masterPub := c.MasterPubKey

	sitePass, err := pio.PromptPass(fmt.Sprintf("Enter password for %s", name))
	if err != nil {
		log.Fatalf("Could not get password for site: %s", err.Error())
	}

	passSealed, err := pc.SealAsym([]byte(sitePass), &masterPub, priv)
	if err != nil {
		log.Fatalf("Could not seal new site password: %s", err.Error())
	}

	si := pio.SiteInfo{
		PubKey:     *pub,
		Name:       name,
		PassSealed: passSealed,
	}

	err = si.AddFile(passSealed, name)
	if err != nil {
		log.Fatalf("Could not save site file: %s", err.Error())
	}
}

// File is used to add a new file entry to the vault.
func File(path, filename string) {
	var c pio.ConfigFile
	pub, priv, err := box.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatalf("Could not generate site key: %s", err.Error())
	}

	config, err := pio.GetConfigPath()
	if err != nil {
		log.Fatalf("Could not get config file name: %s", err.Error())
	}

	// Read the master public key.
	configContents, err := ioutil.ReadFile(config)
	if err != nil {
		log.Fatalf("Could not get config file contents: %s", err.Error())
	}

	err = json.Unmarshal(configContents, &c)
	if err != nil {
		log.Fatalf("Could not unmarshal config file contents: %s", err.Error())
	}

	masterPub := c.MasterPubKey

	fileBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Could not open and read file that is being encrypted: %s", err.Error())
	}

	fileSealed, err := pc.SealAsym([]byte(fileBytes), &masterPub, priv)
	if err != nil {
		log.Fatalf("Could not seal file bytes: %s", err.Error())
	}

	si := pio.SiteInfo{
		PubKey:   *pub,
		Name:     path,
		IsFile:   true,
		FileName: path,
	}

	err = si.AddFile(fileSealed, path)
	if err != nil {
		log.Fatalf("Could not save site file after file insert: %s", err.Error())
	}
}
