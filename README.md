# Chatbot Plugins

## Auth

* Contains plugins to limit access who can send commands to the bot

## Chatlog

* Contains plugins that write errors to a file, they can be stored in clear text or encrypted

## Cryptocurrency

* Contains a plugin that interacts with the coin market cap api

## Help

* Formats the help output when the command /help is called

## HIBP (Have I Been Pwnd)

* Contains a way to interact with the account and password api
* The email plugin will list all the breaches and pastes that email has been seen in
* Password will return the number of times a specific password has been seen

## Media

* Contains plugins that allow for the bot to post images and gisf.
* The giphy plugin will pick and random gif based on user input and post it to chat
* The media plugin will take a url and attempt to upload the img or gif to chat

## Reddit

* Contains a plugin that queries a subreddit and returns the top 10 posts

## Screenshot

* Contains a plugin that will take a screenshot of a website and post it to the chat

## URL Shortener

* Contains a plugin that will return a shortened URL version of the requested URL

## Shodan

* Contains a plugin that queries the Shodan api and returns the organization, ASN, hostnames, and ports from a given IP.

## VirusTotal

* Contains a plugin that queries the VirusTotal api and returns the scan detection results if any exist

## Weather

* Contains a plugin that communicates with the open weather api and returns the weather of a given city.