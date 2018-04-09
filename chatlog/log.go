package chatlog

const (
	//TimeFMT for the log file
	TimeFMT = "2006 Jan 2 15:04:05 UTC"
	//StrFMT of the log file, [Time Stamp] command and or error
	StrFMT = "[%s] %s"
)

var (
	//LogFile to be written to
	LogFile = "request-logs.txt"
)
