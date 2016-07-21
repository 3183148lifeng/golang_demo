package main

import (
	// "flag"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	// "io"
)

type User struct {
	Id   int
	Name string
	Age  int
	Sex  int
}

func init() {

	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", "root:fli123@tcp(127.0.0.1:3306)/fli?charset=utf8")
	orm.RegisterModel(new(User))
	orm.RunSyncdb("default", false, true)
}

func main() {

	o := orm.NewOrm()
	user := &User{
		Name: "fli",
		Age:  18,
		Sex:  1,
	}

	fmt.Println(o.Insert(user))
}
