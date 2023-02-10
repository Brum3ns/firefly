package technique

//[TODO] Version => 1.1 - 2.0

func Detect() {
	/** Run all detection scans
	* 	Result will later be used to adapt fuzzing techniques to avoid false positive/negatives
	 */
}

func RevProxy() (bool, []string) {
	/** Detect reverse proxy used
	* Use common fuzzing payloads that use to behave different
	 */
	var (
		lst []string
		s   = false
	)
	return s, lst
}

func CMS() (bool, string) {
	/** Detect CMS
	* Use common directory and endpoints to detect the CMS
	 */
	var (
		cms string
		s   = false
	)
	return s, cms
}

func Java() bool {
	return true
}

func PHP() bool {
	return true
}

func Nginx() bool {
	return true
}

func Apache() bool {
	return true
}

func Envoy() bool {
	return true
}

func IIS() bool {
	return true
}
