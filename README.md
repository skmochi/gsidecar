# agones-gameserver-sidecar
A container for healthcheck, lifetime check and deschedule AgonesSDKServer

## environment value
|  Key |  Value  | default | Description |
| ---- | ---- | ---- | ---- |
|  ENABLE_HEALTHCHECK | "true" | "true" or "false" | use healthcheck or not |
|  ENABLE_LIFETIMECHECK  | "true" | "true" or "false"  | use lifetime check or not |
|  ENABLE_DESCHEDULECHECK  |  "true" |  "true" or "false"  | use descheduler or not |
|  HEALTHCHECK_DURATION  |  time.Duration e.g. "1s" | "1s" | a duration of healthcheck |
|  LIFETIMECHECK_DURATION  |  time.Duration e.g. "1m"  | "30m" | a duration of lifetime check |
|  DESCHEDULE_DURATION  |  time.Duration e.g. "1h"  | "2h" | a duration of deschedule |

## What is LifetimeCheck?
xxx


## What is Deschedule?
xxx
