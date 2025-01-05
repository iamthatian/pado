// TODO
package state

import (
	// "fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/duckonomy/parkour/utils"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

var k = koanf.New(".")

// TODO: change this to per OS settings
func GetConfig() error {
	// change behavior for macos
	userConfigDir, err := os.UserConfigDir()
	if err != nil {
		return err
	}

	configPath := filepath.Join(userConfigDir, "parkour", "parkour.toml")

	if err := utils.CheckPath(configPath, false); err != nil {
		return err
	}

	if err := k.Load(file.Provider(configPath), toml.Parser()); err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	// blacklist := k.Get("blacklist").([]interface{})
	// for _, j := range blacklist {
	// 	fmt.Println(j == "awesome")
	// }

	return nil
}
