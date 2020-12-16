package dao

import (
	"platform_report/lib"
	"time"
)

type CurrencyType int

var (
	GXZ CurrencyType = 1 // 贡献值
	JB  CurrencyType = 2 // 金币
)

type DzUserAsset struct {
	Id           int64 `gorm:"AUTO_INCREMENT"`
	UserId       int64
	CurrencyType int32
	Amount       int
	InAmount     int
	OutAmount    int
}

type DzUserAssetLog struct {
	Id           int64 `gorm:"AUTO_INCREMENT"`
	AppId        int32
	SubId        int32
	UserId       int64
	CurrencyType int32
	State        int32
	Amount       int
	BeforeAmount int
	AfterAmount  int
	Atom         string
	Remark       string
	CreatedAt    time.Time
}

func GetAsset(uid int64, ct CurrencyType) DzUserAsset {
	var asset DzUserAsset
	lib.GetDbInstance().Raw("select * from dz_user_asset where user_id = ? and currency_type =? limit 1 ", uid, int32(ct)).Scan(&asset)
	return asset
}

func AssetIncrement(appId, subId int, uid int64, ct CurrencyType, amount int, atom, remark string) error {
	var asset DzUserAsset
	tx := lib.GetDbInstance().Begin()
	tx.Debug().Raw("select * from dz_user_asset where user_id = ? and currency_type =? limit 1  FOR UPDATE", uid, int32(ct)).Scan(&asset)
	if asset.UserId == 0 {
		rank := DzUserAsset{UserId: uid, CurrencyType: int32(ct), Amount: amount, InAmount: amount,
			OutAmount: 0}
		e := tx.Debug().Create(&rank).Error
		if e != nil {
			tx.Rollback()
			return e
		}
	} else {
		e := tx.Debug().Exec("update dz_user_asset set amount=amount + ?,in_amount=in_amount + ? where user_id = ? and currency_type = ?", amount, amount, uid, int32(ct)).Error
		if e != nil {
			tx.Rollback()
			return e
		}
	}
	log := DzUserAssetLog{AppId: int32(appId), SubId: int32(subId), UserId: uid,
		CurrencyType: int32(ct), State: 1, Amount: amount,
		BeforeAmount: asset.Amount, AfterAmount: asset.Amount + amount,
		Atom: atom, Remark: remark, CreatedAt: time.Now()}
	e := tx.Debug().Create(&log).Error
	if e != nil {
		tx.Rollback()
		return e
	}
	tx.Commit()
	return nil
}
