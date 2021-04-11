package notifiers

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	log "github.com/sirupsen/logrus"
)

type NotifierOptions map[string]interface{}

func (ndo *NotifierOptions) Scan(val interface{}) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &ndo)
		return nil
	case string:
		json.Unmarshal([]byte(v), &ndo)
		return nil
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}

func (ndo *NotifierOptions) Value() (driver.Value, error) {
	return json.Marshal(ndo)
}

type NotifierDb struct {
	Id      int64
	UserId  int64 `db:"user_id"`
	Name    string
	Type    string
	Options NotifierOptions
}

type RepositoryInterface interface {
	Find(id int64) (Notifier, error)
	FindByUserId(userId int64) ([]Notifier, error)
	FindByMonitorId(monitorId int64) ([]Notifier, error)
}

type Repository struct {
	db *sqlx.DB
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Find(id int64) (Notifier, error) {
	mdb := &NotifierDb{}

	if err := r.db.Get(mdb, "SELECT * FROM notifiers WHERE id = $1", id); err != nil {
		return nil, fmt.Errorf("error looking up monitor id '%d': %s", id, err)
	}

	return TransformNotifierDbToNotifier(mdb)
}

func (r *Repository) FindByUserId(userId int64) ([]Notifier, error) {
	var mdbs []NotifierDb

	if err := r.db.Select(&mdbs, "SELECT * FROM notifiers WHERE user_id = $1", userId); err != nil {
		return nil, fmt.Errorf("error looking up notifiers by user id '%d': %s", userId, err)
	}

	var ret []Notifier
	for _, mdb := range mdbs {
		m, err := TransformNotifierDbToNotifier(&mdb)
		if err != nil {
			log.Errorf("error loading notifier '%d': %s", mdb.Id, err)
			continue
		}
		ret = append(ret, m)
	}

	return ret, nil
}

func (r *Repository) FindByMonitorId(monitorId int64) ([]Notifier, error) {
	var mdbs []NotifierDb

	if err := r.db.Select(&mdbs, `
		SELECT n.* 
			FROM notifiers n 
    		JOIN monitor_notifiers mn on n.id = mn.notifier_id 
			WHERE mn.monitor_id = $1
	`, monitorId); err != nil {
		return nil, fmt.Errorf("error looking up notifiers by monitor id '%d': %s", monitorId, err)
	}

	var ret []Notifier
	for _, mdb := range mdbs {
		m, err := TransformNotifierDbToNotifier(&mdb)
		if err != nil {
			log.Errorf("error loading notifier '%d': %s", mdb.Id, err)
			continue
		}
		ret = append(ret, m)
	}

	return ret, nil
}

func TransformNotifierDbToNotifier(mdb *NotifierDb) (Notifier, error) {
	switch mdb.Type {
	case "slack":
		return NewSlackNotifier(mdb.Name, mdb.Options)
	default:
		return nil, fmt.Errorf("invalid notifier type '%s'", mdb.Type)
	}
}
