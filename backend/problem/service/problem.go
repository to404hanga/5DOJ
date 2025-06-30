package service

import (
	"5DOJ/pkg/transaction"
	"5DOJ/problem/domain"
	"5DOJ/problem/global"
	"5DOJ/problem/model"
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/to404hanga/pkg404/gotools/transform"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type ProblemService struct {
}

var _ IProblemService = (*ProblemService)(nil)

func NewProblemService() *ProblemService {
	return &ProblemService{}
}

func (p *ProblemService) Get(ctx context.Context, pid uint64) (problemView domain.ProblemView, err error) {
	var (
		problemInfo    model.ProblemInfo
		problemContent model.ProblemContent
	)

	eg := &errgroup.Group{}
	eg.Go(func() error {
		return global.MySQL.WithContext(ctx).Where("id = ?", pid).First(&problemInfo).Error
	})
	eg.Go(func() error {
		res := global.MongoDB.Collection(problemContent.TableName()).FindOne(ctx, bson.M{
			"pid": pid,
		})
		if res.Err() != nil {
			return res.Err()
		}
		return res.Decode(&problemContent)
	})
	if err = eg.Wait(); err != nil {
		return
	}

	return domain.ProblemView{
		Id:            problemInfo.Id,
		Title:         problemInfo.Title,
		Level:         problemInfo.Level,
		CreatedBy:     problemInfo.CreatedBy,
		UpdatedBy:     problemInfo.UpdatedBy,
		Enabled:       problemInfo.Enabled,
		TimeLimit:     problemInfo.TimeLimit,
		MemoryLimit:   problemInfo.MemoryLimit,
		TotalScore:    problemInfo.TotalScore,
		TotalTestCase: problemInfo.TotalTestCase,
		CreatedAt:     problemInfo.CreatedAt,
		UpdatedAt:     problemInfo.UpdatedAt,
		Markdown:      problemContent.Markdown,
	}, nil
}

func (p *ProblemService) GetTestCaseList(ctx context.Context, pid uint64) (testCaseList []domain.TestCaseView, err error) {
	var (
		testCases []model.TestCase
		cursor    *mongo.Cursor
	)

	cursor, err = global.MongoDB.Collection(model.TestCase{}.TableName()).Find(ctx, bson.M{
		"pid": pid,
	})
	if err != nil {
		return
	}
	defer cursor.Close(ctx)
	if err = cursor.All(ctx, &testCases); err != nil {
		return
	}

	testCaseList = transform.SliceFromSlice[model.TestCase, domain.TestCaseView](testCases, func(i int, tc model.TestCase) domain.TestCaseView {
		return domain.TestCaseView{
			Id:        tc.Tid,
			Input:     tc.Input,
			Output:    tc.Output,
			Score:     tc.Score,
			CreatedBy: tc.CreatedBy,
			UpdatedBy: tc.UpdatedBy,
			Enabled:   tc.Enabled,
		}
	})
	return
}

func (p *ProblemService) GetList(ctx context.Context, size int, cursorIn uint64, title string) (cursorOut uint64, list []domain.ProblemView, err error) {
	var problemInfos []model.ProblemInfo
	query := global.MySQL.WithContext(ctx).Where("id > ?", cursorIn)
	if title != "" {
		query = query.Where("title LIKE ?", "%"+title+"%")
	}
	query.Limit(size)
	err = query.Find(&problemInfos).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return
	}
	if len(problemInfos) == 0 {
		return cursorIn, list, nil
	}

	problemContents := make([]model.ProblemContent, 0, len(problemInfos))
	col := global.MongoDB.Collection(model.ProblemContent{}.TableName())
	for _, problemInfo := range problemInfos {
		res := col.FindOne(ctx, bson.M{
			"pid": problemInfo.Id,
		})
		if res.Err() != nil {
			return cursorIn, list, res.Err()
		}
		var problemContent model.ProblemContent
		if err = res.Decode(&problemContent); err != nil {
			return cursorIn, list, err
		}
		problemContents = append(problemContents, problemContent)
	}

	if len(problemInfos) != len(problemContents) {
		return cursorIn, list, errors.New("problemInfo 与 problemContent 长度不匹配")
	}
	cursorOut = problemInfos[len(list)-1].Id
	for idx := range problemInfos {
		list = append(list, domain.ProblemView{
			Id:            problemInfos[idx].Id,
			Title:         problemInfos[idx].Title,
			Level:         problemInfos[idx].Level,
			CreatedBy:     problemInfos[idx].CreatedBy,
			UpdatedBy:     problemInfos[idx].UpdatedBy,
			Enabled:       problemInfos[idx].Enabled,
			TimeLimit:     problemInfos[idx].TimeLimit,
			MemoryLimit:   problemInfos[idx].MemoryLimit,
			TotalScore:    problemInfos[idx].TotalScore,
			TotalTestCase: problemInfos[idx].TotalTestCase,
			CreatedAt:     problemInfos[idx].CreatedAt,
			UpdatedAt:     problemInfos[idx].UpdatedAt,
			Markdown:      problemContents[idx].Markdown,
		})
	}
	return
}

