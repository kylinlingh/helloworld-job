package entity

type TaskRecord struct {
	FindingID    string  `json:"finding_id" gorm:"column:finding_id"`
	TaskId       string  `json:"task_id" gorm:"column:celery_task_id"`
	Repo         string  `json:"repo" gorm:"column:repo"`
	Description  string  `json:"description" gorm:"column:description"`
	RuleID       string  `json:"rule_id" gorm:"column:rule_id"`
	File         string  `json:"file" gorm:"column:file"`
	Match        string  `json:"match" gorm:"column:match"`
	Match2Redact string  `json:"match_redact" gorm:"column:match_redact"`
	Secret       string  `json:"secret" gorm:"column:secret"`
	Entropy      float32 `json:"entropy" gorm:"column:entropy"`
	Redact       bool    `json:"redact" gorm:"column:redact"`
	ScanTime     string  `json:"scan_time" gorm:"column:scan_time"`
}
