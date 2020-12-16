package dao

import (
	"platform_report/lib"
	"time"
)

type ReportBoard struct {
	BoardId      string    `json:"board_id"`
	BoardName    string    `json:"board_name"`
	Layout       int       `json:"layout"`
	CreateUserId int       `json:"create_user_id"`
	InTime       time.Time `json:"in_time"`
}


func (c *PlatformReport) AddReportBoard(boardId, boardName string, createUserId int) error {
	report := new(ReportBoard)
	report.BoardId = boardId
	report.BoardName = boardName
	report.CreateUserId = createUserId
	report.InTime = time.Now()
	_, err := lib.InitXormMySql().Insert(report)
	return err
}

func (c *PlatformReport) ModifyBoardName(boardId, boardName string) error {
	_, err := lib.InitXormMySql().SQL("update report_board set board_name=? where board_id=?", boardName, boardId).Execute()
	return err
}

func (c *PlatformReport) ModifyLayout(boardId string, layout int) error {
	_, err := lib.InitXormMySql().SQL("update report_board set layout=? where board_id=?", layout, boardId).Execute()
	return err
}

func (c *PlatformReport) RemoveBoardByID(boardId string) error {
	_, err := lib.InitXormMySql().Exec("delete from report_board where board_id=?", boardId)
	return err
}

func (c *PlatformReport) GetBoardDefines() (datas []*ReportBoard, err error) {
	err = lib.InitXormMySql().Table("report_board").OrderBy("in_time").Find(&datas)
	return
}

//func (c *PlatformReport) Test(reportId, sql string, args map[string]interface{}) {
//	_ = os.Remove("./sql/" + reportId + ".stpl")
//	file, _ := os.OpenFile("./sql/"+reportId+".stpl", os.O_RDWR|os.O_CREATE, 0766);
//	_, _ = file.WriteString(sql)
//	_ = file.Close()
//	//_ = lib.InitXormMySql().SqlTemplate.AddSqlTemplate(reportId+".stpl", sql)
//	results, err := lib.InitXormMySql().SqlTemplateClient(reportId+".stpl", &args).Query().List()
//	fmt.Println(results)
//	if err != nil {
//		fmt.Println(err.Error())
//	}
//
//}