func (p *ProblemService) Create(ctx context.Context, title string, level int, createdBy string, timeLimit, memoryLimit int, markdown string) (pid uint64, err error) {
	// eg := &errgroup.Group{}
	// txSQL := global.MySQL.WithContext(ctx).Begin()
	// var sess mongo.Session
	// if sess, err = global.MongoDB.Client().StartSession(); err != nil {
	// 	txSQL.Rollback()
	// 	return
	// }
	// defer func() {
	// 	sess.EndSession(ctx)
	// }()

	// sess.StartTransaction()
	// defer func() {
	// 	if rec := recover(); rec != nil {
	// 		txSQL.Rollback()
	// 		sess.AbortTransaction(ctx)
	// 		pid = 0
	// 		err = fmt.Errorf("%v", rec)
	// 		return
	// 	}
	// 	if err != nil {
	// 		txSQL.Rollback()
	// 		sess.AbortTransaction(ctx)
	// 		pid = 0
	// 		return
	// 	}
	// 	txSQL.Commit()
	// 	sess.CommitTransaction(ctx)
	// }()

	// eg.Go(func() error {
	// 	problemInfo := model.ProblemInfo{
	// 		Title:       title,
	// 		Level:       int8(level),
	// 		CreatedBy:   createdBy,
	// 		UpdatedBy:   createdBy,
	// 		Enabled:     true,
	// 		TimeLimit:   timeLimit,
	// 		MemoryLimit: memoryLimit,
	// 	}
	// 	err := txSQL.WithContext(ctx).Create(&problemInfo).Error
	// 	if err != nil {
	// 		return err
	// 	}
	// 	pid = problemInfo.Id
	// 	return nil
	// })
	// eg.Go(func() error {
	// 	problemContent := model.ProblemContent{
	// 		Pid:      pid,
	// 		Markdown: markdown,
	// 	}
	// 	_, err := global.MongoDB.Collection(model.ProblemContent{}.TableName()).InsertOne(ctx, problemContent)
	// 	if err != nil {
	// 		return err
	// 	}
	// 	return nil
	// })

	// err = eg.Wait()
	// return

	prepare := func() bool {
		return false
	}
	mongodb := func() error {
		problemContent := model.ProblemContent{
			Pid:      pid,
			Markdown: markdown,
		}
		_, err := global.MongoDB.Collection(model.ProblemContent{}.TableName()).InsertOne(ctx, problemContent)
		if err != nil {
			return err
		}
		return nil
	}
	mysql := func(txSQL *gorm.DB) error {
		problemInfo := model.ProblemInfo{
			Title:       title,
			Level:       int8(level),
			CreatedBy:   createdBy,
			UpdatedBy:   createdBy,
			Enabled:     true,
			TimeLimit:   timeLimit,
			MemoryLimit: memoryLimit,
		}
		err := txSQL.WithContext(ctx).Create(&problemInfo).Error
		if err != nil {
			return err
		}
		pid = problemInfo.Id
		return nil
	}
	if err = transaction.TransactionWithMongoDBAndMySQL(ctx, prepare, mongodb, mysql); err != nil {
		pid = 0
		return
	}
	return
}

