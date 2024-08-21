package models

import (
	"context"
	"github.com/coreos/go-systemd/dbus"
	"github.com/linxlib/godeploy/base/models"
	"github.com/linxlib/godeploy/pkgs/golang-iis/iis"
	"github.com/linxlib/godeploy/pkgs/golang-iis/iis/applicationpools"
	cp "github.com/otiai10/copy"
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

// ServiceTypeString returns the string representation of the given ServiceType.
func ServiceTypeString(t ServiceType) string {
	// Switch the ServiceType to get the string representation.
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

// CheckExistColumns returns a map of columns and their values that are non-zero for a given Service.
//
// The map contains the "name" column and its value from the Service.
//
// Returns:
//   - A map[string]any containing the "name" column and its value.
func (s *Service) CheckExistColumns() map[string]any {
	// Create a map to store the column and its value.
	m := make(map[string]any)
	// Add the "name" column and its value to the map.
	m["name"] = s.Name
	// Return the map.
	return m
}

// Stop stops the service.
//
// Returns true if the service is stopped successfully, false otherwise.
func (s *Service) Stop() bool {
	// Check the type of the service.
	switch s.ServiceType {
	case Systemd:
		// Create a systemd connection.
		conn, _ := dbus.NewSystemdConnectionContext(context.Background())
		defer conn.Close()
		// Stop the systemd unit.
		_, err := conn.StopUnitContext(context.Background(), s.ServiceName, "replace", nil)
		if err != nil {
			// If there is an error, return false.
			return false
		}
		// Return true if the service is stopped successfully.
		return true
	case Directory:
		// Directories are always stopped.
		return true
	case IISSite:
		// Stop the IIS site.
		// Create an IIS client.
		iisClient, err := iis.NewClient()
		if err != nil {
			// If there is an error, return false.
			return false
		}
		// Get the IIS site.
		site, err := iisClient.Websites.Get(s.ServiceName)
		if err != nil {
			// If there is an error, return false.
			return false
		}
		err = iisClient.Websites.Stop(site.Name)
		if err != nil {
			// If there is an error, return false.
			return false
		}
		// Stop the IIS application pool.
		err = iisClient.AppPools.Stop(site.ApplicationPool)
		if err != nil {
			// If there is an error, return false.
			return false
		}
		// Return true if the site is stopped successfully.
		return true
	case ConsoleApp:
		// Console apps are always stopped.
	case WindowsService:
		// Windows services are always stopped.
	default:
		// Unknown service type, assume stopped.
		return true
	}
	// Return true if the service is stopped successfully.
	return true
}

// Status returns the status of the service.
//
// Returns a ServiceStatus struct containing the status, PID, command line, and error message.
func (s *Service) Status() ServiceStatus {
	switch s.ServiceType {
	case Systemd: // Systemd service
		// Create a systemd connection.
		conn, err := dbus.NewSystemdConnectionContext(context.Background())
		if err != nil {
			return ServiceStatus{
				Status:   Stopped,
				PID:      -1,
				CmdLine:  "systemctl status " + s.ServiceName,
				ErrorMsg: err.Error(),
			}
		}
		defer conn.Close()
		// Get the status of the service.
		uss, err := conn.ListUnitsByNamesContext(context.Background(), []string{s.ServiceName})
		if err != nil {
			return ServiceStatus{
				Status:   Stopped,
				PID:      -1,
				CmdLine:  "systemctl status " + s.ServiceName,
				ErrorMsg: "ListUnitsByNames:" + err.Error(),
			}
		}
		if len(uss) <= 0 {
			return ServiceStatus{
				Status:   Stopped,
				PID:      -1,
				CmdLine:  "systemctl status " + s.ServiceName,
				ErrorMsg: "ListUnitsByNames: unit not found",
			}
		}
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
			Status:   s1,
			PID:      int(pid),
			CmdLine:  "systemctl status " + s.ServiceName,
			ErrorMsg: "",
		}
	case Directory: // Directory service
		return ServiceStatus{
			Status:  Running,
			PID:     -1,
			CmdLine: s.ServiceName,
		}
	case IISSite: // IIS site
		iisClient, err := iis.NewClient()
		if err != nil {
			return ServiceStatus{
				Status:  Stopped,
				PID:     -1,
				CmdLine: "IIS site:" + s.ServiceName,
			}
		}
		site, err := iisClient.Websites.Get(s.ServiceName)
		if err != nil {
			return ServiceStatus{
				Status:   Stopped,
				PID:      -1,
				CmdLine:  "IIS AppPool:" + s.ServiceName,
				ErrorMsg: "FindWebsite:" + err.Error(),
			}
		}
		pool, err := iisClient.AppPools.Get(site.ApplicationPool)
		if err != nil {
			return ServiceStatus{
				Status:   Stopped,
				PID:      -1,
				CmdLine:  "IIS AppPool:" + s.ServiceName,
				ErrorMsg: "FindAppPool:" + err.Error(),
			}
		}
		var s1 Status
		if pool.State == applicationpools.StateStarted {
			s1 = Running
		} else {
			s1 = Stopped
		}
		if site.State == "Stopped" {
			s1 = Stopped
		} else {
			s1 = Running
		}

		return ServiceStatus{
			Status:  s1,
			PID:     -1,
			CmdLine: "IIS AppPool:" + pool.Name,
		}

	case ConsoleApp: // Console app

	case WindowsService: // Windows service

	default:

	}
	return ServiceStatus{}
}

// OverwriteFrom Overwrites the service files from the given directory.
//
// Parameter dir is the source directory to copy files from.
// Returns true if the copy operation is successful, false otherwise.
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

// Start starts the service.
//
// Returns true if the service is started successfully, false otherwise.
func (s *Service) Start() bool {
	switch s.ServiceType {
	case Systemd:
		// Connect to the systemd DBus API.
		conn, _ := dbus.NewSystemdConnectionContext(context.Background())
		defer conn.Close()

		// Start the systemd unit.
		_, err := conn.StartUnitContext(context.Background(), s.ServiceName, "replace", nil)
		if err != nil {
			// If the unit does not exist, do not return an error.
			if strings.Contains(err.Error(), "Unit ") && strings.Contains(err.Error(), " not found") {
				return false
			}
			return false
		}
		return true
	case Directory:
		// Directories are always started.
		return true
	case IISSite:
		// Start the IIS site.
		iisClient, err := iis.NewClient()
		if err != nil {
			return false
		}
		site, err := iisClient.Websites.Get(s.ServiceName)
		if err != nil {
			return false
		}

		err = iisClient.AppPools.Start(site.ApplicationPool)
		if err != nil {
			return false
		}

		err = iisClient.Websites.Start(site.Name)
		if err != nil {
			return false
		}
		return true
	case ConsoleApp:
	case WindowsService:
	default:
		// Unknown service type, assume started.
		return true
	}
	return true
}

// Restart restarts the service.
//
// Returns true if the service is restarted successfully, false otherwise.
func (s *Service) Restart() bool {
	switch s.ServiceType {
	case Systemd:
		// Restart the systemd unit.
		conn, _ := dbus.NewSystemdConnectionContext(context.Background())
		defer conn.Close()

		_, err := conn.RestartUnitContext(context.Background(), s.ServiceName, "replace", nil)
		if err != nil {
			return false
		}
		return true
	case Directory:
		// Directories are always restarted.
		return true
	case IISSite:
		// Restart the IIS site.
		iisClient, err := iis.NewClient()
		if err != nil {
			return false
		}
		site, err := iisClient.Websites.Get(s.ServiceName)
		if err != nil {
			return false
		}
		_ = iisClient.Websites.Stop(site.Name)
		_ = iisClient.AppPools.Stop(site.ApplicationPool)
		time.Sleep(time.Second * 2)
		err = iisClient.AppPools.Start(site.ApplicationPool)
		if err != nil {
			return false
		}
		err = iisClient.Websites.Start(site.Name)
		if err != nil {
			return false
		}
		return true
	case ConsoleApp:
	case WindowsService:
	default:
		return true
	}
	return true
}

type ServiceStatus struct {
	Status   Status `json:"status"`
	PID      int    `json:"pid"`
	CmdLine  string `json:"cmdline"`
	ErrorMsg string `json:"error_msg"`
}
