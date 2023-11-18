package lib

import (
	"context"

	"github.com/jackc/pgx/v5"
)

type PostgresClient struct {
	conn *pgx.Conn
}

func NewPostgresClient(dsn string) (*PostgresClient, error) {
	conn, err := pgx.Connect(context.Background(), dsn)
	if err != nil {
		return nil, err
	}
	return &PostgresClient{conn: conn}, nil
}

func (c *PostgresClient) InsertStatsFromLegacy(tx pgx.Tx, stats StatsLegacy) error {
	_, err := tx.Exec(
		context.Background(),
		`
INSERT INTO heartbeat.stats
	(total_visits, longest_absence, server_start_time, _id)
	VALUES ($1, $2, $3, 0)
ON CONFLICT (_id)
DO UPDATE SET
	total_visits = EXCLUDED.total_visits,
	longest_absence = EXCLUDED.longest_absence,
	server_start_time = EXCLUDED.server_start_time
;
`,
		stats.TotalVisits,
		stats.GetLongestAbsence(),
		stats.GetServerStartTime(),
	)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}
	return nil
}

func (c *PostgresClient) InsertBeats(tx pgx.Tx, beats []Beat) error {
	_, err := tx.Exec(
		context.Background(),
		`
INSERT INTO heartbeat.beats
	(device, time_stamp)
SELECT
	x.device, x.time_stamp
FROM jsonb_to_recordset($1::jsonb) AS x(
	device BIGINT,
	time_stamp TIMESTAMP WITH TIME ZONE
);
`,
		beats,
	)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}
	return nil
}

func (c *PostgresClient) UpdateNumBeats(tx pgx.Tx, ds []Device) error {
	_, err := tx.Exec(
		context.Background(),
		`
UPDATE heartbeat.devices SET
	num_beats = x.total_beats
FROM jsonb_to_recordset($1::jsonb) AS x(
	id BIGINT, total_beats BIGINT
)
WHERE heartbeat.devices.id = x.id;
`,
		ds,
	)
	if err != nil {
		tx.Rollback(context.Background())
		return err
	}
	return nil
}

func (c *PostgresClient) BeginTransaction() (pgx.Tx, error) {
	return c.conn.Begin(context.Background())
}

func (c *PostgresClient) Close() {
	c.conn.Close(context.Background())
}
