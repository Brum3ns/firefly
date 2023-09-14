package design

import (
	"fmt"
	"strconv"
)

// Store all design values
type Design struct {
	Color
	Icons
	Debug
	Detect
	Status
	Behavior
}

type Color struct {
	WHITE       string
	WHITEBG     string
	BLACK       string
	GREY        string
	GREYLIGHT   string
	GREYBG      string
	RED         string
	REDLIGHT    string
	REDBG       string
	ORANGE      string
	ORANGELIGHT string
	ORANGEBG    string
	YELLOW      string
	GREEN       string
	GREENLIGHT  string
	BLUE        string
	BLUELIGHT   string
	PINK        string
	PURPEL      string
}

type Icons struct {
	PLUS     string
	AWARE    string
	NEGATIVE string
	POSSIBLE string
}

type Status struct {
	OK       string
	SUCCESS  string
	INFO     string
	FAIL     string
	ERROR    string
	WARNING  string
	CRITICAL string
}

type Detect struct {
	Firm      string
	Certain   string
	Tentative string
}

type Debug struct {
	DEBUG   string
	PAYLOAD string
	INPUT   string
	EXAMPLE string
}

type Behavior struct {
	NONE           string
	DIFF           string
	TIME           string
	REFLECT        string
	TRANSFORMATION string
	PATTERN        string
}

// static [Icons] and design variables:
var (
	// Colors:
	COLOR = &Color{
		WHITE:       "\033[0m",
		WHITEBG:     "\033[47m",
		BLACK:       "\033[30m",
		GREY:        "\033[90m",
		GREYLIGHT:   "\033[1;90m",
		GREYBG:      "\033[1;40m",
		RED:         "\033[31m",
		REDLIGHT:    "\033[1:31m",
		REDBG:       "\033[1;41m",
		ORANGE:      "\033[33m",
		ORANGELIGHT: "\033[1;33m",
		ORANGEBG:    "\033[43m",
		YELLOW:      "\033[93m",
		GREEN:       "\033[32m",
		GREENLIGHT:  "\033[1;32m",
		BLUE:        "\033[34m",
		BLUELIGHT:   "\033[36m",
		PINK:        "\033[1;35m",
		PURPEL:      "\033[35m",
	}

	DETECT = &Detect{
		Certain:   (COLOR.REDBG + ("Certain") + COLOR.WHITE),
		Firm:      (COLOR.ORANGEBG + ("Firm") + COLOR.WHITE),
		Tentative: (COLOR.GREYBG + ("Tentative") + COLOR.WHITE),
	}

	ICON = &Icons{
		PLUS:     ("[" + COLOR.GREENLIGHT + ("+") + COLOR.WHITE + "]"),
		AWARE:    ("[" + COLOR.ORANGELIGHT + ("!") + COLOR.WHITE + "]"),
		NEGATIVE: ("[" + COLOR.REDLIGHT + ("-") + COLOR.WHITE + "]"),
		POSSIBLE: ("[" + COLOR.ORANGELIGHT + ("?") + COLOR.WHITE + "]"),
	}

	STATUS = &Status{
		OK:       ("[" + COLOR.GREEN + ("OK") + COLOR.WHITE + "]"),
		SUCCESS:  ("[" + COLOR.GREENLIGHT + ("OK") + COLOR.WHITE + "]"),
		INFO:     ("[" + COLOR.BLUE + ("INF") + COLOR.WHITE + "]"),
		FAIL:     ("[" + COLOR.RED + ("FAI") + COLOR.WHITE + "]"),
		WARNING:  ("[" + COLOR.ORANGE + ("WAR") + COLOR.WHITE + "]"),
		ERROR:    (COLOR.REDBG + ("ERROR") + COLOR.WHITE),
		CRITICAL: (COLOR.REDBG + ("CRITICAL") + COLOR.WHITE),
	}

	DEBUG = &Debug{
		DEBUG:   ("[" + COLOR.BLUELIGHT + ("DEBUG") + COLOR.WHITE + "]"),
		PAYLOAD: ("[" + COLOR.GREEN + ("PAYLOAD") + COLOR.WHITE + "]"),
		INPUT:   ("[" + COLOR.BLUELIGHT + ("INPUT") + COLOR.WHITE + "]"),
		EXAMPLE: (COLOR.ORANGE + ("Exemple") + COLOR.WHITE),
	}

	BEHAVIOR = &Behavior{
		NONE:           ("----"),
		TRANSFORMATION: (COLOR.REDBG + ("Tfmt") + COLOR.WHITE),
		DIFF:           (COLOR.REDBG + ("Diff") + COLOR.WHITE),
		TIME:           (COLOR.ORANGEBG + ("Time") + COLOR.WHITE),
		REFLECT:        (COLOR.GREYBG + ("Reflect") + COLOR.WHITE),
		PATTERN:        (COLOR.ORANGEBG + ("Pattern") + COLOR.WHITE),
	}

	// Colorized HTTP status codes:
	STATUSCODE_COLOR = map[int]string{
		//(Default Status code colors)
		1: (COLOR.GREYLIGHT + "{CODE}" + COLOR.WHITE),
		2: (COLOR.GREENLIGHT + "{CODE}" + COLOR.WHITE),
		3: (COLOR.BLUELIGHT + "{CODE}" + COLOR.WHITE),
		4: (COLOR.PURPEL + "{CODE}" + COLOR.WHITE),
		5: (COLOR.PINK + "{CODE}" + COLOR.WHITE),
		//(100)
		100: (COLOR.GREY + "100" + COLOR.WHITE),
		//(200)
		200: (COLOR.GREEN + "200" + COLOR.WHITE),
		//(300)
		301: (COLOR.BLUELIGHT + "301" + COLOR.WHITE),
		302: (COLOR.BLUE + "302" + COLOR.WHITE),
		//(404)
		400: (COLOR.PURPEL + "400" + COLOR.WHITE),
		404: (COLOR.GREY + "404" + COLOR.WHITE),
		403: (COLOR.RED + "403" + COLOR.WHITE),
		429: (COLOR.REDBG + "429" + COLOR.WHITE),
		//(500)
		500: (COLOR.PINK + "500" + COLOR.WHITE),
		501: (COLOR.PINK + "501" + COLOR.WHITE),
		502: (COLOR.YELLOW + "502" + COLOR.WHITE),
		503: (COLOR.ORANGELIGHT + "503" + COLOR.WHITE),
	}
)

