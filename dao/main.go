package dao

import (
	"net/url"

	"bitbucket.org/optiisolutions/go-common/db"
	"bitbucket.org/optiisolutions/go-template-service/datamodels"
)

type (
	Template interface {
		CreateTemplate(shardID int64, template *datamodels.Template) (isDuplicate bool, err error)
		GetTemplateById(shardID int64, id int64) (template *datamodels.Template, isFound bool, err error)
		ListTemplates(shardID int64, filters url.Values) (templates []datamodels.Template, total int, hasMore bool, err error)
		UpdateTemplate(shardID int64, template *datamodels.Template) (err error)
		DeleteTemplateById(shardID int64, id int64) (isFound bool, err error)
	}

	DataAccessObject interface {
		db.BaseDataAccessObject //this is need only if you plain to use transactions or some base additional features
		Template
	}
)
