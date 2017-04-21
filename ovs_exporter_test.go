package main

import (
	"os"
	"testing"
	"github.com/Sirupsen/logrus"
	"github.com/joatmon08/ovs_exporter/utils"
)

func CleanupTesting() {
	ids, err := utils.GetAllContainerIDs()
	if err != nil {
		logrus.Error(err)
		panic(err)
	}
	if len(ids) == 0 {
		return
	}
	logrus.Info("deleting all existing containers")
	for _, id := range ids {
		if err := utils.DeleteContainer(id); err != nil {
			logrus.Error(err)
		}
	}
}

func TestMain(m *testing.M) {
	if !testing.Short() {
		CleanupTesting()
	}
	result := m.Run()
	os.Exit(result)
}
