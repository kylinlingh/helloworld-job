package entity

import (
	"gorm.io/gorm"
	"helloworld/internal/utils"
	v1 "helloworld/pkg/meta/v1"
)

type TaskRecord struct {
	v1.ObjectMeta
	File     string `json:"file" gorm:"column:file"`
	ScanTime string `json:"scan_time" gorm:"column:scan_time"`
}

type TaskRecordList struct {
	v1.ListMeta `json:",inline"`
	Records     []*TaskRecord `json:"records"`
}

// AfterCreate Hooks 在记录插入数据库之后，生成并更新到数据库的 instanceID字段
func (t *TaskRecord) AfterCreate(tx *gorm.DB) error {
	t.InstanceID = utils.GetInstanceID(t.ID, "record-")
	return tx.Save(t).Error
}
