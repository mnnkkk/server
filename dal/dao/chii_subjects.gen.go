// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package dao

const TableNameSubject = "chii_subjects"

// Subject mapped from table <chii_subjects>
type Subject struct {
	ID          uint32       `gorm:"column:subject_id;type:mediumint(8) unsigned;primaryKey;autoIncrement:true"`
	TypeID      uint8        `gorm:"column:subject_type_id;type:smallint(6) unsigned;not null"`
	Name        string       `gorm:"column:subject_name;type:varchar(80);not null"`
	NameCN      string       `gorm:"column:subject_name_cn;type:varchar(80);not null"`
	UID         string       `gorm:"column:subject_uid;type:varchar(20);not null"` // isbn / imdb
	Creator     uint32       `gorm:"column:subject_creator;type:mediumint(8) unsigned;not null"`
	Dateline    uint32       `gorm:"column:subject_dateline;type:int(10) unsigned;not null"`
	Image       string       `gorm:"column:subject_image;type:varchar(255);not null"`
	Platform    uint16       `gorm:"column:subject_platform;type:smallint(6) unsigned;not null"`
	Infobox     string       `gorm:"column:field_infobox;type:mediumtext;not null"`
	Summary     string       `gorm:"column:field_summary;type:mediumtext;not null"`            // summary
	Field5      string       `gorm:"column:field_5;type:mediumtext;not null"`                  // author summary
	Volumes     uint32       `gorm:"column:field_volumes;type:mediumint(8) unsigned;not null"` // 卷数
	Eps         uint32       `gorm:"column:field_eps;type:mediumint(8) unsigned;not null"`
	Wish        uint32       `gorm:"column:subject_wish;type:mediumint(8) unsigned;not null"`
	Done        uint32       `gorm:"column:subject_collect;type:mediumint(8) unsigned;not null"`
	Doing       uint32       `gorm:"column:subject_doing;type:mediumint(8) unsigned;not null"`
	OnHold      uint32       `gorm:"column:subject_on_hold;type:mediumint(8) unsigned;not null"` // 搁置人数
	Dropped     uint32       `gorm:"column:subject_dropped;type:mediumint(8) unsigned;not null"` // 抛弃人数
	Series      bool         `gorm:"column:subject_series;type:tinyint(1) unsigned;not null"`
	SeriesEntry uint32       `gorm:"column:subject_series_entry;type:mediumint(8) unsigned;not null"`
	IdxCn       string       `gorm:"column:subject_idx_cn;type:varchar(1);not null"`
	Airtime     uint8        `gorm:"column:subject_airtime;type:tinyint(1) unsigned;not null"`
	Nsfw        bool         `gorm:"column:subject_nsfw;type:tinyint(1);not null"`
	Ban         uint8        `gorm:"column:subject_ban;type:tinyint(1) unsigned;not null"`
	Fields      SubjectField `gorm:"foreignKey:subject_id;references:field_sid" json:"fields"`
}

// TableName Subject's table name
func (*Subject) TableName() string {
	return TableNameSubject
}
