package updater

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
)

type ScheduleSubParams struct {
	Enabled         bool `json:"enabled"`
	IntervalSeconds int  `json:"interval_seconds"`
}

type ScheduleParams struct {
	ScheduleSubParams
	LastUpdate      time.Time `json:"last_update"`
	NextUpdate      time.Time `json:"next_update"`
}

type PUTRequest struct {
	Enabled         bool `json:"enabled"`
	IntervalSeconds int  `json:"interval_seconds"`
}

type TriggerResponse struct {
	UpdatedCount int       `json:"updated_count"`
	Timestamp    time.Time `json:"timestamp"`
}

func GETScheduleParamsHandler(u *Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		updateTime := u.GetUpdateTime()
		lastUpdate := u.GetLastUpdate()
		nextUpdate := lastUpdate.Add(time.Duration(updateTime))

		resp := ScheduleParams{
			ScheduleSubParams: ScheduleSubParams{
				Enabled: u.IsEnabled(),
				IntervalSeconds: updateTime,
    		},
			LastUpdate: lastUpdate,
			NextUpdate: nextUpdate,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}
}

func PUTScheduleParamsHandler(u *Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var putRequest PUTRequest
		
		err := json.NewDecoder(r.Body).Decode(&putRequest)
		if err != nil || !(putRequest.IntervalSeconds <= 3600 && putRequest.IntervalSeconds >= 10) {
			log.Println("Error during JSON parsing: ", err)
			http.Error(w, `Bad Request`, http.StatusBadRequest)
			return
		}
		
		u.RestartUpdating(putRequest.IntervalSeconds)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(putRequest)
	}
}

func POSTScheduleTriggerHandler(u *Updater) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cnt, err := u.Update()
		if err != nil {
			http.Error(w, `Server error`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(TriggerResponse{
			UpdatedCount: cnt,
			Timestamp: time.Now(),
		})
	}
}