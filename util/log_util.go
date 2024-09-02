package util

import (
	"time"
)

func LogRequest(requestString string, uuid string, proses string) string {
	currentTime := time.Now()
	return "[Start][RequestId]= " + uuid + ", [Proses]= " + proses + ", [Time]= " + currentTime.Format("2006-01-02 15:04:05.000000") + ", [Request]= " + requestString
}

func LogResponse(responseString string, uuid string, proses string) string {
	currentTime := time.Now()
	return "[Stop][RequestId]= " + uuid + ", [Proses]= " + proses + ",  [Time]= " + currentTime.Format("2006-01-02 15:04:05.000000") + ", [Response]= " + responseString
}

func LogError(errorString string, uuid string, proses string) string {
	currentTime := time.Now()
	return "[Error][RequestId]= " + uuid + ", [Proses]= " + proses + ", [Time]= " + currentTime.Format("2006-01-02 15:04:05.000000") + ", [Error]= " + errorString
}
