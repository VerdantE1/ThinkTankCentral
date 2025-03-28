package api

import "ThinkTankCentral/service"

type ApiGroup struct {
	BaseApi
	UserApi
	ImageApi
	ArticleApi
}

var ApiGroupApp = new(ApiGroup)

var baseService = service.ServiceGroupApp.BaseService
var userService = service.ServiceGroupApp.UserService
var jwtService = service.ServiceGroupApp.JwtService
var qqService = service.ServiceGroupApp.QQService
var imageService = service.ServiceGroupApp.ImageService
var articleService = service.ServiceGroupApp.ArticleService
