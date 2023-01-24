package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
	"github.com/sirupsen/logrus"

	agonesv1 "agones.dev/agones/pkg/apis/agones/v1"
	sdk "agones.dev/agones/sdks/go"
)

type Env struct {
	EnableHealthcheck   bool `envconfig:"ENABLE_HEALTHCHECK" default:"true"`
	EnableLifetimecheck bool `envconfig:"ENABLE_LIFETIMECHECK" default:"true"`
	EnableDeschedule    bool `envconfig:"ENABLE_DESCHEDULECHECK" default:"true"`

	HealthcheckDuration   time.Duration `envconfig:"HEALTHCHECK_DURATION" default:"1s"`
	LifetimecheckDuration time.Duration `envconfig:"LIFETIMECHECK_DURATION" default:"30m"`
	DescheduleDuration    time.Duration `envconfig:"DESCHEDULE_DURATION" default:"2h"`
}

func main() {
	logrus.Debug("setting environment variables...")
	var env Env
	err := envconfig.Process("", &env)
	if err != nil {
		log.Fatal(err.Error())
	}

	logrus.Debug("starting sidecar processes...")
	s, err := sdk.NewSDK()
	if err != nil {
		log.Fatalf(err.Error())
	}
	if env.EnableHealthcheck {
		go healthCheck(s, env.HealthcheckDuration)
	}
	if env.EnableLifetimecheck {
		go lifetimeCheck(s, env.LifetimecheckDuration)
	}
	if env.EnableDeschedule {
		go deschedule(s, env.DescheduleDuration)
	}
	select {} // keep main goroutine
}

func healthCheck(s *sdk.SDK, duration time.Duration) {
	logrus.Debug("healthcheck called")
	for range time.Tick(duration) {
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

// if expire gameserver's lifetime, shutdown it
func lifetimeCheck(s *sdk.SDK, duration time.Duration) {
	logrus.Debug("check lifetime called")
	for range time.Tick(duration) {
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
func deschedule(s *sdk.SDK, duration time.Duration) {
	logrus.Debug("deschedule called")
	for range time.Tick(duration) {
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
