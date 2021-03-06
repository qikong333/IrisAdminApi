package main

import (
	"testing"

	"github.com/kataras/iris/v12"
)

// 登陆成功
//func TestUserLoginSuccess(t *testing.T) {
//
//	oj := map[string]string{
//		"username": rc.TestData.UserName,
//		"password": rc.TestData.Pwd,
//	}
//	login(t,  oj, iris.StatusOK, true, "登陆成功", nil)
//}

// 输入不存在的用户名登陆
func TestUserLoginWithErrorName(t *testing.T) {
	oj := map[string]string{
		"username": "err_user",
		"password": rc.TestData.Pwd,
	}

	login(t, oj, iris.StatusOK, false, "用户不存在")
}

// 输入错误的登陆密码
func TestUserLoginWithErrorPwd(t *testing.T) {

	oj := map[string]string{
		"username": rc.TestData.UserName,
		"password": "admin",
	}
	login(t, oj, iris.StatusOK, false, "用户名或密码错误")
}

// 输入登陆密码格式错误
func TestUserLoginWithErrorFormtPwd(t *testing.T) {
	oj := map[string]string{
		"username": rc.TestData.UserName,
		"password": "123",
	}

	login(t, oj, iris.StatusOK, false, "密码格式错误")
}

// 输入登陆密码格式错误
func TestUserLoginWithErrorFormtUserName(t *testing.T) {

	oj := map[string]string{
		"username": "df",
		"password": "123",
	}

	login(t, oj, iris.StatusOK, false, "用户名格式错误")
}
