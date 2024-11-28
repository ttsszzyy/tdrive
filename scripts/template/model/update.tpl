
func (m *default{{.upperStartCamelObject}}Model) Update(ctx context.Context, {{if .containsIndexCache}}newData{{else}}data{{end}} *{{.upperStartCamelObject}}, s ...sqlx.Session) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne(ctx, newData.{{.upperStartCamelPrimaryKey}}, s...)
	if err!=nil{
		return err
	}

{{end}}	{{.keys}}
    _, {{if .containsIndexCache}}err{{else}}err:{{end}}= m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
		if s != nil{
        	return s[0].ExecCtx(ctx,query, {{.expressionValues}})
        }
		return conn.ExecCtx(ctx, query, {{.expressionValues}})
	}, {{.keyValues}}){{else}}query := fmt.Sprintf("update %s set %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, {{.lowerStartCamelObject}}RowsWithPlaceHolder)
    if s != nil{
    	 _,err:=s[0].ExecCtx(ctx,query, {{.expressionValues}})
    }
    _,err:=m.conn.ExecCtx(ctx, query, {{.expressionValues}}){{end}}
	return err
}

func (m *default{{.upperStartCamelObject}}Model) FindOneByBuilder(ctx context.Context, sbs ...squirrel.SelectBuilder)(*{{.upperStartCamelObject}}, error) {
    var sb squirrel.SelectBuilder
    if sbs == nil {
        sb = squirrel.Select()
    } else {
        sb = sbs[0]
    }

    // query rows
    query, args, err := sb.From(m.table).Columns({{.lowerStartCamelObject}}Rows).Limit(1).ToSql()
    if err != nil {
        return nil, err
    }
    var ret {{.upperStartCamelObject}}
    if err := m.QueryRowNoCacheCtx(ctx, &ret, query, args...); err != nil {
        return nil, err
    }

    return &ret, nil
}

func (m *default{{.upperStartCamelObject}}Model) ListPage(ctx context.Context, page, size int64, sbs ...squirrel.SelectBuilder)([]*{{.upperStartCamelObject}}, int64, error) {
    var sb squirrel.SelectBuilder
    if sbs == nil {
        sb = squirrel.Select()
    } else {
        sb = sbs[0]
    }

    // count builder
    cb := sb.Column("Count(*) as count").From(m.table)
    query, args, err := cb.ToSql()
    if err != nil {
        return nil, 0, err
    }

    var total int64
    if err := m.QueryRowNoCacheCtx(ctx, &total, query, args...); err != nil {
        return nil, 0, err
    }

    // query rows
    if page <= 0 {
        page = 1
    }
    query, args, err = sb.From(m.table).Columns({{.lowerStartCamelObject}}Rows).Offset(uint64((page - 1)* size)).Limit(uint64(size)).ToSql()
    if err != nil {
        return nil, 0, err
    }
    var list []*{{.upperStartCamelObject}}
    if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
        return nil, 0, err
    }

    return list, total, nil
}

func (m *default{{.upperStartCamelObject}}Model) List(ctx context.Context, sbs ...squirrel.SelectBuilder)([]*{{.upperStartCamelObject}}, error) {
    var sb squirrel.SelectBuilder
    if sbs == nil {
        sb = squirrel.Select()
    } else {
        sb = sbs[0]
    }

    // query rows
    query, args, err := sb.From(m.table).Columns({{.lowerStartCamelObject}}Rows).ToSql()
    if err != nil {
        return nil, err
    }
    var list []*{{.upperStartCamelObject}}
    if err := m.QueryRowsNoCacheCtx(ctx, &list, query, args...); err != nil {
        return nil, err
    }

    return list, nil
}

func (m *default{{.upperStartCamelObject}}Model) Count(ctx context.Context, cbs ...squirrel.SelectBuilder)(int64, error) { 
    var cb squirrel.SelectBuilder
    if cbs == nil {
        cb = squirrel.Select()
    } else {
        cb = cbs[0]
    }

    // count builder
    query, args, err := cb.Column("Count(*) as count").From(m.table).ToSql()
    if err != nil {
        return 0, err
    }

    var total int64
    if err := m.QueryRowNoCacheCtx(ctx, &total, query, args...); err != nil {
        return 0, err
    }

    return total, err
}

