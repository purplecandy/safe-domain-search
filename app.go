package main

import (
	"context"
	"fmt"
	"net"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
)

// App struct
type App struct {
	ctx context.Context
}

type CheckResult struct {
	Status  string `json:"status"`  // passed, failed, skipped, error
	Details string `json:"details"` // human-readable summary
}

type DomainCheckResponse struct {
	Domain      string                 `json:"domain"`
	IsAvailable bool                   `json:"isAvailable"`
	Checks      map[string]CheckResult `json:"checks"`
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) CheckDomain(domain string, forceAll bool) (DomainCheckResponse, error) {
	result := DomainCheckResponse{
		Domain: domain,
		Checks: make(map[string]CheckResult),
	}

	// Step 1: DNS Check
	ips, err := net.LookupIP(domain)
	if err == nil && len(ips) > 0 {
		result.Checks["dns"] = CheckResult{
			Status:  "failed",
			Details: fmt.Sprintf("Domain resolves to IP: %v", ips[0]),
		}
		if !forceAll {
			result.IsAvailable = false
			result.Checks["ns"] = CheckResult{Status: "skipped", Details: "DNS check already confirmed registration"}
			result.Checks["whois"] = CheckResult{Status: "skipped", Details: "DNS check already confirmed registration"}
			return result, nil
		}
	} else {
		result.Checks["dns"] = CheckResult{
			Status:  "failed",
			Details: "No IP address found",
		}
	}

	// Step 2: NS Check
	nsRecords, err := net.LookupNS(domain)
	if err == nil && len(nsRecords) > 0 {
		result.Checks["ns"] = CheckResult{
			Status:  "failed",
			Details: fmt.Sprintf("Domain has name servers like %s", nsRecords[0].Host),
		}
		if !forceAll {
			result.IsAvailable = false
			result.Checks["whois"] = CheckResult{Status: "skipped", Details: "NS check already confirmed registration"}
			return result, nil
		}
	} else {
		result.Checks["ns"] = CheckResult{
			Status:  "failed",
			Details: "No name servers found",
		}
	}

	// Step 3: WHOIS Check
	raw, err := whois.Whois(domain)
	if err != nil {
		result.Checks["whois"] = CheckResult{
			Status:  "error",
			Details: fmt.Sprintf("WHOIS error: %v", err),
		}
		result.IsAvailable = false
		return result, nil
	}

	parsed, err := whoisparser.Parse(raw)
	if err != nil || parsed.Registrar.Name == "" {
		result.Checks["whois"] = CheckResult{
			Status:  "passed",
			Details: "No WHOIS record found",
		}
		result.IsAvailable = true
	} else {
		result.Checks["whois"] = CheckResult{
			Status:  "failed",
			Details: fmt.Sprintf("Registered via %s", parsed.Registrar.Name),
		}
		result.IsAvailable = false
	}

	return result, nil
}
