package dao

import (
	"platform_report/lib"
	"time"
)

type ReportFieldDefines struct {
	ActionId int       `json:"action_id" xorm:"not null pk autoincr INT(11)"`
	Uuid     string    `json:"uuid"`
	Name     string    `json:"name"`
	InTime   time.Time `json:"in_time"`
}

type ReportField struct {
	ParentId     int       `json:"parent_id" xorm:"-"`
	ActionId     int       `json:"action_id" xorm:"INT(11)"`
	FieldId      int       `json:"-" xorm:"not null pk autoincr INT(11)"`
	Types        string    `json:"types" xorm:"not null default '' comment('分类') VARCHAR(200)"`
	Field        string    `json:"field" xorm:"not null comment('字段名') VARCHAR(50)"`
	Name         string    `json:"name" xorm:"not null comment('标识名') VARCHAR(50)"`
	FormType     string    `json:"form_type" xorm:"not null comment('字段类型') VARCHAR(20)"`
	DefaultValue string    `json:"default_value" xorm:"not null default '' comment('默认值') VARCHAR(255)"`
	MaxLength    int       `json:"max_length" xorm:"not null default 0 comment(' 字数上限') INT(4)"`
	IsPrimary    int       `json:"is_primary" xorm:"default 0 comment('是否主键 1.是') TINYINT(1)"`
	IsIdentity   int       `json:"is_identity" xorm:"default 0 TINYINT(1)"`
	IsUnique     int       `json:"is_unique" xorm:"not null default 0 comment('是否唯一（1是，0否）') TINYINT(1)"`
	IsNull       int       `json:"is_null" xorm:"not null default 0 comment('是否必填（1是，0否）') TINYINT(1)"`
	IsShow       int       `json:"is_show" xorm:"not null default 1 comment('是否显示 1.显示') TINYINT(1)"`
	InputTips    string    `json:"input_tips" xorm:"not null default '' comment('输入提示') VARCHAR(100)"`
	Setting      int       `json:"setting" xorm:"comment('设置') INT(4)"`
	Value        string    `json:"value" xorm:"-"`
	OrderId      int       `json:"order_id" xorm:"not null default 0 comment('排序ID') INT(4)"`
	CreateTime   time.Time `json:"-" xorm:"not null comment('创建时间') DATETIME"`
}

func NewReportFieldInvoke() *ReportField {
	return new(ReportField)
}

func (c *ReportField) GetActions() (ret []*ReportFieldDefines, err error) {
	err = lib.InitXormMySql().Find(&ret)
	return
}

func (c *ReportField) GetFieldsByActionId(actionId int) (ret []*ReportField, err error) {
	err = lib.InitXormMySql().Where("action_id = ?", actionId).Find(&ret)
	return
}

func (c *ReportField) GetFieldsByReportId(reportId int) (ret []*ReportField, err error) {
	err = lib.InitXormMySql().SQL(`SELECT
					  rfd.*
					FROM report_field_relation rfr
					  INNER JOIN report_field  rfd
						ON rfr.action_id = rfd.action_id
					WHERE rfr.report_id = ? ORDER BY rfd.order_id `, reportId).Find(&ret)
	return
}

func (c *ReportField) GetActionsByReportId(reportId int) (*ReportFieldDefines, error) {
	item := new(ReportFieldDefines)
	_, err := lib.InitXormMySql().SQL(`SELECT
							  rfd.*
							FROM report_field_relation rfr
							  INNER JOIN report_field_defines  rfd
								ON rfr.action_id = rfd.action_id
							WHERE rfr.report_id = ? `, reportId).Get(item)
	//fmt.Println(5555555555555555)
	return item, err
}

func (c *ReportField) CreateActions(uuid, name string) int {
	item := new(ReportFieldDefines)
	item.Name = name
	item.Uuid = uuid
	item.InTime = time.Now()
	_, _ = lib.InitXormMySql().Insert(item)
	return item.ActionId
}

func (c *ReportField) UpdActions(uuid, name string) error {
	item := new(ReportFieldDefines)
	item.Name = name
	//item.Uuid = uuid
	_, err := lib.InitXormMySql().Where("uuid = ?", uuid).Cols("name").Update(item)
	return err
}

func (c *ReportField) UnBind(reportId, actionId int) error {
	_, err := lib.InitXormMySql().Exec(`DELETE FROM report_field_relation WHERE report_id=? AND action_id=?`, reportId, actionId)
	return err
}

func (c *ReportField) Bind(reportId, actionId int) error {
	_, _ = lib.InitXormMySql().Exec(`
		DELETE FROM report_field_relation WHERE report_id=?
				`, reportId)
	_, err := lib.InitXormMySql().Exec(`
		INSERT INTO report_field_relation (report_id, action_id)
				  VALUES (?,?)
				`, reportId, actionId)
	return err
}

func (c *ReportField) Save(actionID int, fields []ReportField) error {
	session := lib.InitXormMySql().NewSession()
	defer session.Clone()
	err := session.Begin()
	_, err = session.SQL("delete from report_field where action_id = ?", actionID).Execute()
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_, err = session.InsertMulti(fields)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_ = session.Commit()
	return nil
}

func (c *ReportField) Delete(actionId int) error {
	_, _ = lib.InitXormMySql().Exec(`
		DELETE FROM report_field_relation WHERE action_id=?
				`, actionId)
	_, err := lib.InitXormMySql().Exec(`
		DELETE FROM report_field_defines WHERE action_id=?
				 `, actionId)
	return err
}
