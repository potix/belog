package appender

var (
	levelToSyslogLevel = map[LogLevel]string{
		LLEmergency: LLEmergency,
		LLAlert:     LLAlert,
		LLCritical:  LLCritical,
		LLError:     LLError,
		LLWarning:   LLWarn,
		LLInfo:      LLInfo,
		LLDebug:     LLDebug,
		LLTrace:     LLDebug,
	}
)
