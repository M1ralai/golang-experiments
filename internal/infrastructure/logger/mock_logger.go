package logger

// MockLogger is a no-op logger for testing purposes
type MockLogger struct {
	InfoLogs  []LogEntry
	ErrorLogs []LogEntry
}

type LogEntry struct {
	Message string
	Details map[string]interface{}
	Error   error
}

func NewMockLogger() *MockLogger {
	return &MockLogger{
		InfoLogs:  make([]LogEntry, 0),
		ErrorLogs: make([]LogEntry, 0),
	}
}

func (l *MockLogger) Info(msg string, details map[string]interface{}) {
	l.InfoLogs = append(l.InfoLogs, LogEntry{
		Message: msg,
		Details: details,
	})
}

func (l *MockLogger) Error(msg string, err error, details map[string]interface{}) {
	l.ErrorLogs = append(l.ErrorLogs, LogEntry{
		Message: msg,
		Details: details,
		Error:   err,
	})
}
