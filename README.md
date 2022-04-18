[![Go Report Card](https://goreportcard.com/badge/github.com/abenz1267/related)](https://goreportcard.com/report/github.com/abenz1267/related)

# Related - Create files based on individual definitions or groups

Related helps with common file-creation-based tasks. You can predefine single types as well as groups. Scripts (lua, binary, javascript [needs node]) can be executed pre/post file creation as well as pre/post groups.

## Installation

`go install github.com/abenz1267/related@latest`

## Usage

Place a "(.)related.json" file in your project folder. Content example:

```json
{
    "types": [
        {
            "name": "component",
            "template": "typescript/NextFuncComponent.tmpl",
            "path": "./",
            "pre": "typescript/MyCustomJS.lua",
            "post": "typescript/MyCustomPostScript.lua",
            "suffix": ".tsx"
        },
        {
            "name": "cssmodule",
            "path": "./styles/",
            "suffix": ".module.css"
        }
    ],
    "groups": [
        {
            "name": "component",
            "types": ["component", "cssmodule"],
            "pre": "typescript/MyCustomGroupPreScript",
            "post": "typescript/MyCustomGroupPostScript"
        }
    ]
}
```

If no template is provided, the file will be empty.

### Commands

| Command                                | Function                                                                                    |
| -------------------------------------- | ------------------------------------------------------------------------------------------- |
| `list <scripts or templates> <parent>` | Lists all available templates or scripts, grouped by parent-folder. The parent is optional. |
| `<type or group> <name> <filename>`    | Creates the file(s) based on the type or group provided                                     |

### Custom Files

Custom files like templates or scripts can be either placed in a `.related` folder near the config or in a `related` folder in your users config directory. Templates must be placed inside `templates` and scripts inside `scripts`. You can nest further.

Related will prioritize project-level files over config ones.

### Templates

Templates are vanilla Golang templates, only data passed is the filename you provided. It can be accessed via `{{.}}` in the template.

### Scripts

The following types are executable: lua scripts, javascript (via node), and binaries.

You can execute scripts by settings pre- and post-scripts in the type or group definition. Related will look for the script and execute it according to the lifecycle.

Passed command-line arguments:

1. current working directory
2. path
3. filename
4. extension
