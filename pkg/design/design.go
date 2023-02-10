package design

import (
	"flag"
	"fmt"
	"strings"
	"time"
)

//static [Icons] and design variables:
var (
	Plus     = "[\033[1;32m+\033[0m]"
	BeAware  = "[\033[1;33m!\033[0m]"
	Negative = "[\033[1;31m-\033[0m]"
	Possible = "[\033[1;33m?\033[0m]"

	OK       = "[\033[32mOK\033[0m]"
	Success  = "[\033[1;32mOK\033[0m]"
	Info     = "[\033[34mINF\033[0m]"
	Fail     = "[\033[31mFAI\033[0m]"
	Warning  = "[\033[33mWAR\033[0m]"
	Critical = "[\033[1;41mCRITICAL\033[0m]"
	Debug    = "[\033[36mDEBUG\033[0m]"
	Payload  = "[\033[32mPAYLOAD\033[0m]"
	Example  = "\033[33mExemple\033[0m:"

	Null      = "--------"
	Reflect   = "\033[130;100mReflect\033[0m"
	RespError = "\033[130;41mRespErr\033[0m"
	XSS       = "\033[_COLOR_mXSS\033[0m"

	White       = "\033[0m"
	WhiteBG     = "\033[47m"
	Black       = "\033[30m"
	Grey        = "\033[90m"
	GreyLight   = "\033[1;90m"
	Red         = "\033[31m"
	RedLight    = "\033[1:31m"
	Orange      = "\033[33m"
	OrangeLight = "\033[1;33m"
	Yellow      = "\033[93m"
	Green       = "\033[32m"
	GreenLight  = "\033[1;32m"
	Blue        = "\033[34m"
	BlueLight   = "\033[36m"
	Pink        = "\033[1;35m"
	Purpel      = "\033[35m"

	BG_Red = "\033[1;41m"
)

func InfoBanner(show bool) {
	fmt.Println(strings.Repeat("_", 64), "\n\r")

	if show {
		flag.VisitAll(func(f *flag.Flag) {

			lenName := len(fmt.Sprintf("%v", f.Value))
			Name := fmt.Sprintf("%v", f.Name)
			var value string
			if lenName > 0 && len(Name) > 2 {
				if Name == "raw" {
					value = "HTTP Raw"
				} else {
					value = fmt.Sprintf("%v", f.Value)
				}
				fmt.Printf(" :: %s\r\t\t\t: %s\n", strings.Title(f.Name), value)
			}
		})
		fmt.Println(strings.Repeat("_", 64), "\n\r")
		time.Sleep(2000 * time.Millisecond)
	}
}

func Colortxt(s, c string, clean bool) string {
	/** Change color of text
	 */
	var txt string
	switch c {
	case "r":
		c = Red
	case "lr":
		c = RedLight
	case "br":
		c = BG_Red
	case "g":
		c = Green
	case "lg":
		c = GreenLight
	case "b":
		c = Blue
	case "lb":
		c = BlueLight
	case "o":
		c = Orange
	case "lo":
		c = OrangeLight
	case "y":
		c = Yellow
	case "p":
		c = Pink
	case "P":
		c = Purpel
	case "B":
		c = Black

	default:
		c = White
	}

	if clean {
		txt = (c + s + White)
	} else {
		txt = ("(" + (c + s + White) + ")")
	}

	return txt
}

func Disclaimer() {
	fmt.Println(BeAware + " Stay ethical. The creator of the tool is not responsible for any misuse or damage.")
}

func Color_boolean(b bool) string {
	var booleanColor string

	if b {
		booleanColor = GreenLight + fmt.Sprint(b) + White
	} else {
		booleanColor = Blue + fmt.Sprint(b) + White
	}

	return booleanColor
}

func Color_StatusCode(status_code string) string {
	//Check recived *status codes* and add color to it:

	switch string(status_code)[0:1] {
	case "2": //Status: 200
		status_code = fmt.Sprintf(Green + status_code + White)
	case "3": //Status: 300
		switch status_code {
		case "301":
			status_code = fmt.Sprintf(BlueLight + status_code + White)
		case "302":
			status_code = fmt.Sprintf(Blue + status_code + White)
		default:
			status_code = fmt.Sprintf(BlueLight + status_code + White)
		}
	case "1": //Status: 100
		status_code = fmt.Sprintf(Grey + status_code + White)
	case "4": //Status: 400
		switch status_code {
		case "404":
			status_code = fmt.Sprintf(Grey + status_code + White)
		case "403":
			status_code = fmt.Sprintf(BG_Red + status_code + White)
		case "429":
			status_code = fmt.Sprintf(BG_Red + status_code + White)
		case "400":
			status_code = fmt.Sprintf(Purpel + status_code + White)
		default:
			status_code = fmt.Sprintf(Blue + status_code + White)
		}
	case "5": //Status: 500
		switch status_code {
		case "502":
			status_code = fmt.Sprintf(Yellow + status_code + White)
		case "503":
			status_code = fmt.Sprintf(OrangeLight + status_code + White)
		case "504":
			status_code = fmt.Sprintf(Yellow + status_code + White)
		default:
			status_code = fmt.Sprintf(OrangeLight + status_code + White)
		}
	}

	return status_code
}

func Color_Time(respTime float64) string {
	//Check recived *response time* and add color to it if it's odd from the original responses:

	if respTime > 4 {
		return fmt.Sprintf("\033[31m%.3f\033[0m", respTime)
	}

	return fmt.Sprintf("%.3f", respTime)
}
