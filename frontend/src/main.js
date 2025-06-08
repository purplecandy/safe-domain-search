import { CheckDomain } from "../wailsjs/go/main/App";

document.getElementById("search-form").addEventListener("submit", async (e) => {
  e.preventDefault();

  console.log("Form submitted");

  const input = document.getElementById("domain-input");
  const force = document.getElementById("force-check").checked;

  const results = document.getElementById("results");
  const domain = input.value.trim();

  console.log("Domain to check:", domain);
  if (!domain) return;

  // Reset and show loading
  results.innerHTML = `
    <article>
      <header><strong>Checking "${domain}"...</strong></header>
      <progress></progress>
    </article>
  `;

  try {
    const res = await CheckDomain(domain, force);

    console.log("CheckDomain result:", res);
    // Build a result card
    const statusBadge = res.isAvailable
      ? `<mark class="secondary">✅ Available</mark>`
      : `<mark class="contrast">❌ Not Available</mark>`;

    const checksHtml = Object.entries(res.checks)
      .map(([key, check]) => {
        const color =
          {
            passed: "secondary",
            failed: "contrast",
            skipped: "muted",
            error: "warning",
          }[check.status] || "muted";

        return `
        <li>
          <strong>${key.toUpperCase()}</strong>: 
          <mark class="${color}">${check.status}</mark> – ${check.details}
        </li>
      `;
      })
      .join("");

    results.innerHTML = `
      <article>
        <header>
          <h3>${res.domain}</h3>
          ${statusBadge}
        </header>
        <ul>${checksHtml}</ul>
      </article>
    `;
  } catch (err) {
    results.innerHTML = `
      <article>
        <header><strong>Error</strong></header>
        <p>Something went wrong. Please try again.</p>
        <pre>${err.message || err}</pre>
      </article>
    `;
    console.error(err);
  }
});
