// Package sync 是 git-sync-core 的库入口：配置加载与 Service 构造。
// 详细子域见 model / service / executor / dao / lock。
package sync

import (
	"github.com/yi-nology/git-sync-core/model"
	"github.com/yi-nology/git-sync-core/service"
)

type Service = service.Service

type Config = model.Config

func NewService(cfg *Config) (*Service, error) {
	return service.NewService(cfg)
}

func LoadConfig(path string) (*Config, error) {
	return model.LoadConfig(path)
}
