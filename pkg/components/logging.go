package components

import (
    "fmt"
    "time"
)

type Logger struct {
    LogFile string
}

func (l *Logger) LogItem(workflowName string, message string) {
    fmt.Printf("%s::%s::%s\n", workflowName, time.Now().Format(time.RFC3339), message)
}