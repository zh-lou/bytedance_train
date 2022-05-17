package main

import (
	"fmt"

	// "github.com/jinzhu/gorm/dialects/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type User struct {
	UUID    string
	Name    string
	Age     int
	Version int
}

func main() {
	// db, err := gorm.Open(mysql.New(mysql.Config{
	// 	DSN: "root:hallo2014@tcp(127.0.0.1:3306)/hw513?charset=utf8&parseTime=True&loc=Local",
	// 	// DefaultStringSize:         256,   // string 类型字段的默认长度
	// 	// DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	// 	// DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
	// 	// DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	// 	// SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	// }), &gorm.Config{})

	// dsn指的是数据库的配置用户名、密码、通信协议、数据库名字
	dsn := "root:hallo2014@tcp(127.0.0.1:3306)/hw513?charset=utf8&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	// DefaultStringSize:         256,   // string 类型字段的默认长度
	// DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
	// DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
	// DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
	// SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置

	if err != nil {
		fmt.Println(err)
		return
	}

	// 创建表 自动迁移 （把结构体和数据表进行对应）
	err = db.AutoMigrate(&User{})

	if err != nil {
		fmt.Println(err)
		return
	}

	// Create
	user1 := User{"01", "byte", 22, 1}
	user2 := User{"02", "dance", 22, 1}
	db.Create(user1)
	db.Create(user2)

	// Read
	var user User
	db.First(&user, 1) // 查询获得第一条数据保存到user
	fmt.Printf("user:%#v\n", user)

	var user_ User
	db.First(&user_, "Name = ?", "dance") // 查找 名字为dance的数据
	fmt.Printf("user:%#v\n", user)

	// Update
	db.Model(&user_).Where("Name = ?", "dance").Update("Version", 2) //- 将user中的版本改为2
	fmt.Printf("user:%#v\n", user_)

	// Delete
	db.Delete(&user, 1) // 删除第一条数据
}
