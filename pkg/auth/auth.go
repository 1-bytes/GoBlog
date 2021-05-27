package auth

import (
	"GoBlog/app/models/user"
	"GoBlog/pkg/encrypts"
	"GoBlog/pkg/password"
	"GoBlog/pkg/route"
	"GoBlog/pkg/session"
	"encoding/hex"
	"errors"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

// _getUID 获取当前用户的UID 如果不存在则返回空
func _getUID() string {
	_uid := session.Get("uid")
	uid, ok := _uid.(string)
	if ok && len(uid) > 0 {
		return uid
	}
	return ""
}

// User 获取登录用户信息
func User() user.User {
	uid := _getUID()
	if len(uid) > 0 {
		_user, err := user.Get(uid)
		if err != nil {
			return _user
		}
	}
	return user.User{}
}

// Attempt 尝试登录
func Attempt(email string, password string) error {
	// 1. 根据 Email 获取用户
	_user, err := user.GetByEmail(email)

	// 2. 如果出现错误
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("账号不存在或密码错误")
		}
		return errors.New("内部错误，请稍后再试")
	}

	// 3. 匹配密码
	if !_user.ComparePassword(password) {
		return errors.New("账号不存在或密码错误")
	}

	// 4. 登录用户，保存会话
	session.Put("uid", _user.GetStringID())
	return nil
}

// Login 登录指定用户
func Login(_user user.User) {
	session.Put("uid", _user.GetStringID())
}

// Logout 退出用户
func Logout() {
	session.Forget("uid")
}

// Check 检测是否登录
func Check() bool {
	return len(_getUID()) > 0
}

// CreateVerifyCode 生成一个长度为 6 的验证码
func CreateVerifyCode() string {
	expires := strconv.FormatInt(time.Now().Unix()+1800, 10)
	verifyCode := password.GetRandomPassword(6, []rune("0123456789"))
	session.Put("VerifyCode", expires+verifyCode)
	return verifyCode
}

// CheckVerifyCode 判断验证码是否正确
func CheckVerifyCode(vcode string) bool {
	// 从 session 当中获取邮箱验证码
	_verifyCode := session.Get("VerifyCode")
	verifyCode, ok := _verifyCode.(string)
	if !ok && len(verifyCode) < 11 {
		return false
	}
	// 前十位是秒级时间戳，剩下的是验证码
	expires := verifyCode[:10]
	verifyCode = verifyCode[10:]
	timeStr := strconv.FormatInt(time.Now().Unix(), 10)
	if timeStr > expires {
		return false
	}
	if verifyCode == "" || vcode != verifyCode {
		return false
	}
	session.Forget("VerifyCode")
	return true
}

// GetLostPasswordURL 获取重置密码链接
func GetLostPasswordURL(r *http.Request, email string) (url string, err error) {
	expires := strconv.FormatInt(time.Now().Unix()+1800, 10)
	// 数据加密
	aes := encrypts.Aes{
		Key: []byte("PNKZIJWMFCHQDTYLRSEOUGVBXARGHBWH"),
	}
	if !user.HasUserByEmail(email) {
		err := errors.New("邮箱错误或账号不存在")
		return "", err
	}
	origData := []byte(expires + email)
	token := hex.EncodeToString(aes.AesEncryptCFB(origData))
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	url = scheme + "://" + r.Host + route.Name2URL("auth.repassword", "token", token)
	return url, nil
}

// GetLostPasswordEmailByToken 根据 TOKEN 获取
func GetLostPasswordEmailByToken(token string) (email string, err error) {
	var encrypted []byte
	if encrypted, err = hex.DecodeString(token); err != nil {
		err := errors.New("令牌错误或已过期")
		return "", err
	}
	aes := encrypts.Aes{
		Key: []byte("PNKZIJWMFCHQDTYLRSEOUGVBXARGHBWH"),
	}

	decrypted := string(aes.AesDecryptCFB(encrypted))
	timeStr := strconv.FormatInt(time.Now().Unix(), 10)
	expires := decrypted[:10]
	if timeStr > expires {
		err := errors.New("令牌错误或已过期")
		return "", err
	}
	email = decrypted[10:]
	return email, nil
}
