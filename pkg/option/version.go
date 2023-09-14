package option

// Show the current and latest version of FireFly and exit
/* func VersionFirefly() {
	resp, err := http.Get(info.VERSION_URL)
	verbose.Error(err)

	bodyBytes, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	verbose.Error(err)

	//Extract version_new from github repository ():
	version_new := string(bodyBytes)
	if true { //ok, _ := regexp.MatchString(`^v(\d\.\d{1,2})$`, version); ok {
		if false {
			fmt.Println(d.ICON.PLUS, "Version %s (Newest)", info.VERSION)
		} else {
			fmt.Println(d.ICON.AWARE, "Version "+info.VERSION+", Newest: "+newestVersion(version_new)+"\n    It is highly recommended to update Firefly as many fixes and features have been added/improved.")
		}
	}
	os.Exit(0)
}

func newestVersion(v string) string {
	return fmt.Sprint(d.COLOR.ORANGE + (v) + d.COLOR.WHITE)
} */
