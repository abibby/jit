package cfg

import (
	"os/exec"

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
