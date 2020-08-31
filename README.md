# TwitchCategoryChangeDetector

To compile you need to install GO from https://golang.org/dl/  
Then open a command prompt in the source directory and type "go build main.go"  
It will then tell you which dependencies are missing.  
You need to install them each like this as an example: "go get github.com/faiface/beep"  
Then you can rerun "go build main.go" and your executable should be ready to use  

Flags you can use with the executable:  
-s  
  accepts the name of the twitch channel (ex: -s xqcow) default is xqcow  
-t  
  accepts the interval in which to recheck the stream category in seconds (ex: -t 10) default is 10  
