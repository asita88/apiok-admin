package models

import (
	"apiok-admin/app/packages"
	"apiok-admin/app/validators"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"
)

func accessLogWhereClause(param *validators.AccessLogList) (string, []interface{}) {
	var parts []string
	var args []interface{}
	if param.StartTime > 0 {
		parts = append(parts, "CAST(NULLIF(TRIM(`time`),'') AS UNSIGNED) >= ?")
		args = append(args, param.StartTime)
	}
	if param.EndTime > 0 {
		parts = append(parts, "CAST(NULLIF(TRIM(`time`),'') AS UNSIGNED) <= ?")
		args = append(args, param.EndTime)
	}
	q := strings.TrimSpace(param.Query)
	if len(q) != 0 {
		search := "%" + q + "%"
		parts = append(parts, "(url LIKE ? OR `request` LIKE ? OR server_name LIKE ? OR remote_addr LIKE ? OR x_forwarded_for LIKE ? OR referer LIKE ? OR user_agent LIKE ? OR upstream_addr LIKE ? OR request_id LIKE ? OR block_reason LIKE ? OR block_rule LIKE ? OR `method` LIKE ? OR CAST(`status` AS CHAR) LIKE ?)")
		args = append(args, search, search, search, search, search, search, search, search, search, search, search, search, search)
	}
	if len(parts) == 0 {
		return "", nil
	}
	return strings.Join(parts, " AND "), args
}

func applyAccessLogListFilters(tx *gorm.DB, param *validators.AccessLogList) *gorm.DB {
	clause, args := accessLogWhereClause(param)
	if clause != "" {
		tx = tx.Where(clause, args...)
	}
	return tx
}

func accessLogStartOfUTCDay(t time.Time) time.Time {
	t = t.UTC()
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
}

func AccessLogDailyTableName(day time.Time) string {
	return fmt.Sprintf("logs-%s", day.UTC().Format("2006-01-02"))
}

func accessLogTableNames(param *validators.AccessLogList) []string {
	now := time.Now().UTC()
	var startDay, endDay time.Time
	switch {
	case param.StartTime > 0 && param.EndTime > 0:
		startDay = accessLogStartOfUTCDay(time.Unix(param.StartTime, 0))
		endDay = accessLogStartOfUTCDay(time.Unix(param.EndTime, 0))
	case param.StartTime > 0:
		startDay = accessLogStartOfUTCDay(time.Unix(param.StartTime, 0))
		endDay = accessLogStartOfUTCDay(now)
	case param.EndTime > 0:
		d := accessLogStartOfUTCDay(time.Unix(param.EndTime, 0))
		startDay = d
		endDay = d
	default:
		d := accessLogStartOfUTCDay(now)
		startDay = d
		endDay = d
	}
	if endDay.Before(startDay) {
		startDay, endDay = endDay, startDay
	}
	var names []string
	for d := startDay; !d.After(endDay); d = d.AddDate(0, 0, 1) {
		names = append(names, AccessLogDailyTableName(d))
	}
	return names
}

func filterExistingAccessLogTables(db *gorm.DB, candidates []string) []string {
	if len(candidates) == 0 {
		return nil
	}
	type nameRow struct {
		Name string `gorm:"column:TABLE_NAME"`
	}
	var found []nameRow
	err := db.Raw(
		"SELECT TABLE_NAME FROM information_schema.TABLES WHERE TABLE_SCHEMA = DATABASE() AND TABLE_NAME IN ?",
		candidates,
	).Scan(&found).Error
	if err != nil {
		return candidates
	}
	if len(found) == 0 {
		return nil
	}
	exist := make(map[string]struct{}, len(found))
	for _, r := range found {
		exist[r.Name] = struct{}{}
	}
	out := make([]string, 0, len(found))
	for _, c := range candidates {
		if _, ok := exist[c]; ok {
			out = append(out, c)
		}
	}
	return out
}

func resolveAccessLogTables(db *gorm.DB, param *validators.AccessLogList) []string {
	return filterExistingAccessLogTables(db, accessLogTableNames(param))
}

func buildAccessLogUnionAllSQL(tables []string, whereClause string) string {
	var b strings.Builder
	for i, t := range tables {
		if i > 0 {
			b.WriteString(" UNION ALL ")
		}
		b.WriteString("(SELECT * FROM `")
		b.WriteString(t)
		b.WriteString("`")
		if whereClause != "" {
			b.WriteString(" WHERE ")
			b.WriteString(whereClause)
		}
		b.WriteString(")")
	}
	return b.String()
}

func accessLogFromTables(db *gorm.DB, tables []string, param *validators.AccessLogList) *gorm.DB {
	if len(tables) == 0 {
		return db.Table("(SELECT NULL AS id WHERE 1=0) AS access_union")
	}
	if len(tables) == 1 {
		return applyAccessLogListFilters(db.Table("`"+tables[0]+"`"), param)
	}
	clause, args := accessLogWhereClause(param)
	unionSQL := buildAccessLogUnionAllSQL(tables, clause)
	rep := make([]interface{}, 0, len(args)*len(tables))
	for range tables {
		rep = append(rep, args...)
	}
	return db.Table("(?) AS access_union", db.Raw(unionSQL, rep...))
}

func AccessLogListPage(param *validators.AccessLogList) (list []map[string]interface{}, total int, listError error) {
	db := packages.GetDb()
	tables := resolveAccessLogTables(db, param)
	if len(tables) == 0 {
		return nil, 0, nil
	}
	tx := accessLogFromTables(db, tables, param)

	countError := ListCount(tx, &total)
	if countError != nil {
		listError = countError
		return
	}

	tx2 := accessLogFromTables(db, tables, param).Order("id desc")
	listError = ListPaginate(tx2, &list, &param.BaseListPage)

	return
}
