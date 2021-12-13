package dc

import (
	"go-backend/dc/client"
	"go-backend/dc/dao"
	"go-backend/util/mysql_client"
)

var STUDao = &stuDao{
	c: client.DefaultMysqlClient,
}

type stuDao struct {
	c *mysql_client.Client
}

func (d *stuDao) Add(stu *dao.Student) error {
	return d.c.Create(&stu).Error
}
