# Purpose
PB simply stands for "Project Browser" and is intended specifically to be a tool to aide in finding projects on your filesystem.

# Flexbility
PB is not just for simply browsing projects in a single directory.  It is capable of managing projects from various seperate folders on your filesystem and can even traverse subdirectories if configured to.  PB is also a useful tool to combine with something like tmux sessions which can allow you to use PB to open several different projects in their own tmux session.

## Config
PB looks for config settings in a file called settings.json located in the pb subdirectory of your OS's config folder.  On linux this is usually ~/.config/pb/settings.json
<br><br>
The following options are available:

```json
{
    "comment": "Where to find project folders",
    "Sources": [
        "/home/user/repos",
        "/home/user/repos/php/wordpress"
    ],

    "comment": "How many subdirectory levels to traverse before collecting the project folders.  Value of 0 fetches directories of the pattern /home/user/repos/* and 1 the pattern /home/user/repos/*/*",
    "SourceTraversalDepth": 1,

    "comment": "The command to execute when a project is opened.  Two variables are provided in the form of $projectName and $projectPath which contain the values described",
    "ProjectOpenCommand": "tmux new-session -s $projectName -A -c $projectPath nvim .",

    "comment1": "How many subdirectories should the user be allowed to traverse before executing the project open command",
    "comment2": "Value of 0 will cause the first item selected to open as a project. Value of 1 will traversal one directory down before opening as project.  Items can be opened as directories regardless of depth by using the hotkey",
    "DefaultOpenDepth": 0,

    "comment": "Whether or not the absolute path of the project should be displayed alongside its name"
    "DisplayAbsolutePath": true
}
```
