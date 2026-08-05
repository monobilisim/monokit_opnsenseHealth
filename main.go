// opnsenseHealth plugin for monokit2 — OPNsense health: web UI TLS certificate validity and expiry.
//
// Rule of thumb: plugins import only monokit_lib — nothing else from the project.
package main

import (
	"os"

	lib "github.com/monobilisim/monokit_lib"
)

// Populated at build time via -ldflags "-X 'main.version=<version>'".
var version string

// Plugin name — used to scope Zulip alarms / Redmine issues to this plugin.
var pluginName string = "opnsenseHealth"

// Alarm states passed to lib.SendZulipAlarm / Redmine helpers.
var up string = "up"
var down string = "down"

// Config files this plugin needs from /etc/mono (in addition to the shared
// global.yml, which lib.InitConfig always loads). Leave empty if you only need
// global.yml; otherwise add your own (e.g. "plugin.yml") and wire it into
// lib's InitConfig switch and config types.
var configFiles []string = []string{"opnsense.yml"}

func main() {
	// Handles --version, --help and other common flags, then exits.
	if len(os.Args) > 1 {
		lib.HandleCommonPluginArgs(os.Args, version, configFiles)
		return
	}

	if err := lib.InitConfig(configFiles...); err != nil {
		panic("Failed to initialize config: " + err.Error())
	}

	logger, err := lib.InitLogger()
	if err != nil {
		panic("Failed to initialize logger: " + err.Error())
	}

	lib.InitializeDatabase()

	logger.Info().Msg("Starting " + pluginName + " plugin...")

	// TODO: implement your plugin's health checks here. For alarms, use the
	// up/down states, e.g.:
	//   lib.SendZulipAlarm("something is wrong", pluginName, "moduleName", down)
	//   lib.SendZulipAlarm("recovered", pluginName, "moduleName", up)
	_ = up
	_ = down
}
