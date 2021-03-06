package orm

import (
	mydb "filestore-hsz/service/dbproxy/conn"
	"time"
	"log"
)

// 更新用户文件表
func OnUserFileUploadFinished(username, filehash, filename string, filesize int64) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"insert ignore into tbl_user_file (`user_name`, `file_sha1`, `file_name`, " +
			"`file_size`, `upload_at`, `status`) values (?,?,?,?,?,0)")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(username, filehash, filename, filesize, time.Now())
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

// 批量获取用户文件信息
func QueryUserFileMetas(username string, status, limit int64) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1, file_name, file_size, upload_at," +
			"last_update from tbl_user_file where user_name=? and status=? limit ?")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, status, limit)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	var userFiles []TableUserFile
	for rows.Next() {
		ufile := TableUserFile{}
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize,
			&ufile.UploadAt, &ufile.LastUpdated)
		if err != nil {
			log.Println(err.Error())
			break
		}
		userFiles = append(userFiles, ufile)
	}
	res.Suc = true
	res.Data = userFiles
	return
}

// 删除文件（标记删除）
func DeleteUserFile(username, filehash, filename string, status int64) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"update tbl_user_file set status=? where user_name=? and file_sha1=? and file_name=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(status, username, filehash, filename)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

// 文件重命名
func RenameFileName(username, filehash, filename, filenameOld string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"update tbl_user_file set file_name=? where user_name=? and file_sha1=? and status=0 and file_name=? limit 1")
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(filename, username, filehash, filenameOld)
	if err != nil {
		log.Println(err.Error())
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	res.Suc = true
	return
}

// 获取用户单个文件信息
func QueryUserFileMeta(username string, filehash string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"select file_sha1, file_name, file_size, upload_at," +
			"last_update from tbl_user_file where user_name=? and file_sha1=? and status=0 limit 1")
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	}
	defer stmt.Close()

	rows, err := stmt.Query(username, filehash)
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	}

	ufile := TableUserFile{}
	if rows.Next() {
		err = rows.Scan(&ufile.FileHash, &ufile.FileName, &ufile.FileSize,
			&ufile.UploadAt, &ufile.LastUpdated)

		if err != nil {
			log.Println(err.Error())
			res.Suc = false
			res.Msg = err.Error()
			return
		}
	}
	res.Suc = true
	res.Data = ufile
	return
}

// 文件是否已经上传
func IsUserFileUploaded(username string, filehash string) (res ExecResult) {
	stmt, err := mydb.DBConn().Prepare(
		"select 1 from tbl_user_file where user_name=? and file_sha1=? and status=0 limit 1")
	rows, err := stmt.Query(username, filehash)
	if err != nil {
		res.Suc = false
		res.Msg = err.Error()
		return
	} else if rows == nil || !rows.Next() {
		res.Suc = true
		res.Data = map[string]bool{
			"exists": false,
		}
		return
	}
	res.Suc = true
	res.Data = map[string]bool{
		"exists": true,
	}
	return
}













