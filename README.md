# Twitch Category Change Detector
This program simply checks the current category of a twitch streamer and notifies you as soon as it changes.  

To compile you need to install GO from https://golang.org/dl/  
Then open a command prompt in the source directory and type "go build main.go"  
It will then tell you which dependencies are missing.  
You need to install them each like this as an example: "go get github.com/faiface/beep"  
Then you can rerun "go build main.go" and your executable should be compiled without errors.  
At last the executable needs an audio file in the same directory called "juntos.ogg" in order to notify the user of a category change.  
I got mine from here: https://notificationsounds.com/message-tones/juntos-607  


Flags you can use with the executable:  
-s  
  accepts the name of the twitch channel (ex: -s xqcow) default is xqcow  
-t  
  accepts the interval in which to recheck the stream category in seconds (ex: -t 10) default is 10  
