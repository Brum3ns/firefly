package functions

import (
	"fmt"
	"io/ioutil"
	"strings"

	st "github.com/Brum3ns/firefly/pkg/storage"
)

/*
* Create and add all "[G]rep" wordlists and add them to each unique List.
* The reason why is because will later be used for payload adjustment.
*
* [Example Senario]
*
* If FireFly detects an error in the response that was a match from
* the list -> [wl.LstG_errorPython]. It will remember that specific
* hostname that gave an python error in "x" insert point. It will
* then adjust some payloads specific for that parameter and add more
* deep testing specific for python.
 */
func LstGSetup(wl *st.Wordlists) {
	/** Setup all lists within it's own category
	* Add all lists into the storage (storage.Wordlists struct{})
	 */

	typs := []string{
		"cache",
		"cms",
		"dbms",
		"error_mariadb",
		"error_php",
		"error_postgresql",
		"error_python",
		"errors",
		"extension",
		"filename",
		"infoleak",
		"pattern",
		"patternJavascript",
		"technology",
		"webservice",
		"other",
	}

	for i := 0; i < len(typs); i++ {
		t := typs[i]
		f := fmt.Sprintf("Grep_%s.txt", t)
		wl.MG_patterns[t] = FileToLst(st.StaticFilePath + f)

	}
}

func WordlistFolder_Declare(replace, encode, tamper string, wordlistFolder string, wl *st.Wordlists) {
	/** Collect wordlist files from a selected folder
	* Return the collected wordlists
	 */

	//Read wordlists filename in given/default folder:
	wl_file, err := ioutil.ReadDir(wordlistFolder)
	IFError("p", err)

	//If it's a file then add it to an array that holds the wordlist filenames for later use:
	for _, wl_file := range wl_file {
		if !wl_file.IsDir() {

			wl_pathToFile := string(wordlistFolder + wl_file.Name())

			//Check if the technique for the wordlist should be used:
			attackMethod := wl_pathToFile[strings.LastIndex(wl_pathToFile, "/")+1:]
			split_toTag := strings.Split(strings.Split(attackMethod, "_")[0], "/")
			checkUse := split_toTag[len(split_toTag)-1]

			if InLst(wl.UseTechniques, checkUse) {
				StringToArray_Attack(replace, encode, tamper, wl_pathToFile, wl, true)
				wl.Lst_default = append(wl.Lst_default, wl_pathToFile)
			}
		}
	}
}
