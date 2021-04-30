package controllers

import (
	"GoBlog/app/models/user"
	"GoBlog/app/requests"
	"GoBlog/pkg/auth"
	"GoBlog/pkg/email"
	"GoBlog/pkg/view"
	"fmt"
	"net/http"
)

// AuthController 处理静态页面
type AuthController struct{}

// Register 注册页面
func (*AuthController) Register(w http.ResponseWriter, _ *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.register")
}

// DoRegister 处理注册逻辑
func (*AuthController) DoRegister(w http.ResponseWriter, r *http.Request) {
	// 1. 初始化数据
	_user := user.User{
		Name:            r.PostFormValue("name"),
		Email:           r.PostFormValue("email"),
		Password:        r.PostFormValue("password"),
		PasswordConfirm: r.PostFormValue("password_confirm"),
		VerifyCode:      r.PostFormValue("verify_code"),
	}
	// 2. 表单规则
	errs := requests.ValidateRegistrationForm(_user)

	if len(errs) > 0 {
		// 3. 表单验证不通过，重新显示表单
		view.RenderSimple(w, view.D{
			"Errors": errs,
			"User":   _user,
		}, "auth.register")
	} else {
		// 4. 验证成功，创建数据
		err := _user.Create()
		if _user.ID > 0 && err == nil {
			// 登录用户并跳转到首页
			auth.Login(_user)
			http.Redirect(w, r, "/", http.StatusFound)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, "创建用户失败，请联系管理员")
		}
	}
}

// SendVerifyCode 处理发送验证码逻辑
func (*AuthController) SendVerifyCode(w http.ResponseWriter, r *http.Request) {
	emailAddress := r.PostFormValue("email")
	var server email.SMTPServer
	if err := server.SendEmail("verifyEmail", emailAddress, auth.CreateVerifyCode()); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "发送邮件错误，请稍后再试")
		return
	}
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

// Login 显示登录表单
func (*AuthController) Login(w http.ResponseWriter, _ *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.login")
}

// DoLogin 处理登录表单逻辑
func (*AuthController) DoLogin(w http.ResponseWriter, r *http.Request) {
	// 1. 初始化表单数据
	emailAddress := r.PostFormValue("email")
	password := r.PostFormValue("password")

	// 2. 尝试登录
	if err := auth.Attempt(emailAddress, password); err == nil {
		http.Redirect(w, r, "/", http.StatusFound)
	} else {
		view.RenderSimple(w, view.D{
			"Error":    err.Error(),
			"Email":    emailAddress,
			"Password": password,
		}, "auth.login")
	}
}

// Logout 退出登录
func (*AuthController) Logout(w http.ResponseWriter, r *http.Request) {
	auth.Logout()
	http.Redirect(w, r, "/", http.StatusFound)
}

// LostPassword 找回密码页面
func (*AuthController) LostPassword(w http.ResponseWriter, r *http.Request) {
	view.RenderSimple(w, view.D{}, "auth.lostPassword")
}

// DoLostPassword 处理找回密码逻辑
func (*AuthController) DoLostPassword(w http.ResponseWriter, r *http.Request) {
	// 初始化数据
	emailAddress := r.PostFormValue("email")
	_user := user.User{
		Email: emailAddress,
	}
	// 表单验证
	errs := requests.ValidateLostPasswordForm(_user)
	if len(errs) > 0 {
		view.RenderSimple(w, view.D{
			"Errors": errs,
			"User":   _user,
		}, "auth.lostPassword")
		return
	}
	url := auth.GetLostPasswordURL(r, emailAddress)

	var server email.SMTPServer
	if err := server.SendEmail("lostPassword", emailAddress, url); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(w, "发送邮件错误，请稍后再试")
		return
	}
}
