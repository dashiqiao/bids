package e

import "github.com/jinzhu/gorm"

var MsgFlags = map[int]string{
	SUCCESS:        "成功",
	ERROR:          "fail",
	INVALID_PARAMS: "请求参数错误",

	ERROR_EXIST:       "已存在该对象名称",
	ERROR_EXIST_FAIL:  "获取已存在对象失败",
	ERROR_NOT_EXIST:   "该对象不存在",
	ERROR_GET_S_FAIL:  "获取所有对象失败",
	ERROR_COUNT_FAIL:  "统计对象失败",
	ERROR_ADD_FAIL:    "新增对象失败",
	ERROR_EDIT_FAIL:   "修改对象失败",
	ERROR_DELETE_FAIL: "删除对象失败",
	ERROR_EXPORT_FAIL: "导出对象失败",
	ERROR_IMPORT_FAIL: "导入对象失败",
	ERROR_NEED_PARAM:  "缺少参数",
	ERROR_NAME_EXIST:  "名称已经存在",
	ERROR_NOT_DATA:    "暂无数据",
	ERROR_NOT_FIND:    "数据不存在",
	ERROR_PARAM_RANGE: "参数范围错误",

	ERROR_AUTH_CHECK_TOKEN_FAIL:    "Token鉴权失败",
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT: "Token已超时",
	ERROR_AUTH_TOKEN:               "Token生成失败",
	ERROR_AUTH:                     "Token错误",
	ErrAuth:                        "权限错误",
	ERROR_PARAM_QUESTION_ERROR:     "问题状态不对",
	ERROR_PARAM_ERROR:              "参数错误",
	ErrParam:                       "参数错误",
	ERROR_EXIST_RECORD:             "记录已存在,请勿重复添加",
	ERROR_ZSK_NO_RELATION:          "问题与分类至少没有归属关系",
	ERROR_ZSK_NOT_ALLOW_DEl:        "问题与分类至少要保持一种归属关系",
	ErrorRecordNotFound:            gorm.ErrRecordNotFound.Error(),
	ErrOverOneDay:                  "超过24小时不可删除",
	ErrUploadPath:                  "上传路径创建错误",
	ERROR_RPC_CLIENT_FAIL:          "rpc连接失败！",
	ERROR_RULE_AUTH_FAIL:           "权限验证失败！",
	ERROR_DEL_FAIL_RELATION:        "删除失败，存在关联数据",
	ERROR_NO_AUTH:                  "无权限！",
}

func GetMsg(code int) string {
	msg, ok := MsgFlags[code]
	if ok {
		return msg
	}

	return MsgFlags[ERROR]
}
