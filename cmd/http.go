package cmd

import (
	"fmt"
	"github.com/nova2018/easygin/app"
	"github.com/nova2018/easygin/http"
	"github.com/nova2018/goutils"
	"github.com/spf13/cobra"
)

func AddServer(r http.Router, handler goutils.RecoveryHandle) {
	httpCommand := &cobra.Command{
		Use:   "http",
		Short: "启动http服务",
		Long: `启动方式:
   ./app http --config=./config/test.toml 
   ./app http -c=./config/test.toml
注: --config/-c:代表配置文件`,
		Run: daemon(httpServer(r, handler)),
	}
	rootCmd.AddCommand(httpCommand)
}

func httpServer(r http.Router, handler goutils.RecoveryHandle) func(cmd *cobra.Command, args []string) {
	return func(cmd *cobra.Command, args []string) {
		fmt.Println("cfgFile-->", cfgFile)
		runServerWithPanic(func(server *app.Server) error {
			server.Router(r)
			return nil
		}, handler)
	}
}

func runServerWithPanic(fn func(*app.Server) error, handler goutils.RecoveryHandle) {
	s := app.NewServerWithFlags(rootFlag)

	err := fn(s)
	if err != nil {
		panic(err)
	}

	s.RunWithPanic(handler)
}
