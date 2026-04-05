package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/charmbracelet/bubbletea"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"gopkg.in/op/go-logging.v1"

	"github.com/arch-err/yqp/tui/bubbles/yqplayground"
	"github.com/arch-err/yqp/tui/theme"
)

func init() {
	// Silence yq's go-logging output to prevent stderr noise in the TUI.
	backend := logging.NewLogBackend(&discardWriter{}, "", 0)
	logging.SetBackend(backend)
}

type discardWriter struct{}

func (*discardWriter) Write(p []byte) (int, error) { return len(p), nil }

var rootCmd = &cobra.Command{
	Version: "0.1.0",
	Use:     "yqp [query]",
	Short:   "yqp is a TUI to explore yq",
	Long: `yqp is a terminal user interface (TUI) for exploring the yq command line utility.

You can use it to run yq queries interactively. If no query is provided, the interface will prompt you for one.

The command accepts an optional query argument which will be executed against the input YAML.
You can provide the input YAML either through a file or via standard input (stdin).`,
	Args:         cobra.MaximumNArgs(1),
	SilenceUsage: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		query := ""
		if len(args) == 1 {
			query = strings.TrimSpace(args[0])
		}

		configTheme := viper.GetString(configKeysName.themeName)
		if !cmd.Flags().Changed(flagsName.theme) {
			flags.theme = configTheme
		}
		themeOverrides := viper.GetStringMapString(configKeysName.themeOverrides)

		styleOverrides := viper.GetStringMapString(configKeysName.styleOverrides)
		yqtheme, defaultTheme := theme.GetTheme(flags.theme, styleOverrides)

		if !defaultTheme && configTheme == flags.theme && len(themeOverrides) > 0 {
			chromaTypes := make(map[string]chroma.TokenType)
			for tokenType, short := range chroma.StandardTypes {
				chromaTypes[short] = tokenType
			}

			builder := yqtheme.ChromaStyle.Builder()
			for k, v := range themeOverrides {
				builder.Add(chromaTypes[k], v)
			}
			style, err := builder.Build()
			if err == nil {
				yqtheme.ChromaStyle = style
			}
		}

		if isInputFromPipe() {
			stdin, err := streamToBytes(os.Stdin)
			if err != nil {
				return err
			}
			bubble, err := yqplayground.New(stdin, "STDIN", query, yqtheme)
			if err != nil {
				return err
			}
			p := tea.NewProgram(bubble, tea.WithAltScreen())
			m, err := p.Run()
			if err != nil {
				return err
			}
			if m, ok := m.(yqplayground.Bubble); ok && m.ExitMessage != "" {
				return errors.New(m.ExitMessage)
			}
			return nil
		}

		// get the file
		file, e := getFile()
		if e != nil {
			return e
		}
		defer file.Close()

		// read the file
		data, err := os.ReadFile(flags.filepath)
		if err != nil {
			return err
		}

		// get file info so we can get the filename
		fi, err := os.Stat(flags.filepath)
		if err != nil {
			return err
		}

		bubble, err := yqplayground.New(data, fi.Name(), query, yqtheme)
		if err != nil {
			return err
		}
		p := tea.NewProgram(bubble, tea.WithAltScreen())

		m, err := p.Run()
		if err != nil {
			return err
		}
		if m, ok := m.(yqplayground.Bubble); ok && m.ExitMessage != "" {
			return errors.New(m.ExitMessage)
		}
		return nil
	},
}

func initConfig() {
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
		if err := viper.ReadInConfig(); err != nil {
			fmt.Fprintf(os.Stderr, "Config file %s was unable to be read: %v\n", viper.ConfigFileUsed(), err)
		}
		return
	}
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)

	viper.AddConfigPath(home)
	viper.SetConfigName(".yqp")
	viper.SetConfigType("yaml")

	if err := viper.ReadInConfig(); err != nil {
		var errFileNotFound viper.ConfigFileNotFoundError
		if !errors.As(err, &errFileNotFound) {
			fmt.Fprintf(os.Stderr, "Default config file %s was unable to be read: %v\n", viper.ConfigFileUsed(), err)
		}
	}
}

var flags struct {
	filepath, theme string
}

var flagsName = struct {
	file, fileShort, theme, themeShort string
}{
	file:       "file",
	fileShort:  "f",
	theme:      "theme",
	themeShort: "t",
}

var configKeysName = struct {
	themeName      string
	themeOverrides string
	styleOverrides string
}{
	themeName:      "theme.name",
	themeOverrides: "theme.chromaStyleOverrides",
	styleOverrides: "theme.styleOverrides",
}

var cfgFile string

func Execute() error {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "path to config file (default is $HOME/.yqp.yaml)")

	rootCmd.Flags().StringVarP(
		&flags.filepath,
		flagsName.file,
		flagsName.fileShort,
		"", "path to the input YAML file")

	rootCmd.Flags().StringVarP(
		&flags.theme,
		flagsName.theme,
		flagsName.themeShort,
		"", "yqp theme")

	return rootCmd.Execute()
}
