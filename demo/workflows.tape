Set FontSize 20
Set Width 1000
Set Height 700
Set TypingSpeed 0.1s
Output ./outputs/gifs/workflows.gif
Output ./validation/workflow_final_result.ascii
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
Type@100ms "Nothing in the stdout buffer yet"
Sleep 1s
Enter
Type@100ms "Let's run our workflow!"
Sleep 3s
Escape
Type ":qa!"
Sleep 0.5s
Enter
Sleep 1.5s
Tab
Screenshot ./outputs/screenshots/workflow_view.png
Sleep 3.5s
Enter
Type "v"
Sleep 2s
Type "jo"
Enter 2
Type@100ms "Nice!"
Sleep 1s
Type@100ms " Capitalized, duplicated and trimmed!"
Sleep 1s
Enter 2
Type@100ms "And now we have the workflow saved as a file"
Sleep 1s
Enter 2
Type@100ms "This can be used like a script or added to a chain"
Sleep 5s
