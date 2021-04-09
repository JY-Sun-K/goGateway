package controller

import (
	"encoding/json"
	"fmt"
	"github.com/e421083458/golang_common/lib"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"go_gateway/dao"
	"go_gateway/dto"
	"go_gateway/middleware"
	"go_gateway/public"
)

type AdminController struct {

}

func AdminRegister(group *gin.RouterGroup)  {
	adminLogin :=&AdminController{}
	group.GET("/admin_info",adminLogin.AdminInfo)
	group.POST("/change_pwd",adminLogin.ChangePwd)
}

// AdminInfo godoc
// @Summary 管理员信息
// @Description 管理员信息
// @Tags 管理员接口
// @ID /admin/admin_info
// @Accept  json
// @Produce  json
// @Success 200 {object} middleware.Response{data=dto.AdminInfoOutput} "success"
// @Router /admin/admin_info [get]
func (adminInfo *AdminController)AdminInfo(c *gin.Context)  {
	session := sessions.Default(c)
	sessionInfo:= session.Get(public.AdminSessionInfoKey)
	adminSessionInfo :=&dto.AdminSessionInfo{}
	if err:=json.Unmarshal([]byte(fmt.Sprint(sessionInfo)),adminSessionInfo);err!=nil{
		middleware.ResponseError(c,2000,err)
		return
	}
	out:=&dto.AdminInfoOutput{
		ID:           adminSessionInfo.ID,
		Name:         adminSessionInfo.UserName,
		LoginTime:    adminSessionInfo.LoginTime,
		Avatar:       "https://editorial.designtaxi.com/editorial-images/news-totorobics27062016/7.gif",
		Introduction: "I am a super administrator",
		Roles:        []string{"admin"},
	}
	middleware.ResponseSuccess(c,out)
}

// ChangePwd godoc
// @Summary 管理员密码修改
// @Description 管理员密码修改
// @Tags 管理员接口
// @ID /admin/change_pwd
// @Accept  json
// @Produce  json
// @Param body body dto.ChangePwdInput true "body"
// @Success 200 {object} middleware.Response{data=string} "success"
// @Router /admin/change_pwd [post]
func (adminInfo *AdminController)ChangePwd(c *gin.Context)  {
	params := &dto.ChangePwdInput{}
	if err:= params.BindValidParam(c);err!=nil{
		middleware.ResponseError(c,2000,err)
		return
	}

	session := sessions.Default(c)
	sessionInfo:= session.Get(public.AdminSessionInfoKey)
	adminSessionInfo :=&dto.AdminSessionInfo{}

	if err:=json.Unmarshal([]byte(fmt.Sprint(sessionInfo)),adminSessionInfo);err!=nil{
		middleware.ResponseError(c,2001,err)
		return
	}
	tx,err:=lib.GetGormPool("default")
	if err != nil {
		middleware.ResponseError(c,2002,err)
		return
	}



	admin := &dao.Admin{}


	admin,err= admin.Find(c,tx,&dao.Admin{
		UserName: adminSessionInfo.UserName,
	})
	if err != nil {
		middleware.ResponseError(c,2003,err)
		return
	}
	saltPassword := public.GenSaltPassword(admin.Salt,params.Password)
	admin.Password=saltPassword
	err=admin.Save(c,tx)
	if err != nil {
		middleware.ResponseError(c,2004,err)
		return
	}


	middleware.ResponseSuccess(c,"")
}
