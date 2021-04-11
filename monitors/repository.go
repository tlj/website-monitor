package monitors

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"time"
	"website-monitor/content_checkers"
	"website-monitor/schedule"
)

type MonitorHeaders map[string]string

func (mh *MonitorHeaders) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &mh)
		return nil
	case string:
		json.Unmarshal([]byte(v), &mh)
		return nil
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

type MonitorDb struct {
	Id                          int64
	UserId                      int64 `db:"user_id"`
	Name                        string
	Type                        string
	Url                         string
	DisplayUrl                  sql.NullString `db:"display_url"`
	Headers                     MonitorHeaders
	ExpectedStatusCode          int `db:"expected_status_code"`
	Interval                    string
	IntervalVariationPercentage int    `db:"interval_variation_percentage"`
	ScheduleDays                string `db:"schedule_days"`
	ScheduleHours               string `db:"schedule_hours"`
}

type RepositoryInterface interface {
	Find(id int64) (*Monitor, error)
	All() ([]*Monitor, error)
	FindByUserId(userId int64) ([]*Monitor, error)
}

type Repository struct {
	db          *sqlx.DB
	checkerRepo *content_checkers.Repository
}

func NewRepository(db *sqlx.DB, checkerRepo *content_checkers.Repository) *Repository {
	return &Repository{
		db:          db,
		checkerRepo: checkerRepo,
	}
}

func (r *Repository) Find(id int64) (*Monitor, error) {
	mdb := &MonitorDb{}

	if err := r.db.Get(mdb, "SELECT * FROM monitors WHERE id = $1", id); err != nil {
		return nil, fmt.Errorf("error looking up monitor id '%d': %s", id, err)
	}

	mon, err := r.transformMonitorDbToMonitor(mdb)
	if err != nil {
		return nil, err
	}

	if err := r.applyContentChecksToMonitors([]*Monitor{mon}); err != nil {
		return nil, err
	}

	return mon, nil
}

func (r *Repository) All() ([]*Monitor, error) {
	var mdbs []MonitorDb

	if err := r.db.Select(&mdbs, "SELECT * FROM monitors"); err != nil {
		return nil, fmt.Errorf("error looking up monitors: %s", err)
	}

	var ret []*Monitor
	for _, mdb := range mdbs {
		m, err := r.transformMonitorDbToMonitor(&mdb)
		if err != nil {
			return nil, err
		}
		ret = append(ret, m)
	}

	if err := r.applyContentChecksToMonitors(ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Repository) applyContentChecksToMonitors(mons []*Monitor) error {
	var monitorIds []int64
	monitorIdx := map[int64]int{}

	for k, m := range mons {
		monitorIdx[m.Id] = k
		monitorIds = append(monitorIds, m.Id)
	}

	if len(monitorIds) > 0 {
		checks, err := r.checkerRepo.FindByMonitorIds(monitorIds)
		if err != nil {
			return err
		}
		for monitorId, c := range checks {
			mons[monitorIdx[monitorId]].ContentChecks = c
		}
	}

	return nil
}

func (r *Repository) FindByUserId(userId int64) ([]*Monitor, error) {
	var mdbs []MonitorDb

	if err := r.db.Select(&mdbs, "SELECT * FROM monitors WHERE user_id = $1", userId); err != nil {
		return nil, fmt.Errorf("error looking up monitors by user id '%d': %s", userId, err)
	}

	var ret []*Monitor
	for _, mdb := range mdbs {
		m, err := r.transformMonitorDbToMonitor(&mdb)
		if err != nil {
			return nil, err
		}
		ret = append(ret, m)
	}

	if err := r.applyContentChecksToMonitors(ret); err != nil {
		return nil, err
	}

	return ret, nil
}

func (r *Repository) transformMonitorDbToMonitor(mdb *MonitorDb) (*Monitor, error) {
	m := Monitor{
		Id:                 mdb.Id,
		Name:               mdb.Name,
		Type:               MonitorType(mdb.Type),
		UserId:             mdb.UserId,
		Url:                mdb.Url,
		DisplayUrl:         mdb.Url,
		Headers:            mdb.Headers,
		ExpectedStatusCode: mdb.ExpectedStatusCode,
	}

	if m.DisplayUrl == "" {
		m.DisplayUrl = m.Url
	}

	if m.Headers == nil {
		m.Headers = make(map[string]string)
	}

	if _, ok := m.Headers["Referer"]; !ok {
		m.Headers["Referer"] = m.DisplayUrl
	}

	i, _ := time.ParseDuration(mdb.Interval)
	d := schedule.Days{}
	d.FromString(mdb.ScheduleDays)
	h := schedule.Hours{}
	h.FromString(mdb.ScheduleHours)

	m.Scheduler = schedule.NewScheduler(i, &mdb.IntervalVariationPercentage, h, d)

	return &m, nil
}
