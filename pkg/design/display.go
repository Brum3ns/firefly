package design

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Brum3ns/firefly/pkg/storage"
)

func DisplayInfo(count int, result storage.Collection) {
	/** Display information on the screen
	*	Return information about the total process to the user
	 */

	var (
		reMatch        string
		banner_avgDiff string
		diff_status    string
		location       string
		transformation string
		statusCode     = strconv.FormatInt(int64(result.StatusCode), 10)
	)

	//Response was skipped/filtered:
	if result.Skip {
		return
	}

	//Diff average amount banner:
	banner := ""
	for keyItem, amount := range result.AvgAmountDiff {
		//Include amount diff:
		if amount > 0 {
			banner += fmt.Sprintf("%s:(%s+%d%s) ", keyItem, BG_Red, amount, White)
		}
	}
	if len(banner) > 0 {
		banner_avgDiff = (Plus + "AvgDiff: " + banner)
	}

	//Redirect 300 status code the add location header to the verbose:
	if statusCode[:1] == "3" {
		location = (" => " + Purpel + (result.Headers.Get("location")) + White)
	}

	//Check if user regex input got a hit:
	if result.RegexMatch {
		reMatch = "Rx:" + Color_boolean(result.RegexMatch)
	}

	//Payload transformation discovered:
	if result.Tfmt_ok {
		transformation = (Plus + "Transformation:[" + GreenLight + result.Tfmt_display + White + "]")
	}

	//Display information to CLI:
	display := fmt.Sprintf("┌(%d)%s| %s [%s] - %s M:%s, S:%s%s, T:%s,  H:%s, WC:%s, LC:%s, BL:%s, CL:%s, CT:%s %v\n"+
		"\r└> %s %s",

		//Format variables for display:
		count,
		(result.Icon + result.HasErr),
		Purpel+(result.Url)+White,          //URL
		Orange+(result.PayloadClear)+White, //Payload
		reMatch,
		Purpel+(result.Method)+White, //[M] : HTTP Method
		Color_StatusCode(statusCode), //[S] : Status code
		location,                     //[_] : Location header (Redirect)
		Color_Time(result.RespTime),  //[T] : Response time
		Blue+(fmt.Sprint(len(result.Headers)))+White,    //[H] : Header count
		BlueLight+(fmt.Sprint(result.WordCount))+White,  //[WC] : Word count
		BlueLight+(fmt.Sprint(result.LineCount))+White,  //[LC] : Line count
		Pink+(fmt.Sprint(result.BodyByteSize))+White,    //[BL] : Body byte size
		Purpel+(fmt.Sprint(result.ContentLength))+White, //Content length
		Purpel+(result.ContentType)+White,               //Contept type
		diff_status,
		/*========[Behavior information]========*/
		transformation,
		banner_avgDiff)

	fmt.Println(display)
}

/**Display each payload*/
func DisplayPayload(id, wl int, p string) {
	/**Display result of every payload and exit*/
	id++
	fmt.Printf("%v[%d] %s\n", Payload, id, p)
	if wl == id {
		os.Exit(0)
	}
}
