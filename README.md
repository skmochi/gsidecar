# gsidecar
gsidecar is a container for healthcheck, lifetime check and deschedule Agones SDKServer and GameGerver.  
This works as a sidecar of GameServer.

## How to use
An example of Fleet.
```yaml
apiVersion: "agones.dev/v1"
kind: Fleet
metadata:
  name: test
  namespace: test
spec:
  template:
    spec:
      ports:
      - name: default
        containerPort: 7654
      # if there is more than one container, specify which one is the game server
      container: app
      template:
        spec:
          containers:
          - name: app
            image: APPLICATION_IMAGE
          - name: gsidecar
            image: skmochi/gsidecar:latest
            env:
            - name: ENABLE_HEALTHCHECK
              value: "true"
            - name: HEALTHCHECK_DURATION
              value: "5s"
```

## Get image from DockerHub
https://hub.docker.com/r/skmochi/gsidecar

## Environment value
|  Key |  type of Value | example of Value | default Value | Description |
| ---- | ---- | ---- | ---- | ---- |
|  ENABLE_HEALTHCHECK | bool | "true" | "true" | use healthcheck or not |
|  ENABLE_LIFETIMECHECK  | bool | "false" | "true" | use lifetime check or not |
|  ENABLE_DESCHEDULECHECK  | bool | "true" |  "true" | use descheduler or not |
|  HEALTHCHECK_DURATION  |  time.Duration | "1s" | "1s" | a duration of healthcheck |
|  LIFETIMECHECK_DURATION  |  time.Duration | "1m"  | "30m" | a duration of lifetime check |
|  DESCHEDULE_DURATION  |  time.Duration | "3h"  | "1h" | a duration of deschedule |


## What is LifetimeCheck?
This automatically shutdowns gameserver when the time set in the annotation "agones.dev/sdk-lifetime" is exceeded.  
The annotation must be set in Unixtime.  
This works every LIFETIMECHECK_DURATION time.


## What is Deschedule?
Therefore Gameserver is not evictable, their placement would be scatter.  
This option shutdown gameserver whose state is not "Allocated" automatically.  
This works every DESCHEDULE_DURATION time.
You should use this option with "Packed" storategy.
