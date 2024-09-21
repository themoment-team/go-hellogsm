package internal

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

var SafeApplicationProperties ApplicationProperties

type ApplicationProperties struct {
	Mysql MysqlProperties   `yaml:"mysql"`
	API   APIInfoProperties `yaml:"api"`
}

type APIInfoProperties struct {
	RelayAPI APIInfo `yaml:"relay-api"`
}

type APIInfo struct {
	URL string `yaml:"url"`
	Key string `yaml:"key"`
}

type MysqlProperties struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func InitApplicationProperties(activeProfile AppProfile) {
	applicationYamlName := getApplicationYamlRelativePath(activeProfile)
	yamlFile, err := os.ReadFile(applicationYamlName)

	if err != nil {
		panic(fmt.Sprintf("%s 를 찾을 수 없습니다.", applicationYamlName))
	}

	var applicationProperties ApplicationProperties
	err = yaml.Unmarshal(yamlFile, &applicationProperties)
	if err != nil {
		panic(err)
	} else {
		printApplicationProperties(applicationProperties)
	}

	// 전역 변수로 사용 가능하도록 한다.
	SafeApplicationProperties = applicationProperties
}

func printApplicationProperties(applicationProperties ApplicationProperties) {
	printMysqlInfo(applicationProperties)
	printAPI(applicationProperties)
}

func printAPI(properties ApplicationProperties) {
	log.Println(fmt.Sprintf("api info found : %s", properties.API.RelayAPI))
}

func printMysqlInfo(applicationProperties ApplicationProperties) {
	log.Println(fmt.Sprintf("mysql info found : %s / %s / %s / %s / %s",
		applicationProperties.Mysql.Host,
		applicationProperties.Mysql.Port,
		applicationProperties.Mysql.Username,
		applicationProperties.Mysql.Password,
		applicationProperties.Mysql.Database,
	))
}

// application-{activeProfile}.yml 을 /resource 에서 가져온다.
func getApplicationYamlRelativePath(activeProfile AppProfile) string {
	return fmt.Sprintf("../resources/application-%s.yaml", activeProfile.Value)
}
