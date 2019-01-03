package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/influxdata/influxdb/client/v2"
)

// queryDB convenience function to query the database
func queryDB(clnt client.Client, db string, cmd string) (res []client.Result, err error) {
	q := client.Query{
		Command:  cmd,
		Database: db,
	}
	if response, err := clnt.Query(q); err == nil {
		if response.Error() != nil {
			return res, response.Error()
		}
		res = response.Results
	} else {
		return res, err
	}
	return res, nil
}

func doEvery(d time.Duration, f func(time.Time)) {
	f(time.Now())
	for x := range time.Tick(d) {
		f(x)
	}
}

func doReading(_ time.Time) {
	influxURL, found := os.LookupEnv("INFLUXDB_URL")
	if !found {
		log.Fatal("INFLUXDB_URL was not found in the environment")
	}
	influxUser, found := os.LookupEnv("INFLUXDB_USER")
	influxPass, found := os.LookupEnv("INFLUXDB_PASS")
	influxDB, found := os.LookupEnv("INFLUXDB_DB")
	if !found {
		log.Fatal("INFLUXDB_DB was not found in the environment")
	}
	c, err := client.NewHTTPClient(client.HTTPConfig{
		Addr:     influxURL,
		Username: influxUser,
		Password: influxPass,
	})
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	_, err = queryDB(c, influxDB, fmt.Sprintf("CREATE DATABASE %s", influxDB))
	if err != nil {
		log.Fatal(err)
	}

	apcupsdURL, found := os.LookupEnv("APCUPSD_URL")
	if !found {
		log.Fatal("APCUPSD_URL was not found in the environment")
	}
	upsStatus, err := getUpsStatus(apcupsdURL)
	if err != nil {
		log.Fatal(err)
	}

	bp, err := createBatchPoints(upsStatus)
	if err != nil {
		log.Fatal(err)
	}
	bp.SetDatabase(influxDB)

	// Write the batch
	if err := c.Write(bp); err != nil {
		log.Fatal(err)
	}

	// Close client resources
	if err := c.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	interval, found := os.LookupEnv("READING_INTERVAL")
	if !found {
		log.Fatal("READING_INTERVAL was not found in the environment")
	}
	i, err := strconv.Atoi(interval)
	if err != nil {
		log.Fatal(err)
	}
	intervalSeconds := time.Duration(1000*i) * time.Millisecond
	doEvery(intervalSeconds, doReading)
}
