package requests

import (
	"GoBlog/pkg/model"
	"errors"
	"fmt"
	"github.com/thedevsaddam/govalidator"
	"strings"
)

// init 方法会在初始化时执行
func init() {
	govalidator.AddCustomRule("not_exists", func(field string, rule string, message string, value interface{}) error {
		rng := strings.Split(strings.TrimPrefix(rule, "not_exists:"), ",")

		tableName := rng[0]
		dbFiled := rng[1]
		val := value.(string)

		var count int64
		model.DB.Table(tableName).Where(dbFiled+"= ?", val).Count(&count)

		if count != 0 {
			if message != "" {
				return errors.New(message)
			}
			return fmt.Errorf("%v 已被占用", val)
		}
		return nil
	})
}
