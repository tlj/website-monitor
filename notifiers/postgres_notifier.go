package notifiers

import (
	"context"
	"fmt"
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"time"
	"website-monitor/result"
)

type PgLog struct {
	Id            int64
	Name          string
	DisplayUrl    string
	MatchesChecks bool
	CreatedAt     time.Time
}

type PostgresNotifier struct {
	name   string
	dbConn *pg.DB
}

func NewPostgresNotifier(name string, options map[string]string) (*PostgresNotifier, error) {
	mn := &PostgresNotifier{
		name: name,
	}

	opts, err := pg.ParseURL(options["url"])
	if err != nil {
		return nil, fmt.Errorf("error parsing Postgres URL: %v", err)
	}

	mn.dbConn = pg.Connect(opts)
	if err := mn.dbConn.Ping(context.Background()); err != nil {
		return nil, fmt.Errorf("error pinging Postgres: %v", err)
	}

	err = mn.dbConn.Model((*PgLog)(nil)).CreateTable(&orm.CreateTableOptions{
		IfNotExists: true,
	})
	if err != nil {
		return nil, fmt.Errorf("error creating table PgLog: %v", err)
	}

	return mn, nil
}

func (mn *PostgresNotifier) Name() string {
	return mn.name
}

func (mn *PostgresNotifier) Notify(name, displayUrl string, result *result.Results) error {
	if mn.dbConn == nil {
		return fmt.Errorf("no postgres database connected")
	}

	pl := PgLog{
		Name:          name,
		DisplayUrl:    displayUrl,
		MatchesChecks: result.AllTrue(),
		CreatedAt:     time.Now(),
	}

	_, err := mn.dbConn.Model(&pl).Insert()
	if err != nil {
		return fmt.Errorf("error inserting log for %s: %v", name, err)
	}

	return nil
}
