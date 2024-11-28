
func (m *default{{.upperStartCamelObject}}Model) tableName() string {
	return m.table
}

var {{.upperStartCamelObject}}SoftDelete bool

func init() {
    tp := reflect.TypeOf({{.upperStartCamelObject}}{})
	for i := 0; i < tp.NumField(); i++ {
		if tp.Field(i).Tag.Get("db") == "deleted_time" {
			{{.upperStartCamelObject}}SoftDelete = true
			return
		}
	}
}
