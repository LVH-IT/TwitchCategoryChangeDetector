# Twitch Category Change Detector
This package mostly uses native go packages, except for these external dependencies:
* [beep](https://github.com/faiface/beep), which is licenced under an [MIT-like](https://github.com/faiface/beep/blob/master/LICENSE) license. This package is used for playing audio files.  
* [DiscordGo](https://github.com/bwmarrin/discordgo), which is licenced under a [BSD 3-Clause](https://github.com/bwmarrin/discordgo/blob/master/LICENSE) license. This package is used for communicating with discord.

## What does it do and why?
This program simply checks the current category of a twitch streamer and notifies you as soon as it changes, so you don't have to manually check until there is a category you like to watch. It is also way more resource-friendly than having a browser open in the background to check categories.


## Usage
### Supported Operating Systems
Linux and Windows only  

### How to setup
There are no prerequisites when using the precompiled binary file, which you can download from the releases section.  
Though if you want you can also compile the code yourself as explained below.  

The executable needs an audio file in its directory called "juntos.ogg" to notify the user of a category change, but you can also change the audio filename in the config.json file.  
Your audio file has to use the Vorbis codec, however (usually .ogg files) or the program will not work.  
I got my sound from here: https://notificationsounds.com/message-tones/juntos-607  

In the config file, you have to set a client id and a client secret to use the twitch API.  
To do so you have to create a twitch application: https://dev.twitch.tv/console/apps  
You can call the application whatever you want. For "OAuth Redirect URLs" just use "http://localhost" or something, as this value is irrelevant. For the Category, you should choose "Application Integration"  
You will then get a client id and a client secret for the application you created. You need to put both of these into your config.json file.

You are now done and you can use the quickstart script to start monitoring a twitch channel. When it asks you whether to get a new Token, just type "y" to get one. In case it fails to get a Token, go to the "Solving errors" section down below.

### Discord bot
Under development

### Flags you can use
* **-s** :   accepts the name of the twitch channel (ex: -s xqcow), default is xqcow  
* **-t** :   accepts the interval in which to recheck the stream category in seconds (ex: -t 10), default is 10  
* **-dcbot** :   starts the app as a discord bot. Notifications are sent through discord and will not appear in the command line  

### Config changes you can make
* **BearerToken** : You can either manually get one, as described [here](#automatic-bearer-token-retrieval-fails), or let the application do it for you  
* **ClientID** : You have to fill this in yourself for now, as described here: [How to setup](#how-to-setup)
* **ClientSecret** : You have to also fill this in yourself for now, as described here: [How to setup](#how-to-setup)
* **SoundFile** : This has to be the filename of your audio file. Default: "juntos.ogg"  
* **UseCategoryWhitelist** : If you only want to be notified when a streamer changes to a specific category of your choice, set this to true and also change the next config entry  
* **Categories** : Here you can list the categories you want to whitelist. Default: ["Watch Parties","Just Chatting"]  
* **NotifyOnOfflineTitleChange** : Notifies you when the stream title changes while the stream is offline  
* **NotifyOnOnlineTitleChange** : Notifies you when the stream title changes while the stream is online  
* **DiscordBotToken** : A bot token for Discord is needed in order to use this app as a discord bot  

## Compilation
### Prerequisites for compilation
Go >= 1.15 (https://golang.org/dl/)  
You'll get the rest when trying to compile  


### How to compile
Open a command prompt in the source directory and you should be able to install all dependencies by executing this command inside the source folder: 
```sh
go get -d ./...
```
Before compiling on Linux you need a few dependencies (libasound2-dev and pkg-config), which you can install like this: 
```sh
apt install libasound2-dev pkg-config
```
Then simply type:
```sh
go build
```
It will then try to compile and tell you whether there are dependencies that are still missing.
If so, you need to install them each like this: 
```sh
go get github.com/faiface/beep
```
Then rerun "go build" and your executable should be compiled in the source directory.

## Solving errors
### Automatic Bearer Token retrieval fails
Error Message: "Could not obtain a new Bearer Access Token. Please try again or get one manually"  
If the application somehow fails to automatically get a bearer Token, you can manually get one and put it into the config.json file.  
To do that simply open a command prompt and type:  
```sh
curl -X POST "https://id.twitch.tv/oauth2/token?client_id=PUTYOURCLIENTIDHERE&client_secret=PUTYOURCLIENTSECRETHERE&grant_type=client_credentials"  
```
Just replace "PUTYOURCLIENTIDHERE" and "PUTYOURCLIENTSECRETHERE" with the actual information you got before.  
You will then find your bearer token in the response and you can put it into the config.json file.  
Just keep in mind that your bearer token will expire after the amount of time given in the curl response and you will have to get a new one, either automatically by letting the application do it or by doing it manually again in case it fails.  