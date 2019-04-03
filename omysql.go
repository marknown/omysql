// Package omysql is own mysql
package omysql

import (
	"crypto/md5"
	"fmt"
	"log"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

// Config mysql 的配置信息
type Config struct {
	DataSourceName  string // 数据库连接信息
	SingularTable   bool   // 取消默认的 user 结构体到 users 的差异
	LogMode         bool   // 是否开始日志模式，可以看日志
	MaxIdleConns    int    // SetMaxIdleConns设置连接池中的最大闲置连接数
	MaxOpenConns    int    // SetMaxOpenConns设置与数据库建立连接的最大数目
	ConnMaxLifetime int    // 这个值小于mysql的默认值会自动重连 mysql show global variables like 'wait_timeout'
}

// 包内变量，存储实例相关对象
var packageOnce = map[string]*sync.Once{}
var packageInstance = map[string]*gorm.DB{}
var packageMutex = &sync.Mutex{}

// GetInstance 根据配置信息初始化gorm对象 只初始化一次
func GetInstance(config Config) *gorm.DB {
	packageMutex.Lock()
	defer packageMutex.Unlock()

	md5byte := md5.Sum([]byte(config.DataSourceName))
	md5key := fmt.Sprintf("%x", md5byte)

	// 如果有值直接返回
	if v, ok := packageInstance[md5key]; ok {
		// fmt.Println("return direct")
		return v
	}

	// 如果once 不存在
	if _, ok := packageOnce[md5key]; !ok {
		var once = &sync.Once{}
		var db *gorm.DB
		var err error
		once.Do(func() {
			db, err = gorm.Open("mysql", config.DataSourceName)
			// 多go程，不关闭
			// defer db.Close()

			if nil != err {
				log.Println(err.Error())
			}

			// 是否记录日志
			db.LogMode(config.LogMode)
			// 取消默认的 user 结构体到 users 的差异
			db.SingularTable(config.SingularTable)

			// SetMaxIdleConns设置连接池中的最大闲置连接数
			db.DB().SetMaxIdleConns(config.MaxIdleConns)
			// SetMaxOpenConns设置与数据库建立连接的最大数目
			db.DB().SetMaxOpenConns(config.MaxOpenConns)
			// 这个值小于mysql的默认值会自动重连 mysql show global variables like 'wait_timeout'
			db.DB().SetConnMaxLifetime(time.Duration(int64(config.ConnMaxLifetime)) * time.Second)

			packageInstance[md5key] = db
			packageOnce[md5key] = once
			// fmt.Printf("init %p %v\n", db, db)
		})

		return db
	}

	return nil
}
