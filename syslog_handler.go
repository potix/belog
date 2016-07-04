package belog

import (
	"log/syslog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	facilityMap = map[string]syslog.Priority{
		"KERN":     syslog.LOG_KERN,
		"USER":     syslog.LOG_USER,
		"MAIL":     syslog.LOG_MAIL,
		"DAEMON":   syslog.LOG_DAEMON,
		"AUTH":     syslog.LOG_AUTH,
		"SYSLOG":   syslog.LOG_SYSLOG,
		"LPR":      syslog.LOG_LPR,
		"NEWS":     syslog.LOG_NEWS,
		"UUCP":     syslog.LOG_UUCP,
		"CRON":     syslog.LOG_CRON,
		"AUTHPRIV": syslog.LOG_AUTHPRIV,
		"FTP":      syslog.LOG_FTP,
		"LOCAL0":   syslog.LOG_LOCAL0,
		"LOCAL1":   syslog.LOG_LOCAL1,
		"LOCAL2":   syslog.LOG_LOCAL2,
		"LOCAL3":   syslog.LOG_LOCAL3,
		"LOCAL4":   syslog.LOG_LOCAL4,
		"LOCAL5":   syslog.LOG_LOCAL5,
		"LOCAL6":   syslog.LOG_LOCAL6,
		"LOCAL7":   syslog.LOG_LOCAL7,
	}
)

//SyslogHandler is handler of syslog
type SyslogHandler struct {
	network  string
	addr     string
	tag      string
	facility syslog.Priority
	writer   *syslog.Writer
	mutex    *sync.RWMutex
}

//Open is open syslog
func (h *SyslogHandler) Open() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	writer, err := syslog.Dial(h.network, h.addr, h.facility, h.tag)
	if err != nil {
		go h.reopenSyslog()
	} else {
		h.writer = writer
	}
}

//Write is output to syslog
func (h *SyslogHandler) Write(loggerName string, logEvent LogEvent, formattedLog string) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	if h.writer == nil {
		// statistics
		return
	}
	switch logEvent.LogLevelNum() {
	case LogLevelEmerg:
		if err := h.writer.Emerg(formattedLog); err != nil {
			// statistics
		}
	case LogLevelAlert:
		if err := h.writer.Alert(formattedLog); err != nil {
			// statistics
		}
	case LogLevelCrit:
		if err := h.writer.Crit(formattedLog); err != nil {
			// statistics
		}
	case LogLevelError:
		if err := h.writer.Err(formattedLog); err != nil {
			// statistics
		}
	case LogLevelWarn:
		if err := h.writer.Warning(formattedLog); err != nil {
			// statistics
		}
	case LogLevelNotice:
		if err := h.writer.Notice(formattedLog); err != nil {
			// statistics
		}
	case LogLevelInfo:
		if err := h.writer.Info(formattedLog); err != nil {
			// statistics
		}
	case LogLevelTrace:
		fallthrough
	case LogLevelDebug:
		if err := h.writer.Debug(formattedLog); err != nil {
			// statistics
		}
	default:
		// statistics
	}
}

//Flush is call nothing to do
func (h *SyslogHandler) Flush() {
}

//Close is close syslog
func (h *SyslogHandler) Close() {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if err := h.writer.Close(); err != nil {
		// statistics
	}
	h.writer = nil
}

//SetNetworkAndAddr is set netowrk type and remote addr
func (h *SyslogHandler) SetNetworkAndAddr(network string, addr string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	if h.network != network || h.addr != addr {
		if err := h.writer.Close(); err != nil {
			// statistics
		}
		h.writer = nil
		writer, err := syslog.Dial(h.network, h.addr, h.facility, h.tag)
		if err != nil {
			go h.reopenSyslog()
		} else {
			h.writer = writer
		}
	}
	h.network = network
	h.addr = addr
}

//SetTag is set tag
func (h *SyslogHandler) SetTag(tag string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	h.tag = tag
}

//SetFacility is set facility
func (h *SyslogHandler) SetFacility(facility string) {
	h.mutex.Lock()
	defer h.mutex.Unlock()
	fac, ok := facilityMap[facility]
	if !ok {
		h.facility = syslog.LOG_LOCAL0
	}
	h.facility = fac
}

func (h *SyslogHandler) reopenSyslog() {
	// retry Open
	time.Sleep(time.Second)
	h.Open()
}

//NewSyslogHandler is create SyslogHandler
func NewSyslogHandler() (syslogHandler *SyslogHandler) {
	return &SyslogHandler{
		network:  "",
		addr:     "",
		tag:      filepath.Base(os.Args[0]),
		facility: syslog.LOG_LOCAL0,
		mutex:    new(sync.RWMutex),
	}
}

func init() {
	RegisterHandler("SyslogHandler", func() (handler Handler) {
		return NewSyslogHandler()
	})
}
