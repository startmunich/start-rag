import os

from slack_bolt import App
from slack_bolt.adapter.socket_mode import SocketModeHandler


# Install the Slack app and get xoxb- token in advance
app = App(token=os.environ["SLACK_BOT_TOKEN"])

@app.event("app_mention")
def event_test(event, say):
    print(event)

if __name__ == "__main__":
    SocketModeHandler(app, os.environ["SLACK_APP_TOKEN"]).start()