package service

import (
	"errors"
	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
	"platform_report/dao"
	"strings"
)

type GraphService struct {
}

func NewGraphService() *GraphService {
	return new(GraphService)
}

func (c *GraphService) GetGraphConditions(graphID string) []Conditions {
	key := "graph_support"
	var supports []*dao.ReportGraphSupport
	if x, found := kvcache.Get(key); found {
		supports = x.([]*dao.ReportGraphSupport)
	}
	if supports == nil {
		supports, _ = dao.NewPlatformReport().GetGraphConditions()
		if supports != nil {
			kvcache.Set(key, supports, cache.DefaultExpiration)
		}
	}
	conditions := NewReportService().GetReportConditions()
	items := make([]Conditions, 0)
	for _, val := range supports {
		if val.GraphId == graphID {
			for _, cond := range conditions {
				if val.CondId == cond.CondId {
					items = append(items, Conditions{ParentId: cond.ParentId, CondId: cond.CondId, Name: cond.CondName, Type: cond.CondType, Val: cond.CondVal, Holder: cond.Holder, IsChar: cond.IsChar, Conduct: cond.Conduct})
				}
			}
		}
	}
	return items
}

func (c *GraphService) AddGraphRelation(uuid, graphId string) error {
	reportInvoke := NewReportService()
	defines, _ := reportInvoke.GetReportDefinesByUuid(uuid)
	if defines == nil {
		return errors.New("报表ID异常，数据未找到")
	}
	invoke := dao.NewPlatformReport()
	err := invoke.AddGraphRelation(defines.ReportId, graphId)
	reportInvoke.CC()
	return err
}

func (c *GraphService) RemoveGraphRelation(uuid, graphId string) error {
	reportInvoke := NewReportService()
	defines, _ := reportInvoke.GetReportDefinesByUuid(uuid)
	if defines == nil {
		return errors.New("报表ID异常，数据未找到")
	}
	invoke := dao.NewPlatformReport()
	err := invoke.RemoveGraphRelation(defines.ReportId, graphId)
	reportInvoke.CC()
	return err
}

func (c *GraphService) AddGraph(graphName, sql string, gtype, mk, auto int, x, y1, y2, tip string, supports []int, datasource int) error {
	if graphName == "" {
		return errors.New("请输入图表名称")
	}
	reportInvoke := NewReportService()
	uuid := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	invoke := dao.NewPlatformReport()
	graphSupport := make([]dao.ReportGraphSupport, 0, len(supports))
	for _, item := range supports {
		graphSupport = append(graphSupport, dao.ReportGraphSupport{GraphId: uuid, CondId: item})
	}
	err := invoke.AddGraph(uuid, graphName, mk, auto, sql, gtype, x, y1, y2, tip, graphSupport, datasource)
	reportInvoke.CC()
	return err
}

func (c *GraphService) ModifyGraph(graphID, graphName, sql string, gtype, mk, auto int, x, y1, y2, tip string, supports []int, datasource int) error {
	if graphName == "" {
		return errors.New("请输入图表名称")
	}
	reportInvoke := NewReportService()
	invoke := dao.NewPlatformReport()
	graphSupport := make([]dao.ReportGraphSupport, 0, len(supports))
	for _, item := range supports {
		graphSupport = append(graphSupport, dao.ReportGraphSupport{GraphId: graphID, CondId: item})
	}
	err := invoke.ModifyGraph(graphID, graphName, sql, gtype, auto, mk, x, y1, y2, tip, graphSupport, datasource)
	reportInvoke.CC()
	return err
}

func (c *GraphService) RemoveGraph(graphID string) error {
	reportInvoke := NewReportService()
	invoke := dao.NewPlatformReport()
	err := invoke.RemoveGraph(graphID)
	reportInvoke.CC()
	return err
}

func (c *GraphService) GetGraphList(reportID string) []*dao.ReportGraph {
	reportInvoke := NewReportService()
	defines, _ := reportInvoke.GetReportDefinesByUuid(reportID)
	if defines == nil {
		return nil
	}

	key := "report_graph" + reportID
	if x, found := kvcache.Get(key); found {
		return x.([]*dao.ReportGraph)
	}

	datas, _ := dao.NewPlatformReport().GetGraph(defines.ReportId)
	if datas == nil {
		return nil
	}
	kvcache.Set(key, datas, cache.DefaultExpiration)
	return datas
}

func (c *GraphService) GetAllGraphList() []*dao.ReportGraph {
	key := "all_graph"
	if x, found := kvcache.Get(key); found {
		return x.([]*dao.ReportGraph)
	}
	//fmt.Println(222222222222222)
	datas, _ := dao.NewPlatformReport().GetAllGraph()
	if datas == nil {
		return nil
	}
	kvcache.Set(key, datas, cache.DefaultExpiration)
	return datas
}

func (c *GraphService) GetAllGraphListByCondition(uuid string) []*dao.ReportGraph {
	reportInvoke := NewReportService()
	defines, _ := reportInvoke.GetReportDefinesByUuid(uuid)
	if defines == nil {
		return nil
	}
	graphs, _ := dao.NewPlatformReport().GetAllGraphByCondition(defines.ReportId)
	return graphs
}

func (c *GraphService) GetGraphByGraphId(graphId string) *dao.ReportGraph {
	ls := c.GetAllGraphList()
	for _, val := range ls {
		if val.GraphId == graphId {
			return val
		}
	}
	return nil
}

func (c *GraphService) GetGraphByID(reportID, graphID string) *dao.ReportGraph {
	graphs := c.GetGraphList(reportID)
	for _, val := range graphs {
		if val.GraphId == graphID {
			return val
		}
	}
	return nil
}

func (c *GraphService) LinkGraph(reportID, graphID string) error {
	reportInvoke := NewReportService()
	defines, _ := reportInvoke.GetReportDefinesByUuid(reportID)
	if defines == nil {
		return errors.New("404")
	}
	err := dao.NewPlatformReport().AddGraphRelation(defines.ReportId, graphID)
	NewReportService().CC()
	return err
}

func (c *GraphService) UnLinkGraph(reportID, graphID string) error {
	reportInvoke := NewReportService()
	defines, _ := reportInvoke.GetReportDefinesByUuid(reportID)
	if defines == nil {
		return errors.New("404")
	}
	err := dao.NewPlatformReport().RemoveGraphRelation(defines.ReportId, graphID)
	NewReportService().CC()
	return err
}
