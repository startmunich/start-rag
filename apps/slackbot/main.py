import os

from dotenv import load_dotenv, dotenv_values
from slack_bolt import App
from slack_bolt.adapter.socket_mode import SocketModeHandler

load_dotenv(dotenv_path=".env.local")

# Install the Slack app and get xoxb- token in advance
app = App(token=os.environ["SLACK_APP_TOKEN"])

@app.event("app_mention")
def event_test(event, say):
    user = event["user"]
    channel = event["channel"]
    text = event["text"]
    print("Received Message\n | channel: " + channel + ", user: " + user + ", text: " + text)
    # TODO: Lang Magic and say answer
    return

if __name__ == "__main__":
    SocketModeHandler(app, os.environ["SLACK_BOT_TOKEN"]).start()