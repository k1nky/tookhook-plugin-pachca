package main

import (
	"context"
	"os"
	"strings"

	hcplugin "github.com/hashicorp/go-plugin"
	"github.com/k1nky/tookhook-plugin-pachca/internal/options"
	"github.com/k1nky/tookhook-plugin-pachca/internal/pachca"
	"github.com/k1nky/tookhook/pkg/logger"
	"github.com/k1nky/tookhook/pkg/plugin"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	DefaultLogLevel = "debug"
)

type Plugin struct {
	log *logger.Logger
}

func (p Plugin) Validate(ctx context.Context, r plugin.Handler) error {
	opts, err := options.New(r.Options)
	if err != nil {
		p.log.Errorf("validate: %v", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}
	if err := opts.Validate(); err != nil {
		p.log.Errorf("validate: %v", err)
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

func (p Plugin) Health(ctx context.Context) error {
	return nil
}

func (p Plugin) Forward(ctx context.Context, r plugin.Handler, data []byte) ([]byte, error) {
	opts, err := options.New(r.Options)
	if err != nil {
		p.log.Errorf("forward: %v", err)
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	chat := strings.Split(opts.Chat, "/")
	pch := pachca.NewPachca(opts.Token)
	m := pachca.MessagePayload{Message: pachca.Message{
		EntityType:       chat[0],
		EntityId:         chat[1],
		Content:          string(data),
		DisplayName:      opts.DisplayName,
		DisplayAvatarUrl: opts.DisplayAvatarUrl,
	}}
	response, err := pch.Send(m)
	p.log.Debugf("forward to %s with response: %s", chat, string(response))
	if err != nil {
		err = status.Error(codes.Unavailable, err.Error())
	}
	return response, err
}

func main() {
	log := newLogger()
	hcplugin.Serve(&hcplugin.ServeConfig{
		HandshakeConfig: plugin.Handshake,
		Plugins: map[string]hcplugin.Plugin{
			"grpc": &plugin.GRPCPlugin{Impl: &Plugin{
				log: log,
			}},
		},

		GRPCServer: hcplugin.DefaultGRPCServer,
	})
}

func newLogger() *logger.Logger {
	logLevel := os.Getenv("TOOKHOK_PLUGIN_PACHCA_LOG_LEVEL")
	if logLevel == "" {
		logLevel = DefaultLogLevel
	}
	l := logger.New("pachca")
	if err := l.SetLevel(DefaultLogLevel); err != nil {
		l.Errorf("invalid log level: %v", err)
	}
	return l
}
