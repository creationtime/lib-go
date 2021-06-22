package log

import (
	"context"
	"encoding/json"
	"github.com/lights-T/lib-go/util/strings"
	"github.com/micro/go-micro/v2/client"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/server"
)

type logWrapper struct {
	client.Client
}

func (l *logWrapper) Call(ctx context.Context, req client.Request, rsp interface{}, opts ...client.CallOption) error {
	err := l.Client.Call(ctx, req, rsp)
	reqBuf, _ := json.Marshal(req.Body())
	rspBuf, _ := json.Marshal(rsp)
	if err != nil {
		logger.Errorf("client call service:%s endpoint:%s request:%s rsp:%s error:%s", req.Service(), req.Endpoint(), strings.BytesToString(reqBuf), strings.BytesToString(rspBuf), err.Error())
		return err
	}
	logger.Infof("client call service:%s endpoint:%s request:%s rsp:%s", req.Service(), req.Endpoint(), strings.BytesToString(reqBuf), strings.BytesToString(rspBuf))
	return nil
}

func NewClientWrapper(c client.Client) client.Client {
	return &logWrapper{c}
}

func NewHandlerWrapper(fn server.HandlerFunc) server.HandlerFunc {
	return func(ctx context.Context, req server.Request, rsp interface{}) error {
		err := fn(ctx, req, rsp)
		reqBuf, _ := json.Marshal(req.Body())
		rspBuf, _ := json.Marshal(rsp)
		if err != nil {
			logger.Errorf("handler service:%s endpoint:%s request:%s rsp:%s error:%s", req.Service(), req.Endpoint(), strings.BytesToString(reqBuf), strings.BytesToString(rspBuf), err.Error())
			return err
		}
		logger.Infof("handler service:%s endpoint:%s request:%s rsp:%s", req.Service(), req.Endpoint(), strings.BytesToString(reqBuf), strings.BytesToString(rspBuf))
		return nil
	}
}
