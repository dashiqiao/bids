package dao

import (
	"platform_report/lib"
	"time"
)

type ExamineStatus int
type ExamineState int

const (
	Wait     ExamineStatus = 0
	Allow    ExamineStatus = 1
	NotAllow ExamineStatus = 2

	Up   ExamineState = 0
	Down ExamineState = 1
)

type ExamineObject struct {
	Id            int
	ExamineStatus int
	UserId        int64
	CreatedAt     time.Time
}

type DzTipOffAuditing struct {
	Id        int `gorm:"AUTO_INCREMENT"`
	AppId     int
	SubId     int
	Comment   string
	State     int
	UserId    int
	UserName  string
	CreatedAt time.Time
}

func GetTopicObject(id int) ExamineObject {
	var object ExamineObject
	lib.GetDbInstance().Raw("select id,examine_status,user_id,created_at from dz_topic where id =?", id).Scan(&object)
	return object
}

func GetDynamicObject(id int) ExamineObject {
	var object ExamineObject
	lib.GetDbInstance().Raw("select id,examine_status,user_id,created_at from dz_dynamic where id =?", id).Scan(&object)
	return object
}

func GetStrategyObject(id int) ExamineObject {
	var object ExamineObject
	lib.GetDbInstance().Raw("select id,examine_status,user_id,created_at from dz_strategy where id =?", id).Scan(&object)
	return object
}

func ExamineTopic(id int, status ExamineStatus, msg string) error {
	return lib.GetDbInstance().Exec("UPDATE dz_topic SET examine_status =?,examine_msg = ?,updated_at=NOW() WHERE id= ?",
		status, msg, id).Error
}

func ExamineDynamic(id int, status ExamineStatus, msg string) error {
	return lib.GetDbInstance().Exec("UPDATE dz_dynamic SET examine_status =?,examine_msg = ?,updated_at=NOW() WHERE id= ?",
		status, msg, id).Error
}

func ExamineStrategy(id int, status ExamineStatus, msg string) error {
	return lib.GetDbInstance().Exec("UPDATE dz_strategy SET examine_status =?,examine_msg = ?,updated_at=NOW() WHERE id= ?",
		status, msg, id).Error
}

func ShelfTopic(id int, status ExamineState) error {
	return lib.GetDbInstance().Exec("UPDATE dz_topic SET state =?,updated_at=NOW() WHERE id= ?",
		status, id).Error
}

func ShelfDynamic(id int, status ExamineState) error {
	return lib.GetDbInstance().Exec("UPDATE dz_dynamic SET state =?,updated_at=NOW() WHERE id= ?",
		status, id).Error
}

func ShelfStrategy(id int, status ExamineState) error {
	return lib.GetDbInstance().Exec("UPDATE dz_strategy SET state =?,updated_at=NOW() WHERE id= ?",
		status, id).Error
}

func ShelfAnswer(id int, status ExamineState) error {
	if status == Up {
		return lib.GetDbInstance().Exec("UPDATE dz_topic_answer SET deleted_at=NULL WHERE id= ?",
			status, id).Error
	} else {
		return lib.GetDbInstance().Exec("UPDATE dz_topic_answer SET deleted_at=NOW() WHERE id= ?",
			status, id).Error
	}

}

func AddJbAuditing(appId, subId int, comment string, state, userId int, userName string) error {
	item := DzTipOffAuditing{AppId: appId, SubId: subId, Comment: comment, State: state, UserId: userId, UserName: userName,
		CreatedAt: time.Now()}
	return lib.GetDbInstance().Create(&item).Error
}

func ExamineJbById(id, state int) error {
	return lib.GetDbInstance().Exec(`
			UPDATE dz_tip_off SET state = ? WHERE id = ?
			`, state, id).Error
}

func ExamineJbByCategory(appId, subId, state, category int) error {
	return lib.GetDbInstance().Exec(`
			UPDATE dz_tip_off SET state = ? WHERE app_id = ? AND sub_id= ? AND category=?;
			`, state, appId, subId, category).Error
}
