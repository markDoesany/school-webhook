import os
import sys
import requests
import json
from collections import Counter
from datetime import datetime, timedelta

def send_slack_message(text):
    webhook_url = os.environ["SLACK_WEBHOOK_URL"]
    payload = {"text": text}
    headers = {"Content-Type": "application/json"}
    response = requests.post(webhook_url, data=json.dumps(payload), headers=headers)
    if response.status_code != 200:
        raise Exception(f"Slack request failed: {response.status_code}, {response.text}")

def notify_commit():
    commit_message = os.environ.get("COMMIT_MESSAGE", "No message")
    author = os.environ.get("COMMIT_AUTHOR", "Unknown")
    commit_url = os.environ.get("COMMIT_URL", "")
    text = (
        f"🚀 New commit pushed to *main*!\n"
        f"*Author:* {author}\n"
        f"*Message:* {commit_message}\n"
        f"*Link:* {commit_url}"
    )
    send_slack_message(text)

def leaderboard():
    repo = os.environ["GITHUB_REPOSITORY"]  # e.g. "username/repo"
    token = os.environ["GITHUB_TOKEN"]

    # Time window: last 7 days
    now = datetime.utcnow()
    since_dt = now - timedelta(days=7)
    since = since_dt.isoformat() + "Z"

    url = f"https://api.github.com/repos/{repo}/commits?since={since}"
    headers = {"Authorization": f"token {token}"}
    resp = requests.get(url, headers=headers)
    commits = resp.json()

    authors = [c["commit"]["author"]["name"] for c in commits if "commit" in c]
    leaderboard = Counter(authors).most_common()

    # Format date range for display
    start_str = since_dt.strftime("%b %d, %Y")
    end_str = now.strftime("%b %d, %Y")

    if not leaderboard:
        text = f"📊 No commits found in *{repo}* between {start_str} and {end_str}."
    else:
        lines = [
            f"{i+1}. {author} – {count} commits"
            for i, (author, count) in enumerate(leaderboard)
        ]
        text = (
            f"🏆 Weekly Commit Champions for *{repo}*\n"
            f"As of {end_str} (covering {start_str} → {end_str}):\n"
            + "\n".join(lines)
        )

    send_slack_message(text)

if __name__ == "__main__":
    mode = sys.argv[1] if len(sys.argv) > 1 else "notify"
    if mode == "notify":
        notify_commit()
    elif mode == "leaderboard":
        leaderboard()
    else:
        print(f"Unknown mode: {mode}")