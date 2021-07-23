package qoinhelper

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

func LoggerInfo(message string) {
	fmt.Println(colorBlue)
	log.Println(colorBlue, "========== Start Info Message ==========")
	log.Println(colorBlue, "INFO => "+message+".")
	log.Println(colorBlue, "========== End Of Info Message ==========")
	fmt.Println(colorBlue)
}

func LoggerWarning(message string) {
	fmt.Println(colorYellow)
	log.Println(colorYellow, "========== Start Info Message ==========")
	log.Println(colorYellow, "INFO => "+message+".")
	log.Println(colorYellow, "========== End Of Info Message ==========")
	fmt.Println(colorYellow)
}

func LoggerSuccess(message string) {
	fmt.Println(colorGreen)
	log.Println(colorGreen, "========== Start Info Message ==========")
	log.Println(colorGreen, "INFO => "+message+".")
	log.Println(colorGreen, "========== End Of Info Message ==========")
	fmt.Println(colorGreen)
}

func LoggerError(err error) {
	if err != nil {
		fmt.Println(colorRed)
		log.Println(colorRed, "========== Start Error Message ==========")
		log.Println(colorRed, "ERROR => "+err.Error()+".")
		log.Println(colorRed, "========== End Of Error Message ==========")
		fmt.Println(colorRed)
	}
}

func LoggerDebug(msg string) {
	fmt.Println(colorPurple)
	log.Println(colorPurple, "========== Start Error Message ==========")
	log.Println(colorPurple, "ERROR => "+msg+".")
	log.Println(colorPurple, "========== End Of Error Message ==========")
	fmt.Println(colorPurple)
}