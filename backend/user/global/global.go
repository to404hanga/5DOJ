package global

import (
	"github.com/to404hanga/pkg404/grpcx"
	"github.com/to404hanga/pkg404/logger"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gorm.io/gorm"
)

var (
	MySQL      *gorm.DB
	L          logger.Logger
	Etcd       *clientv3.Client
	GrpcServer *grpcx.Server
)
