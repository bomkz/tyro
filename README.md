# TYRO
 Telemetry Yield Real-time Observations
## Installation
(NOTE: binaries not available until permission is gotten from image authors, until then, you can build your own, except you'll have to use your own images)
Download precompiled binaries off the release tab in Github, or if you're ✨special✨, spend the next hour and a half trying to figure out why it is not compiling on your side. (looking at you, future me)
### Warning, you /might/ receive a false positive detection by Windows Defender or whatever antivirus you use due to this being an uncommonly downloaded program. I can promise you this is not a virus, but realistically, you should always follow your gut instinct, don't download something and run it because you were promised its not a virus, if you feel uneasy, verify it yourself by using the vast array of utilities online to scan files and run them in VMs to scan their behaviour.

## Compile 
Requirements:
- Golang
- Hopes and dreams

`git clone https://github.com/AngelFluffyOokami/jamcat-mach.git`

`cd .\tyro`

`go mod tidy`

Open discordrp.go and replace APPLICATION_ID with one of your own from https://discord.com/developers, or if not, use the mine:

1220960048704913448

if using your own APPLICATION_ID, go to https://discord.com/developers and create a new application, name it, and then go to rich presence art asset section, there you will upload a few images, they have to strictly have the following names given to them after upload, but before saving:

vtolvr
f45a
fa26b
av42c
ef24g
t55
ah94

`go build`

