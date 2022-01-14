package qoingohelper

import (
	"fmt"
	"log"
)

const (
	colorRed    string = "\033[31m"
	colorGreen  string = "\033[32m"
	colorYellow string = "\033[33m"
	colorBlue   string = "\033[34m"
	colorPurple string = "\033[35m"
	colorCyan   string = "\033[36m"
	colorWhite  string = "\033[37m"
)

//this is for logger info level
func LoggerInfo(message interface{}) {
	fmt.Println(colorCyan)
	log.Println(colorCyan, "========== Start Info Message ==========")
	log.Println(colorCyan, "INFO => "+fmt.Sprint(message)+".")
	log.Println(colorCyan, "========== End Of Info Message ==========")
}

//this is for logger warning level
func LoggerWarning(message interface{}) {
	fmt.Println(colorYellow)
	log.Println(colorYellow, "========== Start Warning Message ==========")
	log.Println(colorYellow, "INFO => "+fmt.Sprint(message)+".")
	log.Println(colorYellow, "========== End Of Info Message ==========")
}

//this is for logger success level
func LoggerSuccess(message interface{}) {
	fmt.Println(colorGreen)
	log.Println(colorGreen, "========== Start Success Message ==========")
	log.Println(colorGreen, "INFO => "+fmt.Sprint(message)+".")
	log.Println(colorGreen, "========== End Of Info Message ==========")
}

//this is for logger error level
func LoggerError(err error) {
	if err != nil {
		fmt.Println(colorRed)
		log.Println(colorRed, "========== Start Error Message ==========")
		log.Println(colorRed, "ERROR => "+err.Error()+".")
		log.Println(colorRed, "========== End Of Error Message ==========")
	}
}

//this is for logger debug level
func LoggerDebug(msg interface{}) {
	fmt.Println(colorPurple)
	log.Println(colorPurple, "========== Start Debug Message ==========")
	log.Println(colorPurple, "ERROR => "+fmt.Sprint(msg)+".")
	log.Println(colorPurple, "========== End Of Error Message ==========")
}
