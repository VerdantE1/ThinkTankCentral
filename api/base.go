package api

import (
	"ThinkTankCentral/global"
	"ThinkTankCentral/model/response"
	"github.com/gin-gonic/gin"
	"github.com/mojocn/base64Captcha"
	"go.uber.org/zap"
)

type BaseApi struct {
}

func (baseAPI *BaseApi) Init() {

}

var store = base64Captcha.DefaultMemStore

// Captcha 生成数字验证码
func (baseApi *BaseApi) Captcha(c *gin.Context) {
	// 创建验证码驱动，设置验证码格式
	driver := base64Captcha.NewDriverDigit(
		global.Config.Captcha.Height,
		global.Config.Captcha.Width,
		global.Config.Captcha.Length,
		global.Config.Captcha.MaxSkew,
		global.Config.Captcha.DotCount,
	)

	//生成验证码
	captcha := base64Captcha.NewCaptcha(driver, store) // 创建验证码对象
	id, b64s, _, err := captcha.Generate()             // 生成验证码
	if err != nil {
		global.Log.Error("Failed to generate captcha:", zap.Error(err))
		response.FailWithMessage("Failed to generate captcha", c)
		return
	}

	response.OkWithData(response.Captcha{
		CaptchaID: id,
		PicPath:   b64s,
	}, c)
}

//// SendEmailVerificationCode 发送邮箱验证码
//func (baseApi *BaseApi) SendEmailVerificationCode(c *gin.Context) {
//	var req request.SendEmailVerificationCode
//	err := c.ShouldBindJSON(&req)
//	if err != nil {
//		response.FailWithMessage(err.Error(), c)
//		return
//	}
//	if store.Verify(req.CaptchaID, req.Captcha, true) {
//		err = baseService.SendEmailVerificationCode(c, req.Email)
//		if err != nil {
//			global.Log.Error("Failed to send email:", zap.Error(err))
//			response.FailWithMessage("Failed to send email", c)
//			return
//		}
//		response.OkWithMessage("Successfully sent email", c)
//		return
//	}
//	response.FailWithMessage("Incorrect verification code", c)
//}

// QQLoginURL 返回 QQ 登录链接
func (baseApi *BaseApi) QQLoginURL(c *gin.Context) {
	url := global.Config.QQ.QQLoginURL()
	response.OkWithData(url, c)
}
