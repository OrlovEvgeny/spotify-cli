package main

import (
	"fmt"
	"spotify/internal/back"
	"spotify/internal/device"
	"spotify/internal/login"
	"spotify/internal/next"
	"spotify/internal/p"
	"spotify/internal/pause"
	"spotify/internal/play"
	"spotify/internal/playlist"
	"spotify/internal/queue"
	"spotify/internal/repeat"
	"spotify/internal/save"
	"spotify/internal/shuffle"
	"spotify/internal/status"
	"spotify/internal/unsave"
	"spotify/internal/update"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func main() {
	// TODO: https://github.com/spf13/viper/pull/1064
	viper.AddConfigPath("$HOME")
	viper.SetConfigName(".spotify-cli")
	viper.SetConfigType("json")
	_ = viper.SafeWriteConfig()
	_ = viper.ReadInConfig()

	root := &cobra.Command{
		Use:               "spotify",
		Short:             "Spotify for the terminal 🎵",
		Version:           buildVersion(version, commit, date),
		PersistentPreRunE: promptUpdate,
	}

	root.AddCommand(back.NewCommand())
	root.AddCommand(device.NewCommand())
	root.AddCommand(login.NewCommand())
	root.AddCommand(next.NewCommand())
	root.AddCommand(p.NewCommand())
	root.AddCommand(pause.NewCommand())
	root.AddCommand(play.NewCommand())
	root.AddCommand(playlist.NewCommand())
	root.AddCommand(queue.NewCommand())
	root.AddCommand(repeat.NewCommand())
	root.AddCommand(save.NewCommand())
	root.AddCommand(shuffle.NewCommand())
	root.AddCommand(status.NewCommand())
	root.AddCommand(unsave.NewCommand())
	root.AddCommand(update.NewCommand())

	// Hide help command
	root.SetHelpCommand(&cobra.Command{Hidden: true})

	_ = root.Execute()
}

// Sets ldflags by goreleaser https://goreleaser.com/customization/build/ (default values)
var (
	version = "dev"
	commit  string
	date    string
)

func buildVersion(version, commit, date string) string {
	result := version
	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}
	if date != "" {
		result = fmt.Sprintf("%s\nbuilt at: %s", result, date)
	}
	return result
}

func promptUpdate(cmd *cobra.Command, _ []string) error {
	if time.Now().Unix() < viper.GetInt64("prompt_update_timer") {
		return nil
	}

	isUpdated, err := update.IsUpdated(cmd)
	if err != nil {
		return err
	}
	if !isUpdated {
		cmd.Println("Use 'spotify update' to get the latest version.")
	}

	// Wait one day before the next prompt
	const day int64 = 24 * 60 * 60
	viper.Set("prompt_update_timer", time.Now().Unix()+day)

	return viper.WriteConfig()
}
