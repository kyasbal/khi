package logger

import "github.com/GoogleCloudPlatform/kubernetes-history-inspector/pkg/common"

type LogThrottleStatus = int

var (
	StatusNoThrottle         LogThrottleStatus = 0
	StatusJustBeforeThrottle LogThrottleStatus = 1
	StatusThrottled          LogThrottleStatus = 2
)

// LogThrottler is an interface to decide if the new record should be printed or not.
// It mainly for reducing log volume addressing same kind log.
type LogThrottler interface {
	// ThrottleStatus returns true only when it should be handled. This must be called once for a single log entry.
	ThrottleStatus(logKind string) LogThrottleStatus
}

type ConstantLogThrottle struct {
	counter         *common.ConcurrentCounter
	MaxCountPerKind int
}

func NewConstantLogThrottle(maxCountPerKind int) LogThrottler {
	return ConstantLogThrottle{
		counter:         common.NewDefaultConcurrentCounter(common.NewSuffixShardingProvider(16, 1)),
		MaxCountPerKind: maxCountPerKind,
	}
}

func (c ConstantLogThrottle) ThrottleStatus(logKind string) LogThrottleStatus {
	if logKind == "" {
		return StatusNoThrottle
	}
	cnt := c.counter.Incr(logKind)
	if cnt == c.MaxCountPerKind {
		return StatusJustBeforeThrottle
	} else if cnt > c.MaxCountPerKind {
		return StatusThrottled
	} else {
		return StatusNoThrottle
	}
}
