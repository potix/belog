package formatter

type DateTimeFormatter interface {
	DateTimeFormat(log *belog.Log) (dateTime string)
}

type Formatter interface {
	SetDateTimeFormatter(dateTimeFormatter DateTimeFormatter)
	Format(log belog.Log) (dateTime string)
}
