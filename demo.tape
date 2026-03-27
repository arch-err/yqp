# yqp demo — record with: just generate-gif
# Uses a fake kubectl wrapper to simulate a real K8s workflow.

Output demo.gif

Set Shell "bash"
Set FontSize 14
Set Width 1200
Set Height 700
Set Theme "Catppuccin Mocha"
Set Padding 10
Set TypingSpeed 40ms

# Simulate piping kubectl output into yqp
Type "./demo/kubectl | ./yqp"
Sleep 300ms
Enter
Sleep 1.5s

# Query 1: Get all container images
Type ".spec.template.spec.containers[].image"
Sleep 300ms
Enter
Sleep 1.5s

# Query 2: Find containers missing memory limits
Ctrl+A
Type ".spec.template.spec.containers[] | select(.resources.limits.memory == null) | .name"
Sleep 300ms
Enter
Sleep 1.5s

# Query 3: Show labels
Ctrl+A
Type ".metadata.labels"
Sleep 300ms
Enter
Sleep 1.5s

# Toggle input panel off for clean output view
Ctrl+T
Sleep 1s

# Toggle back
Ctrl+T
Sleep 1s

# Quit
Ctrl+C
