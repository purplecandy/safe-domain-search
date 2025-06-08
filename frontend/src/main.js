import { Greet } from "../wailsjs/go/main/App";

document.getElementById("search-form").addEventListener("submit", async (e) => {
  e.preventDefault();

  const input = document.getElementById("domain-input");
  const domain = input.value;

  try {
    const response = await Greet(domain);
    document.getElementById("results").textContent = response;
  } catch (err) {
    document.getElementById("results").textContent = "Something went wrong.";
    console.error(err);
  }
});
