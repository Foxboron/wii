wii
===

Webfrontend for [iii]{https://github.com/Foxboron/iii}.  

Allows the use of a simple webfront-end for viewing logs or adding webhooks

```
Usage of ./wii:
  -a string
    	Authentication. Example: user:pass (default none)
  -i string
    	IRC directory
  -p string
    	Port of the server (default "8003")
```


Example: 
```
curl -u user:pass --data "msg=Hello World" "localhost:8003/irr.hackint.org/channel/buffer"
curl -u user:pass "localhost:8003/irc.hackint.org/channel/buffer"
```
