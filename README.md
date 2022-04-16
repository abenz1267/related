# Related - Create files based on individual definitions or groups

## Installation

`go install github.com/abenz1267/related@latest`

## Usage

Place a ".related.json" file in your project folder. Content example:

```json
{
    "types": [
        {
            "name": "component",
            "template": "typescript/NextFuncComponent",
            "path": "./",
            "pre": "typescript/MyCustomPreScript",
            "post": "typescript/MyCustomPostScript",
            "extension": "tsx"
        },
        {
            "name": "cssmodule",
            "path": "./styles/",
            "extension": "module.css"
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

### Templates

Templates are embedded. You can create your own templates by placing them into your config folder. F.e. on Linux `~/config/related/templates/<parent>/<name>.tmpl`.

You can overwrite the default templates by simply placing a copy in your config folder. Related will always prioritize custom templates over default ones.

### Scripts

Scripts must be placed in your config folder. F.e. on Linux `~/config/related/scripts/<parent>/<name>.lua`.

You can execute \*.lua files by settings pre- and post-scripts in the type or group definition. Related will look for the \*.lua script and execute itaccording to the lifecycle. A global variable named `Name` will hold the name you provided to the initial command, f.e. `related group component MyComponent` will add `"MyComponent"` as a global variable to the script.