func (p *ProblemService) Update(ctx context.Context, pid uint64, title string, level int, updatedBy string, timeLimit, memoryLimit int, markdown string) (err error) {
	// eg := &errgroup.Group{}
	// txSQL := global.MySQL.WithContext(ctx).Begin()
	// var sess mongo.Session
	// if sess, err = global.MongoDB.Client().StartSession(); err != nil {
	// 	txSQL.Rollback()
	// 	return
	// }
	// defer func() {
	// 	sess.EndSession(ctx)
	// }()

	// sess.StartTransaction()
	// defer func() {
	// 	if rec := recover(); rec != nil {
	// 		txSQL.Rollback()
	// 		sess.AbortTransaction(ctx)
	// 		err = fmt.Errorf("%v", rec)
	// 		return
	// 	}
	// 	if err != nil {
	// 		txSQL.Rollback()
	// 		sess.AbortTransaction(ctx)
	// 		return
	// 	}
	// 	txSQL.Commit()
	// 	sess.CommitTransaction(ctx)
	// }()

	// updates := map[string]any{}
	// if title != "" {
	// 	updates["title"] = title
	// }
	// if level > 0 && level < 4 {
	// 	updates["level"] = level
	// }
	// if timeLimit > 0 {
	// 	updates["time_limit"] = timeLimit
	// }
	// if memoryLimit > 0 {
	// 	updates["memory_limit"] = memoryLimit
	// }
	// updates["updated_by"] = updatedBy
	// // 如果 MySQL 仅更新 updated_by 字段，且 MongoDB 没有更新的内容，则不更新
	// if len(updates) == 1 && markdown == "" {
	// 	return nil
	// }

	// eg.Go(func() error {
	// 	return txSQL.WithContext(ctx).Model(&model.ProblemInfo{}).
	// 		Where("id = ?", pid).
	// 		Updates(updates).Error
	// })
	// eg.Go(func() error {
	// 	if markdown == "" {
	// 		return nil
	// 	}

	// 	res, err := global.MongoDB.Collection(model.ProblemContent{}.TableName()).UpdateOne(ctx, bson.M{
	// 		"pid": pid,
	// 	}, bson.M{
	// 		"$set": bson.M{
	// 			"markdown": markdown,
	// 		},
	// 	})
	// 	if err != nil {
	// 		return err
	// 	}
	// 	if res.ModifiedCount == 0 {
	// 		return errors.New("PID 错误")
	// 	}
	// 	return nil
	// })

	// err = eg.Wait()
	// return

	updates := map[string]any{}
	prepare := func() bool {
		if title != "" {
			updates["title"] = title
		}
		if level > 0 && level < 4 {
			updates["level"] = level
		}
		if timeLimit > 0 {
			updates["time_limit"] = timeLimit
		}
		if memoryLimit > 0 {
			updates["memory_limit"] = memoryLimit
		}
		updates["updated_by"] = updatedBy
		// 如果 MySQL 仅更新 updated_by 字段，且 MongoDB 没有更新的内容，则不更新
		if len(updates) == 1 && markdown == "" {
			return true
		}
		return false
	}
	mongodb := func() error {
		if markdown == "" {
			return nil
		}
		res, err := global.MongoDB.Collection(model.ProblemContent{}.TableName()).UpdateOne(ctx, bson.M{
			"pid": pid,
		}, bson.M{
			"$set": bson.M{
				"markdown": markdown,
			},
		})
		if err != nil {
			return err
		}
		if res.ModifiedCount == 0 {
			return errors.New("PID 错误")
		}
		return nil
	}
	mysql := func(txSQL *gorm.DB) error {
		return txSQL.WithContext(ctx).Model(&model.ProblemInfo{}).
			Where("id = ?", pid).
			Updates(updates).Error
	}
	if err = transaction.TransactionWithMongoDBAndMySQL(ctx, prepare, mongodb, mysql); err != nil {
		return
	}
	return
}

func (p *ProblemService) Enable(ctx context.Context, pid uint64, updatedBy string) (err error) {
	return global.MySQL.WithContext(ctx).Model(&model.ProblemInfo{}).
		Where("id = ?", pid).
		Updates(map[string]any{
			"updated_by": updatedBy,
			"enabled":    true,
		}).Error
}

func (p *ProblemService) Disable(ctx context.Context, pid uint64, updatedBy string) (err error) {
	return global.MySQL.WithContext(ctx).Model(&model.ProblemInfo{}).
		Where("id = ?", pid).
		Updates(map[string]any{
			"updated_by": updatedBy,
			"enabled":    false,
		}).Error
}

func (p *ProblemService) AppendTestCase(ctx context.Context, pid uint64, input, output string, score int, createdBy string) (tid string, err error) {
	prepare := func() bool {
		return false
	}
	mongodb := func() error {
		tid := fmt.Sprintf("%s:%s", pid, uuid.New().String())
		_, err := global.MongoDB.Collection(model.TestCase{}.TableName()).InsertOne(ctx, bson.M{
			"tid":       tid,
			"pid":       pid,
			"input":     input,
			"output":    output,
			"score":     score,
			"createdBy": createdBy,
			"updatedBy": createdBy,
			"enabled":   true,
		})
		return err
	}
	mysql := func(txSQL *gorm.DB) error {
		return txSQL.WithContext(ctx).Exec("UPDATE problem_info SET total_score = total_score + ?, total_test_case = total_test_case + 1 WHERE id = ?", score, pid).Error
	}
	if err = transaction.TransactionWithMongoDBAndMySQL(ctx, prepare, mongodb, mysql); err != nil {
		return
	}
	return
}

