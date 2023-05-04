package server

import (
	"wework-vkm/internal/initialization"
	"log"

	internalHandler "wework-vkm/internal/handlers"

	"github.com/gin-gonic/gin"
	"wework-vkm/pkg/openai"
	weworkHandler "wework-vkm/pkg/wework/handlers"

	"github.com/spf13/cobra"
)

var cmd = &cobra.Command{
	Use:   "server",
	Short: "run wework webhook server",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := cmd.Flags().GetString("config")		

		if err != nil {
			log.Println(err)
			return
		}

		initialization.InitRoleList()
		config := initialization.LoadConfig(cfg)

		initialization.LoadWeworkClient(*config)

		gpt := openai.NewChatGPT(*config)
		internalHandler.InitHandlers(gpt, *config)

		r := gin.Default()

		// 存储应用配置
		r.Use(func (c *gin.Context) {
			c.Keys = make(map[string]interface{})
			c.Keys["config"] = config
			c.Keys["handler"] = internalHandler.HandleMessage
			c.Next()
		})
		
		r.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})

		r.GET("/webhook/event/wework", weworkHandler.HandleGinRequest)
		r.POST("/webhook/event/wework", weworkHandler.HandleGinRequest)

		err = initialization.StartServer(*config, r)
		if err != nil {
			log.Fatalf("failed to start server: %v", err)
		}
	},
}

func Register(rootCmd *cobra.Command) error {
	rootCmd.AddCommand(cmd)
	return nil
}
