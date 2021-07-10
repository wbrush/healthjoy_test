package subscriber

import (
	"cloud.google.com/go/pubsub"
	"context"
	"encoding/json"
	"errors"
	"github.com/sirupsen/logrus"
	"github.com/wbrush/go-template-service/configuration"
	"github.com/wbrush/go-template-service/dao"
	"strings"
)

type (
	// Sub receives messages from pubsub
	Sub struct {
		config *configuration.Config

		dao dao.DataAccessObject
		ps  messaging.PublisherSubscriber
	}
)

// NewSub initializes a new instance of PubSub Subscriber with needed fields, but doesn't subscribe at all
func NewSub(cfg *configuration.Config, dao dao.DataAccessObject, ps messaging.PublisherSubscriber) *Sub {
	s := &Sub{
		config: cfg,
		dao:    dao,
		ps:     ps,
	}

	return s
}

// GracefulStop shuts down the server without interrupting any
// active connections.
func (s *Sub) GracefulStop(ctx context.Context) error {
	return s.ps.Close()
}

// Title returns the title.
func (s *Sub) Title() string {
	return "pubsub Subscriber"
}

// Run starts the PubSub subscriptions to a topics.
func (s *Sub) Run() {
	var err error

	for _, topic := range s.config.GCP.GetSubscriptionTopicsList() {
		if topic == "" {
			continue
		}

		topic := strings.Trim(topic, "\" ")

		go func(topicName string) {
			logrus.Infof("Subscribing to a %s topic...", topicName)
			err = s.ps.Subscribe(topicName,
				s.ps.MakeSubscriptionName(
					configuration.ServiceName,
					string(s.config.ServiceParams.Environment),
					topicName),
				s.receiveMessage)
			if err != nil {
				logrus.Errorf("can't subscribe by %s on topic %s: %v", s.config.GCP.ServicePubTopic, topicName, err)
			}
		}(topic)
	}

	select {}
}

func (s *Sub) receiveMessage(m *pubsub.Message) error {
	var messageData messaging.ChangeMessage

	err := json.Unmarshal(m.Data, &messageData)
	if err != nil {
		logrus.Errorf("wrong message data (not included BaseMessage?): %s", err.Error())
		return err
	}

	media, _ := messageData.Filter.GetMedia()
	err = errors.New("Received message that has no defined message handler:" + media)

	return err
}
