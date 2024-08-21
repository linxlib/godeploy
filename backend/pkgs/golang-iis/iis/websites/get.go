package websites

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/linxlib/godeploy/pkgs/golang-iis/iis/helpers"
)

type getWebsiteResponse struct {
	Name                string                   `json:"name"`
	ApplicationPool     string                   `json:"applicationPool"`
	PhysicalPath        string                   `json:"physicalPath"`
	State               string                   `json:"state"`
	WebsiteStartsOnBoot bool                     `json:"serverAutoStart"`
	Limits              getWebsiteLimitsResponse `json:"limits"`
}

type getWebsiteLimitsResponse struct {
	MaxBandwidth int64 `json:"maxBandwidth"`
}

func (c *WebsitesClient) Get(name string) (*Website, error) {
	commands := fmt.Sprintf(`
Import-Module WebAdministration
Get-Website -Name %q | ConvertTo-Json -Compress
  `, name)

	stdout, _, err := c.Run(commands)
	if err != nil {
		return nil, fmt.Errorf("Error retrieving Website: %+v", err)
	}

	var site getWebsiteResponse
	if out := stdout; out != nil && *out != "" {
		v := *out
		err := json.Unmarshal([]byte(v), &site)
		if err != nil {
			return nil, fmt.Errorf("Error unmarshalling Website %q: %+v", name, err)
		}
	}

	if site.Name == "" {
		return nil, fmt.Errorf("Website %q was not found", name)
	}

	website := Website{
		Name:                         site.Name,
		ApplicationPool:              site.ApplicationPool,
		PhysicalPath:                 helpers.FixPowerShellPath(site.PhysicalPath),
		State:                        site.State,
		StartsOnBoot:                 site.WebsiteStartsOnBoot,
		MaxBandwidthPerSecondInBytes: site.Limits.MaxBandwidth,
	}
	return &website, nil
}

func (c *WebsitesClient) GetAll() ([]*Website, error) {
	commands := `
Import-Module WebAdministration
Get-Website | ConvertTo-Json -Compress
  `

	stdout, _, err := c.Run(commands)
	if err != nil {
		return nil, fmt.Errorf("error retrieving Website: %+v", err)
	}

	var sites []getWebsiteResponse
	var site getWebsiteResponse
	if out := stdout; out != nil && *out != "" {
		v := *out
		if !strings.HasPrefix(v, "[") {
			err := json.Unmarshal([]byte(v), &site)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling Websites: %+v", err)
			}
			sites = append(sites, site)
		} else {
			err := json.Unmarshal([]byte(v), &sites)
			if err != nil {
				return nil, fmt.Errorf("error unmarshalling Websites: %+v", err)
			}
		}

	}

	if len(sites) == 0 {
		return nil, fmt.Errorf("nothing not found")
	}
	websites := make([]*Website, len(sites))
	for i, site := range sites {
		websites[i] = &Website{
			Name:                         site.Name,
			ApplicationPool:              site.ApplicationPool,
			PhysicalPath:                 helpers.FixPowerShellPath(site.PhysicalPath),
			State:                        site.State,
			StartsOnBoot:                 site.WebsiteStartsOnBoot,
			MaxBandwidthPerSecondInBytes: site.Limits.MaxBandwidth,
		}
	}

	return websites, nil
}
