package cmd

import (
	"fmt"
	"strings"

	"github.com/gdamore/tcell"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/derricw/siggo/model"
)

func init() {
	cfgCmd.AddCommand(cfgColorCmd)
	rootCmd.AddCommand(cfgCmd)
}

var cfgCmd = &cobra.Command{
	Use:   "cfg",
	Short: "configure siggo",
	Long: `Example:
    $ siggo cfg`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := model.GetConfig()
		if err != nil {
			log.Fatalf("couldn't load config: %s\n", err)
		}
		cfg.Print()
	},
}

var cfgColorCmd = &cobra.Command{
	Use:   "color",
	Short: "sets or prints the color for a contact",
	Long: `Accepts W3C color names or hex format.
	Example:
    $ siggo cfg color "John Smith" DeepSkyBlue
    $ siggo cfg color "John Smith" "#FF00FF"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatalf("color config requires at least a contact name")
		} else if len(args) == 1 {
			// show color for contact
			contactName := args[0]
			cfg, err := model.GetConfig()
			if err != nil {
				log.Fatalf("couldn't load current config: %s", err)
			}
			if color, ok := cfg.ContactColors[contactName]; ok {
				// found contact
				fmt.Printf("%s: %s\n", contactName, color)
			} else {
				log.Fatalf("contact '%s' has no color configuration", contactName)
			}
		} else if len(args) == 2 {
			// set color for contact
			contactName := args[0]
			colorName := strings.ToLower(args[1])
			cfg, err := model.GetConfig()
			if err != nil {
				log.Fatalf("couldn't load current config: %s", err)
			}
			if cfg.UserNumber == "" {
				log.Fatalf("no user phone number configured @ %s", model.DefaultConfigPath())
			}
			// make sure contact exists?
			// make sure color exists
			color := tcell.GetColor(colorName)
			if color == -1 {
				log.Fatalf("color is not valid W3C color: %s", colorName)
			}
			// set color and save config
			if cfg.ContactColors == nil {
				cfg.ContactColors = make(map[string]string)
			}
			cfg.ContactColors[contactName] = colorName
			err = cfg.Save()
			if err != nil {
				log.Fatal(err)
			}
		} else {
			log.Fatalf("too many args")
		}
	},
}