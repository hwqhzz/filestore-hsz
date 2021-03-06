package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"context"
	userProto "filestore-hsz/service/account/proto"
	"log"
	cfg "filestore-hsz/config"
)

// 查询批量的文件元信息
func FileQueryHandler(c *gin.Context) {
	limitCnt, _ := strconv.Atoi(c.Request.FormValue("limit"))
	username := c.Request.FormValue("username")

	rpcResp, err := userCli.UserFiles(context.TODO(), &userProto.ReqUserFile{
		Username: username,
		Limit:    int32(limitCnt),
		Status:   int32(0),
	}, cfg.RpcOpts)

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(rpcResp.FileData) <= 0 {
		rpcResp.FileData = []byte("[]")
	}

	c.Data(http.StatusOK, "application/json", rpcResp.FileData)
}

// 批量查询已删除文件
func FileQueryDeletedHandler(c *gin.Context)  {

	limitCnt, _ := strconv.Atoi(c.Request.FormValue("limit"))
	username := c.Request.FormValue("username")
	rpcResp, err := userCli.UserFiles(context.TODO(), &userProto.ReqUserFile{
		Username: username,
		Limit:    int32(limitCnt),
		Status:   int32(1),
	}, cfg.RpcOpts)
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}

	if len(rpcResp.FileData) <= 0 {
		rpcResp.FileData = []byte("[]")
	}

	c.Data(http.StatusOK, "application/json", rpcResp.FileData)
}

// 更新元信息接口(重命名)
func FileMetaUpdateHandler(c *gin.Context)  {

	opType := c.Request.FormValue("op")
	fileSha1 := c.Request.FormValue("filehash")
	newFileName := c.Request.FormValue("newfilename")
	username := c.Request.FormValue("username")
	oldFilename := c.Request.FormValue("oldfilename")

	if opType != "0" || len(newFileName) < 1 {
		c.Status(http.StatusForbidden)
		return
	}

	rpcResp, err := userCli.UserFileRename(context.TODO(), &userProto.ReqUserFileRename{
		Username: username,
		Filehash: fileSha1,
		NewFileName: newFileName,
		OldFilename: oldFilename,
	})
	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	if len(rpcResp.FileData) <= 0 {
		rpcResp.FileData = []byte("[修改文件名成功]")
	}
	c.Data(http.StatusOK, "application/json", rpcResp.FileData)
}

// 删除用户文件信息接口
func FileDeleteHandler(c *gin.Context) {
	fileSha1 := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")
	filename := c.Request.FormValue("filename")
	// context.TODO返回一个非nil的空上下文。代码应该使用上下文。当不清楚要使用哪个上下文或者上下文还不可用(因为周围的函数还没有扩展到接受上下文参数)时，可以使用TODO。
	rpcResp, err := userCli.UserFileDelete(context.TODO(), &userProto.ReqUserFileDelete{
		Username: username,
		Filehash: fileSha1,
		Filename: filename,
		Status:   int32(1),
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	if len(rpcResp.FileData) <= 0 {
		rpcResp.FileData = []byte("[]")
	}
	c.Data(http.StatusOK, "application/json", rpcResp.FileData)
}

// 恢复用户文件接口
func FileRecoverHandler(c *gin.Context) {
	fileSha1 := c.Request.FormValue("filehash")
	username := c.Request.FormValue("username")
	filename := c.Request.FormValue("filename")
	// context.TODO返回一个非nil的空上下文。代码应该使用上下文。当不清楚要使用哪个上下文或者上下文还不可用(因为周围的函数还没有扩展到接受上下文参数)时，可以使用TODO。
	rpcResp, err := userCli.UserFileDelete(context.TODO(), &userProto.ReqUserFileDelete{
		Username: username,
		Filehash: fileSha1,
		Filename: filename,
		Status:   int32(0),
	})

	if err != nil {
		log.Println(err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	if len(rpcResp.FileData) <= 0 {
		rpcResp.FileData = []byte("[]")
	}
	c.Data(http.StatusOK, "application/json", rpcResp.FileData)
}
