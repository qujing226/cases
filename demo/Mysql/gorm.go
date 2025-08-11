package Mysql

import (
	"context"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"time"
)

// IncrReadCnt 增加阅读数
// clause.OnConflict 表示在冲突时更新 "Upsert"
func IncrReadCnt(ctx context.Context, id int) error {
	dao, _ := gorm.Open(mysql.Open(""))
	err := dao.WithContext(ctx).Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "id"}}, // 明确冲突列
		DoUpdates: clause.Assignments(map[string]interface{}{
			"read_cnt": gorm.Expr("`read_cnt` + 1"),
			"utime":    time.Now().UnixMilli(),
		}),
	}).Create(&struct {
		Id      int
		ReadCnt int
		Utime   int64
	}{
		Id:      id,
		ReadCnt: 1,
		Utime:   time.Now().UnixMilli(),
	}).Error
	return err
}

// SelectReadCnt 查询阅读数
func SelectReadCnt(ctx context.Context, id int) (int, error) {
	// 优化 SQL，以命中索引直接在索引上就计算出来总数，而不需要回表。
	dao, _ := gorm.Open(mysql.Open(""))
	var cnt int
	err := dao.WithContext(ctx).Model(&struct {
		ReadCnt int
	}{}).Select("COUNT(*)").Where("id > ?", 0).Scan(&cnt).Error
	return cnt, err
}
