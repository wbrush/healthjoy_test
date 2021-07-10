package setup

import (
	"github.com/wbrush/go-common/messaging"
	"github.com/wbrush/go-template-service/configuration"
	"github.com/wbrush/go-template-service/dao"
)

func StartUp(cfg *configuration.Config, dao dao.DataAccessObject, ps messaging.PublisherSubscriber) {
	//TODO: do any staff that should be done on startup
}
