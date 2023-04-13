package mongodb

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/spf13/cast"
	"github.com/sunmi-OS/gocore/v2/conf/viper"
	"github.com/sunmi-OS/gocore/v2/glog"
	"github.com/sunmi-OS/gocore/v2/utils"
	"github.com/sunmi-OS/gocore/v2/utils/closes"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var dbMap sync.Map
var closeOnce sync.Once

// NewDB 初始化db
func NewDB(dbname string) {
	var (
		orm *mongo.Database
		err error
	)
	oldConn, _ := dbMap.Load(dbname)
	err = utils.Retry(func() error {
		orm, err = openDB(dbname)
		if err != nil {
			glog.ErrorF("UpdateDB(%s) error:%+v", dbname, err)
			return err
		}
		return nil
	}, 5, 3*time.Second)
	if err != nil {
		panic(err)
	}
	dbMap.Delete(dbname)
	dbMap.Store(dbname, orm)
	if oldConn != nil {
		db, _ := oldConn.(*mongo.Database)
		if db != nil {
			err = db.Client().Disconnect(context.Background())
			if err != nil {
				panic(err)
			}
		}
	}
	closeOnce.Do(func() {
		closes.AddShutdown(closes.ModuleClose{
			Name:     "MongoDB Close",
			Priority: closes.MongoPriority,
			Func:     Close,
		})
	})
}

// GetDB 获取实例
func GetDB(dbname string) *mongo.Database {
	v, ok := dbMap.Load(dbname)
	if ok {
		return v.(*mongo.Database)
	}
	panic("The database is not initialized")
}

func Close() {
	dbMap.Range(func(dbName, orm interface{}) bool {
		glog.WarnF("close mongodb %s", dbName)
		dbMap.Delete(dbName)
		db, _ := orm.(*mongo.Database)
		if db != nil {
			_ = db.Client().Disconnect(context.Background())
		}
		return true
	})
}

func openDB(dbname string) (*mongo.Database, error) {
	// ["host:port" ...]
	dbEndpoint := viper.GetEnvConfig(dbname + ".Endpoint").SliceString()
	dbName := viper.GetEnvConfig(dbname + ".Name").String()
	dbUser := viper.GetEnvConfig(dbname + ".User").String()
	dbPasswd := viper.GetEnvConfig(dbname + ".Passwd").String()
	replicaSet := viper.GetEnvConfig(dbname + ".ReplicaSet").String()
	maxPoolSize := viper.GetEnvConfig(dbname + ".MaxPoolSize").Int64()
	minPoolSize := viper.GetEnvConfig(dbname + ".MinPoolSize").Int64()
	dbEndpointStr := ""
	if len(dbEndpoint) > 1 {
		dbEndpointStr = strings.Join(dbEndpoint, ",")
	} else {
		dbEndpointStr = dbEndpoint[0]
	}
	// single-node or sharded-clusters
	uri := fmt.Sprintf("mongodb://%s:%s@%s/admin",
		dbUser,
		dbPasswd,
		dbEndpointStr,
	)
	if replicaSet != "" {
		// replica set instance
		uri += "?replicaSet=" + replicaSet
	}
	opts := options.Client().ApplyURI(uri)
	opts.SetMaxPoolSize(cast.ToUint64(maxPoolSize))
	opts.SetMinPoolSize(cast.ToUint64(minPoolSize))
	client, err := mongo.Connect(context.TODO(), opts)
	if err != nil {
		return nil, err
	}
	return client.Database(dbName), nil
}
