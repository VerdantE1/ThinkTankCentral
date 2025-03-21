package api

import (
	"ThinkTankCentral/global"
	"ThinkTankCentral/model/database"
	"ThinkTankCentral/model/request"
	"ThinkTankCentral/model/response"
	"ThinkTankCentral/utils"
	"errors"
	"time"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"go.uber.org/zap"
)

type UserApi struct {
}

// Register 用户注册
// @Summary 用户注册接口
// @Description 通过邮箱验证码进行用户注册
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param body body request.Register true "注册请求参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /user/register [post]
// 实现逻辑：
// 1. 校验验证码、邮箱一致性
// 2. 验证邮箱验证码有效性
// 3. 创建用户记录
// 4. 生成JWT令牌
// Register 注册
func (userApi *UserApi) Register(c *gin.Context) {
	var req request.Register
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	session := sessions.Default(c)
	// 两次邮箱一致性判断
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("This email doesn't match the email to be verified", c)
		return
	}

	// 获取会话中存储的邮箱验证码
	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("Invalid verification code", c)
		return
	}

	// 判断邮箱验证码是否过期
	savedTime := session.Get("expire_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("The verification code has expired, please resend it", c)
		return
	}

	u := database.User{Username: req.Username, Password: req.Password, Email: req.Email}

	user, err := userService.Register(u)
	if err != nil {
		global.Log.Error("Failed to register user:", zap.Error(err))
		response.FailWithMessage("Failed to register user", c)
		return
	}

	// 注册成功后，生成 token 并返回
	userApi.TokenNext(c, user)
}

// Login 登录接口，根据不同的登录方式调用不同的登录方法
// Login 用户登录接口
// @Summary 用户登录入口
// @Description 支持邮箱/QQ两种登录方式
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param flag query string true "登录方式(email/qq)"
// @Success 200 {object} response.Response{data=response.Login}
// @Failure 400 {object} response.Response
// @Router /user/login [post]
// 实现逻辑：
// 1. 根据flag参数路由到对应登录方式
// 2. 邮箱登录需校验图形验证码
// 3. QQ登录需校验授权码
func (userApi *UserApi) Login(c *gin.Context) {
	switch c.Query("flag") {
	case "email":
		userApi.EmailLogin(c)
	case "qq":
		userApi.QQLogin(c)
	default:
		userApi.EmailLogin(c)
	}
}

// EmailLogin 邮箱登录
func (userApi *UserApi) EmailLogin(c *gin.Context) {
	var req request.Login
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	// 校验验证码
	if store.Verify(req.CaptchaID, req.Captcha, true) {
		u := database.User{Email: req.Email, Password: req.Password}
		user, err := userService.EmailLogin(u)
		if err != nil {
			global.Log.Error("Failed to login:", zap.Error(err))
			response.FailWithMessage("Failed to login", c)
			return
		}

		// 登录成功后生成 token
		userApi.TokenNext(c, user)
		return
	}
	response.FailWithMessage("Incorrect verification code", c)
}

// QQLogin QQ登录
func (userApi *UserApi) QQLogin(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		response.FailWithMessage("Code is required", c)
		return
	}

	// 获取访问令牌
	accessTokenResponse, err := qqService.GetAccessTokenByCode(code)
	if err != nil || accessTokenResponse.Openid == "" {
		global.Log.Error("Invalid code", zap.Error(err))
		response.FailWithMessage("Invalid code", c)
		return
	}

	// 根据访问令牌进行QQ登录
	user, err := userService.QQLogin(accessTokenResponse)
	if err != nil {
		global.Log.Error("Failed to login:", zap.Error(err))
		response.FailWithMessage("Failed to login", c)
		return
	}

	// 登录成功后生成 token
	userApi.TokenNext(c, user)
}

