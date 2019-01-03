package main

import (
	"fmt"
	"log"
	"os"

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

func main() {
	influxURL, found := os.LookupEnv("INFLUXDB_URL")
	if !found {
		log.Fatal("INFLUXDB_URL was not found in the environment")
	}
	influxUser, found := os.LookupEnv("INFLUXDB_USER")
	if !found {
		log.Fatal("INFLUXDB_USER was not found in the environment")
	}
	influxPass, found := os.LookupEnv("INFLUXDB_PASS")
	if !found {
		log.Fatal("INFLUXDB_PASS was not found in the environment")
	}
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
