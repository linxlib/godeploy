package models

import (
	"github.com/coreos/go-systemd/dbus"
	"github.com/linxlib/godeploy/base/models"
	"time"
)

// ServiceType
// @Enum
type ServiceType = int

const (
	IISSite ServiceType = iota
	ConsoleApp
	WindowsService
	Systemd
	Directory
)

// Status
// @Enum
type Status = int

const (
	Running Status = iota
	Stopped
	NotFound
)

func ServiceTypeString(t ServiceType) string {
	switch t {
	case IISSite:
		return "IISSite"
	case ConsoleApp:
		return "ConsoleApp"
	case WindowsService:
		return "WindowsService"
	case Systemd:
		return "Systemd"
	case Directory:
		return "Directory"
	default:
		return "Unknown"
	}
}

// @Body
type Service struct {
	*models.BaseModel `gorm:"embedded"`
	Name              string      `json:"name"`
	ServiceName       string      `json:"service_name"`
	Arg               string      `json:"arg"`
	RealPath          string      `json:"real_path"`
	ListenUrl         string      `json:"listen_url"`
	ServiceType       ServiceType `json:"service_type"`
	ServerID          uint        `json:"server_id"`
	LastDeployTime    time.Time   `json:"last_deploy_time"`
}

func (s *Service) Stop() bool {
	return false
}
func (s *Service) Status() ServiceStatus {
	switch s.ServiceType {
	case Systemd:
		conn, _ := dbus.NewSystemdConnection()
		defer conn.Close()
		uss, _ := conn.ListUnitsByNames([]string{s.ServiceName})
		status := uss[0].ActiveState
		pid := uss[0].JobId
		var s1 Status
		switch status {
		case "active":
			s1 = Running
		default:
			s1 = Stopped
		}
		return ServiceStatus{
			Status:  s1,
			PID:     int(pid),
			CmdLine: s.ServiceName,
		}
	case Directory:
		return ServiceStatus{
			Status:  Running,
			PID:     -1,
			CmdLine: s.ServiceName,
		}
	case IISSite:

	case ConsoleApp:
	case WindowsService:
	default:

	}
	return ServiceStatus{}
}
func (s *Service) OverwriteFrom(dir string) bool {
	return false
}
func (s *Service) Start() bool {
	return false
}

type ServiceStatus struct {
	Status  Status `json:"status"`
	PID     int    `json:"pid"`
	CmdLine string `json:"cmdline"`
}
