package config

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"path"
	"runtime"
)

type Conf struct {
	//db
	DbHost         string `yaml:"db_host"`
	DbPort         string `yaml:"db_port"`
	DbUsername     string `yaml:"db_username"`
	DbPassword     string `yaml:"db_password"`
	DbDatabase     string `yaml:"db_database"`
	DbPrefix       string `yaml:"db_prefix"`
	DbMaxIdleConns int    `yaml:"db_max_idle_conns"`
	DbMaxOpenConns int    `yaml:"db_max_open_conns"`
	//redis
	RedisHost     string `yaml:"redis_host"`
	RedisPort     string `yaml:"redis_port"`
	RedisPassword string `yaml:"redis_password"`
	RedisDB       int    `yaml:"redis_db"`
	//log
	LogPath         string `yaml:"log_path"`
	LogFileName     string `yaml:"log_file_name"`
	LogMaxAge       int    `yaml:"log_max_age"`
	LogRotationTime int    `yaml:"log_rotation_time"`
	//jwt
	JwtSecret string `yaml:"jwtSecret"`

	PageSize  int               `yaml:"page_size"`
	ChatAk    string            `yaml:"chat_ak"`
	ChatSk    string            `yaml:"chat_sk"`
	ExcelTemp map[string]string `yaml:"excel_temp"`

	Port string `yaml:"port"`

	RpcHost   string `yaml:"rpc_host"`
	ReportRpc string `yaml:"report_rpc"`

	UpyunOprator string `yaml:"upyun_oprator"`
	UpyunSecret  string `yaml:"upyun_secret"`

	CkHost     string `yaml:"click_house_host"`
	CkPort     string `yaml:"click_house_port"`
	CkUsername string `yaml:"click_house_username"`
	CkPassword string `yaml:"click_house_password"`
	CkDatabase string `yaml:"click_house_database"`

	PrestoHost     string `yaml:"presto_host"`
	PrestoPort     string `yaml:"presto_port"`
	PrestoUsername string `yaml:"presto_username"`
	PrestoPassword string `yaml:"presto_password"`

	DxzXzqCourseUrl          string `yaml:"dxz_xzq_course_url"`            // 学长圈获取课程地址
	DxzXzqLessonUrl          string `yaml:"dxz_xzq_lesson_url"`            // 学长圈获取课节地址
	DxzXzqEssayUrl           string `yaml:"dxz_xzq_essay_url"`             // 学长圈获取文章地址
	DxzXzqExhibitionUrl      string `yaml:"dxz_xzq_exhibition_url"`        // 学长圈获取展会地址
	DxzXzqThemeUrl           string `yaml:"dxz_xzq_theme_url"`             // 学长圈获取主题地址
	DxzXzqFirstLabelListUrl  string `yaml:"dxz_xzq_first_label_list_url"`  // 学长圈一级类目地址
	DxzXzqSecondLabelListUrl string `yaml:"dxz_xzq_second_label_list_url"` // 学长圈二级类目地址
	DxzXzqMultiPlayUrl       string `yaml:"dxz_xzq_multi_play_url"`        // 学长圈获取多码率视频播放地址

}

type tomlConfig struct {
	Servers map[string]server
}

type server struct {
	Catalog string
	Schema  string
}

func (c *Conf) GetConf() *Conf {
	//编译前使用， 编译后的可执行文件getCurrentPath() 返回为开发模式下的目录
	yamlFile, err := ioutil.ReadFile("./config/conf.yaml")
	//开发模式使用
	if err != nil {
		yamlFile, err = ioutil.ReadFile(getCurrentPath() + "/conf.yaml")
		if err != nil {
			log.Printf("yamlFile02.Get err   #%v ", err)
		}
	}

	err = yaml.Unmarshal(yamlFile, c)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}
	return c
}

func getCurrentPath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Dir(filename)
}
