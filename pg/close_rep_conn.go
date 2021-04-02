package main

func closeRepConn() {
	err := repConn.Close()
	if err != nil {
		log.WithError(err).Errorln("db.Close()")
		return
	}
	log.Infoln("db closed connection")
}
