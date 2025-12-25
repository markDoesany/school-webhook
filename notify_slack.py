import requests
import json
import os

def send_slack_message(commit_message, author, commit_url):
    webhook_url = os.environ["SLACK_WEBHOOK_URL"]
    payload = {
        "text": f"🚀 New commit pushed to *main*!\n"
                f"*Author:* {author}\n"
                f"*Message:* {commit_message}\n"
                f"*Link:* {commit_url}"
    }
    headers = {"Content-Type": "application/json"}
    response = requests.post(webhook_url, data=json.dumps(payload), headers=headers)
    if response.status_code != 200:
        raise Exception(f"Slack request failed: {response.status_code}, {response.text}")

if __name__ == "__main__":
    # GitHub provides commit info via environment variables
    commit_message = os.environ.get("COMMIT_MESSAGE", "No message")
    author = os.environ.get("COMMIT_AUTHOR", "Unknown")
    commit_url = os.environ.get("COMMIT_URL", "")
    send_slack_message(commit_message, author, commit_url)
