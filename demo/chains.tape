Set FontSize 20
Set Width 1000
Set Height 700
Output chains.gif
Set Shell zsh
Env VISUAL "vim"

Sleep 1s
Type "launcher"
Sleep 1s
Enter
Sleep 2s
Type "v"
Sleep 2s
Type "i"
Type@100ms "In the other demo,"
Sleep 0.5s
Type@100ms " we ran the scripts one at a time"
Sleep 1s
Enter
Type@100ms "Now lets chain them together first"
Sleep 1s
Type@100ms " and then run them all"
Sleep 3s
Escape
Type ":qa!"
Sleep 0.5s
Enter
Sleep 2s
Type "a"
Sleep 0.5s
Type@150ms "echo hello world!"
Sleep 1s
Enter
Sleep 2s
Type "j"
Sleep 0.5s
Type "a"
Sleep 1s
Type "j"
Sleep 0.5s
Type "a"
Sleep 1s
Type "j"
Sleep 0.5s
Type "a"
Sleep 2s
Screenshot screenshots/chain_view.png
Type "R"
Sleep 1s
Type "v"
Sleep 2s
Type "jo"
Enter 2
Type@100ms "Nice!"
Sleep 1.5s
Type@100ms " Capitalized, duplicated and trimmed!"
Sleep 1s
Enter 2
Type@100ms "Just like last time!"
Sleep 1s
Enter
Type@100ms "But now we have the chain in memory"
Sleep 1s
Enter 2
Type@100ms "Lets save this chain as a workflow to be used later"
Sleep 3s
Escape
Type ":wq"
Sleep 0.5s
Enter
Sleep 2s
Type "W"
Sleep 2s
Type@100ms "Echo, capitalize, duplicate and trim"
Sleep 2s
Enter
Sleep 2s
Tab
Sleep 3s
Type "v"
Sleep 2s
Type@200ms "VGxi"
Sleep 1s
Type@100ms "Now this workflow can be used like a script!"
Sleep 1s
Enter
Type@100ms "Run it with the launcher or chain further with its base scripts"
Sleep 5s
