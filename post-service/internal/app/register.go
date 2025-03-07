package app

import (
	"errors"
	"posts/internal/pkg/config"
	"posts/internal/usecase/kafka"
)

func Registries(k_handler *KafkaHandler, cfg *config.Config) error {
	brokers := []string{cfg.KafkaUrl}
	kcm := kafka.NewKafkaConsumerManager()

	if err := kcm.RegisterConsumer(brokers, "log-update", "log-u", k_handler.LogUpdate()); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			return errors.New("consumer for topic 'log-update' already exists")
		} else {
			return errors.New("error registering consumer:" + err.Error())
		}
	}

	if err := kcm.RegisterConsumer(brokers, "log-delete", "log-d", k_handler.LogDelete()); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			return errors.New("consumer for topic 'log-delete' already exists")
		} else {
			return errors.New("error registering consumer:" + err.Error())
		}
	}

	if err := kcm.RegisterConsumer(brokers, "post-update", "post-u", k_handler.PostUpdate()); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			return errors.New("consumer for topic 'post-update' already exists")
		} else {
			return errors.New("error registering consumer:" + err.Error())
		}
	}

	if err := kcm.RegisterConsumer(brokers, "post-delete", "post-d", k_handler.PostDelete()); err != nil {
		if err == kafka.ErrConsumerAlreadyExists {
			return errors.New("consumer for topic 'post-delete' already exists")
		} else {
			return errors.New("error registering consumer:" + err.Error())
		}
	}

	return nil
}
