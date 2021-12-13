package dao

type Student struct {
	ID        uint64 `gorm:"column:id"`
	StuId     string `gorm:"column:stu_id"`
	FirstName string `gorm:"column:first_name"`
	LastName  string `gorm:"column:last_name"`
	Email     string `gorm:"column:email"`
	PN        string `gorm:"column:phone_number"`
}

func (s Student) TableName() string {
	return "student"
}
