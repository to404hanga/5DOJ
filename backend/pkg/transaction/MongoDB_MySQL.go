package transaction

import (
	"5DOJ/problem/global"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

// prepare func() bool，预处理语句，如果要在预处理处中断，返回 nil
func TransactionWithMongoDBAndMySQL(ctx context.Context, prepare func() bool, mongodb func() error, mysql func(txSQL *gorm.DB) error) (err error) {
	eg := &errgroup.Group{}
	txSQL := global.MySQL.WithContext(ctx).Begin()
	var sess mongo.Session
	if sess, err = global.MongoDB.Client().StartSession(); err != nil {
		txSQL.Rollback()
		return
	}
	defer func() {
		sess.EndSession(ctx)
	}()

	sess.StartTransaction()
	defer func() {
		if rec := recover(); rec != nil {
			txSQL.Rollback()
			sess.AbortTransaction(ctx)
			err = fmt.Errorf("%v", rec)
			return
		}
		if err != nil {
			txSQL.Rollback()
			sess.AbortTransaction(ctx)
			return
		}
		txSQL.Commit()
		sess.CommitTransaction(ctx)
	}()

	if ok := prepare(); ok {
		return
	}

	eg.Go(func() error {
		return mysql(txSQL)
	})
	eg.Go(mongodb)

	err = eg.Wait()
	return
}
