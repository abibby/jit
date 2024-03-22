package cfg

import (
	"log"
	"os/exec"
	"strconv"

	"github.com/spf13/viper"
)

func GetString(key string) string {
	suffix := "_command"
	if viper.IsSet(key) {
		return viper.GetString(key)

	} else if viper.IsSet(key + suffix) {
		newKey := key[:len(key)-len(suffix)]
		b, err := exec.Command("sh", "-c", viper.GetString(key+suffix)).Output()
		value := string(b[:len(b)-1])
		if err == nil {
			viper.Set(newKey, value)
		}

		return value
	}
	return ""
}

func GetInt(key string) int {
	s := GetString(key)
	i, err := strconv.Atoi(s)
	if err != nil {
		log.Print(err)
	}
	return i
}
func GetIntDefault(key string, defaultValue int) int {
	i := GetInt(key)
	if i == 0 {
		return defaultValue
	}
	return i
}
