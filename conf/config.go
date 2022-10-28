package conf

import (
	"context"
	sql "database/sql"
	"fmt"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

//全局config实例对象也就是程序在内存中的配置对象
//程序内都通过读取此对象来获取配置
//运行时仅初始化一次，避免反复初始化
//该对象什么时候被初始化？——>加载配置时也就是执行LoadConfigFromToml和LoadConfigFromEnv函数时
//为了不被程序在运行时恶意修改，设置成首字母小写的私有变量
var config *Config

//由于config全局配置为私有变量，单独提供下面的获取配置函数C
func C() *Config {
	return config
}

// Config 应用配置
//通过封装一个对象，来与外部配置(本项目使用的配置文件、环境变量)进行对接
type Config struct {
	App   *App   `toml:"app"` //toml文件中的[app]项
	Log   *Log   `toml:"log"`
	MySQL *MySQL `toml:"mysql"`
}

//获取有默认值的Config实例函数
//保证了即使不提供任何配置的情况下也能运行并发现问题如提示默认值不正确
func NewDefaultConfig() *Config {
	return &Config{
		App:   NewDefaultApp(),
		Log:   NewDefaultLog(),
		MySQL: NewDefaultMysql(),
	}
}

//App配置对象
type App struct {
	//同时支持配置文件和环境变量
	Name string `toml:"name" env:"APP_NAME"`
	Host string `toml:"host" env:"APP_HOST"`
	Port string `toml:"port" env:"APP_PORT"`
}

//配置默认值
func NewDefaultApp() *App {
	return &App{
		Name: "demo",
		Host: "127.0.0.1",
		Port: "8050",
	}
}

// MySQL配置对象
type MySQL struct {
	Host     string `toml:"host" env:"MYSQL_HOST"`
	Port     string `toml:"port" env:"MYSQL_PORT"`
	UserName string `toml:"username" env:"MYSQL_USERNAME"`
	Password string `toml:"password" env:"MYSQL_PASSWORD"`
	Database string `toml:"database" env:"MYSQL_DATABASE"`
	//下面的参数是针对数据库连接的优化配置，通常不用配置，使用默认值就好
	//控制当前mysql打开的最大连接数
	MaxOpenConn int `toml:"max_open_conn" env:"MYSQL_MAX_OPEN_CONN"`
	//控制mysql复用，比如设置成5表示最多运行5个来复用
	MaxIdleConn int `toml:"max_idle_conn" env:"MYSQL_MAX_IDLE_CONN"`
	//mysql连接的生命周期，与mysql配置相关，必须小于mysql的配置
	//一个conn用12小时换一个conn，保证一定的可用性
	MaxLifeTime int `toml:"max_life_time" env:"MYSQL_MAX_LIFE_TIME"`
	//连接最多允许存活多久
	MaxIdleTime int `toml:"max_idle_time" env:"MYSQL_MAX_idle_TIME"`
	lock        sync.Mutex
}

//获取mysql默认配置实例函数 用于生成数据库sql.DB实例
func NewDefaultMysql()  *MySQL {
	return &MySQL{
		Host:        "192.168.0.206",
		Port:        "3306",
		UserName:    "root",
		Password:    "Password1",
		Database:    "api_demo",
		MaxOpenConn: 200,
		MaxIdleConn: 100,
	}
}

//mysql连接池配置方法，返回sql.DB,为了防止运行时修改，不对外提供，使用下面的GetDB方法获取
//*sql.DB中的属性freeConn维护着一个连接池对象[]*driverConn，确保里面的连接都是可用的，定期检查conn健康
//如果某一个conn失效，会调用driverConn.resetSession方法重置，清空driverConn结构体数据
//这样避免了driverConn结构体的内存申请和释放成本
func (m *MySQL) getDBconn() (*sql.DB, error) {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&multiStatements=true", m.UserName, m.Password, m.Host, m.Port, m.Database)
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("connect to mysql<%s> error, %s", dsn, err.Error())
	}
	db.SetMaxOpenConns(m.MaxOpenConn)
	db.SetMaxIdleConns(m.MaxIdleConn)
	//创建超时时间的context
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	//程序退出前也要执行上下文取消函数
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("ping mysql<%s> error, %s", dsn, err.Error())
	}
	return db, nil
}

//初始化DB实例方法
//在加载全局配置config实例后，使用conf.C().MySQL.GetDB()方法，获取DB时，动态判断再初始化
func (m *MySQL) GetDB() *sql.DB {
	//加锁，锁住临界值？
	m.lock.Lock()
	defer m.lock.Unlock()
	//如果db未初始化，先初始化
	if db == nil {
		conn, err := m.getDBconn()
		if err != nil {
			panic(err)
		}
		db = conn
	}
	//全局变量一定存在
	return db
}

//全局DB实例，同全局配置config一样，运行时仅初始化一次
var db *sql.DB

//日志配置对象
type Log struct {
	Level   string    `toml:"level" env:"LOG_LEVEL"`
	Format  LogFormat `toml:"format" env:"LOG_FORMAT"`
	To      LogTo     `toml:"to" env:"LOG_TO"`
	PathDir string    `toml:"path_dir" env:"LOG_PATH_DIR"`
}

//log默认配置函数
func NewDefaultLog() *Log {
	return &Log{
		//日志级别支持debug,info,warn,error
		Level: "info",
		//日志格式支持json和text
		Format: TextFormat,
		//日志输出支持stdout和file
		To: ToStdout,
	}
}
