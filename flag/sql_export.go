package flag

import (
	"ThinkTankCentral/global"
	"fmt"
	"os"
	"os/exec"
	"time"
)

/**
 * 使用 mysqldump 工具将 MySQL 数据库导出为 .sql 文件，适合备份数据
 */

// SQLExport 导出 MySQL 数据
func SQLExport() error {
	//获取Mysql配置文件
	mysql := global.Config.Mysql

	//生成导出文件名，mysql_20250316
	timer := time.Now().Format("20060102")
	sqlPath := fmt.Sprintf("mysql_%s.sql", timer)

	//构造命令
	//通过 docker exec 进入 MySQL 容器，使用 mysqldump 工具导出指定数据库的数据
	cmd := exec.Command("docker", "exec", "mysql", "mysqldump", "-u"+mysql.Username, "-p"+mysql.Password, mysql.DBName)

	//创建输出文件
	//创建一个新的 .sql 文件，用于保存导出的数据
	outFile, err := os.Create(sqlPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	//执行导出并写入文件
	cmd.Stdout = outFile
	return cmd.Run()
}
