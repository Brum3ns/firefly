package runner

/* func (r *Runner) AnalyzeBasic() {

	//Calculate the average normal response time:
	r.vresp.VR_RespTime = math.Round((calculate.Sum(r.vresp.VR_AvgRespTime)/4)*100) / 100

	//Check histeric service quickly:
	var serverCheck = ""
	if len(r.vresp.Lstv_headers["server"]) > 1 {
		serverCheck = r.vresp.Lstv_headers["server"]
		fmt.Println(d.Success, "Heuristic checks shows possible service [server]:", r.vresp.Lstv_headers["server"])

	} else if len(r.vresp.Lstv_headers["x-powered-by"]) > 1 {
		serverCheck = r.vresp.Lstv_headers["x-powered-by"]
		fmt.Println(d.Success, "Heuristic checks shows possible [x-powered-by]:", r.vresp.Lstv_headers["x-powered-by"])
	}

	//If a service is found. Check if FireFly got known attack techniques for it
	if len(serverCheck) > 1 {

		for _, service := range r.wordlist.Lst_Experience {
			if strings.Contains(serverCheck, service) {
				fmt.Println(d.Info, fmt.Sprintf("FireFly has known attack techniques for: %s - (%s)", service, serverCheck))

			}
		}
	}

	//Check dynamic content from the test requests:
	if len(r.vresp.Lstv_responeData) <= 1 {
		fmt.Println(d.Info, "Target appear to be static or we dealing with a cached response If the response is cached some cache techniques should be used in major of the test cases.")

	} else {
		fmt.Println(d.Warning, "Target appears to be dynamic. Manual checks are extra important")
	}
}

func AnalyzeDefault(lst []string, result storage.Response, vRes *storage.VerifyResponse) {

	//Header change detection <===========[START HERE]
	/*for header, value := range vRes.Lstv_headers {
		for _, cheader := range lst {
			fmt.Println(cheader, "::", header, "|", value)
		}
	}* /

	//Response Time:
	if result.R_RespTime >= vRes.VR_RespTime+1 {
		fmt.Println(d.Info, "Response time *over* average")

	} else if result.R_RespTime <= vRes.VR_RespTime/2 {
		fmt.Println(d.Info, "Response times *under* average")
	}

	if result.R_Status == 429 {
		fmt.Println(d.Info, "Seems like FireFly getting rate-limited (Blocked)")
	}
}
*/
