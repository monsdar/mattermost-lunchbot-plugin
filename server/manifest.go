// This file is automatically generated. Do not modify it manually.

package main

import (
	"strings"

	"github.com/mattermost/mattermost-server/v5/model"
)

var manifest *model.Manifest

const manifestStr = `
{
  "id": "com.nilsbrinkmann.lunchbot",
  "name": "Lunchbot Plugin",
  "description": "This plugin pairs random users for lunch",
  "homepage_url": "https://github.com/monsdar/mattermost-lunchbot-plugin",
  "support_url": "https://github.com/monsdar/mattermost-lunchbot-plugin/issues",
  "release_notes_url": "https://github.com/monsdar/mattermost-lunchbot-plugin/releases",
  "version": "1.7.0",
  "min_server_version": "5.12.0",
  "server": {
    "executables": {
      "linux-amd64": "server/dist/plugin-linux-amd64",
      "darwin-amd64": "server/dist/plugin-darwin-amd64",
      "windows-amd64": "server/dist/plugin-windows-amd64.exe"
    },
    "executable": ""
  }
}
`

func init() {
	manifest = model.ManifestFromJson(strings.NewReader(manifestStr))
}
