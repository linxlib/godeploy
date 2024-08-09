package models

import (
	"context"
	"github.com/coreos/go-systemd/dbus"
	"github.com/linxlib/godeploy/base/models"
	cp "github.com/otiai10/copy"
	"github.com/tombuildsstuff/golang-iis/iis"
	"github.com/tombuildsstuff/golang-iis/iis/applicationpools"
	"os"
	"strings"
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

func (s *Service) CheckExistColumns() map[string]any {
	m := make(map[string]any)
	m["name"] = s.Name
	return m
}

func (s *Service) Stop() bool {
	return true
}
func (s *Service) Status() ServiceStatus {
	switch s.ServiceType {
	case Systemd:
		conn, _ := dbus.NewSystemdConnectionContext(context.Background())
		defer conn.Close()
		uss, _ := conn.ListUnitsByNamesContext(context.Background(), []string{s.ServiceName})
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
		iisClient, err := iis.NewClient()
		if err != nil {
			return ServiceStatus{
				Status:  Stopped,
				PID:     -1,
				CmdLine: s.ServiceName,
			}
		}
		site, err := iisClient.Websites.Get(s.ServiceName)
		if err != nil {
			return ServiceStatus{
				Status:  Stopped,
				PID:     -1,
				CmdLine: s.ServiceName,
			}
		}
		pool, err := iisClient.AppPools.Get(site.ApplicationPool)
		if err != nil {
			return ServiceStatus{
				Status:  Stopped,
				PID:     -1,
				CmdLine: s.ServiceName,
			}
		}
		var s1 Status
		if pool.State == applicationpools.StateStarted {
			s1 = Running
		} else {
			s1 = Stopped
		}

		return ServiceStatus{
			Status:  s1,
			PID:     -1,
			CmdLine: pool.Name,
		}

	case ConsoleApp:

	case WindowsService:

	default:

	}
	return ServiceStatus{}
}
func (s *Service) OverwriteFrom(dir string) bool {
	opt := cp.Options{
		Skip: func(info os.FileInfo, src, dest string) (bool, error) {
			return strings.HasSuffix(src, ".git") || strings.HasSuffix(src, ".idea"), nil
		},
		OnDirExists: func(src, dest string) cp.DirExistsAction {
			return cp.Replace
		},
	}
	return cp.Copy(dir, s.RealPath, opt) == nil
}
func (s *Service) Start() bool {
	return true
}

type ServiceStatus struct {
	Status  Status `json:"status"`
	PID     int    `json:"pid"`
	CmdLine string `json:"cmdline"`
}
