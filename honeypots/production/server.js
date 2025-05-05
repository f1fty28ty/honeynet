const fs = require("fs");
const path = require("path");

const LOG_FILE = "/app/logs/production.log";

function now() {
  return new Date().toISOString();
}

function randomIP() {
  return `${rand(11, 223)}.${rand(0, 255)}.${rand(0, 255)}.${rand(1, 254)}`;
}

function randomUser() {
  return `attacker${rand(1000, 9999)}`;
}

function rand(min, max) {
  return Math.floor(Math.random() * (max - min + 1)) + min;
}

function log(entry) {
  entry["@timestamp"] = now();
  entry["honeypot"] = "production";
  entry["event_source"] = "simulated";
  const line = JSON.stringify(entry);
  console.log(line); // Docker stdout
  fs.appendFileSync(LOG_FILE, line + "\n"); // Filebeat log
}

// Example simulation
function simulate() {
  const ip = randomIP();
  const user = randomUser();

  const urls = ["/", "/404", "/submit", "/login"];
  for (let i = 0; i < 3; i++) {
    const path = urls[Math.floor(Math.random() * urls.length)];
    const method = Math.random() < 0.5 ? "GET" : "POST";
    log({
      event_type: "http_request",
      command: `${method} ${path}`,
      response: "HTTP 404",
      username: user,
      ip: ip
    });
  }
}

setInterval(simulate, 8000);