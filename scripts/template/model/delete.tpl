func (m *default{{.upperStartCamelObject}}Model) Delete(ctx context.Context, {{.lowerStartCamelPrimaryKey}} {{.dataType}}, s ...sqlx.Session) error {
	{{if .withCache}}{{if .containsIndexCache}}data, err:=m.FindOne(ctx, {{.lowerStartCamelPrimaryKey}}, s...)
	if err!=nil{
		return err
	}

	{{end}}	{{.keys}}

	var query string
    if {{.upperStartCamelObject}}SoftDelete {
        query = fmt.Sprintf("update %s set `deleted_time` = %v where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table, time.Now().Unix())
    } else {
        query = fmt.Sprintf("delete from %s where {{.originalPrimaryKey}} = {{if .postgreSql}}$1{{else}}?{{end}}", m.table)
    }
    
     _, err {{if .containsIndexCache}}={{else}}:={{end}} m.ExecCtx(ctx, func(ctx context.Context, conn sqlx.SqlConn) (result sql.Result, err error) {
		if s!=nil{
        	return s[0].ExecCtx(ctx,query, {{.lowerStartCamelPrimaryKey}})
        }
		return conn.ExecCtx(ctx, query, {{.lowerStartCamelPrimaryKey}})
	}, {{.keyValues}}){{else}}
		_,err:= s[0].ExecCtx(ctx,query, {{.lowerStartCamelPrimaryKey}})
        	return err
        }
		_,err:=m.conn.ExecCtx(ctx, query, {{.lowerStartCamelPrimaryKey}}){{end}}
	return err
}