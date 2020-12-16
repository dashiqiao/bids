package dao

import (
	"platform_report/lib"
	"time"
)

type ReportTargetDefines struct {
	Tip        string  `json:"tip"`
	TargetId   string  `json:"target_id"`
	Icon       string  `json:"icon"`
	EngName    string  `json:"eng_name"`
	ChsName    string  `json:"chs_name"`
	ReportId   string  `json:"report_id"`
	ReportName string  `json:"report_name" xorm:"-"`
	Value      float64 `json:"value" xorm:"-"`
	PK         float64 `json:"pk" xorm:"-"`
	Formula    string  `json:"formula"`
}

type DxzDataIndexDetails struct {
	Id          int       `json:"id" xorm:"not null pk autoincr comment('指标id') INT(11)"`
	IndexValue  float64   `json:"index_value" xorm:"comment('指标值') DECIMAL(16,2)"`
	IndexEnName string    `json:"index_en_name" xorm:"comment('指标英文名称') VARCHAR(50)"`
	IndexChName string    `json:"index_ch_name" xorm:"comment('指标中文名称') VARCHAR(50)"`
	ReportDate  time.Time `json:"report_date" xorm:"comment('数据日期') DATE"`
	EtlTime     time.Time `json:"etl_time" xorm:"comment('etl时间') DATETIME"`
}

type ReportBoardCard struct {
	BoardId   string    `json:"board_id"`
	CardId    string    `json:"card_id"`
	CardName  string    `json:"card_name"`
	Source    int       `json:"source"`
	Targets   string    `json:"targets"`
	GraphId   string    `json:"graph_id"`
	ReportId  string    `json:"-"`
	Number    int       `json:"number"`
	InTime    time.Time `json:"-"`
	GraphType int       `json:"graph_type" xorm:"-"`
	GraphName string    `json:"graph_name" xorm:"-"`
}

type ReportTarget struct {
}

func NewReportTarget() *ReportTarget {
	return new(ReportTarget)
}

func (c *ReportTarget) GetDataIndex(cols []string, time1, time2 string) (items []DxzDataIndexDetails, err error) {
	err = lib.InitXormMySql().In("index_en_name", cols).Where("report_date between ? and ?", time1, time2).Find(&items)
	return
}

//func (c *ReportTarget) GetTargetList() ([]ReportTargetDefines, error) {
//	var items []ReportTargetDefines
//	err := lib.InitXormMySql().Find(&items)
//	return items, err
//}

func (c *ReportTarget) Get(name string, p, limit int) ([]string, int64, error) {
	var
	(
		ids   []string
		err   error
		count int64
	)

	if name == "" {
		err = lib.InitXormMySql().Table(ReportTargetDefines{}).Cols("target_id").Asc("_id").Limit(limit, (p-1)*limit).Find(&ids)
		count, _ = lib.InitXormMySql().Table(ReportTargetDefines{}).Count()
	} else {
		err = lib.InitXormMySql().Table(ReportTargetDefines{}).Where("chs_name like ?", "%"+name+"%").Cols("target_id").Asc("_id").Limit(limit, (p-1)*limit).Find(&ids)
		count, _ = lib.InitXormMySql().Table(ReportTargetDefines{}).Where("chs_name like ?", "%"+name+"%").Count()
	}
	return ids, count, err
}

func (c *ReportTarget) GetMulti(ids []string) (ret []*ReportTargetDefines, err error) {
	err = lib.InitXormMySql().In("target_id", ids).Find(&ret)
	return
}

func (c *ReportTarget) SaveTarget(targetID, icon, eng, chs, link, formula, tip string) error {
	report := new(ReportTargetDefines)
	report.Tip = tip
	report.TargetId = targetID
	report.Icon = icon
	report.EngName = eng
	report.ChsName = chs
	report.ReportId = link
	report.Formula = formula
	_, err := lib.InitXormMySql().Insert(report)
	return err
}

func (c *ReportTarget) ModifyTarget(targetID, icon, eng, chs, formula, link, tip string) error {
	_, err := lib.InitXormMySql().SQL("update report_target_defines set icon=?,eng_name=?,chs_name=?,report_id=?,formula= ?,tip=? where target_id=?", icon, eng, chs, link, formula, tip, targetID).Execute()
	return err
}

func (c *ReportTarget) RemoveTarget(targetID string) error {
	_, err := lib.InitXormMySql().Exec("delete from report_target_defines where target_id=?", targetID)
	return err
}

func (c *ReportTarget) StepTwo(boardId string, cards []ReportBoardCard) error {
	session := lib.InitXormMySql().NewSession()
	defer session.Clone()
	err := session.Begin()
	_, err = session.SQL("delete from report_board_card where board_id= ? and  source in (2,3)", boardId).Execute()
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_, err = session.InsertMulti(cards)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_ = session.Commit()
	return nil
}

func (c *ReportTarget) StepOne(cardId, cardName, boardID, targets string) error {
	session := lib.InitXormMySql().NewSession()
	defer session.Clone()
	err := session.Begin()
	_, err = session.SQL("delete from report_board_card where board_id= ? and source = 1 ", boardID).Execute()
	if err != nil {
		_ = session.Rollback()
		return err
	}
	card := new(ReportBoardCard)
	card.BoardId = boardID
	card.CardId = cardId
	card.CardName = cardName
	card.Source = 1
	card.Targets = targets
	card.Number = 1
	card.InTime = time.Now()
	_, err = session.Insert(card)
	if err != nil {
		_ = session.Rollback()
		return err
	}
	_ = session.Commit()
	return nil
}

func (c *ReportTarget) GetCards(boardId string) (ret []*ReportBoardCard, err error) {
	err = lib.InitXormMySql().Table("report_board_card").Where("board_id=?", boardId).OrderBy("number").Find(&ret)
	return
}
