package flag

import (
	"ThinkTankCentral/global"
	"os"
	"strings"
)

/**
 *从 .sql 文件中读取并逐条执行 SQL 语句，将数据导入数据库
 */

// SQLImport 导入 MySQL 数据
func SQLImport(sqlPath string) (errs []error) {
	//读取文件内容
	byteData, err := os.ReadFile(sqlPath)
	if err != nil {
		return append(errs, err)
	}

	// 分割 SQL 语句
	sqlList := strings.Split(string(byteData), ";")
	for _, sql := range sqlList {
		// 去除字符串开头和结尾的空白符
		sql = strings.TrimSpace(sql)
		if sql == "" {
			continue
		}
		// 执行sql语句
		err = global.DB.Exec(sql).Error
		if err != nil {
			errs = append(errs, err)
			continue
		}
	}
	return nil
}
