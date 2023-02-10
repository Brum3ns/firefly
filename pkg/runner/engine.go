package runner

import (
	"fmt"
	"os"

	"github.com/Brum3ns/firefly/pkg/functions"
	G "github.com/Brum3ns/firefly/pkg/functions/globalVariables"
	"github.com/Brum3ns/firefly/pkg/storage"
)

type EiTransformation struct {
	OK                     bool
	Transformation         map[string]string
	Transformation_display string

	Status bool
	ID     int
}

type EiDifference struct {
	OK        bool
	Amount    int
	AvgAmount map[string]int
	Diff      map[string]map[string]string //Target -> Payload, Amount, Diff(s)
	Reflect   map[string]map[int][]string  //Name, Amount, Value(s)
	Dynamic   map[string]map[int][]string  //Name, Amount, Value(s)

	WordCount int
	LineCount int

	Status bool
	ID     int
}

type EiPatterns struct {
	OK       bool
	Patterns map[string]map[int][]string //Type, Amount, Pattern(s)

	Status bool
	ID     int
}

type EiErrors struct {
	OK         bool
	Errs       map[string]map[int][]string //Type, Amount, Error(s)
	ErrsAmount map[string]int

	Icon     string //DELETE? Replaced with "OK"
	HasErr   string
	CountErr int

	Status bool
	ID     int
}

/* type EiCalc struct {
	Calculations map[string]map[int][]string //Type, Amount, Calculation(s)

	LineCount int
	WordCount int

	Status bool
	ID     int
} */

type EiDynamic struct {
	Calculations map[string]map[int][]string //Type, Amount, Calculation(s)

	Status bool
	ID     int
}

//[x4] - [chan]nels amount (Aka: tasks that will run):
var EiProcessCount = 4

/**Engine that do all the hevy enumiration tasks */
func (r *Runner) Engine(RespID int, RespResult storage.Response, EiResults chan<- storage.Collection) {
	/*=============================================================\
	|						FireFly Engine						   |
	\==============================================================/
	+---------------------------------------------------------------------------------+
	| Action | Take response result from target & memory storage setup 			      |
	+---------------------------------------------------------------------------------+
	| Usage  | Use the collected data for later processes & update the output result. |
	+---------------------------------------------------------------------------------+
	| Return | Return gathered data to update the storage, verbose & output	   		  |
	+---------------------------------------------------------------------------------+
	* Engine spinning up threads & goroutines (Fuzz process)
	* -> [Request]
	* -> [Analyze]
	* -> [Verbose]
	* -> [Output]
	*/

	//Check filter result and decided if FireFly should continue (ignore verify requests)
	if !r.verify && (len(G.Lst_mFilter) > 0 || len(G.Lst_mMatch) > 0) {
		if r.Filter(RespResult) {
			EiResults <- storage.Collection{
				Skip: true,
			}
			return
		}
	}

	//Setup memory storage (structs)
	var (
		eiT = &EiTransformation{}
		//eiC = &EiCalc{}
		eiD = &EiDifference{}
		eiE = &EiErrors{}
		eiP = &EiPatterns{}
	)

	//Setup [chan]nels for each task & make a "defer" that close the channel once it's no longer needed & done:
	c_transformation := make(chan EiTransformation, 1)
	c_difference := make(chan EiDifference, 1)
	c_errors := make(chan EiErrors, 1)
	c_patterns := make(chan EiPatterns, 1)
	//c_calcs := make(chan EiCalc, 1)

	//Close all [chan]nels when they are finished:
	defer close(c_transformation)
	defer close(c_difference)
	defer close(c_errors)
	defer close(c_patterns)
	//defer close(c_calcs)

	//Start internal tasks:
	go func(rslt chan<- EiTransformation) {
		eiT = r.Transformation(RespID, RespResult)
		rslt <- *eiT
	}(c_transformation)

	go func(rslt chan<- EiDifference) {
		eiD = r.Difference(RespID, RespResult)
		rslt <- *eiD

	}(c_difference)

	go func(rslt chan<- EiPatterns) {
		eiP = r.Pattern(RespID)
		rslt <- *eiP
	}(c_patterns)

	go func(rslt chan<- EiErrors) {
		eiE = r.Error(RespID, RespResult)
		rslt <- *eiE
	}(c_errors)

	//DELETE?
	/* go func(rslt chan<- EiCalc) {
		eiC = r.Calcs(RespID, RespResult)
		rslt <- *eiC
	}(c_calcs) */

	//Monitor all the goroutine processes at once and control each status (stability):
	for i := 0; i < EiProcessCount; i++ {
		select {

		//[Task] - Calculation of response data
		/* case _, ok := <-c_calcs:
		eiC.Status = ok */

		//[Task] - Error detection
		case _, ok := <-c_errors:
			eiE.Status = ok

		//[Task] - Pattern gathering
		case _, ok := <-c_patterns:
			eiP.Status = ok

		//[Task] - Differences detection
		case _, ok := <-c_difference:
			eiD.Status = ok

		//[Task] - Transformation detection
		case _, ok := <-c_transformation:
			eiT.Status = ok
		}
	}

	done := struct {
		l_id     []int
		l_status []bool
	}{
		l_id:     []int{eiE.ID, eiP.ID, eiD.ID, eiT.ID},
		l_status: []bool{eiE.Status, eiP.Status, eiD.Status, eiT.Status},
	}

	ok, msg := Check(RespID, done)
	functions.IFFail(msg)

	if ok {
		var icon string

		//Collect all data from all tasks and return it to the runner for future data handling:
		EiResults <- storage.Collection{
			ThreadID:     RespResult.ThreadID,
			ID:           RespID,
			Tag:          RespResult.Tag,
			Payload:      RespResult.Payload,
			PayloadClear: functions.PayloadClearPattern(RespResult.Payload),

			Url:          RespResult.Url,
			UrlNoPayload: RespResult.UrlNoPayload,
			Method:       RespResult.Method,
			StatusCode:   RespResult.Status,

			ContentLength: RespResult.ContentLength,
			ContentType:   RespResult.ContentType,
			Headers:       RespResult.Headers,

			RespTime: RespResult.RespTime,

			Body:         RespResult.Body,
			BodyByteSize: RespResult.BodyByteSize,

			WordCount: eiD.WordCount,
			LineCount: eiD.LineCount,

			AvgAmountDiff: eiD.AvgAmount,
			//RespDiff:      eiD.Diff,
			RespErr:       eiE.Errs,
			RespErrAmount: eiE.ErrsAmount,

			Transformation: eiT.Transformation,
			Tfmt_display:   eiT.Transformation_display,
			Tfmt_ok:        eiT.OK,

			//Filter: m_Filter,

			//HeadersMatch: ,
			//RegexMatch: ,
			//Valid: ,

			//Technologies:     ,
			//RegexMatch:     , */

			Icon:   icon,
			HasErr: eiE.HasErr,

			ErrMsg: RespResult.ErrMsg,
			Errors: RespResult.Error,
			Skip:   RespResult.Skip,
			//TaskStatus: m_Ready,
			VerifyProcess: r.verify,
			Status:        true,
		}
	} else {
		fmt.Println(functions.IFFail("eiID"))
		os.Exit(0)
	}

}

/**Check so all tasks are done*/
func Check(id int, data struct {
	l_id     []int
	l_status []bool
}) (bool, string) {

	msg := "taskcheck"
	if len(data.l_id) != len(data.l_status) {
		return false, msg
	}

	for i := 0; i < len(data.l_id); i++ {
		if data.l_id[i] != id {
			return false, msg
		} else if !data.l_status[i] {
			return false, msg
		}
	}
	return true, ""
}
