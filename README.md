# gsidecar
A container for healthcheck, lifetime check and deschedule AgonesSDKServer.  
This works as a sidecar of GameServer.

## Get image from DockerHub
https://hub.docker.com/repository/docker/skmochi/agones-sidecar/tags?page=1&ordering=last_updated

## Environment value
|  Key |  Value  | default | Description |
| ---- | ---- | ---- | ---- |
|  ENABLE_HEALTHCHECK | "true" or "false" | "true" | use healthcheck or not |
|  ENABLE_LIFETIMECHECK  | "true" or "false" | "true" | use lifetime check or not |
|  ENABLE_DESCHEDULECHECK  |  "true" or "false" |  "true" | use descheduler or not |
|  HEALTHCHECK_DURATION  |  time.Duration e.g. "1s" | "1s" | a duration of healthcheck |
|  LIFETIMECHECK_DURATION  |  time.Duration e.g. "1m"  | "30m" | a duration of lifetime check |
|  DESCHEDULE_DURATION  |  time.Duration e.g. "1h"  | "2h" | a duration of deschedule |


## What is LifetimeCheck?
This automatically shutdowns gameserver when the time set in the annotation "agones.dev/sdk-lifetime" is exceeded.  
The annotation must be set in Unixtime.  
This works every LIFETIMECHECK_DURATION time.


## What is Deschedule?
Therefore Gameserver is not evictable, their placement would be scatter.  
This option shutdown gameserver whose state is not "Allocated" automatically.  
This works every DESCHEDULE_DURATION time.
