package subscriber

import (
	"reflect"
	"strconv"

	"bitbucket.org/optiisolutions/go-common/datamodels"

	"encoding/json"

	"bitbucket.org/optiisolutions/go-common/messaging"
	"cloud.google.com/go/pubsub"
	"github.com/sirupsen/logrus"
)

func (s *Sub) AddNewShard(m *pubsub.Message) error {
	logrus.Infof("received message: %s", string(m.Data))

	var (
		messageData messaging.ChangeMessage
		//shard         *datamodels.Shard
		err error
	)

	err = json.Unmarshal(m.Data, &messageData)
	if err != nil {
		logrus.Errorf("wrong message data: %s", err.Error())
		return err
	}

	if messageData.Action != messaging.RecordUpdateNew && messageData.Action != messaging.RecordUpdateSync {
		logrus.Warnf("unsupported message action: %s", string(messageData.Action))
		return nil //not an error, at all
	}

	media, err := messageData.Filter.GetMedia()
	if err != nil {
		logrus.Errorf("cannot get message media filter: %s", err.Error())
		return err
	}

	if media != reflect.TypeOf(datamodels.Property{}).Name() {
		logrus.Warnf("unsupported message media: %s", string(messageData.Action))
		return nil //not an error, at all
	}

	shardIdStr, err := messageData.Filter.GetId()
	if err != nil {
		logrus.Errorf("cannot get message Id filter: %s", err.Error())
		return err
	}

	shardId, err := strconv.ParseInt(shardIdStr, 10, 64)
	if err != nil {
		logrus.Errorf("wrong message Id (not an int?) [%s]: %s", shardIdStr, err.Error())
		return err
	}

	err = s.dao.AddNewShard(&datamodels.Shard{
		ShardId:      shardId,
		PropertyName: "", //leaved empty at the moment (TODO: need to add additional http-request to get this?)
	})
	if err != nil {
		logrus.Errorf("cannot add new shard: %s", err.Error())
		return err
	}

	logrus.Infof("New Shard %d was created successfully", shardId)

	return nil
}
