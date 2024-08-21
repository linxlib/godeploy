package applicationpools

// ApplicationPool is a reference to an Application Pool within IIS
type ApplicationPool struct {
	Name             string
	FrameworkVersion ManagedFrameworkVersion
	// MaxCPUPerInterval is the amount of (1/1000's) of % CPU allocated per interval (5s)
	MaxCPUPerInterval int64
	AutoStart         bool
	StartMode         StartMode
	State             State
}

// ManagedFrameworkVersion is the version of the .net Framework used in the Application Pool
type ManagedFrameworkVersion string

const (
	ManagedFrameworkVersionFour ManagedFrameworkVersion = "v4.0"
	ManagedFrameworkVersionTwo  ManagedFrameworkVersion = "v2.0"
	ManagedFrameworkVersionNone ManagedFrameworkVersion = ""
)

// StartMode is the start mode for the Application Pool
type StartMode string

const (
	StartModeAlwaysRunning StartMode = "AlwaysRunning"
	StartModeOnDemand      StartMode = "OnDemand"
)

// State returns the current running state of the Application Pool
type State string

const (
	StateStarted State = "Started"
	StateStopped State = "Stopped"
)
