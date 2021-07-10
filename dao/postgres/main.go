package postgres

import (
	"github.com/sirupsen/logrus"
	"github.com/wbrush/go-template-service/configuration"
	"github.com/wbrush/go-template-service/services/api"
	"strconv"
)

type (
	PgDAO struct {
		db.BasePgDAO
		cfg *configuration.Config
	}
)

func NewPgDao(cfg *configuration.Config) (*PgDAO, error) {
	d := PgDAO{
		BasePgDAO: db.NewBasePgDAO(cfg.DbMigrationPath),
		cfg:       cfg,
	}

	err := d.Init(&cfg.DbParams)
	if err != nil {
		return &d, err /// returning dao pointer (not nil) here is to process cluster init errors
	}

	return &d, nil
}

//this is used only for integration tests
var dao *PgDAO

func GetPgDao() *PgDAO {
	return dao
}
func SetPgDao(d *PgDAO) {
	dao = d
}

func (d *PgDAO) Close() {
	if d.BaseDB == nil {
		return
	}

	err := d.BaseDB.Close()
	if err != nil {
		logrus.Fatalf("cannot close a base DB connection: %s", err.Error())
	}
}

func (d *PgDAO) buildSelfPath(id int64) (path string) {
	path = d.cfg.BaseUri
	if d.cfg.ServiceParams.Port != "" {
		path = path + ":" + d.cfg.ServiceParams.Port
	}
	path = path + configuration.APIBasePath + configuration.APIVersion + api.TemplatePath + "/" + strconv.Itoa(int(id))
	return
}