// Token 认证机制:
// 通过 JWT 实现用户认证和授权。
// 提供访问令牌和刷新令牌，确保安全和灵活性。
// 多点登录拦截:
// 防止同一账户在多个地点同时登录，提高系统安全性。
// 通过 Redis 管理用户的登录状态，实现高效的登录会话管理。
// TokenNext 令牌生成核心逻辑
// @Description 生成访问令牌和刷新令牌，处理多点登录拦截
// 实现逻辑：
//  1. 生成JWT访问令牌(有效期2小时)
//  2. 生成刷新令牌(有效期7天)
//  3. 多点登录检测处理：
//     a) 未开启多点登录时直接返回令牌
//     b) 已存在旧令牌时加入黑名单
//     c) 更新Redis中的最新令牌
//
// 错误处理：
// - 用户冻结状态返回403
// - JWT生成失败返回500
func (userApi *UserApi) TokenNext(c *gin.Context, user database.User) {
	// 检查用户是否被冻结
	if user.Freeze {
		response.FailWithMessage("The user is frozen, contact the administrator", c)
		return
	}

	baseClaims := request.BaseClaims{
		UserID: user.ID,
		UUID:   user.UUID,
		RoleID: user.RoleID,
	}

	j := utils.NewJWT()

	// 创建访问令牌
	accessClaims := j.CreateAccessClaims(baseClaims)
	accessToken, err := j.CreateAccessToken(accessClaims)
	if err != nil {
		global.Log.Error("Failed to get accessToken:", zap.Error(err))
		response.FailWithMessage("Failed to get accessToken", c)
		return
	}

	// 创建刷新令牌
	refreshClaims := j.CreateRefreshClaims(baseClaims)
	refreshToken, err := j.CreateRefreshToken(refreshClaims)
	if err != nil {
		global.Log.Error("Failed to get refreshToken:", zap.Error(err))
		response.FailWithMessage("Failed to get refreshToken", c)
		return
	}

	// 是否开启了多地点登录拦截
	if !global.Config.System.UseMultipoint {
		// 设置刷新令牌并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
		return
	}

	// 检查 Redis 中是否已存在该用户的 JWT
	if jwtStr, err := jwtService.GetRedisJWT(user.UUID); errors.Is(err, redis.Nil) {
		// 不存在就设置新的
		if err := jwtService.SetRedisJWT(refreshToken, user.UUID); err != nil {
			global.Log.Error("Failed to set login status:", zap.Error(err))
			response.FailWithMessage("Failed to set login status", c)
			return
		}

		// 设置刷新令牌并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
	} else if err != nil {
		// 出现错误处理
		global.Log.Error("Failed to set login status:", zap.Error(err))
		response.FailWithMessage("Failed to set login status", c)
	} else {
		// Redis 中已存在该用户的 JWT，将旧的 JWT 加入黑名单，并设置新的 token
		var blacklist database.JwtBlacklist
		blacklist.Jwt = jwtStr
		if err := jwtService.JoinInBlacklist(blacklist); err != nil {
			global.Log.Error("Failed to invalidate jwt:", zap.Error(err))
			response.FailWithMessage("Failed to invalidate jwt", c)
			return
		}

		// 设置新的 JWT 到 Redis
		if err := jwtService.SetRedisJWT(refreshToken, user.UUID); err != nil {
			global.Log.Error("Failed to set login status:", zap.Error(err))
			response.FailWithMessage("Failed to set login status", c)
			return
		}

		// 设置刷新令牌并返回
		utils.SetRefreshToken(c, refreshToken, int(refreshClaims.ExpiresAt.Unix()-time.Now().Unix()))
		c.Set("user_id", user.ID)
		response.OkWithDetailed(response.Login{
			User:                 user,
			AccessToken:          accessToken,
			AccessTokenExpiresAt: accessClaims.ExpiresAt.Unix() * 1000,
		}, "Successful login", c)
	}
}

// ForgotPassword 找回密码
// ForgotPassword 密码找回
// @Summary 通过邮箱验证码重置密码
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param body body request.ForgotPassword true "密码找回参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /user/forgot-password [post]
// 实现逻辑：
// 1. 校验邮箱验证码有效性
// 2. 更新用户密码为加密后的新密码
// 安全机制：
// - 会话中存储的验证码5分钟有效
// - 密码使用bcrypt加密存储
func (userApi *UserApi) ForgotPassword(c *gin.Context) {
	var req request.ForgotPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	session := sessions.Default(c)

	// 两次邮箱一致性判断
	savedEmail := session.Get("email")
	if savedEmail == nil || savedEmail.(string) != req.Email {
		response.FailWithMessage("This email doesn't match the email to be verified", c)
		return
	}

	// 获取会话中存储的邮箱验证码
	savedCode := session.Get("verification_code")
	if savedCode == nil || savedCode.(string) != req.VerificationCode {
		response.FailWithMessage("Invalid verification code", c)
		return
	}

	// 判断邮箱验证码是否过期
	savedTime := session.Get("expire_time")
	if savedTime.(int64) < time.Now().Unix() {
		response.FailWithMessage("The verification code has expired, please resend it", c)
		return
	}

	err = userService.ForgotPassword(req)
	if err != nil {
		global.Log.Error("Failed to retrieve the password:", zap.Error(err))
		response.FailWithMessage("Failed to retrieve the password", c)
		return
	}
	response.OkWithMessage("Successfully retrieved", c)
}

