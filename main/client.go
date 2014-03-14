package main

import (
	"encoding/json"
	"flag"
	"github.com/golang/glog"
	"github.com/gyokuro/tally/proto"
	"math"
	"strconv"
	"strings"
	"time"
)

var (
	timestamp  = flag.String("timestamp", "", "Event timestamp")
	eventType  = flag.String("type", "event", "Event type")
	context    = flag.String("context", "", "Event context")
	source     = flag.String("source", "", "Event source")
	lat        = flag.Float64("lat", 0., "Event location:latitude")
	lon        = flag.Float64("lon", 0., "Event location:longitude")
	attributes = flag.String("attributes", "", "Event attributes, {key=value;}+")
)

var nanoseconds = math.Pow10(9)

func unix_timestamp(secs float64) string {
	t := time.Unix(int64(secs/nanoseconds), int64(secs*nanoseconds))
	return t.Format(time.RFC3339Nano)
}

func to_seconds(nanos int64) float64 {
	return float64(nanos) / nanoseconds
}

func to_geojson(loc *Tally.Location) []float64 {
	return []float64{*loc.Lon, *loc.Lat}
}

func format_json(event *Tally.Event) (bytes []byte, err error) {
	payload := map[string]interface{}{
		"@timestamp": unix_timestamp(*event.Timestamp),
		"@type":      event.Type,
		"@source":    event.Source,
		"@context":   event.Context,
		"@location":  to_geojson(event.Location),
	}
	for _, attr := range event.Attributes {
		if attr.BoolValue != nil {
			payload[*attr.Key] = attr.BoolValue
		} else if attr.IntValue != nil {
			payload[*attr.Key] = attr.IntValue
		} else if attr.DoubleValue != nil {
			payload[*attr.Key] = attr.DoubleValue
		} else if attr.StringValue != nil {
			payload[*attr.Key] = attr.StringValue
		}
	}
	return json.Marshal(payload)
}

func parse_attribute(key string, value string) *Tally.Attribute {
	attr := Tally.Attribute{
		Key: &key,
	}
	if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
		attr.DoubleValue = &floatValue
	} else if intValue, err2 := strconv.ParseInt(value, 10, 64); err2 == nil {
		attr.IntValue = &intValue
	} else if boolValue, err3 := strconv.ParseBool(value); err3 == nil {
		attr.BoolValue = &boolValue
	} else {
		attr.StringValue = &value
	}
	return &attr
}

func main() {
	flag.Parse()

	var event = Tally.Event{
		Type:    eventType,
		Source:  source,
		Context: context,
		Location: &Tally.Location{
			Lon: lon,
			Lat: lat,
		},
	}

	if *timestamp == "" {
		now := to_seconds(time.Now().UnixNano())
		event.Timestamp = &now
	}

	for i, p := range strings.Split(*attributes, ";") {
		kv := strings.Split(p, ":")
		if len(kv) == 2 {
			glog.Infof("i=%d Key=%s, Value=%s", i, kv[0], kv[1])
			event.Attributes = append(event.Attributes, parse_attribute(kv[0], kv[1]))
		}
	}

	if json, err := format_json(&event); err == nil {
		glog.Infof("JSON2 = %s", json)
	}
}
