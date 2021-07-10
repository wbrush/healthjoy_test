package setup

import (
	"bitbucket.org/optiisolutions/go-common/messaging"
	"bitbucket.org/optiisolutions/go-template-service/configuration"
	"bitbucket.org/optiisolutions/go-template-service/dao"
)

func StartUp(cfg *configuration.Config, dao dao.DataAccessObject, ps messaging.PublisherSubscriber) {
	//TODO: do any staff that should be done on startup
}