// UserCard 获取用户卡片信息
// UserCard 获取用户名片信息
// @Summary 获取用户公开信息
// @Tags 用户模块
// @Accept json
// @Produce json
// @Param uuid query string true "用户UUID"
// @Success 200 {object} response.Response{data=response.UserCard}
// @Router /user/card [get]
// 返回字段：
// - UUID 用户唯一标识
// - Username 用户名
// - Avatar 头像URL
// - Address 所在地
// - Signature 个性签名
func (userApi *UserApi) UserCard(c *gin.Context) {
	var req request.UserCard
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	userCard, err := userService.UserCard(req)
	if err != nil {
		global.Log.Error("Failed to get card:", zap.Error(err))
		response.FailWithMessage("Failed to get card", c)
		return
	}
	response.OkWithData(userCard, c)
}

// Logout 登出
// Logout 用户登出
// @Summary 清除登录状态
// @Tags 用户模块
// @Produce json
// @Success 200 {object} response.Response
// @Router /user/logout [post]
// 实现逻辑：
// 1. 清除Refresh Token Cookie
// 2. 将JWT加入黑名单
// 3. 删除Redis中的令牌记录
func (userApi *UserApi) Logout(c *gin.Context) {
	userService.Logout(c)
	response.OkWithMessage("Successful logout", c)
}

// UserResetPassword 修改密码
// UserResetPassword 用户修改密码
// @Summary 登录状态下修改密码
// @Tags 用户模块
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body request.UserResetPassword true "密码修改参数"
// @Success 200 {object} response.Response
// @Failure 400 {object} response.Response
// @Router /user/reset-password [post]
// 安全机制：
// - 需校验原始密码
// - 修改成功后强制登出
func (userApi *UserApi) UserResetPassword(c *gin.Context) {
	var req request.UserResetPassword
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	err = userService.UserResetPassword(req)
	if err != nil {
		global.Log.Error("Failed to modify:", zap.Error(err))
		response.FailWithMessage("Failed to modify, orginal password does not match the current account", c)
		return
	}
	response.OkWithMessage("Successfully changed password, please log in again", c)
	userService.Logout(c)
}

// UserInfo 获取个人信息
// UserInfo 获取用户详细信息
// @Summary 获取当前登录用户完整信息
// @Tags 用户模块
// @Security ApiKeyAuth
// @Produce json
// @Success 200 {object} response.Response{data=database.User}
// @Router /user/info [get]
// 返回字段：
// - Email 注册邮箱
// - RoleID 角色权限
// - Freeze 冻结状态
// - CreateAt 注册时间
func (userApi *UserApi) UserInfo(c *gin.Context) {
	userID := utils.GetUserID(c)
	user, err := userService.UserInfo(userID)
	if err != nil {
		global.Log.Error("Failed to get user information:", zap.Error(err))
		response.FailWithMessage("Failed to get user information", c)
		return
	}
	response.OkWithData(user, c)
}

// UserChangeInfo 修改个人信息
// UserChangeInfo 修改个人信息
// @Summary 更新用户基本信息
// @Tags 用户模块
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body request.UserChangeInfo true "用户信息参数"
// @Success 200 {object} response.Response
// @Router /user/change-info [post]
// 可修改字段：
// - Avatar 头像URL
// - Address 所在地
// - Signature 个性签名
func (userApi *UserApi) UserChangeInfo(c *gin.Context) {
	var req request.UserChangeInfo
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	req.UserID = utils.GetUserID(c)
	err = userService.UserChangeInfo(req)
	if err != nil {
		global.Log.Error("Failed to change user information:", zap.Error(err))
		response.FailWithMessage("Failed to change user information", c)
		return
	}
	response.OkWithMessage("Successfully changed user information", c)
}

// UserWeather 获取天气
// UserWeather 获取用户天气
// @Summary 根据IP地址获取天气信息
// @Tags 用户模块
// @Produce json
// @Success 200 {object} response.Response{data=string}
// @Router /user/weather [get]
// 实现逻辑：
// 1. 通过客户端IP获取地理位置
// 2. 查询高德天气API
// 3. 缓存结果1小时
// 返回示例：
// "广东-广州 晴 28℃ 东南风3级 湿度65%"
func (userApi *UserApi) UserWeather(c *gin.Context) {
	ip := c.ClientIP()
	weather, err := userService.UserWeather(ip)
	if err != nil {
		global.Log.Error("Failed to get user weather", zap.Error(err))
		response.FailWithMessage("Failed to get user weather", c)
		return
	}
	response.OkWithData(weather, c)
}

