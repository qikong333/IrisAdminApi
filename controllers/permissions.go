package controllers

import (
	"errors"
	"io"
	"os"
	"time"

	"IrisAdminApi/files"
	"IrisAdminApi/models"
	"IrisAdminApi/tools"
	"IrisAdminApi/transformer"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/kataras/iris/v12"
	gf "github.com/snowlyg/gotransformer"
)

/**
* @api {get} /admin/permissions/:id 根据id获取权限信息
* @apiName 根据id获取权限信息
* @apiGroup Permissions
* @apiVersion 1.0.0
* @apiDescription 根据id获取权限信息
* @apiSampleRequest /admin/permissions/:id
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission
 */
func GetPermission(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	permission := models.GetPermissionById(id)

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, permTransform(permission), "操作成功"))
}

/**
* @api {post} /admin/permissions/ 新建权限
* @apiName 新建权限
* @apiGroup Permissions
* @apiVersion 1.0.0
* @apiDescription 新建权限
* @apiSampleRequest /admin/permissions/
* @apiParam {string} name 权限名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func CreatePermission(ctx iris.Context) {

	aul := new(models.PermissionRequest)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		_, _ = ctx.JSON(errorData(err))
	} else {
		err := validate.Struct(aul)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			//for _, err := range err.(validator.ValidationErrors) {
			//	fmt.Println()
			//	fmt.Println(err.Namespace())
			//	fmt.Println(err.Field())
			//	fmt.Println(err.Type())
			//	fmt.Println(err.Param())
			//	fmt.Println()
			//}
		} else {
			u := models.CreatePermission(aul)
			ctx.StatusCode(iris.StatusOK)
			if u.ID == 0 {
				_, _ = ctx.JSON(ApiResource(false, u, "操作失败"))
			} else {
				_, _ = ctx.JSON(ApiResource(true, u, "操作成功"))
			}
		}
	}
}

/**
* @api {post} /admin/permissions/:id/update 更新权限
* @apiName 更新权限
* @apiGroup Permissions
* @apiVersion 1.0.0
* @apiDescription 更新权限
* @apiSampleRequest /admin/permissions/:id/update
* @apiParam {string} name 权限名
* @apiParam {string} display_name
* @apiParam {string} description
* @apiParam {string} level
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func UpdatePermission(ctx iris.Context) {
	aul := new(models.PermissionRequest)

	if err := ctx.ReadJSON(&aul); err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		_, _ = ctx.JSON(errorData(err))
	} else {
		err := validate.Struct(aul)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			//for _, err := range err.(validator.ValidationErrors) {
			//	fmt.Println()
			//	fmt.Println(err.Namespace())
			//	fmt.Println(err.Field())
			//	fmt.Println(err.Type())
			//	fmt.Println(err.Param())
			//	fmt.Println()
			//}
		} else {
			id, _ := ctx.Params().GetInt("id")
			uid := uint(id)

			u := models.UpdatePermission(aul, uid)
			ctx.StatusCode(iris.StatusOK)
			if u.ID == 0 {
				_, _ = ctx.JSON(ApiResource(false, u, "操作失败"))
			} else {
				_, _ = ctx.JSON(ApiResource(true, u, "操作成功"))
			}
		}
	}
}

/**
* @api {delete} /admin/permissions/:id/delete 删除权限
* @apiName 删除权限
* @apiGroup Permissions
* @apiVersion 1.0.0
* @apiDescription 删除权限
* @apiSampleRequest /admin/permissions/:id/delete
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func DeletePermission(ctx iris.Context) {
	id, _ := ctx.Params().GetUint("id")
	models.DeletePermissionById(id)
	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, nil, "删除成功"))
}

/**
* @api {post} /admin/permissions/import 导入权限
* @apiName 导入权限
* @apiGroup ImportPermission
* @apiVersion 1.0.0
* @apiDescription 导入权限
* @apiSampleRequest /admin/permissions/import
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func ImportPermission(ctx iris.Context) {

	file, info, err := ctx.FormFile("file")
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, err.Error(), "导入失败"))
	}

	fullPath, err := files.GetUploadFileUPath(file, info, "excel")
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, err.Error(), "导入失败"))
	}

	out, err := os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, err.Error(), "导入失败"))
	}

	if out == nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, errors.New("excel 文件上次失败"), "导入失败"))
		return
	}

	defer out.Close()

	_, _ = io.Copy(out, file)

	f, err := excelize.OpenFile(fullPath)
	if err != nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, err.Error(), "导入失败"))
	}

	if f == nil {
		ctx.StatusCode(iris.StatusOK)
		_, _ = ctx.JSON(ApiResource(false, errors.New("excel 文件打开失败"), "导入失败"))
		return
	}

	// Excel 导入行数据转换
	// 获取 Sheet1 上所有单元格
	rows := f.GetRows("Sheet1")
	titles := map[string]string{"0": "Name", "1": "DisplayName", "2": "Description", "3": "Act"}
	for roI, row := range rows {
		if roI > 0 {
			// 将数组  转成对应的 map
			m := models.PermissionRequest{}
			x := gf.NewXlxsTransform(&m, titles, row, "", time.RFC3339, nil)
			err := x.XlxsTransformer()
			if err != nil {
				ctx.StatusCode(iris.StatusOK)
				_, _ = ctx.JSON(ApiResource(false, err.Error(), "导入失败"))
			}

			models.CreatePermission(&m)
		}
	}

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, nil, "导入成功"))
}

/**
* @api {get} /permissions 获取所有的权限
* @apiName 获取所有的权限
* @apiGroup Permissions
* @apiVersion 1.0.0
* @apiDescription 获取所有的权限
* @apiSampleRequest /permissions
* @apiSuccess {String} msg 消息
* @apiSuccess {bool} state 状态
* @apiSuccess {String} data 返回数据
* @apiPermission null
 */
func GetAllPermissions(ctx iris.Context) {
	offset := tools.ParseInt(ctx.URLParam("offset"), 1)
	limit := tools.ParseInt(ctx.URLParam("limit"), 20)
	name := ctx.FormValue("name")
	orderBy := ctx.FormValue("orderBy")

	permissions := models.GetAllPermissions(name, orderBy, offset, limit)

	ctx.StatusCode(iris.StatusOK)
	_, _ = ctx.JSON(ApiResource(true, permsTransform(permissions), "操作成功"))
}

func permsTransform(perms []*models.Permission) []*transformer.Permission {
	var rs []*transformer.Permission
	for _, perm := range perms {
		r := permTransform(perm)
		rs = append(rs, r)
	}
	return rs
}

func permTransform(perm *models.Permission) *transformer.Permission {
	r := &transformer.Permission{}
	g := gf.NewTransform(r, perm, time.RFC3339)
	_ = g.Transformer()
	return r
}
