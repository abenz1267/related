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
            "extension": "tsx"
        },
        {
            "name": "cssmodule",
            "path": "./styles/",
            "extension": "module.css"
        }
    ],
    "groups": [{ "name": "component", "types": ["component", "cssmodule"] }]
}
```

If no template is provided, the file will be empty.

### Commands

| Command                             | Function                                                                         |
| ----------------------------------- | -------------------------------------------------------------------------------- |
| `listtemplates <parent>`            | Lists all available templates, grouped by parent-folder. The parent is optional. |
| `<type or group> <name> <filename>` | Creates the file(s) based on the type or group provided                          |

### Templates

Templates are embedded. You can create your own templates by placing them into your config folder. F.e. on Linux `~/config/related/templates/<parent>/<name>.tmpl`.

You can overwrite the default templates by simply placing a copy in your config folder. Related will always prioritize custom templates over default ones.
