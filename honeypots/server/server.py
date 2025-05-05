import json
import random
import time
from datetime import datetime

LOG_FILE = "/app/logs/server.log"


def now():
    return datetime.utcnow().isoformat() + "Z"


def random_ip():
    return (
        f"{random.randint(11, 223)}."
        f"{random.randint(0, 255)}."
        f"{random.randint(0, 255)}."
        f"{random.randint(1, 254)}"
    )


def random_user():
    return f"attacker{random.randint(1000, 9999)}"


def log(entry):
    entry["@timestamp"] = now()
    entry["honeypot"] = "server"
    entry["event_source"] = "simulated"
    json_line = json.dumps(entry)
    print(json_line, flush=True)            # For Filebeat Docker input
    with open(LOG_FILE, "a") as f:
        f.write(json_line + "\n")           # For Filebeat file input


def simulate_connection():
    ip = random_ip()
    user = random_user()

    log({
        "event_type": "connection",
        "command": "",
        "response": "connected to fake SQL honeypot",
        "username": user,
        "ip": ip
    })

    fake_cmds = [
        "SELECT * FROM users;",
        "SHOW TABLES;",
        "DROP TABLE accounts;",
        "UPDATE users SET admin=1 WHERE id=1;",
        "SELECT password FROM credentials;"
    ]

    for cmd in random.sample(fake_cmds, k=random.randint(2, 4)):
        log({
            "event_type": "command",
            "command": cmd,
            "response": "syntax error",
            "username": user,
            "ip": ip
        })
        time.sleep(random.uniform(0.5, 1.5))

    log({
        "event_type": "disconnect",
        "command": "",
        "response": "client disconnected",
        "username": user,
        "ip": ip
    })


if __name__ == "__main__":
    while True:
        simulate_connection()
        time.sleep(random.uniform(5, 10))