package testcommon

import (
	"encoding/json"
	"flag"
	"runtime"
	"strings"
)

func IsTest() bool {
	if !testFlag {
		return false
	}
	for _, arg := range flag.Args() {
		if arg == "-test.v" || arg == "-test.run" || arg == "-test.timeout" || arg == "-test.coverprofile" {
			return true
		}
	}

	buf := make([]byte, 1<<16)
	runtime.Stack(buf, false)
	stack := string(buf)
	return strings.Contains(stack, "testing.tRunner")
}

var testFlag = true

func SetTestEnv(f bool) {
	testFlag = f
}

// MockConsulConfig ...
func MockConsulConfig() []byte {
	m := map[string]map[string]interface{}{
		"redis": {
			"address": "127.0.0.1", "port": 63796379, "password": "", "database": 0,
		},
		"mysql": {
			"address": "127.0.0.1", "port": 3306, "username": "root", "password": "",
			"database": "fobrain_test", "charset": "utf8mb4", "log-level": "debug", "slow-time": 15,
		},
		"elastic": {
			"address": "127.0.0.1", "port": 9200, "username": "", "password": "", "sniff": false,
		},
		"logger": {
			"level": "info", "output_console": true, "output_file": true,
			"file_name": "../logs/logs.log",
			"max_size":  64, "max_age": 30, "max_backups": 5, "local_time": true, "compress": true,
		},
	}

	jsonData, _ := json.Marshal(m)
	return jsonData
}
