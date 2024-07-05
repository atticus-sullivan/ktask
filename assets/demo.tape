# Settings:
#   Set FontSize <number>           Set the font size of the terminal
#   Set FontFamily <string>         Set the font family of the terminal
#   Set LetterSpacing <float>       Set the font letter spacing (tracking)
#   Set LineHeight <float>          Set the font line height
#   Set LoopOffset <float>%         Set the starting frame offset for the GIF loop
#   Set Theme <json|string>         Set the theme of the terminal
#   Set Padding <number>            Set the padding of the terminal
#   Set Framerate <number>          Set the framerate of the recording
#   Set PlaybackSpeed <float>       Set the playback speed of the recording
#   Set MarginFill <file|#000000>   Set the file or color the margin will be filled with.
#   Set Margin <number>             Set the size of the margin. Has no effect if MarginFill isn't set.
#   Set BorderRadius <number>       Set terminal border radius, in pixels.
#   Set WindowBar <string>          Set window bar type. (one of: Rings, RingsRight, Colorful, ColorfulRight)
#   Set WindowBarSize <number>      Set window bar size, in pixels. Default is 40.
#   Set TypingSpeed <time>          Set the typing speed of the terminal. Default is 50ms.
#
# Sleep:
#   Sleep <time>                    Sleep for a set amount of <time> in seconds
#
# Type:
#   Type[@<time>] "<characters>"    Type <characters> into the terminal with a
#                                   <time> delay between each character
#
# Keys:
#   Escape[@<time>] [number]        Press the Escape key
#   Backspace[@<time>] [number]     Press the Backspace key
#   Delete[@<time>] [number]        Press the Delete key
#   Insert[@<time>] [number]        Press the Insert key
#   Down[@<time>] [number]          Press the Down key
#   Enter[@<time>] [number]         Press the Enter key
#   Space[@<time>] [number]         Press the Space key
#   Tab[@<time>] [number]           Press the Tab key
#   Left[@<time>] [number]          Press the Left Arrow key
#   Right[@<time>] [number]         Press the Right Arrow key
#   Up[@<time>] [number]            Press the Up Arrow key
#   Down[@<time>] [number]          Press the Down Arrow key
#   PageUp[@<time>] [number]        Press the Page Up key
#   PageDown[@<time>] [number]      Press the Page Down key
#   Ctrl+<key>                      Press the Control key + <key> (e.g. Ctrl+C)
#
# Display:
#   Hide                            Hide the subsequent commands from the output
#   Show                            Show the subsequent commands in the output

Output demo.gif

Require ktask

Set FontSize 22
Set Height 1080
Set Width 1920
Set TypingSpeed 0.1s
Set PlaybackSpeed 0.25 # Make output 2 times slower
Set Theme "Catppuccin Frappe"
Set Shell zsh

Hide
Type "cp demo.ktask.bak demo.ktask"
Enter
Show
Sleep 500ms
Type "ktask kanban demo.ktask"
Sleep 1s
Enter
Sleep 500ms
Type "ll"
Sleep 500ms
Type "hh"
Sleep 2s
Type "e"
Sleep 2s
Type "2x"
Sleep 500ms
Enter 2
Sleep 500ms
Type "j"
Sleep 1.5s
Type "lhkj"
Sleep 500ms
Enter
Type "l"
Sleep 500ms
Type "l"
Sleep 1s
Type "h"
Sleep 500ms
Type "a"
Sleep 500ms
Type "testing"
Enter
Type "test"
Sleep 500ms
Enter
Type "jk"
Sleep 1s
Type "x"
Sleep 500ms
Type "q"
Sleep 1s
Ctrl+D
