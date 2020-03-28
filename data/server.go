package data

import (
	"encoding/json"
	"io"
)

type IsRunning struct {
	Status		bool	`json:"status"`
}

var serverStatus = IsRunning{Status:true}

func (t *IsRunning) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(t)
}

func UpdateStatus(t bool) {
	serverStatus.Status = t
	return
}

func GetStatus() IsRunning{
	return serverStatus
}
