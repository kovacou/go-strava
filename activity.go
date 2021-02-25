// Copyright © 2020 Alexandre KOVAC <contact@kovacou.com>.
//
// Use of this source code is governed by an MIT-style
// license that can be found in the LICENSE file.

package strava

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/kovacou/go-types"
)

// Activity is a representation of an Activity on Strava.com.
type Activity struct {
	ID                 uint64    `json:"id"`
	ExternalID         string    `json:"external_id"`
	UploadId           int64     `json:"upload_id"`
	Name               string    `json:"name"`
	Description        string    `json:"description"`
	Type               string    `json:"type"`
	Distance           float64   `json:"distance"`
	MovingTime         uint64    `json:"moving_time"`
	ElapsedTime        uint64    `json:"elapsed_time"`
	AverageSpeed       float64   `json:"average_speed"`
	AverageCadence     float64   `json:"average_cadence"`
	AverageHeartRate   float64   `json:"average_heartrate"`
	MaxSpeed           float64   `json:"max_speed"`
	MaxHeartRate       float64   `json:"max_heartrate"`
	MaxWatts           float64   `json:"max_watts"`
	Score              float64   `json:"suffer_score"`
	Calories           float64   `json:"calories"`
	TotalElevationGain float64   `json:"total_elevation_gain"`
	HighestElevation   float64   `json:"elev_high"`
	LowestElevation    float64   `json:"elev_low"`
	StartLocation      []float64 `json:"start_latlng"`
	EndLocation        []float64 `json:"end_latlng"`
	DeviceName         string    `json:"device_name"`

	StartAt time.Time `json:"start_date"`

	// Nested
	Splits []Split `json:"splits_metric"`
	Laps   []Lap   `json:"laps"`

	Athlete struct {
		ID uint64 `json:"id"`
	} `json:"athlete"`

	Map struct {
		ID              string `json:"id"`
		Polyline        string `json:"polyline"`
		SummaryPolyline string `json:"summary_polyline"`
	} `json:"map"`
}

// Split is a representation of an Split on Strava.com.
type Split struct {
	Distance                  float64 `json:"distance"`
	ElevationDifference       float64 `json:"elevation_difference"`
	ElapsedTime               uint64  `json:"elapsed_time"`
	MovingTime                uint64  `json:"moving_time"`
	AverageSpeed              float64 `json:"average_speed"`
	AverageGradeAdjustedSpeed float64 `json:"average_grade_adjusted_speed"`
	AverageHeartRate          float64 `json:"average_heartrate"`
	PaceZone                  uint64  `json:"pace_zone"`
}

// Lap is a representation of an Lap on Strava.com
type Lap struct {
	ID                 uint64  `json:"id"`
	Split              uint64  `json:"split"`
	Index              uint64  `json:"lap_index"`
	Name               string  `json:"name"`
	Distance           float64 `json:"distance"`
	ElapsedTime        uint64  `json:"elapsed_time"`
	MovingTime         uint64  `json:"moving_time"`
	AverageSpeed       float64 `json:"average_speed"`
	AverageHeartRate   float64 `json:"average_heartrate"`
	AverageCadence     float64 `json:"average_cadence"`
	MaxSpeed           float64 `json:"max_speed"`
	MaxHeartRate       float64 `json:"max_heartrate"`
	TotalElevationGain float64 `json:"total_elevation_gain"`
	PaceZone           uint64  `json:"pace_zone"`
}

// Activity returns information of the activity for the given id.
func (s *strava) Activity(id uint64) (out Activity, err error) {
	r, err := s.GET(fmt.Sprintf("/activities/%d", id), RequestParams{
		WithBearer: true,
	})

	defer closeHTTPResponse(r)
	if err != nil {
		return
	}

	if r.StatusCode == http.StatusOK {
		b, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(b, &out)
	}
	return
}

// ActivitiesRequest contains parameters to request activities from Strava.
type ActivitiesRequest struct {
	Before  types.Date `json:"before"`
	After   types.Date `json:"after"`
	Page    uint64     `json:"page"`
	PerPage uint64     `json:"per_page"`
}

// Queries convert the request to queries.
func (p ActivitiesRequest) Queries() types.Map {
	m := types.Map{}

	if !p.After.IsZero() {
		m.Set("after", p.After.Unix())
	}

	if !p.Before.IsZero() {
		m.Set("before", p.Before.Unix())
	}

	if p.Page > 0 {
		m.Set("page", p.Page)
	}

	if p.PerPage > 0 {
		m.Set("per_page", p.PerPage)
	}

	return m
}

// Activities returns informations about activities
func (s *strava) Activities(p ActivitiesRequest) (out []Activity, err error) {
	r, err := s.GET("/activities", RequestParams{
		WithBearer: true,
		Queries:    p.Queries(),
	})

	defer closeHTTPResponse(r)
	if err != nil {
		return
	}

	if r.StatusCode == http.StatusOK {
		b, _ := ioutil.ReadAll(r.Body)
		err = json.Unmarshal(b, &out)
	}
	return
}
