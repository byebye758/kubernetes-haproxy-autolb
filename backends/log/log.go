package log

import (
	"github.com/Sirupsen/logrus"
	"os"
)

func Log(logcontent, loglable string) {
	logname := "/var/log/k8s-auto/k8s-auto.log"

	var log = logrus.New()

	file, err := os.OpenFile(logname, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
		log.Out = file
	} else {
		log.Info("ERROR")

	}

	//log.Info(logcontent)

	log.WithFields(logrus.Fields{
		"lable": loglable,
	}).Info(logcontent)

}
