package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/godror/godror"
	"github.com/tealeg/xlsx"
)

var oracleHost, oraclePort, oracleUser, oraclePassword, oracleServiceName *string
var excelFilePath *string
var db *sql.DB

func main() {
	oracleHost = flag.String("h", "127.0.0.1", "oracle host")
	oracleUser = flag.String("u", "", "oracle user name")
	oraclePassword = flag.String("p", "", "oracle password")
	oracleServiceName = flag.String("s", "", "oracle service name")
	oraclePort = flag.String("P", "1521", "oracle port")
	excelFilePath = flag.String("t", "/tmp/databse.xlsx", "export xlsx path and filename")
	flag.Parse()
	// fmt.Println(*oracleHost)

	if *oracleHost == "" || *oracleUser == "" || *excelFilePath == "" || *oracleUser == "" || *oraclePassword == "" || *oracleServiceName == "" {
		flag.PrintDefaults()
		return
	}
	dns := `user="` + *oracleUser + `" password="` + *oraclePassword + `" connectString="` + *oracleHost + `:` + *oraclePort + `/` + *oracleServiceName + `"`
	// db, err := sql.Open("godror", `user="oracle" password="oracle" connectString="127.0.0.1:1521/test1"`)
	db, err := sql.Open("godror", dns)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	// excelFilePath = "./databse.xlsx"
	excelAbsFilePath, err := filepath.Abs(*excelFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	reader := bufio.NewReader(os.Stdin)

	//获取Sql语句
	var sqlStr string
	fmt.Printf("Please Input SQL:\n")
	for {
		tmpSql, _ := reader.ReadString('\n')
		sqlStr = sqlStr + tmpSql
		tmpSql = strings.TrimSpace(tmpSql)
		if tmpSql[len(tmpSql)-1] == ';' {
			break
		}
	}
	sqlStr = strings.TrimSpace(sqlStr)
	//去掉末尾分号，要不报错ORA-00911:无效的字符错误
	sqlStr1 := sqlStr[:len(sqlStr)-1]

	stmt, err := db.Prepare(sqlStr1)
	if err != nil {
		fmt.Println(err)
		return
	}

	defer stmt.Close()

	result, err := stmt.Query()
	if err != nil {
		fmt.Println(err)
		return
	}

	//保存excel
	saveExcelByRows(excelAbsFilePath, result)
}

func saveExcelByRows(excelAbsFilePath string, rows *sql.Rows) error {
	file := xlsx.NewFile()
	sheet, err := file.AddSheet("result")
	if err != nil {
		return err
	}

	columns, err := rows.Columns()
	if err != nil {
		return err
	}

	//写column数据
	columnRow := sheet.AddRow()
	columnLen := len(columns)
	for _, name := range columns {
		cell := columnRow.AddCell()
		cell.Value = name
	}

	scanArgs := make([]interface{}, columnLen)
	values := make([][]byte, columnLen)
	for i := range values {
		scanArgs[i] = &values[i]
	}
	for rows.Next() {
		rows.Scan(scanArgs...)
		row := sheet.AddRow()
		for _, v := range values {
			cell := row.AddCell()
			cell.Value = string(v)
		}
	}

	//保存文件
	err = file.Save(excelAbsFilePath)
	if err != nil {
		return err
	} else {
		return nil
	}

}
