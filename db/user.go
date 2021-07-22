package db

import (
	mydb "filestore/db/mysql"
	"fmt"
)

// create table `tbl_user`(
//     `id` int(11) NOT NULL AUTO_INCREMENT,
//     `user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
//     `user_pwd` varchar(256) NOT NULL DEFAULT '' COMMENT '用户encoded 密码',
//     `email` varchar(64) DEFAULT '' COMMENT '邮箱',
//     `phone` varchar(128) DEFAULT '' COMMENT '手机号',
//     `email_validated` tinyint(1) DEFAULT 0 COMMENT '邮箱是否已验证',
//     `phone_validated` tinyint(1) DEFAULT 0 COMMENT '手机号是否已验证',
//     `signup_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '注册日期',
//     `last_active` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后活跃时间',
//     `profile` text COMMENT '用户熟悉',
//     `status` int(11) NOT NULL DEFAULT '0' COMMENT '账户状态(启用/禁用/锁定/标记删除等)',
//     PRIMARY KEY(`id`),
//     UNIQUE KEY `idx_phone`(`phone`),
//     KEY `idx_status`(`status`)
// )ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8mb4;

// 注册操作
func UserSignup(username string, passwd string) bool {
	stmt, err := mydb.DBConn().Prepare("insert ignore into tbl_user(`user_name`, `user_pwd`) values (?, ?)")
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}
	defer stmt.Close()

	ret, err := stmt.Exec(username, passwd)
	if err != nil {
		fmt.Println("Failed to insert, err:" + err.Error())
		return false
	}

	if rowsAffected, err := ret.RowsAffected(); err == nil && rowsAffected > 0 {
		return true
	}

	return false
}

// 判断登录
func UserSignin(username string, encpwd string) bool {
	stmt, err := mydb.DBConn().Prepare("select * from tbl_user where user_name=? limit 1")

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	rows, err := stmt.Query(username)
	if err != nil {
		fmt.Println(err.Error())
		return false
	} else if rows == nil {
		fmt.Println("username not found:" + username)
		return false
	}

	pRows := mydb.ParseRows(rows)

	if len(pRows) > 0 && string(pRows[0]["user_pwd"].([]byte)) == encpwd {
		return true
	}

	return false
}

// create table tbl_user_token
// (
// 	`id` int(11) NOT NULL AUTO_INCREMENT,
// 	`user_name` varchar(64) NOT NULL DEFAULT '' COMMENT '用户名',
// 	`user_token` char(40) NOT NULL DEFAULT '' COMMENT '用户登录token',
// 	PRIMARY KEY (`id`),
// 	UNIQUE KEY `idx_username` (`user_name`)
// )ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

// 刷新token
func UpdateToken(username string, token string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"replace into tbl_user_token (`user_name`, `user_token`) values(?, ?)",
	)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	defer stmt.Close()

	_, err = stmt.Exec(username, token)

	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	return true
}
