package internal

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"log"
	"os"
)

type ApplicationProperties struct {
	Mysql MysqlProperties `yaml:"mysql"`
}

type MysqlProperties struct {
	Host     string `yaml:"host"`
	Port     string `yaml:"port"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Database string `yaml:"database"`
}

func ExtractApplicationProperties(activeProfile AppProfile) ApplicationProperties {
	// application-{activeProfile}.yml 을 /resource 에서 가져온다.
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
		printMysqlInfo(applicationProperties)
	}

	return applicationProperties
}

func printMysqlInfo(applicationProperties ApplicationProperties) {
	log.Println(fmt.Sprintf("mysql info : %s / %s / %s / %s / %s",
		applicationProperties.Mysql.Host,
		applicationProperties.Mysql.Port,
		applicationProperties.Mysql.Username,
		applicationProperties.Mysql.Password,
		applicationProperties.Mysql.Database,
	))
}

func getApplicationYamlRelativePath(activeProfile AppProfile) string {
	return fmt.Sprintf("../resources/application-%s.yaml", activeProfile.Value)
}
