package dao

import (
	"bytes"
	"fmt"
	"ssq-spider/logger"
	"ssq-spider/model"
	"strings"
)

const (
	sysConfigTable = "sys_config_t"
	ssqNumberTable = "ssq_number_t"
)

func GetOneSysConfig(typeId string, name string) (string, error) {
	db, err := NewMysqlDBClient()
	if err != nil {
		logger.Logger.Error("NewMysqlDBClient error: ", err)
		return "", err
	}

	var sysConfig model.SysConfig

	if err := db.Table(sysConfigTable).Where("type_id=? and name=?",
		typeId, name).First(&sysConfig).Error; err != nil {
		logger.Logger.Error("query one sys config error:", err)
		return "", err
	}

	return sysConfig.Value, nil
}

func UpdateOneSysConfig(typeId string, name string, value string) error {
	db, err := NewMysqlDBClient()
	if err != nil {
		logger.Logger.Error("NewMysqlDBClient error: ", err)
		return err
	}
	return db.Table(sysConfigTable).Where("type_id=? and name=?",
		typeId, name).UpdateColumn("value", value).Error
}

func BulkSaveSSQOpenNumbers(ssqOpenNums *[]model.SSQOpenNumber) error {
	if len(*ssqOpenNums) == 0 {
		return nil
	}

	db, err := NewMysqlDBClient()
	if err != nil {
		logger.Logger.Error("NewMysqlDBClient error: ", err)
		return err
	}

	fields := []string{"open_no", "red_num", "blue_num", "ball_sort"}

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("insert into %s (%s) values ", ssqNumberTable, strings.Join(fields, ",")))

	dataFormat := "(\"%s\", \"%s\", \"%s\", \"%s\")"

	for _, ssqNum := range *ssqOpenNums {
		buffer.WriteString(fmt.Sprintf(dataFormat,
			ssqNum.OpenNo,
			ssqNum.RedNum,
			ssqNum.BlueNum,
			ssqNum.BallSort))
		buffer.WriteString(",")
	}

	SQL := buffer.String()[:len(buffer.String())-1] + ";"
	return db.Exec(SQL).Error
}
