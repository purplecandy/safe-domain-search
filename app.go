package main

import (
	"context"
	"errors"
	"fmt"
	"net"

	"github.com/likexian/whois"
	whoisparser "github.com/likexian/whois-parser"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// App struct
type App struct {
	ctx context.Context
}

type CheckResult struct {
	Status  string `json:"status"`          // passed, failed, skipped, error
	Details string `json:"details"`         // summary of result
	Error   string `json:"error,omitempty"` // optional technical error
}

type DomainCheckResponse struct {
	Domain      string                 `json:"domain"`
	IsAvailable bool                   `json:"isAvailable"`
	Checks      map[string]CheckResult `json:"checks"`
	Error       string                 `json:"error,omitempty"` // general error if any
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

func CheckDNS(domain string) (CheckResult, error) {
	ips, err := net.LookupIP(domain)
	if err != nil {
		var dnsError *net.DNSError
		if errors.As(err, &dnsError) && dnsError.IsNotFound {
			return CheckResult{
				Status:  "passed",
				Details: "Domain does not resolve to any IP address",
			}, nil

		}
		return CheckResult{
			Status:  "error",
			Details: "Unable to lookup IP address",
			Error:   err.Error(),
		}, err
	}

	if len(ips) > 0 {
		return CheckResult{
			Status:  "failed",
			Details: fmt.Sprintf("Domain resolves to IP: %v", ips[0]),
		}, nil
	}

	return CheckResult{
		Status:  "passed",
		Details: "No IP address found for domain",
	}, nil

}

func CheckNS(domain string) (CheckResult, error) {
	nsRecords, err := net.LookupNS(domain)
	if err != nil {
		if dnsErr, ok := err.(*net.DNSError); ok && dnsErr.Err == "no such host" {
			return CheckResult{
				Status:  "passed",
				Details: "Domain has no NS records or is unregistered",
			}, nil
		}
		return CheckResult{
			Status:  "error",
			Details: "Unable to lookup name servers",
			Error:   err.Error(),
		}, err
	}

	if len(nsRecords) == 0 {
		return CheckResult{
			Status:  "passed",
			Details: "No name servers found for domain",
		}, nil
	}

	return CheckResult{
		Status:  "failed",
		Details: fmt.Sprintf("Domain has name servers like %s", nsRecords[0].Host),
	}, nil
}

func CheckWhois(domain string) (CheckResult, error) {
	raw, err := whois.Whois(domain)
	if err != nil {
		return CheckResult{
			Status:  "error",
			Details: "Unable to perform WHOIS lookup",
			Error:   err.Error(),
		}, err
	}

	fmt.Println("Raw WHOIS data:", raw)

	parsed, err := whoisparser.Parse(raw)

	if err != nil {
		fmt.Println("Parsing error:", err)

		if errors.Is(err, whoisparser.ErrNotFoundDomain) {
			return CheckResult{
				Status:  "passed",
				Details: "No WHOIS record found",
			}, nil
		}

		if errors.Is(err, whoisparser.ErrPremiumDomain) {
			return CheckResult{
				Status:  "passed",
				Details: "Domain is premium or restricted",
			}, nil
		}

		if errors.Is(err, whoisparser.ErrReservedDomain) || errors.Is(err, whoisparser.ErrBlockedDomain) {
			return CheckResult{
				Status:  "failed",
				Details: "Domain is reserved or blocked",
			}, nil
		}

		return CheckResult{
			Status:  "error",
			Details: "Failed to parse WHOIS data",
			Error:   err.Error(),
		}, err
	}

	fmt.Println("Parsed WHOIS data:", parsed)

	if parsed.Registrar == nil || parsed.Registrar.Name == "" {
		return CheckResult{
			Status:  "passed",
			Details: "No WHOIS record found",
		}, nil
	}

	return CheckResult{
		Status:  "failed",
		Details: fmt.Sprintf("Registered via %s", parsed.Registrar.Name),
	}, nil
}

func (a *App) CheckDomain(domain string, forceAll bool) (DomainCheckResponse, error) {

	fmt.Println("Checking domain:", domain)
	fmt.Println("Force all checks:", forceAll)

	result := DomainCheckResponse{
		Domain: domain,
		Checks: make(map[string]CheckResult),
	}

	// Step 1: DNS Check

	fmt.Println("Performing DNS check...")
	dnsResult, _ := CheckDNS(domain)

	result.Checks["dns"] = dnsResult

	if dnsResult.Status == "passed" {
		result.IsAvailable = true
		result.Checks["ns"] = CheckResult{Status: "skipped", Details: "DNS check already confirmed registration"}
		result.Checks["whois"] = CheckResult{Status: "skipped", Details: "DNS check already confirmed registration"}

		if !forceAll {
			return result, nil
		}
	}

	fmt.Println("DNS check result:", dnsResult)
	// Step 2: NS Check
	fmt.Println("Performing NS check...")
	nsResult, _ := CheckNS(domain)
	result.Checks["ns"] = nsResult

	if nsResult.Status == "passed" {
		result.IsAvailable = true
		result.Checks["whois"] = CheckResult{Status: "skipped", Details: "NS check already confirmed registration"}
		if !forceAll {
			return result, nil
		}
	}

	fmt.Println("NS check result:", nsResult)

	// Step 3: WHOIS Check
	fmt.Println("Performing WHOIS check...")
	whoisResult, _ := CheckWhois(domain)
	result.Checks["whois"] = whoisResult

	if whoisResult.Status == "passed" {
		result.IsAvailable = true

	}

	return result, nil

}

func (a *App) OpenLink(link string) {
	runtime.BrowserOpenURL(a.ctx, link)
}
