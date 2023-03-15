package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"ohcl/config"
	"ohcl/models/ohcl"
	"reflect"
	"strconv"
	"strings"
)

type DB struct {
	*sql.DB
	error
}

type QueryOptions struct {
	Conditions []Conditions
	Pagination Pagination
	TableName  string
}

type Conditions struct {
	Where string
	Value string
}

type Pagination struct {
	Limit  *int
	OffSet *int
}

// connection Will make connection to sql database
func (d DB) connection(connection string) DB {

	db, err := sql.Open("postgres", connection)
	if err != nil {
		return DB{nil, err}
	}

	if err := db.Ping(); err != nil {
		return DB{nil, err}
	}

	return DB{db, nil}
}

func (d DB) Alive() bool {
	if d.Ping() != nil {
		return false
	}

	return true
}

func (d DB) Error() error {
	return d.error
}

func Db() DB {
	d := new(DB)

	// initialize database.
	dsn, err := config.Get("DATABASE")
	if err != nil {
		d.error = err
	}

	return d.connection(dsn)
}

func (d DB) close() DB {
	if err := d.DB.Close(); err != nil {
		d.error = err
	}

	return d
}

func (d DB) Insert(tableName string, recorder interface{}) DB {
	defer d.close()

	v := reflect.ValueOf(recorder)
	t := v.Type()
	columns := make([]string, 0)
	values := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)

		tag := field.Tag.Get("sql")
		if tag != "" {
			s := strings.Split(tag, ":")
			columns = append(columns, s[0])
			switch s[1] {
			case "unix":
				values = append(values, strconv.FormatInt(v.Field(i).Interface().(int64), 10))
			case "symbol":
				values = append(values, v.Field(i).Interface().(string))
			case "open":
				values = append(values, strconv.FormatFloat(v.Field(i).Interface().(float64), 'f', 6, 64))
			case "high":
				values = append(values, strconv.FormatFloat(v.Field(i).Interface().(float64), 'f', 6, 64))
			case "low":
				values = append(values, strconv.FormatFloat(v.Field(i).Interface().(float64), 'f', 6, 64))
			case "close":
				values = append(values, strconv.FormatFloat(v.Field(i).Interface().(float64), 'f', 6, 64))
			case "created_at":
				values = append(values, strconv.FormatInt(v.Field(i).Interface().(int64), 10))
			default:
				continue
			}
		}
	}

	query := fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s) RETURNING id", tableName,
		strings.Join(columns, ","), strings.Join(values, ","))

	row := d.DB.QueryRow(query)
	if err := row.Err(); err != nil {
		d.error = err
		return d
	}

	return d
}

func (d DB) InsertMany(tableName string, records []ohcl.Model) DB {
	defer d.close()

	if len(records) == 0 {
		return d
	}

	// Get column names from the first record.
	v := reflect.ValueOf(records[0])
	t := v.Type()
	columns := make([]string, 0)
	for i := 0; i < v.NumField(); i++ {
		field := t.Field(i)
		if field.Name != "ID" {
			tag := field.Tag.Get("sql")
			if tag != "" {
				s := strings.Split(tag, ":")
				columns = append(columns, s[0])
			}
		}
	}

	// Build the query string.
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES ", tableName, strings.Join(columns, ",")))

	// Build the values string.
	values := make([]interface{}, 0)
	placeholders := make([]string, 0)
	i := 1
	for _, r := range records {

		// Build a single row's placeholder string
		rowPlaceholder := make([]string, len(columns))
		for j := 0; j < len(columns); j++ {
			rowPlaceholder[j] = fmt.Sprintf("$%d", i)
			i++
		}
		placeholders = append(placeholders, "("+strings.Join(rowPlaceholder, ", ")+")")

		// Append the row's values to the values slice
		values = append(values, r.Unix, r.Symbol, r.Open, r.High, r.Low, r.Close)
	}

	query := sb.String() + strings.Join(placeholders, ", ")
	fmt.Println(query)

	// Execute the query and scan the results.
	rows, err := d.DB.Query(query, values...)
	if err != nil {
		d.error = err
		return d
	}

	defer rows.Close()

	return d
}

func (d DB) GetAll(q QueryOptions, dest interface{}) DB {
	defer d.close()

	v := reflect.ValueOf(dest)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Slice {
		d.error = fmt.Errorf("GetAll: dest parameter must be a pointer to a slice")
		return d
	}

	// Build the query string.
	query := fmt.Sprintf("SELECT * FROM %s", q.TableName)
	values := make([]interface{}, 0)
	if len(q.Conditions) > 0 {
		conditions := make([]string, 0)
		for _, c := range q.Conditions {
			conditions = append(conditions, c.Where+" = ?")
			values = append(values, c.Value)
		}
		query += " WHERE " + strings.Join(conditions, " AND ")
	}
	if q.Pagination.Limit != nil && q.Pagination.OffSet != nil {
		query += fmt.Sprintf(" LIMIT %d OFFSET %d", *q.Pagination.Limit, *q.Pagination.OffSet)
	}

	// Execute the query and scan the results into the destination.
	rows, err := d.DB.Query(query, values...)
	if err != nil {
		d.error = err
		return d
	}
	defer rows.Close()

	for rows.Next() {
		item := reflect.New(v.Elem().Type().Elem())
		itemValue := item.Elem()
		fields := make([]interface{}, itemValue.NumField())
		for i := 0; i < itemValue.NumField(); i++ {
			fields[i] = itemValue.Field(i).Addr().Interface()
		}
		err := rows.Scan(fields...)
		if err != nil {
			d.error = err
			return d
		}
		v.Elem().Set(reflect.Append(v.Elem(), itemValue))
	}

	return d
}

func (d DB) Migration(str string) DB {
	defer d.Close()

	if _, err := d.DB.Exec(str); err != nil {
		d.error = err
	}

	return d
}

func QueryOption() QueryOptions {
	return QueryOptions{}
}

func (q QueryOptions) SetPagination(limit, offset int) QueryOptions {
	q.Pagination.OffSet = &offset
	q.Pagination.Limit = &limit
	return q
}

func (q QueryOptions) SetCondition(where, value string) QueryOptions {
	c := Conditions{
		Where: where,
		Value: value,
	}

	q.Conditions = append(q.Conditions, c)
	return q
}

func (q QueryOptions) SetTableName(name string) QueryOptions {
	q.TableName = name
	return q
}
