module exam-store

go 1.18

require (
	exam-api v0.0.0-00010101000000-000000000000
	github.com/banzaicloud/logrus-runtime-formatter v0.0.0-20190729070250-5ae5475bae5e
	github.com/emicklei/go-restful/v3 v3.9.0
	github.com/lib/pq v1.10.6
	github.com/sirupsen/logrus v1.9.0
)

replace exam-api => ../api-service
