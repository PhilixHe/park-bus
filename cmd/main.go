package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/cobra"
	"log"
	"net/http"
	"os"
	"os/signal"
	conf "park-bus/config"
	"park-bus/cron_job"
	"park-bus/pkg/version"
	"time"
)

const ErrExitCode = 1

func NewVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "print version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			fmt.Println(version.Get().String())
			return nil
		},
	}
}

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "park-bus",
		Short:   "park-bus online binary",
		Version: version.Get().String(),
		Run: func(cmd *cobra.Command, args []string) {
			// 获取配置文件路径
			cfgFile, err := cmd.Flags().GetString("config")
			if err != nil {
				panic(err)
			}

			// 获取文件信息，判断文件是否存在。
			_, err = os.Stat(cfgFile)
			if err != nil {
				if os.IsNotExist(err) {
					fmt.Println(fmt.Errorf("config file not exist: %s", cfgFile))
					cmd.Help()
					return
				}
			}

			// 加载配置文件
			conf.LoadConfig(cfgFile)

			// 启动服务
			runServer()
			return
		},
	}

	cmd.Flags().StringP("config", "f", "./config/config.yaml", "config file (default is ./config/config.yaml)")

	cmd.AddCommand(
		NewVersionCmd(),
	)

	return cmd
}

func runServer() {
	gin.ForceConsoleColor()
	r := gin.New()
	r.Use(
		gin.Recovery(),
		gin.Logger(),
	)

	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	srv := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: r,
	}

	// 启动定时任务
	cron_job.Run()

	go func() {
		// 服务连接
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// 等待中断信号以优雅地关闭服务器（设置 5 秒的超时时间）
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

func main() {
	if err := NewRootCmd().Execute(); err != nil {
		fmt.Println(err.Error())
		os.Exit(ErrExitCode)
	}
}