func NewDesign() *Design {
	return &Design{
		Color:    *COLOR,
		Debug:    *DEBUG,
		Icons:    *ICON,
		Detect:   *DETECT,
		Status:   *STATUS,
		Behavior: *BEHAVIOR,
	}
}

// Colorize the status code and return it as a string
func (d *Design) StatusCode(code int) string {
	if v, ok := STATUSCODE_COLOR[code]; ok {
		return v
	} else if v, ok := STATUSCODE_COLOR[code/100]; ok {
		return v
	}
	return fmt.Sprintf("\033[31m%d\033[0m", code)
}

// Colorize the word count and return it as a string
func (d *Design) WordCount(wordCount int) string {
	return d.Color.BLUELIGHT + strconv.Itoa(wordCount) + d.Color.WHITE
}

// Colorize the word line and return it as a string
func (d *Design) LineCount(lineCount int) string {
	return d.Color.BLUE + strconv.Itoa(lineCount) + d.Color.WHITE
}

// Colorize the Content Type header value and return it as a string
func (d *Design) ContentType(contentType string) string {
	return d.Color.PURPEL + contentType + d.Color.WHITE
}

// Colorize the Content Length and return it as a string
func (d *Design) ContentLength(contentLength int) string {
	return d.Color.PINK + strconv.Itoa(contentLength) + d.Color.WHITE
}

// Colorize the response time and return it as a string
func (d *Design) ResponseTime(time float64) string {
	//Check recived *response time* and add color to it if it's odd from the original responses:
	if time > 7 {
		return fmt.Sprintf("\033[31m%.4f\033[0m", time)
	}
	return fmt.Sprintf("\033[1:38m%.4f\033[0m", time)
}

// Colorize the Content Length and return it as a string
func (d *Design) IsDiff(value int) string {
	v := strconv.Itoa(value)
	if value != 0 {
		return d.Color.REDBG + v + d.Color.WHITE
	}
	return v
}
