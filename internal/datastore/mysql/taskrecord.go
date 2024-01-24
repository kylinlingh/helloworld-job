package mysql

import (
	"context"
	"gorm.io/gorm"
	"helloworld/internal/entity"
	"helloworld/pkg/gormutil"
	v1 "helloworld/pkg/meta/v1"
)

type taskrecord struct {
	db *gorm.DB
}

func (t *taskrecord) Create(ctx context.Context, tr *entity.TaskRecord) error {
	return t.db.Create(&tr).Error
}

func (t *taskrecord) List(ctx context.Context, opts v1.ListOptions) (*entity.TaskRecordList, error) {
	ret := &entity.TaskRecordList{}
	ol := gormutil.Unpointer(opts.Offset, opts.Limit)
	d := t.db.Table("recod").Find(&ret.Records).Offset(ol.Offset).Limit(ol.Limit).Order("id desc").Count(&ret.TotalCount)
	return ret, d.Error
}

func newTaskRecord(ds *mysqlStore) *taskrecord {
	return &taskrecord{
		ds.db,
	}
}
