# GOT - golang template engine

Got is simple server side template engine, it supports rendering
- variables
- if statments
- foreach loops

## Instalation

```bash
go get github.com/Pla9er9/template_engine
```

## Run tests
```bash
go test .
```

## Usage

### Variables
Tag Pattern - `{ variablename }`
You can create struct variable and then use exported properties of struct in template
<br>
```go
variables := map[string]any{
    "user" : User{
        Name: "Alex",
    },
}
```

Example:
```html
<h1 style="color: {color}">{user.Name}</h1>
```

### If statment
Starting tag `{@if variableName}` <br>
Ending tag `{/if}`

⚠️ In if statment you can only pass boolean variables, you cannot write conditions like `{@if age > 18}` in templates, you should define conditions in `.go` file as bool variables <br>

```go
variables := map[string]any{
    "isAdult" : (age >= 18),
}
```

❌ Wrong usage - `{@if age >= 18}` <br>
✅ Correct usage - `{@if isAdult}` <br>

Example:
```html
{@if isAuthorized}
    <h1>You are authorized</h1>
{/if}
```

### Foreach loop
Starting tag `{@foreach array as item}` <br>
Ending tag `{/foreach}`
<br>

Example:
```html
{@foreach array as item}
    <p>{ item }</p>
{/foreach}
```

## Example
`main.go`
```Go
package main

import "github.com/Pla9er9/template_engine"

func main() {
    engine := templateEngine.GetTemplateEngine()
    variables := map[string]any{
        "color" : "red",
        "adult" : true,
        "array" : []string{"Item1", "Item2", "Item3"},
    }

    html, err := engine.RenderTemplateFromFile("templ.got", variables)

	if err != nil {
		panic(err)
	}
}
```

`templ.got`
```html
<!DOCTYPE html>
<html lang="en">
    <head>
        <meta charset="UTF-8" />
    </head>
    <body>
        <h1 style="color: { color }">{ color }</h1>
        {@if adult}
            <h1>You are adult</h1>
        {/if}
        {@foreach array as arr}
            <h2>{arr}</h2>
        {/foreach}
    </body>
</html>
```

## [Mit license](https://github.com/Pla9er9/template_engine/blob/main/LICENSE)