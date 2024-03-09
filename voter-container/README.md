# LAB: Cloud Native Packaging Assignment

First and foremost, I apologize for the late submission. I ran into a couple of errors that took me a few days to figure out. I had somehow fat fingered an "find and replace" to the entire repo and pushed up all the broken code... took me 3 days to figure it out and had to back track once I realized I had made that mistake.

Any who without further hiccups... To run my container setups.

1. From this current directory, **voter-container** or where this **README.md** you are currently reading exists, run `make setup`. This will pull from my personal docker repository on the web [here](docker.io/jhyang/cs-t681-voter-api). Please note that this image is probably built from arm64 base architecture. In the event that it doesn't work you can run `make build-container` first to build the image, then `make init-redis` to set things up.
2. The app can be accessed from `localhost:8080` and I left the redis ui avaliable from own preference `localhost:8001` otherwise everything else is the same.
