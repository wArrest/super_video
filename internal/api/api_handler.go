package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	"github.com/wArrest/unwatermark"
)

type BaseParams struct {
	AccessPwd   string
	AccessIPMap map[string]byte
}

type Body struct {
	Pwd        string `json:"pwd"`
	SourceText string `json:"source_text"`
}
type ApiHandler struct {
	dbPool *sqlx.DB
	BaseParams
}

func NewApiHandler(dbPool *sqlx.DB) *ApiHandler {
	return &ApiHandler{dbPool: dbPool}
}

func (a *ApiHandler) Transform(c *gin.Context) {
	var reqBody Body
	err := c.BindJSON(&reqBody)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	//不是白名单ip内且授权码错误的人没有权限访问
	if _, ok := a.AccessIPMap[c.ClientIP()]; !ok && reqBody.Pwd != a.AccessPwd {
		c.JSON(400, gin.H{
			"message": "没有访问的权限！",
		})
		return
	}
	media := unwatermark.GetMedia(reqBody.SourceText)
	if media == nil {
		c.JSON(400, gin.H{
			"message": "暂不支持的媒体！",
		})
		return
	}
	rUrl, err := media.GetRealLink(reqBody.SourceText)
	if err != nil {
		c.JSON(400, gin.H{
			"message": err.Error(),
		})
		return
	}
	c.JSON(200, gin.H{
		"message": "获取成功",
		"list":    []string{},
		"rUrl":    rUrl,
	})
}
