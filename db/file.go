package db

import (
	"database/sql"
	mydb "filestore/db/mysql"
	"fmt"
)

// create table tbl_file
// (
// 	`id` int(11) NOT NULL AUTO_INCREMENT,
// 	`file_sha1` char(40) NOT NULL DEFAULT '' COMMENT '文件名',
// 	`file_name` varchar(256) NOT NULL DEFAULT '0' COMMENT '文件大小',
// 	`file_size` bigint(20) DEFAULT '0' COMMENT '文件大小',
// 	`flie_addr` varchar(1024) NOT NULL DEFAULT '' COMMENT '文件存储位置',
// 	`create_at` datetime DEFAULT NOW() COMMENT '创建日期',
// 	`update_at` datetime DEFAULT NOW() on update current_timestamp() COMMENT '更新日期',
// 	`status` int(11) NOT NULL DEFAULT '0' COMMENT '状态(可用/禁用/已删除等状态)',
// 	`ext1` int(11) DEFAULT '0' COMMENT '备用字段1',
// 	`ext2` text COMMENT '备用字段2',
// 	PRIMARY KEY (`id`),
// 	UNIQUE KEY `idx_file_hash` (`file_sha1`),
// 	KEY `idx_status` (`status`)
// ) ENGINE=InnoDB DEFAULT CHARSET=utf8;

func OnfileUploadFinished(filehash string, filename string, filesize int64, fileaddr string) bool {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_file(`file_sha1`, `file_name`, `file_size`, `file_addr`, `status`) values (?,?,?,?,1)",
	)
	if err != nil {
		fmt.Println("Failed to prepare statement, err:" + err.Error())
	}
	defer stmt.Close()

	ret, err := stmt.Exec(filehash, filename, filesize, fileaddr)
	if err != nil {
		fmt.Println(err.Error())
		return false
	}

	if rf, err := ret.RowsAffected(); nil == err {
		if rf <= 0 {
			fmt.Printf("File with hash:%s has been uploaded before", filehash)
		}
		return true
	}
	return false
}

type TableFile struct {
	FileHash string
	FileName sql.NullString
	FileSize sql.NullInt64
	FileAddr sql.NullString
}

func GetFileMeta(filehash string) (*TableFile, error) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1, file_addr, file_name, file_size from tbl_file where file_sha1=? and status=1 limit 1",
	)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	defer stmt.Close()

	tfile := TableFile{}

	err = stmt.QueryRow(filehash).Scan(&tfile.FileHash, &tfile.FileAddr, &tfile.FileName, &tfile.FileSize)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return &tfile, nil
}
