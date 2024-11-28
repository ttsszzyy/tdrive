
func (m *default{{.upperStartCamelObject}}Model) Insert(ctx context.Context, data *{{.upperStartCamelObject}}, s ...sqlx.Session) (sql.Result,error) {
	{{if .withCache}}{{.keys}}
    return m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
		if s != nil{
        	return s[0].ExecCtx(ctx,query,{{.expressionValues}})
        }
		return conn.ExecCtx(ctx, query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("insert into %s (%s) values ({{.expression}})", m.table, {{.lowerStartCamelObject}}RowsExpectAutoSet)
        if s != nil {
            return s[0].ExecCtx(ctx,query,{{.expressionValues}})
        }
    return m.conn.ExecCtx(ctx, query, {{.expressionValues}}){{end}}
}
