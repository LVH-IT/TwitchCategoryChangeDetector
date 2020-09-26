# Twitch Category Change Detector
This package mostly uses native go packages, except for these external dependencies:
* [beep](https://github.com/faiface/beep), which is licenced under an [MIT-like](https://github.com/faiface/beep/blob/master/LICENSE) licence. This package is used for playing audio files.

## What does it do and why?
This program simply checks the current category of a twitch streamer and notifies you as soon as it changes, so you don't have to manually check until there is a category you like to watch. It is also way more resource friendly than having a browser open in the background to check categories.


## Usage
Download the latest release and contine with "How to setup". There are no prerequisites when using the already compiled binary file.  
But if you want you can also compile the code yourself as explained below.  


### How to setup
The executable needs an audio file in it's directory called "juntos.ogg" in order to notify the user of a category change, but you can also change the audio filename in the config.json file.  
Your audio file has to use the vorbis codec however (usually .ogg files) or the program will not work.  
I got my sound from here: https://notificationsounds.com/message-tones/juntos-607  

In the config file you have to set a client id and a client secret to use the twitch api.  
To do so you have to create a twitch application: https://dev.twitch.tv/console/apps  
You can call the application whatever you want. For "OAuth Redirect URLs" just use "http://localhost" or something, as this value is irrelevant. For the Category you should choose "Application Integration"  
You will then get a client id and a client secret for the application you created. You need to put both of these into your config.json file.

You are now done and you can use quickstart.bat to start monitoring a twitch channel. When it asks you wether to get a new Token, just type "y" to get one. In case it fails to get a Token go to the "Solving errors" section down below.

### Flags you can use
* **-s** :   accepts the name of the twitch channel (ex: -s xqcow), default is xqcow  
* **-t** :   accepts the interval in which to recheck the stream category in seconds (ex: -t 10), default is 10  

### Config changes you can make
* **BearerToken** : You can either manually get one, as described [here](#automatic-bearer-token-retrieval-fails), or let the application do it for you  
* **ClientID** : You have to fill this in yourself for now, as described here: [How to setup](#how-to-setup)
* **ClientSecret** : You have to also fill this in yourself for now, as described here: [How to setup](#how-to-setup)
* **SoundFile** : This has to be the filename of your audiofile. Default: "juntos.ogg"  
* **UseCategoryWhitelist** : If you only want to be notified when a streamer changes to a specific category of your choice, set this to true and also change the next config entry  
* **Categories** : Here you can list the categories you want to whitelist. Default: ["Watch Parties","Just Chatting"]  

## Compilation
### Prerequisites for compilation
Go 1.15 (https://golang.org/dl/)  
You'll get the rest when trying to compile  


### How to compile
Open a command prompt in the source directory and you should be able to install all dependencies by executing this command inside the source folder: 
```sh
go get -d ./...
```
Then simply type:
```sh
go build
```
It will then try to compile and tell you wether there are dependencies which are still missing.
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
Just replace "PUTYOURCLIENTIDHERE" and "PUTYOURCLIENTSECRETHERE" whith the actual information you got before.  
You will then find your bearer token in the response and you can put it into the config.json file.  
Just keep in mind that your bearer token will expire after the amount of time given in the curl response and you will have to get a new one, either automatically by letting the application do it or by doing it manually again in case it fails.  

## Roadmap
### Potential Features to come:
* Manage multiple streams within one CLI
* GUI
  * All the Features in a GUI
