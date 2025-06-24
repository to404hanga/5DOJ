package model

type User struct {
	Uid             string `gorm:"type:varchar(20);primaryKey;comment:学号"`     // 学号
	Name            string `gorm:"type:varchar(20);index:idx_name;comment:姓名"` // 姓名
	Password        string `gorm:"type:char(60);comment:bcrypt 加密密码"`          // bcrypt 加密密码
	TelephoneNumber string `gorm:"type:varchar(20);comment:手机号码"`              // 手机号码
	Gender          string `gorm:"type:varchar(10);comment:性别"`                // 性别
}
