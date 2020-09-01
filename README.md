# Twitch Category Change Detector

## What does it do?
This program simply checks the current category of a twitch streamer and notifies you as soon as it changes.  


## Usage
Download the latest release and contine with "How to setup". There are no prerequisites when using the already compiled binary file.  
But if you want you can also compile the code yourself as explained below.  


### How to setup
The executable needs an audio file in the same directory called "juntos.ogg" in order to notify the user of a category change, but you can also change the audio filename in the config.json file.  
Your audio file has to use the vorbis codec however (usually .ogg files) or the program will not work.  
I got my sound from here: https://notificationsounds.com/message-tones/juntos-607  

In the config file you have to set a bearer token and a client id to use the twitch api.  
To do so you have to create a twtich application: https://dev.twitch.tv/console/apps  
You will then get a client id (which you need to put into the config.json file) and a client secret for the application you created.  
Then you can get a bearer token that matches your client id by using curl for example.  
To do that simply open a command prompt and type:  

curl -X POST "https://id.twitch.tv/oauth2/token?client_id=PUTYOURCLIENTIDHERE&client_secret=PUTYOURCLIENTSECRETHERE&grant_type=client_credentials"  

Just replace "PUTYOURCLIENTIDHERE" and "PUTYOURCLIENTSECRETHERE" whith the actual information you got before.  
You will then find your bearer token in the response and you can put it into the config.json file.  
Just keep in mind that your bearer token will expire after the amount of time given in the curl response and you will have to get a now one.  

You are now done and you can use quickstart.bat to start monitoring a twitch channel.


### Flags you can use
-s  
  accepts the name of the twitch channel (ex: -s xqcow) default is xqcow  
  
-t  
  accepts the interval in which to recheck the stream category in seconds (ex: -t 10) default is 10  

## Compilation
### Prerequisites for compilation
GO  (https://golang.org/dl/)  
You get the rest when trying to compile  


### How to compile
Open a command prompt in the source directory and simply type "go build"  
It will then try to compile and tell you which dependencies are missing.  
You need to install them each like this: "go get github.com/faiface/beep"  
And rerun "go build" until there are no more errors and your executable should be compiled in the source directory.  
