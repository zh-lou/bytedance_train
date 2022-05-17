package main

import (
	"context"
	"fmt"
	"math/rand"
	dal_model "project0513/task2/dal/model"
	"project0513/task2/dal/query"
	"project0513/task2/model"

	"github.com/google/uuid"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

type User struct {
	UUID    string
	Name    string
	Age     int
	Version int
}

func main() {
	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       "root:hallo2014@tcp(127.0.0.1:3306)/hw513",
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{})
	if err != nil {
		println(err)
		return
	}

	// 迁移 schema
	err = db.AutoMigrate(&User{})

	if err != nil {
		fmt.Println(err)
		return
	}

	//gen自动生成people
	g := gen.NewGenerator(gen.Config{
		OutPath:      "./task2/dal/query",
		ModelPkgPath: "./task2/dal/model", //默认同一路径
	})

	g.UseDB(db)

	//从数据库生成模型捕获表信息，返回 BaseStruct
	peopleTbl := g.GenerateModelAs("users", "People")

	//为指定的结构体或表格生成基础CRUD查询方法，
	g.ApplyBasic(
		peopleTbl,
	)
	//为指定的数据库表实现除基础方法外的相关方法，扩展接口
	g.ApplyInterface(func(Method model.Method) {}, peopleTbl)

	//执行并生成代码
	g.Execute()

	err = generatePeople(db)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = GetMaxVersionCount(db)
	if err != nil {
		fmt.Println(err)
		return
	}

}

// 插入数据
func generatePeople(db *gorm.DB) error {

	q := query.Use(db)
	arr := make([]string, 50) //定义长度为50的string
	for i := range arr {
		u1, err := uuid.NewUUID() //uuid
		if err != nil {
			fmt.Println(err)
			return err
		}
		arr[i] = u1.String()
	}

	for i := 0; i < 100; i++ {
		j := rand.Intn(50)

		u, err := q.People.WithContext(context.Background()).Debug().Where(q.People.UUID.Eq(arr[j])).First()
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				err = q.People.WithContext(context.Background()).Create(&dal_model.People{
					UUID:    arr[j],
					Name:    "2",
					Age:     1,
					Version: 1,
				})
				if err != nil {
					fmt.Printf("Create fail:%v\n", err)
					return err
				}
				continue
			}
			fmt.Printf("Find fail:%v\n", err)
			return err
		}

		_, err = q.People.WithContext(context.Background()).Where(q.People.UUID.Eq(arr[j])).Update(q.People.Version, u.Version+1)
		if err != nil {
			fmt.Printf("Update fail:%v\n", err)

			return err
		}
	}
	return nil
}

// 获取最大版本号
func GetMaxVersionCount(db *gorm.DB) error {
	q := query.Use(db)
	people, err := q.People.WithContext(context.Background()).Debug().GetMaxVersionCount()
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println(people)
	return nil
}
