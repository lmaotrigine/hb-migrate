package main

import (
	"context"
	"errors"
	"flag"
	"log"

	"github.com/lmaotrigine/hb-migrate/lib"
	"github.com/spf13/pflag"
)

var databaseDsn string
var redisDsn string
var serverBaseUrl string
var token string
var password string

func init() {
	pflag.ErrHelp = errors.New("")
	pflag.StringVarP(&password, "redis-password", "p", "", "Redis password")
	pflag.StringVarP(
		&databaseDsn, "database-dsn", "d", "postgresql://postgres@localhost/postgres", "PostgreSQL database `DSN`",
	)
	pflag.StringVarP(&redisDsn, "redis-address", "r", "127.0.0.1:6379", "Redis (ReJSON) `address`")
	pflag.StringVarP(
		&serverBaseUrl, "server-base-url", "s", "http://localhost:6060", "Server base `URL` (of the *new* server)",
	)
	pflag.StringVarP(
		&token,
		"token",
		"t",
		"",
		"Token to use for authentication to the server to add new devices. This is the value of the `secret_key` config variable that you set.",
	)
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
}

func main() {
	pflag.Parse()
	httpClient := lib.NewHttpClient(serverBaseUrl, token)
	redisClient, err := lib.NewRedisClient(redisDsn, password)
	if err != nil {
		log.Fatalln(err)
	}
	postgresClient, err := lib.NewPostgresClient(databaseDsn)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Database dsn: %s, Redis %s, server %s\n", databaseDsn, redisDsn, serverBaseUrl)
	beats, err := redisClient.GetAllBeats()
	if err != nil {
		log.Fatalln(err)
	}
	devices, err := redisClient.GetAllDevices()
	if err != nil {
		log.Fatalln(err)
	}
	new_devs, err := httpClient.MigrateDevices(devices)
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("Migrated %d devices\n", len(new_devs))
	map_ := make(map[string]int64, len(new_devs))
	for _, d := range new_devs {
		map_[d.Name] = d.Id
	}
	new_beats := make([]lib.Beat, 0, len(beats))
	for _, b := range beats {
		new_beats = append(new_beats, b.Migrate(map_))
	}
	stats, err := redisClient.GetStats()
	if err != nil {
		log.Fatalln(err)
	}
	tx, err := postgresClient.BeginTransaction()
	if err != nil {
		log.Fatalln(err)
	}
	err = postgresClient.InsertStatsFromLegacy(tx, *stats)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Inserted stats")
	err = postgresClient.InsertBeats(tx, new_beats)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Inserted beats")
	err = postgresClient.UpdateNumBeats(tx, new_devs)
	if err != nil {
		log.Fatalln(err)
	}
	log.Println("Updated num_beats for devices")
	tx.Commit(context.Background())
	log.Println("Done!")
	postgresClient.Close()
}
