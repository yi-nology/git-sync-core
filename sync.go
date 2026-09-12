// Package sync 是 git-sync-core 的库入口：配置加载与 Service 构造。
//
// 本库**不提供任何 HTTP/Web 服务端**，也不依赖 hertz/gin 等 Web 框架。
// 入站 Webhook 请由壳层收包后转成 service.WebhookPayload 再调用 Service.ReceiveWebhook。
// Config.Server 中的 Host/Port/APIKey 仅由壳层消费；core 自身不会监听端口。
//
// 详细子域见 model / service / executor / dao / lock。
package sync

import (
	"github.com/yi-nology/git-sync-core/model"
	"github.com/yi-nology/git-sync-core/service"
)

type Service = service.Service

type Config = model.Config

// WebhookPayload 协议无关的 Webhook 入站载荷（由壳层构造）。
type WebhookPayload = service.WebhookPayload

func NewService(cfg *Config) (*Service, error) {
	return service.NewService(cfg)
}

func LoadConfig(path string) (*Config, error) {
	return model.LoadConfig(path)
}
