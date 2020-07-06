# Mattermost Lunchbot Plugin
This plugin pairs random users from a channel with each other. They get asked to go to lunch together. The idea is for a team to get to know each other better, not just to stick to your people day in and day out.

```
John: /lunchbot topics add Geocaching
John: /lunchbot topics add playing guitar
Mike: /lunchbot blacklist add George
Mike: /lunchbot topics add Basketball
Mike: /lunchbot

Direct chat to Mike:
LunchBot: Hey Mike! You should meet up with John. You could talk about Basketball.

Direct chat to John:
LunchBot: Hey John! You should meet up with Mike. You could talk about Geocaching.
```


## Why?
In COVID times it's hard to get to know your colleagues by casually chatting by the watercooler. This bot enables these type of random interactions between everyone.

## Features
* Everyone can trigger to get paired up by using `/lunchbot`
* Let users set topics they'd like to talk about using `/lunchbot topics add <topic>`
* Let users blacklist certain users they don't want to get paired with using `/lunchbot blacklist add <username>`

## Contribute
This plugin is based on the [mattermost-plugin-starter-template](https://github.com/mattermost/mattermost-plugin-starter-template). See there on how to set everything up and test the plugin.

## Attributions
The logo is licensed under Creative Commons: `together by Adrien Coquet from the Noun Project`
