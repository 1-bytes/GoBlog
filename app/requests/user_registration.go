package requests

import (
	"GoBlog/app/models/user"
	"github.com/thedevsaddam/govalidator"
	"net/url"
)

// ValidateRegistrationForm 验证表单，返回 errs 长度等于零即通过
func ValidateRegistrationForm(data user.User) url.Values {
	// 1. 定制认证规则
	rules := govalidator.MapData{
		"name":             []string{"required", "alpha_num", "between:3,20", "not_exists:users,name"},
		"email":            []string{"required", "min:4", "max:30", "email", "not_exists:users,email"},
		"password":         []string{"required", "min:6"},
		"password_confirm": []string{"required"},
	}

	// 2. 定制错误信息
	message := govalidator.MapData{
		"name": []string{
			"required:用户名为必填项",
			"alpha_num:用户名格式错误，只允许数字和英文",
			"between:用户名长度需在 3~20 之间",
		},
		"email": []string{
			"required:邮箱为必填项",
			"min:邮箱长度需大于 4",
			"max:邮箱长度需小于 30",
			"email:邮箱格式不正确，请提供有效的邮箱地址",
		},
		"password": []string{
			"required:密码为必填项",
			"min:密码长度需大于6",
		},
		"password_confirm": []string{
			"required:确认密码框为必填项",
		},
	}

	// 3. 配初始化
	opts := govalidator.Options{
		Data:          &data,
		Rules:         rules,
		TagIdentifier: "valid",
		Messages:      message,
	}

	// 4. 开始认证
	errs := govalidator.New(opts).ValidateStruct()
	if data.Password != data.PasswordConfirm {
		errs["password_confirm"] = append(errs["password_confirm"], "两次输入的密码不匹配")
	}
	return errs
}
