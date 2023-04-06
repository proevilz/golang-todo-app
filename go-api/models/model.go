package models

type Todo struct {
	ID        string `gorm:"primaryKey" json:"id"`
	Title     string `gorm:"type:varchar(255);collation:utf8mb4_unicode_ci" json:"title"`
	Completed bool   `json:"completed"`
}
