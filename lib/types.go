package lib

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

type Stats struct {
	TotalVisits     int64 `json:"total_visits"`
	LongestAbsence  int64 `json:"longest_absence"`
	Id              int8  `json:"_id"`
	ServerStartTime int64 `json:"server_start_time"`
}

func (s *Stats) GetLongestAbsence() time.Duration {
	return time.Duration(s.LongestAbsence * int64(time.Second))
}

func (s *Stats) GetServerStartTime() time.Time {
	return time.Unix(s.ServerStartTime, 0)
}

type Beat struct {
	Device    int64     `json:"device"`
	TimeStamp time.Time `json:"time_stamp"`
}

type BeatLegacy struct {
	Device    string `json:"device_name"`
	TimeStamp int64  `json:"timestamp"`
}

type Device struct {
	Id       int64  `json:"id"`
	Name     string `json:"name"`
	Token    string `json:"token,omitempty"`
	NumBeats int64  `json:"total_beats"`
}

type DeviceLegacy struct {
	Name       string `json:"device_name"`
	TotalBeats int64  `json:"total_beats"`
}

type StatsLegacy struct {
	TotalVisits    int64 `json:"total_visits"`
	LongestAbsence int64 `json:"longest_missing_beat"`
	UptimeMillis   int64 `json:"total_uptime_milli"`
}

func (l *DeviceLegacy) Migrate(c *HttpClient) (*Device, error) {
	var body []byte
	req := fasthttp.AcquireRequest()
	body, err := json.Marshal(map[string]string{"name": l.Name})
	if err != nil {
		return nil, fmt.Errorf("error marshalling device %s: %s", l.Name, err)
	}
	req.SetBody(body)
	req.Header.SetMethod("POST")
	req.Header.SetContentType("application/json")
	req.Header.Set("Authorization", c.auth)
	req.SetRequestURI(fmt.Sprintf("%s/api/devices", c.baseUrl))
	res := fasthttp.AcquireResponse()
	if err := c.client.Do(req, res); err != nil {
		return nil, fmt.Errorf("error migrating device %s: %s", l.Name, err)
	}
	fasthttp.ReleaseRequest(req)
	body = res.Body()
	var d Device
	err = json.Unmarshal(body, &d)
	fasthttp.ReleaseResponse(res)
	if err != nil {
		return nil, fmt.Errorf("error unmarshalling device %s: %s", l.Name, err)
	}
	d.NumBeats = l.TotalBeats
	return &d, nil
}

func (b *Beat) GetTimeStamp() time.Time {
	return b.TimeStamp
}

func (s *StatsLegacy) GetServerStartTime() time.Time {
	return time.Now().Add(-time.Duration(s.UptimeMillis * int64(time.Millisecond)))
}

func (s *StatsLegacy) GetLongestAbsence() time.Duration {
	return time.Duration(s.LongestAbsence * int64(time.Second))
}

func (l *BeatLegacy) Migrate(ds map[string]int64) Beat {
	return Beat{Device: ds[l.Device], TimeStamp: time.Unix(l.TimeStamp, 0)}
}
