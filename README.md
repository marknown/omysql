# omysql
omysql is a init for gorm

## Mysql Config
```
type Config struct {
    DataSourceName  string // 数据库连接信息
    SingularTable   bool   // 取消默认的 user 结构体到 users 的差异
    LogMode         bool   // 是否开始日志模式，可以看日志
    MaxIdleConns    int    // SetMaxIdleConns设置连接池中的最大闲置连接数
    MaxOpenConns    int    // SetMaxOpenConns设置与数据库建立连接的最大数目
    ConnMaxLifetime int    // 这个值小于mysql的默认值会自动重连 mysql show global variables like 'wait_timeout'
}

example config

config := &Config{
    DataSourceName : "root:password@tcp(127.0.0.1:3306)/dbname?charset=utf8&parseTime=true&loc=Asia%2FShanghai",
    SingularTable  : true,
    LogMode        : true,
    MaxIdleConns   : 10,
    MaxOpenConns   : 100,
    ConnMaxLifetime: 3600
}
```

## Mysql Instance function
```
func GetInstance(config Config) *gorm.DB {
}
```

## Usage
```
db := omysql.GetInstance(config)
```
