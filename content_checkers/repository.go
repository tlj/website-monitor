package content_checkers

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
)

type CheckerDb struct {
	Id         int64
	MonitorId  int64 `db:"monitor_id"`
	Name       string
	Type       string
	Path       sql.NullString
	Value      string
	IsExpected bool `db:"is_expected"`
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Find(id int64) (ContentChecker, error) {
	mdb := &CheckerDb{}

	if err := r.db.Get(mdb, "SELECT * FROM content_checks WHERE id = $1", id); err != nil {
		return nil, fmt.Errorf("error looking up content checks id '%d': %s", id, err)
	}

	return TransformCheckerDbToContentChecker(mdb)
}

func (r *Repository) FindByMonitorId(monitorId int64) ([]ContentChecker, error) {
	var mdbs []CheckerDb

	if err := r.db.Select(&mdbs, "SELECT * FROM content_checks WHERE monitor_id = $1", monitorId); err != nil {
		return nil, fmt.Errorf("error looking up monitors by content checks id '%d': %s", monitorId, err)
	}

	var ret []ContentChecker
	for _, mdb := range mdbs {
		m, err := TransformCheckerDbToContentChecker(&mdb)
		if err != nil {
			return nil, err
		}
		ret = append(ret, m)
	}

	return ret, nil
}

func (r *Repository) FindByMonitorIds(monitorIds []int64) (map[int64][]ContentChecker, error) {
	var mdbs []CheckerDb

	query, args, _ := sqlx.In("SELECT * FROM content_checks WHERE monitor_id IN (?)", monitorIds)
	query = r.db.Rebind(query)

	if err := r.db.Select(&mdbs, query, args...); err != nil {
		return nil, fmt.Errorf("error looking up monitors by content checks id '%v': %s", monitorIds, err)
	}

	ret := make(map[int64][]ContentChecker, len(monitorIds))
	for _, mdb := range mdbs {
		m, err := TransformCheckerDbToContentChecker(&mdb)
		if err != nil {
			return nil, err
		}
		ret[mdb.MonitorId] = append(ret[mdb.MonitorId], m)
	}

	return ret, nil
}

func TransformCheckerDbToContentChecker(mdb *CheckerDb) (ContentChecker, error) {
	var cc ContentChecker

	switch CheckType(mdb.Type) {
	case RegexCheckType:
		cc = NewRegexChecker(mdb.Name, mdb.Value, mdb.IsExpected)
	case HtmlXpathType:
		cc = NewHtmlXPathChecker(mdb.Name, mdb.Path.String, mdb.Value, mdb.IsExpected)
	case JsonPathType:
		cc = NewJsonPathChecker(mdb.Name, mdb.Path.String, mdb.Value, mdb.IsExpected)
	case HtmlRenderType:
		cc = NewHtmlRenderSelectorChecker(mdb.Name, mdb.Path.String, mdb.Value, mdb.IsExpected)
	default:
		return nil, fmt.Errorf("unsupported contentCheck config: '%s'", mdb.Type)
	}

	return cc, nil
}
