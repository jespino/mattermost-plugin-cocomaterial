package main

import (
	"fmt"
	"net/http"
	"path"
	"path/filepath"
	"sync"

	"github.com/mattermost/mattermost-server/v5/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	cocoEntries    []string
	cocoCategories map[string][]string
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.RequestURI)
	bundlePath, err := p.API.GetBundlePath()
	if err != nil {
		return
	}

	for _, cocoEntry := range p.cocoEntries {
		if path.Base(r.RequestURI) == normalizeName(cocoEntry)+".png" {
			assetFile := filepath.Join(bundlePath, "assets", "coco", cocoEntry+".png")
			http.ServeFile(w, r, assetFile)
		}
	}
}

// See https://developers.mattermost.com/extend/plugins/server/reference/

func (p *Plugin) OnActivate() error {
	p.setCocoEntries()
	if err := p.API.RegisterCommand(createCocoCommand(p.cocoCategories)); err != nil {
		return err
	}
	return nil
}
