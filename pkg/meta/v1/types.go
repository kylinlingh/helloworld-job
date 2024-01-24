package v1

import (
	"encoding/json"
	"gorm.io/gorm"
	"time"
)

// Extend defines a new type used to store extended fields.
type Extend map[string]interface{}

// String returns the string format of Extend.
func (ext Extend) String() string {
	data, _ := json.Marshal(ext)
	return string(data)
}

// Merge merge extend fields from extendShadow.
func (ext Extend) Merge(extendShadow string) Extend {
	var extend Extend

	// always trust the extendShadow in the database
	_ = json.Unmarshal([]byte(extendShadow), &extend)
	for k, v := range extend {
		if _, ok := ext[k]; !ok {
			ext[k] = v
		}
	}

	return ext
}

// ObjectMeta 统一资源数据
type ObjectMeta struct {
	ID uint64 `json:"id,omitempty" gorm:"primary_key;AUTO_INCREMENT;column:id"`

	// InstanceID defines a string type resource identifier,
	// use prefixed to distinguish resource types, easy to remember, Url-friendly.
	InstanceID string `json:"instanceID,omitempty" gorm:"unique;column:instance_id;type:varchar(32)"`

	// Extend store the fields that need to be added, but do not want to add a new table column, will not be stored in db.
	Extend Extend `json:"extend,omitempty" gorm:"-" validate:"omitempty"`

	// ExtendShadow is the shadow of Extend. DO NOT modify directly.
	ExtendShadow string `json:"-" gorm:"column:extend_shadow" validate:"omitempty"`

	CreatedAt time.Time `json:"createdAt,omitempty" gorm:"column:created_at"`

	UpdatedAt time.Time `json:"updatedAt,omitempty" gorm:"column:updated_at"`

	DeletedAt gorm.DeletedAt `json:"deletedAt,omitempty" gorm:"column:deleted_at"`
}

// BeforeCreate run before create database record.
func (obj *ObjectMeta) BeforeCreate(tx *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()

	return nil
}

// BeforeUpdate run before update database record.
func (obj *ObjectMeta) BeforeUpdate(tx *gorm.DB) error {
	obj.ExtendShadow = obj.Extend.String()

	return nil
}

// AfterFind run after find to unmarshal a extend shadown string into metav1.Extend struct.
func (obj *ObjectMeta) AfterFind(tx *gorm.DB) error {
	if err := json.Unmarshal([]byte(obj.ExtendShadow), &obj.Extend); err != nil {
		return err
	}

	return nil
}

type ListOptions struct {
	// LabelSelector is used to find matching REST resources.
	LabelSelector string `json:"labelSelector,omitempty" form:"labelSelector"`

	// FieldSelector restricts the list of returned objects by their fields. Defaults to everything.
	FieldSelector string `json:"fieldSelector,omitempty" form:"fieldSelector"`

	// TimeoutSeconds specifies the seconds of ClientIP type session sticky time.
	TimeoutSeconds *int64 `json:"timeoutSeconds,omitempty"`

	// Offset specify the number of records to skip before starting to return the records.
	Offset *int64 `json:"offset,omitempty" form:"offset"`

	// Limit specify the number of records to be retrieved.
	Limit *int64 `json:"limit,omitempty" form:"limit"`
}

// ListMeta describes metadata that synthetic resources must have, including lists and
// various status objects. A resource may have only one of {ObjectMeta, ListMeta}.
type ListMeta struct {
	TotalCount int64 `json:"totalCount,omitempty"`
}
