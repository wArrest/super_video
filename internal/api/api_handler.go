package api

import (
	"github.com/gin-gonic/gin"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
	"github.com/wArrest/unwatermark"
)

type MEDIA_CHANNEL string

const (
	XIGUA  MEDIA_CHANNEL = "xigua"
	DOUYIN MEDIA_CHANNEL = "douyin"
)

type BaseParams struct {
	AccessPwd   string
	AccessIPMap map[string]byte
}

type Body struct {
	Pwd        string        `json:"pwd"`
	SourceText string        `json:"source_text"`
	Media      MEDIA_CHANNEL `json:"media"`
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
	var result map[string]string
	switch reqBody.Media {
	case DOUYIN:
		douyin := unwatermark.NewDouYin([]string{reqBody.SourceText})
		result = douyin.GetResults()
		if _, ok := result[reqBody.SourceText]; !ok {
			log.Warnf("无法解析的链接：%s", reqBody.SourceText)
			c.JSON(400, gin.H{
				"message": "无法解析的链接",
			})
			return
		}
	default:
		c.JSON(400, gin.H{
			"message": "不支持的链接类型",
		})
		return
	}
	rUrls:=[]string{}
	for _, realUrl := range result {
		if realUrl != "" {
			rUrls = append(rUrls, realUrl)
		}
	}
  if len(rUrls)==0 {
    c.JSON(400, gin.H{
      "message": "获取失败",
      "list":    rUrls,
    })
    return
  }
	c.JSON(200, gin.H{
		"message": "获取成功",
		"list":    rUrls,
	})
}
func (a *ApiHandler) addRecord() {

}
