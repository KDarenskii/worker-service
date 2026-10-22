package cmd

import (
	"github.com/KDarenskii/worker-service/internal/app/builder"
	"github.com/urfave/cli/v2"
)

// SubscribeOrderCreated возвращает CLI команду для запуска consumer'а топика order.created.
func SubscribeOrderCreated() *cli.Command {
	return &cli.Command{
		Name:            "subscribe-order-created",
		Aliases:         []string{"consume-order-created"},
		Usage:           "Запускает consumer топика order.created",
		Action:          cmdSubscribeOrderCreated,
		HideHelpCommand: true,
	}
}

func cmdSubscribeOrderCreated(cCtx *cli.Context) error {
	app := builder.NewBuilder(cCtx)
	app.BuildConfig()
	app.BuildMonitorOpenTelemetry()
	app.BuildConnRedis()
	app.BuildBrokerKafka()
	app.BuildConsumerOrderCreated()
	app.BuildProcHttp()

	app.Run()
	return nil
}
