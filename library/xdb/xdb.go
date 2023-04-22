package xdb

import (
	"context"
	"fmt"
	"gorm.io/gorm/logger"
	"os"
	"time"

	"github.com/casbin/casbin/v2"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/os/glog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

var (
	MasterDB *XDb
	Isexit   = false
)

type XDb struct {
	*gorm.DB
	Enforcer *casbin.Enforcer
}

type Config struct {
	Host      string
	Port      int
	Username  string
	Password  string
	Database  string
	Charset   string
	Collation string
	MaxIdle   string
	MaxConn   string
	Query     string
}

func Init() {
	//fmt.Println("xdb init")
	//fmt.Printf("set env")
	//os.Setenv("ZONEINFO", filepath.Join(getCurrentPath(), "zoneinfo.zip"))
	cfg := Config{}
	ce := g.Cfg().GetStruct("mysql", &cfg)
	if ce != nil {
		panic(ce)
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
		cfg.Username,
		cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
		cfg.Collation,
		cfg.Query,
	)
	show_dsn := fmt.Sprintf(
		"%s:******@tcp(%s:%d)/%s?charset=%s&collation=%s&%s",
		cfg.Username,
		//cfg.Password,
		cfg.Host,
		cfg.Port,
		cfg.Database,
		cfg.Charset,
		cfg.Collation,
		cfg.Query,
	)

	glog.Info("数据库连接DSN:", show_dsn)

	newLogger := New(
		//log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer（日志输出的目标，前缀和日志包含的内容）
		logger.Config{
			SlowThreshold:             200 * time.Millisecond, // 慢 SQL 阈值
			LogLevel:                  logger.Info,            // 日志级别
			IgnoreRecordNotFoundError: true,                   // 忽略ErrRecordNotFound（记录未找到）错误
			Colorful:                  false,                  // 禁用彩色打印
		},
	)

	db, err := gorm.Open(mysql.New(mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string 类型字段的默认长度
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据当前 MySQL 版本自动配置
	}), &gorm.Config{
		Logger: newLogger,
		NamingStrategy: schema.NamingStrategy{
			//TablePrefix:   "vod_", // 表名前缀，`User` 的表名应该是 `t_users`
			SingularTable: true, // 使用单数表名，启用该选项，此时，`User` 的表名应该是 `t_user`
		},
	})
	if err != nil {
		glog.Info(dsn)
		panic(err)
	}
	//a, _ := NewAdapterByDB(db)
	//e, err := casbin.NewEnforcer("config/rbac_model.conf", a)
	//if err != nil {
	//	panic(err)
	//	//todo: 暂时不用rbac
	//}
	go keepAlive(db)
	MasterDB = &XDb{
		DB: db,
		//Enforcer: e,
	}
}

func Default() *XDb {
	if MasterDB == nil {
		panic(gerror.New("Not Open Db"))
	}
	return MasterDB
}

func GetTable() {
	row := MasterDB.Raw(`
	SELECT
		* 
	FROM
		information_schema.TABLES 
	WHERE
		table_schema = 'cc_center';
	`).Row()
	g.Dump(row)
}

type ctxTransactionKey struct{}

func CtxWithTransaction(ctx context.Context, tx *gorm.DB) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}
	return context.WithValue(ctx, ctxTransactionKey{}, tx)
}

func Transaction(ctx context.Context, fc func(txctx context.Context) error) error {
	db := MasterDB.WithContext(ctx)

	return db.Transaction(func(tx *gorm.DB) error {
		txctx := CtxWithTransaction(ctx, tx)
		return fc(txctx)
	})
}

func getCurrentPath() string {
	getwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	return getwd
}

//简单事务封装
func Tx(funcs ...func(db *gorm.DB) error) (err error) {
	//func Tx(funcs ...func(db *XDb) error) (err error) {
	tx := MasterDB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			err = fmt.Errorf("%v", err)
		}
	}()
	for _, f := range funcs {
		err = f(tx)
		if err != nil {
			tx.Rollback()
			return
		}
	}
	err = tx.Commit().Error
	return
}

//數據庫保連保连
func keepAlive(dbc *gorm.DB) {
	for {
		sqlDB, err := dbc.DB()
		if err == nil {
			err = sqlDB.Ping()
			if err != nil {
				glog.Warning("數據庫连接已被断开")
				if !Isexit {
					Isexit = true
					go func() {
						glog.Warning("等待5s退出进程")
						time.Sleep(time.Second * 5)
						os.Exit(-1)
					}()
				}
				break
			}
		}
		time.Sleep(60 * time.Second)
	}
}
