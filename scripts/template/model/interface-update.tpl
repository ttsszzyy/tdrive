FindOneByBuilder(ctx context.Context, sbs ...squirrel.SelectBuilder)(*{{.upperStartCamelObject}}, error)
Update(ctx context.Context, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}, s ...sqlx.Session) error
ListPage(ctx context.Context, page, size int64, sbs ...squirrel.SelectBuilder)([]*{{.upperStartCamelObject}}, int64, error)
List(ctx context.Context, sbs ...squirrel.SelectBuilder)([]*{{.upperStartCamelObject}}, error)
Count(ctx context.Context, cbs ...squirrel.SelectBuilder)(int64, error)