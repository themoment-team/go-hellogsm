# go-hellogsm

## 소개

go-hellogsm 은 www.hellogsm.kr (광주소프트웨어마이스터고 입학지원시스템)의 batch jobs 를 모아둔 레포지토리 입니다.

* [1차 평가 배치 - firstEvaluationJob]()
* [2차 평가 배치 - secondEvaluationJob]()
* [최종 학과 배정 배치 - departmentAssignmentJob]()

## 개발자 가이드

### go-hellogsm 실행하기

```shell
cd ~/cmd
go build main.go
./main -profile local -jobs firstEvaluationJob
```

### 파라미터 소개

* profile: `local`, `stage`, `prod` 로 3가지를 지원하며 하나만 입력한다.
* jobs: `firstEvaluationJob`, `secondEvaluationJob`, `departmentAssignmentJob` 3가지를 지원하며 복수 입력 가능하다.
    * 복수 입력 예시: `-jobs firstEvaluationJob,secondEvaluationJob,fake` 혹여나 잘못 입력했다고 하더라도(`fake`) 무시되니 괜찮다.

### panic 사용과 error 리턴의 usecase (go-hellogsm 이하 process)

* process가 즉시 셧다운 되어야 한다면 `panic` 을 사용한다.
* 그 외 일반적인 예외(에러) 케이스에는 `error` 를 사용한다.

### ApplicationProperties 추가하기

* SpringBoot 와 같이 `application-{env}.yml` 파일을 해석할 수 있도록 구현 되어 있다.
    * 현재는 MySQL 정보만 필요하기 때문에 `~/internal/application_properties.go` 에 MysqlProperties 만 선언 돼 있다.
    * 이런식으로 `application-{env}.yml` 에서 추가적인 정보 해석이 필요한 경우 적절하게 추가하면 된다.

### DB 사용하기

* process 가 실제 비즈니스 로직을 처리하기 전에 `MyDB` 라는 전역 변수를 초기화 하기 때문에 안심하고 `MyDB` 를 호출해 쓰면 된다.
* `gorm` 을 사용한다. `gorm`에 ORM 특화 기능은 제하고 MyBatis(SQL mapper) 스타일로 사용한다. (ORM 을 쓰고 싶었던게 아님.)
* error 발생시 transaction rollback 처리 가이드는 TBD...

## 코딩 컨벤션

해당 레포지토리는 아래 레퍼런스들을 참고해 개발합니다.

* [Go 프로젝트 스탠다드 구조](https://github.com/golang-standards/project-layout/tree/master/internal)
* [Go 네이밍 룰](https://docs.google.com/document/u/2/d/1cBxRMfJm43U25akrLLRj6P4O3TsCk2lqYBeK4D9oCWM/mobilebasic?pli=1)
* [뱅크샐러드 Go 코딩 컨벤션](https://blog.banksalad.com/tech/go-best-practice-in-banksalad/)
* [Code Formatting 은 GoLand 표준](https://www.jetbrains.com/ko-kr/go/)