// UserChart 获取用户图表数据，登录和注册人数
// UserChart 用户增长图表
// @Summary 获取注册/登录统计图表数据
// @Tags 用户模块
// @Produce json
// @Param date query int true "统计天数(最大30天)"
// @Success 200 {object} response.Response{data=response.UserChart}
// @Router /user/chart [get]
// 数据结构：
// - DateList 日期序列
// - RegisterData 每日注册数
// - LoginData 每日登录数
func (userApi *UserApi) UserChart(c *gin.Context) {
	var req request.UserChart
	err := c.ShouldBindQuery(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	data, err := userService.UserChart(req)
	if err != nil {
		global.Log.Error("Failed to get user chart:", zap.Error(err))
		response.FailWithMessage("Failed to user chart", c)
		return
	}
	response.OkWithData(data, c)
}

// UserList 获取用户列表
// UserList 用户列表查询
// @Summary 管理员获取用户列表
// @Tags 用户管理
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param uuid query string false "用户UUID"
// @Param register query int false "注册方式(1:邮箱 2:QQ)"
// @Success 200 {object} response.PageResult
// @Router /user/list [get]
// 权限要求：管理员角色
func (userApi *UserApi) UserList(c *gin.Context) {
	var pageInfo request.UserList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := userService.UserList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get user list:", zap.Error(err))
		response.FailWithMessage("Failed to get user list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}

// UserFreeze 冻结用户
// UserFreeze 冻结用户
// @Summary 管理员冻结用户账户
// @Tags 用户管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body request.UserOperation true "用户操作参数"
// @Success 200 {object} response.Response
// @Router /user/freeze [post]
// 实现逻辑：
// 1. 校验管理员权限
// 2. 更新用户冻结状态
// 安全机制：
// - 操作记录留痕
// - 禁止冻结管理员账户
func (userApi *UserApi) UserFreeze(c *gin.Context) {
	var req request.UserOperation
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.UserFreeze(req)
	if err != nil {
		global.Log.Error("Failed to freeze user:", zap.Error(err))
		response.FailWithMessage("Failed to freeze user", c)
		return
	}
	response.OkWithMessage("Successfully freeze user", c)
}

// UserUnfreeze 解冻用户
// UserUnfreeze 解冻用户
// @Summary 管理员解冻用户账户
// @Tags 用户管理
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param body body request.UserOperation true "用户操作参数"
// @Success 200 {object} response.Response
// @Router /user/unfreeze [post]
// 注意事项：
// - 解冻后用户需重新登录
// - 操作记录关联管理员ID
func (userApi *UserApi) UserUnfreeze(c *gin.Context) {
	var req request.UserOperation
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}
	err = userService.UserUnfreeze(req)
	if err != nil {
		global.Log.Error("Failed to unfreeze user:", zap.Error(err))
		response.FailWithMessage("Failed to unfreeze user", c)
		return
	}
	response.OkWithMessage("Successfully unfreeze user", c)
}

// UserLoginList 获取登录日志列表
// UserLoginList 登录日志查询
// @Summary 获取用户登录记录
// @Tags 用户管理
// @Security ApiKeyAuth
// @Produce json
// @Param page query int false "页码"
// @Param page_size query int false "每页数量"
// @Param uuid query string false "用户UUID"
// @Success 200 {object} response.PageResult
// @Router /user/login-list [get]
// 返回字段：
// - LoginAt 登录时间
// - IP 登录IP
// - Device 登录设备
// - Location 登录地点
func (userApi *UserApi) UserLoginList(c *gin.Context) {
	var pageInfo request.UserLoginList
	err := c.ShouldBindQuery(&pageInfo)
	if err != nil {
		response.FailWithMessage(err.Error(), c)
		return
	}

	list, total, err := userService.UserLoginList(pageInfo)
	if err != nil {
		global.Log.Error("Failed to get user login list:", zap.Error(err))
		response.FailWithMessage("Failed to get user login list", c)
		return
	}
	response.OkWithData(response.PageResult{
		List:  list,
		Total: total,
	}, c)
}
