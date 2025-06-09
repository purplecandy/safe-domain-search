# Safe Domain Search

<p align="center">
  <img src="icon.png" alt="App Icon" width="100" />
</p>

Search domain availability **without getting tracked** or frontrun by registrars.

Most online domain search tools silently log your queries and sometimes share them with 3rd parties. This opens the door to a shady practice called **domain front-running** — where someone (often a registrar) sees your search and registers the domain before you can.

Safe Domain Search protects your query by running all checks **locally**, using Go and Wails, with no analytics or telemetry.

![Not Available Screenshot](screenshot-failed.png)

![Available Screenshot](screenshot-success.png)

## 🔍 What It Does

- Checks if a domain resolves via **DNS**
- Looks for existing **Name Server (NS) records**
- Optionally performs a **WHOIS lookup** for deeper validation
- Returns a structured result with the outcome of each check
- Does **not send your query to third-party services** or tracking tools

## 💡 Why This Matters

> “I checked availability at GoDaddy, and it was free. Hours later, it was listed as a premium domain for $500 — registered the same day by GoDaddy. I waited a week and bought it for the regular fee elsewhere. Never used GoDaddy again.”  
> — _pixelrow on Reddit_

Domain front-running is real, and it's not paranoia — it’s been reported [on Hacker News](https://news.ycombinator.com/item?id=37802959), [on Wikipedia](https://en.wikipedia.org/wiki/Domain_name_front_running), and in domain forums like [NamePros](https://www.namepros.com/threads/how-to-prevent-domain-name-searches-from-being-stolen.1263571/).

## 📦 Install & Run

```bash
# Clone the repo
git clone https://github.com/purplecandy/safe-domain-search
cd safe-domain-search

# Install dependencies
npm install
go mod tidy

# Run locally
wails dev

# Build production binary
wails build
```
