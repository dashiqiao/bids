package service

import (
	"encoding/json"
	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	"platform_report/dao"
	"platform_report/lib"
	"platform_report/pkg/util"
	"strings"
	"time"
)

type TargetService struct {
	prefix string
}

type Cards struct {
	CardName string `json:"card_name"`
	GraphId  string `json:"graph_id"`
	//ReportId string `json:"report_id"`
	No     int `json:"no"`
	Source int `json:"source"` //1.指标  2.图表  3.报表
}

func NewTargetService() *TargetService {
	return &TargetService{prefix: "target"}
}

func (c *TargetService) GetKey(uuid string) string {
	return c.prefix + ":" + uuid
}

func (c *TargetService) GetTargetDefines(name string, p, limit int) ([]*dao.ReportTargetDefines, int64) {
	invoke := dao.NewReportTarget()
	ids, count, _ := invoke.Get(name, p, limit)
	retMap := make([]*dao.ReportTargetDefines, 0, limit)
	notExist := make([]string, 0)
	for _, id := range ids {
		resMap := lib.GetRedis().HGetAll(c.GetKey(id)).Val()
		if len(resMap) > 0 {
			retMap = append(retMap, c.Fill(resMap))
		} else {
			notExist = append(notExist, id)
		}
	}
	if len(notExist) > 0 {
		notExistRes, _ := invoke.GetMulti(notExist)
		if notExistRes != nil {
			retMap = append(retMap, notExistRes...)
			for _, res := range notExistRes {
				lib.GetRedis().HMSet(c.GetKey(res.TargetId), util.StructToMap(res))
			}
		}
	}
	return retMap, count
}

func (c *TargetService) GetTargetDefinesByID(id string) *dao.ReportTargetDefines {
	resMap := lib.GetRedis().HGetAll(c.GetKey(id)).Val()
	if len(resMap) > 0 {
		return c.Fill(resMap)
	}
	items, _ := dao.NewReportTarget().GetMulti([]string{id})
	if len(items) > 0 {
		lib.GetRedis().HMSet(c.GetKey(id), util.StructToMap(items[0]))
		return items[0]
	}
	return nil
}

func (c *TargetService) SaveTargetDefines(icon, eng, chs, link, formula, tip string) error {
	ID := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	err := dao.NewReportTarget().SaveTarget(ID, icon, eng, chs, link, formula, tip)
	//NewReportService().CC()
	return err
}

func (c *TargetService) ModifyTarget(targetID, icon, eng, chs, link, formula, tip string) error {
	err := dao.NewReportTarget().ModifyTarget(targetID, icon, eng, chs, formula, link, tip)
	lib.GetRedis().Del(c.GetKey(targetID))
	return err
}

func (c *TargetService) RemoveTargetByID(targetID string) error {
	err := dao.NewReportTarget().RemoveTarget(targetID)
	lib.GetRedis().Del(c.GetKey(targetID))
	return err
}

func (c *TargetService) SaveBoardTarget(cardName, boardID, targets string) error {
	cardId := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	err := dao.NewReportTarget().StepOne(cardId, cardName, boardID, targets)
	NewReportService().CC()
	return err
}

func (c *TargetService) SaveBoardGraph(boardID string, cards []Cards) error {
	newCards := make([]dao.ReportBoardCard, 0, len(cards))
	for _, c := range cards {
		card := dao.ReportBoardCard{
			BoardId:  boardID,
			CardId:   strings.ReplaceAll(uuid.NewV4().String(), "-", ""),
			CardName: c.CardName,
			Source:   c.Source,
			Targets:  "",
			//GraphId:  c.GraphId,
			//ReportId: c.ReportId,
			Number: c.No,
			InTime: time.Now(),
		}
		if c.Source == 2 {
			card.GraphId = c.GraphId
		} else if c.Source == 3 {
			card.ReportId = c.GraphId
		}

		newCards = append(newCards, card)
	}
	err := dao.NewReportTarget().StepTwo(boardID, newCards)
	NewReportService().CC()
	return err
}

func (c *TargetService) GetCards(boardId string) []*dao.ReportBoardCard {
	key := "dashboard_cards" + boardId
	if x, found := kvcache.Get(key); found {
		return x.([]*dao.ReportBoardCard)
	}
	datas, err := dao.NewReportTarget().GetCards(boardId)
	if err != nil {
		return nil
	}
	kvcache.Set(key, datas, cache.DefaultExpiration)
	return datas
}

func (c *TargetService) GetBoardTarget(boardId string) (string, int, int, []*dao.ReportTargetDefines, []*dao.ReportBoardCard) {
	cards := c.GetCards(boardId)
	var title, targets string
	var number int
	graph := make([]*dao.ReportBoardCard, 0)
	graphInvoke := NewGraphService()
	defineInvoke := NewDefinesService()
	for _, card := range cards {
		if card.Source == 1 {
			title = card.CardName
			targets = card.Targets
			number = card.Number
		}
		if card.Source == 2 {
			gg := graphInvoke.GetGraphByGraphId(card.GraphId)
			if gg != nil {
				card.GraphType = gg.GraphType
				card.GraphName = gg.GraphName

			}
			graph = append(graph, card)
		}
		if card.Source == 3 {
			gg := defineInvoke.GetOne(card.ReportId)
			if gg != nil {
				card.GraphId = card.ReportId
				card.GraphName = gg.ReportName
			}
			graph = append(graph, card)
		}

	}
	var targetArr []string
	_ = json.Unmarshal([]byte(targets), &targetArr)
	//fmt.Println(targetArr)
	defines := make([]*dao.ReportTargetDefines, 0, len(targetArr))
	invoke := NewReportService()
	for i := 0; i < len(targetArr); i++ {
		ret := c.GetTargetDefinesByID(targetArr[i])
		if ret != nil {
			d, _ := invoke.GetReportDefinesByUuid(ret.ReportId)
			if d != nil {
				ret.ReportName = d.ReportName
			}
			defines = append(defines, ret)
		}
	}
	board := NewBoardService().GetBoardDefinesById(boardId)
	layout := 0
	if board != nil {
		layout = board.Layout
	}
	return title, layout, number, defines, graph
}

func (c *TargetService) Fill(resMap map[string]string) *dao.ReportTargetDefines {
	if len(resMap) > 0 {
		item := new(dao.ReportTargetDefines)
		item.ReportId = resMap["report_id"]
		item.TargetId = resMap["target_id"]
		item.Icon = resMap["icon"]
		item.EngName = resMap["eng_name"]
		item.ChsName = resMap["chs_name"]
		item.Tip = resMap["tip"]
		item.Formula = resMap["formula"]
		return item
	}
	return nil
}
