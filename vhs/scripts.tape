Set FontSize 20
Set Width 1000
Set Height 700
Set TypingSpeed 0.1s
Output ./outputs/gifs/scripts.gif
Output ./validation/script_final_result.ascii
Set Shell zsh
Env VISUAL "vim"

Sleep 1s
Type "launcher"
Sleep 1s
Enter
Sleep 1s
Screenshot ./outputs/screenshots/script_view.png
Sleep 1s
Enter
Sleep 0.5s
Type@150ms "echo hello world!"
Sleep 1s
Enter
Type "v"
Sleep 1s
Type@300ms "VVVVVV"
Sleep 1s
Type "o"
Enter 2
Type@100ms "In our launcher buffer"
Sleep 1s
Enter
Type@100ms "We can view and edit the latest stdout"
Enter 2
Sleep 1s
Type@100ms "Let's do something with it!"
Sleep 3s
Escape
Type ":qa!"
Sleep 0.5s
Enter
Sleep 3s
Type "j"
Enter
Sleep 2s
Type "v"
Sleep 1s
Type@300ms "VVVVVV"
Sleep 1s
Type "o"
Enter 2
Type@100ms "Ah look!"
Sleep 1s
Enter 
Type@100ms "The output from `echo` has piped into the capitilization script"
Sleep 1s
Enter 2
Type@100ms "Let's run this through the other two scripts"
Sleep 2s
Escape
Type ":qa!"
Sleep 0.5s
Enter
Sleep 2s
Type@300ms "j"
Enter
Sleep 1s
Type@300ms "j"
Enter
Sleep 1s
Type "v"
Sleep 2s
Type "jo"
Enter 2
Type@100ms "Nice!"
Sleep 1.5s
Type@100ms " Capitalized, duplicated and trimmed!"
Sleep 2s
Enter 2
Type@100ms "By chaining scripts together"
Sleep 0.5s
Enter
Type@100ms "You can start building complex workflows"
Escape
Sleep 5s
