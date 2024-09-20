# Purpose
PB simply stands for "Project Browser" and is intended specifically to be a tool to aide in finding projects on your filesystem.

# Flexbility
PB is not just for simply browsing projects in a single directory.  It is capable of managing projects from various seperate folders on your filesystem and can even traverse subdirectories if configured to.  PB is also a useful tool to combine with something like tmux sessions which can allow you to use PB to open several different projects in their own tmux session.

## Config
PB looks for config settings in a file called settings.json located in the pb subdirectory of your OS's config folder.  On linux this is usually ~/.config/pb/settings.json

### Options
`Source.Path`: The path to start searching for directories in. For example `/home/user/repos`

`Source.TraversalDepth`: Can be any number. `0` matches directoriees by `/home/user/repos/*` while `1` matches directories by `/home/user/repos/*/*`

`ProjectOpenCommand`: Set the command that will be used to open projects. `$projectName` and `$projectPath` are available as variables

`DefaultOpenDepth`: Can be any number. `0` will cause the first selection to open a project.  `1` while allow users to select a folder than select a project to open

`DisplayAbsolutePath`: Determines wether the absolute path will be displayed alongside the project name

### Example

```json

{
    "Sources": [
        {
            "Path": "/home/user/repos",
            "TraversalDepth": 1
        },
        {
            "Path": "/home/user/school",
            "TraversalDepth": 0
        }
    ],
    "ProjectOpenCommand": "tmux new-session -s $projectName -A -c $projectPath 'nvim .'",
    "DefaultOpenDepth": 0,
    "DisplayAbsolutePath": true
}
```
