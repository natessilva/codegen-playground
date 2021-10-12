package main

import "github.com/gobuffalo/plush"

func render(template string, d Definition, name string) (string, error) {
	ctx := plush.NewContext()
	ctx.Set("definition", d)
	ctx.Set("name", name)
	ctx.Set("camelize_down", camelizeDown)
	ctx.Set("json_type", jsonType)
	return plush.Render(template, ctx)
}
