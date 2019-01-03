package main

import (
	"fmt"

	"github.com/mdlayher/apcupsd"

	"github.com/influxdata/influxdb/client/v2"
)

func getUpsStatus(apcupsdURL string) (*apcupsd.Status, error) {
	upsclient, err := apcupsd.Dial("tcp", apcupsdURL)
	if err != nil {
		return nil, err
	}

	status, err := upsclient.Status()
	if err != nil {
		return nil, err
	}
	return status, nil
}

func createBatchPoints(status *apcupsd.Status) (client.BatchPoints, error) {
	// Create a new point batch
	bp, err := client.NewBatchPoints(client.BatchPointsConfig{
		Precision: "s",
	})
	if err != nil {
		return nil, err
	}

	tags := map[string]string{
		"hostname":                status.Hostname,
		"ups_name":                status.UPSName,
		"model":                   status.Model,
		"APC":                     status.APC,
		"version":                 status.Version,
		"cable":                   status.Cable,
		"ups_mode":                status.UPSMode,
		"nominal_input_voltage":   fmt.Sprintf("%f", status.NominalInputVoltage),
		"nominal_battery_voltage": fmt.Sprintf("%f", status.NominalBatteryVoltage),
		"nominal_power":           fmt.Sprintf("%d", status.NominalPower),
		"serial_number":           status.SerialNumber,
	}

	fields := map[string]interface{}{
		"start_time":                     status.StartTime,
		"status":                         status.Status,
		"line_voltage":                   status.LineVoltage,
		"load_percent":                   status.LoadPercent,
		"battery_charge_percent":         status.BatteryChargePercent,
		"time_left":                      status.TimeLeft,
		"minimum_battery_charge_percent": status.MinimumBatteryChargePercent,
		"minimum_time_left":              status.MinimumTimeLeft,
		"maximum_time":                   status.MaximumTime,
		"sense":                          status.Sense,
		"low_transfer_voltage":           status.LowTransferVoltage,
		"high_transfer_voltage":          status.HighTransferVoltage,
		"alarm_del":                      status.AlarmDel,
		"battery_voltage":                status.BatteryVoltage,
		"last_transfer":                  status.LastTransfer,
		"number_transfers":               status.NumberTransfers,
		"x_on_battery":                   status.XOnBattery,
		"time_on_battery":                status.TimeOnBattery,
		"cumulative_time_on_battery":     status.CumulativeTimeOnBattery,
		"x_off_battery":                  status.XOffBattery,
		"last_self_test":                 status.LastSelftest,
		"selftest":                       status.Selftest,
		"status_flags":                   status.StatusFlags,
		"battery_date":                   status.BatteryDate,
		"internal_temp":                  status.InternalTemp,
		"output_voltage":                 status.OutputVoltage,
		"line_frequency":                 status.LineFrequency,
		"driver":                         status.Driver,
	}

	pt, err := client.NewPoint(
		"ups_data",
		tags,
		fields,
		status.Date,
	)
	if err != nil {
		return nil, err
	}
	bp.AddPoint(pt)

	return bp, err
}
