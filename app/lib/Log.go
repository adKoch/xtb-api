package log

import "log"

func Debug(message string) {
	logNonFatal(message, "DEBUG")
}

func Info(message string) {
	logNonFatal(message, "INFO")
}

func Warn(message string) {
	logNonFatal(message, "WARN")
}

func Error(message string) {
	logNonFatal(message, "ERROR")
}

func Fatal(message string) {
	logMessage := createPrintMessage(message, "FATAL")
	log.Fatal(logMessage)
}

func logNonFatal(message string, severity string){
	log.Print(createPrintMessage(message, severity))
}

func createPrintMessage(message string, severity string) string {
	return severity + " " +
	message + "\n"
}