func (p *ProblemService) UpdateTestCase(ctx context.Context, pid uint64, tid, input, output string, score int, updatedBy string) (err error) {
	var old model.TestCase
	prepare := func() bool {
		if score > 0 {
			res := global.MongoDB.Collection(model.TestCase{}.TableName()).FindOne(ctx, bson.M{
				"tid": tid,
			})
			if res.Err() != nil {
				err = res.Err()
				return true
			}
			if err = res.Decode(&old); err != nil {
				return true
			}
		}
		return false
	}
	mongodb := func() error {
		set := bson.M{
			"$set": bson.M{
				"updatedBy": updatedBy,
			},
		}
		if input != "" {
			set["$set"].(bson.M)["inpput"] = input
		}
		if output != "" {
			set["$set"].(bson.M)["output"] = output
		}
		if score > 0 {
			set["$set"].(bson.M)["score"] = score
		}
		if len(set["$set"].(bson.M)) == 1 {
			return nil
		}
		res, err := global.MongoDB.Collection(model.TestCase{}.TableName()).UpdateOne(ctx, bson.M{
			"tid": tid,
		}, set)
		if err != nil {
			return err
		}
		if res.ModifiedCount == 0 {
			return errors.New("TID 错误")
		}
		return nil
	}
	mysql := func(txSQL *gorm.DB) error {
		if score > 0 {
			return txSQL.WithContext(ctx).
				Exec("UPDATE problem_info SET total_score = total_score - ? + ?, updated_by = ?",
					old.Score, score, updatedBy).
				Error
		}
		return nil
	}
	if err = transaction.TransactionWithMongoDBAndMySQL(ctx, prepare, mongodb, mysql); err != nil {
		return
	}
	return
}

func (p *ProblemService) EnableTestCase(ctx context.Context, pid uint64, tid, updatedBy string) (err error) {
	var old model.TestCase
	// 原本为启用状态的用例，不重复启用
	prepare := func() bool {
		res := global.MongoDB.Collection(model.TestCase{}.TableName()).FindOne(ctx, bson.M{
			"tid": tid,
		})
		if res.Err() != nil {
			err = res.Err()
			return true
		}
		if err = res.Decode(&old); err != nil {
			return true
		}
		return old.Enabled
	}
	mongodb := func() error {
		res, err := global.MongoDB.Collection(model.TestCase{}.TableName()).
			UpdateOne(ctx, bson.M{
				"tid": tid,
			}, bson.M{
				"$set": bson.M{
					"updatedBy": updatedBy,
					"enabled":   true,
				},
			})
		if err != nil {
			return err
		}
		if res.ModifiedCount == 0 {
			return errors.New("TID 错误")
		}
		return nil
	}
	mysql := func(txSQL *gorm.DB) error {
		return txSQL.WithContext(ctx).Exec("UPDATE problem_info SET total_score = total_score + ?, total_test_case = total_test_case + 1, updated_by = ?", old.Score, updatedBy).Error
	}
	if err = transaction.TransactionWithMongoDBAndMySQL(ctx, prepare, mongodb, mysql); err != nil {
		return
	}
	return
}

func (p *ProblemService) DisableTestCase(ctx context.Context, pid uint64, tid, updatedBy string) (err error) {
	var old model.TestCase
	// 原本为禁用状态的用例，不重复禁用
	prepare := func() bool {
		res := global.MongoDB.Collection(model.TestCase{}.TableName()).FindOne(ctx, bson.M{
			"tid": tid,
		})
		if res.Err() != nil {
			err = res.Err()
			return true
		}
		if err = res.Decode(&old); err != nil {
			return true
		}
		return !old.Enabled
	}
	mongodb := func() error {
		res, err := global.MongoDB.Collection(model.TestCase{}.TableName()).
			UpdateOne(ctx, bson.M{
				"tid": tid,
			}, bson.M{
				"$set": bson.M{
					"updatedBy": updatedBy,
					"enabled":   false,
				},
			})
		if err != nil {
			return err
		}
		if res.ModifiedCount == 0 {
			return errors.New("TID 错误")
		}
		return nil
	}
	mysql := func(txSQL *gorm.DB) error {
		return txSQL.WithContext(ctx).Exec("UPDATE problem_info SET total_score = total_score - ?, total_test_case = total_test_case - 1, updated_by = ?", old.Score, updatedBy).Error
	}
	if err = transaction.TransactionWithMongoDBAndMySQL(ctx, prepare, mongodb, mysql); err != nil {
		return
	}
	return
}
