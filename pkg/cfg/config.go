package cfg

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"

	"lium-product/es-search/pkg/utils"
	testcommon "lium-product/es-search/tests/common_test"
)

type Cfg struct {
	// Common 通用配置
	Common Common `json:"common"`

	// Redis redis配置
	Redis Redis `json:"redis"`

	// MySql mysql配置
	MySql MySql `json:"mysql"`

	// ElasticSearch ElasticSearch配置
	ElasticSearch ElasticSearch `json:"elastic"`

	// Logger 日志配置
	Logger Logger `json:"logger"`

	// Jwt 配置
	Jwt Jwt `json:"jwt"`
}

var (
	once               sync.Once
	cfgLock            sync.RWMutex
	singleInstanceConf *Cfg
)

func InitLoadCfg() {
	GetInstance()
}

func pf(name string, err error) string {
	return fmt.Sprintf("load %s config failed: %v", name, err)
}

// SetInstance 设置配置实例，测试用
func SetInstance(c *Cfg) {
	singleInstanceConf = c
}

func GetInstance() *Cfg {
	cfgLock.Lock()
	defer cfgLock.Unlock()
	if singleInstanceConf == nil {
		// 加载所有配置
		once.Do(func() { singleInstanceConf = loadCfg() })
	}
	return singleInstanceConf
}

func loadCfg() *Cfg {
	conf := &Cfg{}
	// 读取本地配置文件，这个是加载服务器上的 config.json文件的配置
	onInitPath(conf)
	// 更新默认配置
	RefCfgDefVal(conf)

	return conf
}

func onInitPath(conf *Cfg) {
	pwd, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	paths := []string{
		"./config.json",
		"../config.json",
		"../../config.json",
		"../../../config.json",
		"../../../../config.json",
		"../../../../../config.json",
		filepath.Dir(pwd) + "/../config.json",
	}

	var configPath string
	for i := range paths {
		if ok, _ := utils.PathExists(paths[i]); ok {
			configPath = paths[i]
			break
		}
	}
	fmt.Printf("load config template path : %s \n", configPath)

	// 加载配置文件
	data, err := ioutil.ReadFile(configPath)
	if err != nil {
		panic(pf("file", err))
	}

	if err := json.Unmarshal(data, &conf); err != nil {
		panic(pf("json Unmarshal", err))
	}
	fmt.Printf("load config template : %+v \n", conf)
}

// RefCfgDefVal 系统默认配置
func RefCfgDefVal(confs ...*Cfg) {
	if testcommon.IsTest() {
		configure := testcommon.MockConsulConfig()

		cfgStruct := &Cfg{}
		json.Unmarshal(configure, cfgStruct)
		singleInstanceConf = cfgStruct
	}
}
