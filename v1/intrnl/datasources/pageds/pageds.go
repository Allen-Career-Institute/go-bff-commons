package pageds

import (
	ds "github.com/Allen-Career-Institute/go-bff-commons/v1/framework/datasource"
)

type Handlers interface {
	GetPage() ds.HandlerFunc
}
