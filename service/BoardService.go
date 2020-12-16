package service

import "C"
import (
	"errors"
	"fmt"
	"math"
	"platform_report/dao"
	"platform_report/lib"
	"platform_report/pkg/util"
	"regexp"
	"strconv"
	"strings"

	"github.com/patrickmn/go-cache"
	uuid "github.com/satori/go.uuid"
)

type BoardService struct {
}

func NewBoardService() *BoardService {
	return new(BoardService)
}

func (c *BoardService) AddReportBoard(boardName string, createUserId int) error {
	ID := strings.ReplaceAll(uuid.NewV4().String(), "-", "")
	err := dao.NewPlatformReport().AddReportBoard(ID, boardName, createUserId)
	reportInvoke := NewReportService()
	reportInvoke.CC()
	return err
}

func (c *BoardService) ModifyBoardName(boardId, boardName string) error {
	if boardName == "" {
		return errors.New("请输入看板名称")
	}
	reportInvoke := NewReportService()
	err := dao.NewPlatformReport().ModifyBoardName(boardId, boardName)
	reportInvoke.CC()
	return err
}

func (c *BoardService) ModifyLayout(boardId string, layout int) error {
	reportInvoke := NewReportService()
	err := dao.NewPlatformReport().ModifyLayout(boardId, layout)
	reportInvoke.CC()
	return err
}

func (c *BoardService) RemoveBoardByID(boardId string) error {
	err := dao.NewPlatformReport().RemoveBoardByID(boardId)
	reportInvoke := NewReportService()
	reportInvoke.CC()
	return err
}

func (c *BoardService) GetBoardDefines() []*dao.ReportBoard {
	key := "report_board"
	if x, found := kvcache.Get(key); found {
		return x.([]*dao.ReportBoard)
	}
	datas, err := dao.NewPlatformReport().GetBoardDefines()
	if err != nil {
		return nil
	}
	kvcache.Set(key, datas, cache.DefaultExpiration)
	return datas
}

func (c *BoardService) GetBoardDefinesById(id string) *dao.ReportBoard {
	defines := c.GetBoardDefines()
	for _, d := range defines {
		if d.BoardId == id {
			return d
		}
	}
	return nil
}

func (c *BoardService) GetTargetDashBoard(columns []*dao.ReportTargetDefines, date, dt1, dt2 string) []*dao.ReportTargetDefines {
	var time1, time2 string
	if date != "" {
		if date == "custom" {
			time1, time2 = dt1, dt2
		} else {
			t1, t2 := util.GetTimeByType(date)
			time1 = lib.TimeFormat(t1)
			time2 = lib.TimeFormat(t2)
		}
	}
	factorReg := regexp.MustCompile(`\{[^\}]+\}`)
	flagReg := regexp.MustCompile(`\#.*?\#`)
	//fmt.Println("========", time1, time2)
	cols := make([]string, 0)
	formaleCols := make(map[string]string)
	for i := 0; i < len(columns); i++ {
		columns[i].Value = 0
	}
	for _, col := range columns {
		if col.Formula == "" {
			cols = append(cols, col.EngName)
		} else {
			formaleCols[col.EngName] = col.Formula
			factors := factorReg.FindAllString(col.Formula, -1)
			for _, item := range factors {
				factor := strings.ReplaceAll(item, "{", "")
				factor = strings.ReplaceAll(factor, "}", "")
				cols = append(cols, factor)
			}
		}
	}

	sunMaps := make(map[string]float64)
	items, _ := dao.NewReportTarget().GetDataIndex(cols, time1, time2)

	for _, item := range items {
		if val, ok := sunMaps[item.IndexEnName]; ok {
			sunMaps[item.IndexEnName] = val + math.Ceil(item.IndexValue*100)
		} else {
			sunMaps[item.IndexEnName] = math.Ceil(item.IndexValue * 100)
		}

	}

	for key, _ := range sunMaps {
		sunMaps[key] = sunMaps[key] / 100
	}

	//fmt.Println("=========", sunMaps)
	subDays := util.GetSubDays(util.TimeParse(time1), util.TimeParse(time2))
	//fmt.Println("subDays =====", subDays)

	for key, val := range formaleCols {
		//fmt.Println("11111111111111", key)
		factors := factorReg.FindAllString(val, -1)
		//fmt.Println("11111111111111", factors)
		flags := flagReg.FindAllString(val, -1)
		var n = -1
		var result float64
		for _, item := range factors {
			factor := strings.ReplaceAll(item, "{", "")
			factor = strings.ReplaceAll(factor, "}", "")
			//fmt.Println("11111111111111", factor)
			var currentVal float64
			if factor == "@@day" {
				currentVal = float64(subDays)
			} else {
				currentVal = sunMaps[factor]
			}
			//fmt.Println("222222222222222", currentVal)
			if n >= 0 {
				switch flags[n] {
				case "#addition#":
					result = currentVal + result
				case "#subtraction#":
					result = result - currentVal
				case "#multiplication#":
					result = currentVal * result
				case "#division#":
					if currentVal == 0 {
						result = 0
					} else {
						result, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", result/currentVal), 64)
					}
				}
			} else {
				result = currentVal
			}
			n++
		}
		sunMaps[key] = result
	}

	for i := 0; i < len(columns); i++ {
		if val, ok := sunMaps[columns[i].EngName]; ok {
			columns[i].Value = val
		}
	}

	return columns
}
