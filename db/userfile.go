package db

import (
	mydb "filestore/db/mysql"
	"fmt"
	"time"
)

// create table tbl_user_file
// (
// 	`id` int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
// 	`user_name` varchar(64) NOT NULL,
// 	`file_sha1` varchar(64) NOT NULL DEFAULT '' COMMENT '文件hash',
// 	`file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
// 	`file_name` varchar(256) NOT NULL DEFAULT '' COMMENT '文件名',
// 	`upload_at` datetime DEFAULT CURRENT_TIMESTAMP COMMENT '上传时间',
// 	`last_update` datetime DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '最后修改时间',
// 	`status` int(11) NOT NULL DEFAULT '0' COMMENT '文件状态(0正常1删除2禁用)',
// 	UNIQUE KEY `idx_user_file` (`user_name`, `file_sha1`),
// 	KEY `idx_status` (`status`),
// 	KEY `idx_user_id` (`user_name`)
// )ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

type UserFile struct {
	UserName    string
	FileHash    string
	FileName    string
	FileSize    int64
	UploadAt    string
	LastUpdated string
}

func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user_file (`user_name`, `file_sha1`, `file_name`, `file_size`, `upload_at`) values (?,?,?,?,?)",
	)

	if err != nil {
		return false
	}

	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())

	if err != nil {
		return false
	}

	return true
}

func QueryUserFileMetas(username string, limit int) ([]UserFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1, file_name, file_size, upload_at, last_update from tbl_user_file where user_name=? limit=?",
	)

	if err != nil {
		return nil, err
	}

	defer stmt.Close()

	rows, err := stmt.Query(username, limit)

	if err != nil {
		return nil, err
	}

	var userFiles []UserFile
	for rows.Next() {
		ufile := UserFile{}
		err := rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize, &ufile.UploadAt, &ufile.LastUpdated)

		if err != nil {
			fmt.Println(err.Error())
			break
		}

		userFiles = append(userFiles, ufile)
	}

	return userFiles, nil
}
