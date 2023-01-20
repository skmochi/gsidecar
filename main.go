package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/sirupsen/logrus"

	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	sdk "agones.dev/agones/sdks/go"
)

func main() {
	logrus.Debug("main called")
	s, err := sdk.NewSDK()
	if err != nil {
		log.Fatalf(err.Error())
	}
	go healthCheck(s)
	go checkLifetime(s)
	go deschedule(s)
	select {} // keep main goroutine
}

// do healthcheck every 1 second
func healthCheck(s *sdk.SDK) {
	logrus.Debug("healthcheck called")
	for range time.Tick(time.Second) {
		if err := s.Health(); err != nil {
			log.Fatalf(err.Error())
		}
	}
}

func getStateCertainly(s *sdk.SDK) string {
	logrus.Debug("get state called")
	for range time.Tick(time.Second) {
		gameServer, err := s.GameServer()
		if err != nil {
			continue
		}
		state := gameServer.Status.State
		return state
		// unixtime
		// createdTime := gameServer.ObjectMeta.CreationTimestamp
	}
	return ""
}

func shutdownCertainly(s *sdk.SDK) {
	fmt.Println("terminate called")
	for range time.Tick(time.Second) {
		if err := s.Shutdown(); err != nil {
			log.Fatalf(err.Error())
			continue
		}
		return
	}
}

// get annotation of key "agones.dev/sdk-lifetime"
func getAnnotationLifetimeCertainly(s *sdk.SDK) int64 {
	logrus.Debug("get annotation lifetime called")
	for range time.Tick(time.Second) {
		gameServer, err := s.GameServer()
		if err != nil {
			log.Fatalf(err.Error())
			continue
		}
		annotations := gameServer.ObjectMeta.Annotations
		_lifetime, ok := annotations["agones.dev/sdk-lifetime"]
		if !ok {
			continue
		}
		lifetime, err := strconv.ParseInt(_lifetime, 10, 64)
		if err != nil {
			continue
		}
		return lifetime
	}
	return 0
}

// check gameserver's lifetime every 5 minutes
// if expire gameserver's lifetime, shutdown it
func checkLifetime(s *sdk.SDK) {
	logrus.Debug("check lifetime called")
	for range time.Tick(30 * time.Minute) {
		lifetime := getAnnotationLifetimeCertainly(s)
		d := time.Now().Unix()
		diff := lifetime - d
		if diff < 0 {
			logrus.Debug("lifetime is over, terminating...")
			shutdownCertainly(s)
		}
	}
}

// shutdown gameserver which is not Allocated for deschedule
func deschedule(s *sdk.SDK) {
	logrus.Debug("deschedule called")
	for range time.Tick(time.Hour) {
		state := getStateCertainly(s)
		// do not shutdown when calculating
		allocated := agonesv1.GameServerStateAllocated
		if state == "" || state == string(allocated) {
			continue
		}
		logrus.Debug("descheduling...")
		shutdownCertainly(s)
	}
}
