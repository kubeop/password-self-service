/*
Copyright © 2023 Sonic Ma <EMAIL sonic.ma@outlook.com>
*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"password-self-service/internal/router"
	"password-self-service/pkg/config"
	"password-self-service/pkg/gredis"
	"password-self-service/pkg/logging"
	"password-self-service/pkg/tasks"
	"syscall"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "password-self-service",
	Short: "密码自助平台",
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		logging.Init()
		gredis.Init()
		tasks.Init()

		// 注册所有路由
		routers := router.InitRouter()

		srv := &http.Server{
			Addr:    config.Setting.Server.Addr,
			Handler: routers,
		}

		// Initializing the server in a goroutine so that
		// it won't block the graceful shutdown handling below
		go func() {
			// 服务连接
			if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
				logging.Logger().Sugar().Fatalf("Starting server failed: %s", err)
			}
		}()

		logging.Logger().Sugar().Infof("Starting server on %v, Pid is %v", config.Setting.Server.Addr, syscall.Getpid())

		// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
		quit := make(chan os.Signal, 1)
		// kill (no param) default send syscall.SIGTERM
		// kill -2 is syscall.SIGINT
		// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		logging.Logger().Info("Shutting down server...")

		// The context is used to inform the server it has 5 seconds to finish
		// the request it is currently handling
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			logging.Logger().Sugar().Fatalf("Server forced to shutdown: ", err)
		}
		logging.Logger().Info("Server exiting")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/password-self-service.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".password-self-service" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName("password-self-service")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
	}

	if err := viper.Unmarshal(&config.Setting); err != nil {
		log.Fatalf("Fail to parse config: %v", err)
	}

}
