// Refactor this 'abuse' in the future:
// https://leighmcculloch.com/posts/tool-go-check-no-globals-no-inits/
package notion

import "tele-notion-bot/config"

var ApiRoot string
var ApiVersion string

func init() {
	cfg := config.GetConfig()
	ApiRoot = cfg.GetString("NOTION.API_ROOT")
	ApiVersion = cfg.GetString("NOTION.API_VERSION")
}
