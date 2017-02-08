package log

import (
	"github.com/Sirupsen/logrus"
	"os"
)

func Log(logcontent, logname string) {
	var log = logrus.New()

	file, err := os.OpenFile(logname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("ERROR")

	}

	log.Info(logcontent)

}
