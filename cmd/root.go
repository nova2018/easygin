package cmd

import (
	"github.com/nova2018/easygin/app"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	rootCmd = &cobra.Command{
		Use:   "app",
		Short: "app:代表编译后的二进制文件",
	}
	cfgFile  = ""
	isDaemon = false
	rootFlag *pflag.FlagSet
)

func init() {
	rootFlag = rootCmd.PersistentFlags()
	rootFlag.StringVarP(&cfgFile, "config", "c", app.GetDefaultConfigFile(), "config file (default is ./conf/test.toml)")
	rootFlag.BoolVarP(&isDaemon, "daemon", "d", false, "is daemon")
	cobra.OnInitialize(initCmd)
}

func initCmd() {

}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
