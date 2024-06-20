package cfg

import (
	"os/exec"

	"github.com/spf13/cast"
	"github.com/spf13/viper"
)

func get(key string) any {
	suffix := "_command"
	if viper.IsSet(key) {
		return viper.Get(key)

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

func GetString(key string) string {
	return cast.ToString(get(key))
}
func GetInt(key string) int {
	return cast.ToInt(get(key))
}
func GetStringSlice(key string) []string {
	return cast.ToStringSlice(get(key))
}
func GetIntDefault(key string, defaultValue int) int {
	i := GetInt(key)
	if i == 0 {
		return defaultValue
	}
	return i
}
