package handler

import (
	dbplayer "filestore/db"
	"filestore/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	pwd_salt = "*#890"
)

// 处理用户注册请求
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(data)
		return
	}

	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("Invalid parameter"))
		return
	}

	enc_passwd := util.Sha1([]byte(passwd + pwd_salt))
	suc := dbplayer.UserSignup(username, enc_passwd)
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}

}

// 登录接口
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	username := r.Form.Get("username")
	password := r.Form.Get("password")

	encPasswd := util.Sha1([]byte(password + pwd_salt))
	pwdChecked := dbplayer.UserSignin(username, encPasswd)

	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}

	token := GenToken(username)

	upRes := dbplayer.UpdateToken(username, token)

	if !upRes {
		w.Write([]byte("FAILED"))
		return
	}

	w.Write([]byte("http://" + r.Host + "/static/view/home.html"))
}

func GenToken(username string) string {
	// 40位 md5(username+timestamp+token_salt)+timestamp[:8]
	ts := fmt.Sprintf("%x", time.Now().Unix())

	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))

	return tokenPrefix + ts[:8]
}
