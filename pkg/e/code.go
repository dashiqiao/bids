package e

const (
	SUCCESS        = 200
	ERROR          = 500
	INVALID_PARAMS = 400

	ERROR_OK                       = 0
	ERROR_EXIST                    = 10001
	ERROR_EXIST_FAIL               = 10002
	ERROR_NOT_EXIST                = 10003
	ERROR_GET_S_FAIL               = 10004
	ERROR_COUNT_FAIL               = 10005
	ERROR_ADD_FAIL                 = 10006
	ERROR_EDIT_FAIL                = 10007
	ERROR_DELETE_FAIL              = 10008
	ERROR_EXPORT_FAIL              = 10009
	ERROR_IMPORT_FAIL              = 10010
	ERROR_PARAM_ERROR              = 10011
	ERROR_PARAM_QUESTION_ERROR     = 10012
	ERROR_EXIST_RECORD             = 10013
	ERROR_NEED_PARAM               = 10014
	ERROR_NAME_EXIST               = 10015
	ERROR_NOT_DATA                 = 10016
	ERROR_PARAM_RANGE              = 10017
	ErrorRecordNotFound            = 10018
	ErrOverOneDay                  = 10019
	ErrAuth                        = 10020
	ErrParam                       = 10021
	ErrUploadPath                  = 100022
	ERROR_AUTH_CHECK_TOKEN_FAIL    = 20001
	ERROR_AUTH_CHECK_TOKEN_TIMEOUT = 20002
	ERROR_AUTH_TOKEN               = 20003
	ERROR_AUTH                     = 20004
	ERROR_ZSK_NO_RELATION          = 20005
	ERROR_ZSK_NOT_ALLOW_DEl        = 20006
	ERROR_DEL_FAIL_RELATION        = 20007
	ERROR_NOT_FIND                 = 20008

	ERROR_RPC_CLIENT_FAIL = 30001
	ERROR_RULE_AUTH_FAIL  = 30002
	ERROR_NO_AUTH         = 30003
)
