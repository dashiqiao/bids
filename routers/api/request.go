package api

type UserRole struct {
	Menus    string      `json:"Menus"`
	Rules    string      `json:"Rules"`
	Title    string      `json:"Titile"`
	Type     int64       `json:"Type"`
	Types    int64       `json:"Types"`
	RuleList []AdminRule `json:"RuleList"`
	MenuList []AdminMenu `json:"MenuList"`
}
type AdminRule struct {
	Id        int    `json:"id"`
	Types     int    `json:"types"`      // 0系统设置 1工作台 2客户管理 3项目管理 4人力资源 5财务管理 6商业智能
	Title     string `json:"title"`      // 名称
	Name      string `json:"name"`       // 定义
	Func      string `json:"func"`       // 方法 配合权限管理
	Router    string `json:"router"`     // 路由 配合菜单管理
	Level     int    `json:"level"`      // 级别 1模块 2控制器 3操作  1 1级菜单 2 2级菜单 3 3级菜单
	Pid       int    `json:"pid"`        // 父id，默认0
	Status    int    `json:"status"`     // 状态，1启用，0禁用
	Dtime     int64  `json:"dtime"`      // 删除时间
	CompanyId int    `json:"company_id"` // 公司id
}

// 权限规则表
type AdminMenu struct {
	Id int `json:"id"`
	// Types     int    `json:"-"`          // 0系统设置 1工作台 2客户管理 3项目管理 4人力资源 5财务管理 6商业智能
	Title string `json:"title"` // 名称
	// Name      string `json:"-"`          // 定义
	// Func      string `json:"-"`          // 方法 配合权限管理
	Router    string `json:"router"`     // 路由 配合菜单管理
	Level     int    `json:"level"`      // 级别 1模块 2控制器 3操作  1 1级菜单 2 2级菜单 3 3级菜单
	Pid       int    `json:"pid"`        // 父id，默认0
	Sort      int    `json:"sort"`       // 排序
	Pic       string `json:"pic"`        // 图标
	Status    int    `json:"status"`     // 状态，1启用，0禁用
	Dtime     int64  `json:"dtime"`      // 删除时间
	CompanyId int    `json:"company_id"` // 公司id
}
