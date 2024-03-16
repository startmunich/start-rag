import os
import json

from dotenv import load_dotenv, dotenv_values
from slack_bolt import App
from slack_bolt.adapter.socket_mode import SocketModeHandler

from apps.slackbot import ai_functionalitites

load_dotenv(dotenv_path=".env")

# load dict with registered users
with open("registered_users.json", "r") as f:
    registered_users = json.load(f)

# Install the Slack app and get xoxb- token in advance
app = App(token=os.environ["SLACK_APP_TOKEN"])

@app.event("app_mention")
def event_test(event, say):
    user = event["user"]
    channel = event["channel"]
    text = event["text"]
    # check if user is key in registered_users
    if user not in registered_users:
        say("You are not registered to use this bot")
        return
    
    print("Received Message\n | channel: " + channel + ", user: " + user + ", text: " + text)
    
    # Receive message from event, split and remove app mention "@StartGPT"
    text_split = text.split("<@U066N58KTUP>")
    query = text_split[1].strip() if len(text_split) > 1 else ""
    
    say(f"<@{user}>" + " " + ai_functionalitites.get_answer(query))

if __name__ == "__main__":
    SocketModeHandler(app, os.environ["SLACK_BOT_TOKEN"]).start()