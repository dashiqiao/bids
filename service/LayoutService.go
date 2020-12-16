package service

import (
	"github.com/patrickmn/go-cache"
	"platform_report/dao"
)

type LayoutService struct {
}

func NewLayoutService() *LayoutService {
	return new(LayoutService)
}

func (c *LayoutService) Get() ([]*dao.ReportLayout, error) {
	key := "layout"
	if x, found := kvcache.Get(key); found {
		return x.([]*dao.ReportLayout), nil
	} else {
		ret, _ := dao.NewReportLayout().Get()
		if ret != nil {
			kvcache.Set(key, ret, cache.DefaultExpiration)
			return ret, nil
		}
	}
	return nil, nil
}

func (c *LayoutService) GetOne(id int) (*dao.ReportLayout, error) {
	datas, _ := c.Get()
	if datas == nil {
		return nil, nil
	}
	for _, data := range datas {
		if data.LayoutId == id {
			return data, nil
		}
	}
	return nil, nil
}

func (c *LayoutService) Save(title, info string) (int, error) {
	layout := dao.NewReportLayout()
	layout.Title = title
	layout.Info = info
	id, err := layout.Save()
	if id > 0 {
		kvcache.Flush()
	}
	return id, err
}

func (c *LayoutService) Delete(id int) error {
	layout := dao.NewReportLayout()
	err := layout.Delete(id)
	if err == nil {
		kvcache.Flush()
	}
	return err
}

func (c *LayoutService) ModifyInfo(id int, info string) error {
	layout := dao.NewReportLayout()
	maps := make(map[string]interface{})
	maps["info"] = info
	err := layout.Modify(id, maps)
	if err == nil {
		kvcache.Flush()
	}
	return err
}

func (c *LayoutService) ModifyTitle(id int, title string) error {
	layout := dao.NewReportLayout()
	maps := make(map[string]interface{})
	maps["title"] = title
	err := layout.Modify(id, maps)
	if err == nil {
		kvcache.Flush()
	}
	return err
}
