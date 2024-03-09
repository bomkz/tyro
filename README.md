# jamcat-mach
 Jet Audio and Music Control Access Terminal-Media Access and Control Hub (WIndows Media Controls for VTOL VR without modloader)

## Installation
Download precompiled binaries off the release tab in Github, or if you're ✨special✨, spend the next hour and a half trying to figure out why it is not compiling on your side. (looking at you, future me)
### Warning, you /might/ receive a false positive detection by Windows Defender or whatever antivirus you use. This is due to the way I pause/play media, by sending virtual keystrokes, which CAN be seen as malicious by AV programs. I've submitted a false positive report to MSFT, and hopefully they accept it.

## Compile 
Requirements:
- Golang
- Mingw
- CGO enabled (check your goenv)
- Hopes and dreams

`git clone https://github.com/AngelFluffyOokami/jamcat-mach.git`

`cd .\jamcat-mach`

`go mod tidy`

`go build`


## YOU
YES YOU.
You ever wished you could control other music players from within VTOL VR? I'm sure you have, just to be let down when you found out you had to use a modloader.
#### Well you're in luck!
JAMCAT-MACH let's you do just that!*
(*does not support pausing playback)

JAMCAT-MACH places three audio files in VTOL VR's RadioMusic folder, then reads logs in realtime to check when VTOL VR loads in an MP3 file. 
It then uses some quick maffs with the file names to determine whether you pressed play, skip, or rewind, then sends a keyboard input to windows.

JAMCAT-MACH also does some fancy logic to figure out where steam is installed, and then from there, where VTOL VR is installed as well. So you could have steam installed in a secondary drive, then VTOL VR installed in a third drive, and JAMCAT-MACH /should/ be able to figure out where the fuck everything is. But no promises. Just file a Github issue, or shoot me a DM at discord (@f45a) if you encounter an issue with this, I'll try to fix it, but no promises, and much less any implied warranties. 

Just a caution, please backup anything in the RadioMusic folder if you put your own music there. 

#### Don't be stupid.
#### Don't close program with TaskMan unless necessary.
#### Don't mess with log files while program is running.
#### Don't add, modify, or remove files in RadioMusic folder while program is running.

(I did add a way to handle these situations except the first one, however, I cannot guarantee it will work 100 percent of the time, so just d o n ' t please)
