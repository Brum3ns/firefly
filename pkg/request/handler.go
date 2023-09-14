package request

import (
	"crypto/md5"
	"encoding/hex"
	"net/http"
	"sync"
	"time"

	"github.com/Brum3ns/firefly/pkg/payloads"
)

type Handler struct {
	VerifyMode  bool
	Delay       int
	Threads     int
	Worker      worker
	Client      *http.Client
	TaskStorage *TaskStorage
	WaitGroup   sync.WaitGroup

	JobQueue   chan RequestProperties
	WorkerPool chan chan RequestProperties
}

// worker represents the worker that executes the job
type worker struct {
	Delay      int
	client     *http.Client
	jobChannel chan RequestProperties
	workerPool chan chan RequestProperties
}

type TaskStorage struct {
	URLs            []string
	Methods         []string
	Schemes         []string
	Payloads        map[string][]string //(tag|wordlist)
	PostData        string
	InsertPoint     string
	Headers         [][2]string
	RandomUserAgent bool
}

// Start the handler for the workers by giving the tasks to preform and the amount of workers.
func NewHandler(ClientProperties *ClientProperties, task *TaskStorage, threads int, delay int, verifyMode bool) *Handler {
	return &Handler{
		Delay:       delay,
		VerifyMode:  verifyMode,
		Threads:     threads,
		TaskStorage: task,
		JobQueue:    make(chan RequestProperties),
		WorkerPool:  make(chan chan RequestProperties, threads),
		Client:      setClient(ClientProperties),
	}
}

// Start all the workers and assign tasks (jobs) to the request workers
func (h *Handler) Run(listener chan<- Result, RequestAmount chan<- int, done chan<- bool) {
	//Start the amount of workers related to the amount of given threads:
	var result = make(chan Result)
	for i := 0; i < h.Threads; i++ {
		h.Worker = newRequestWorker(h.Client, h.WorkerPool, h.Delay)
		h.Worker.spawnRequestWorker(result)
	}

	//Listen for new jobs from the queue and send it to the job channel for the workers to handle it:
	go func() {
		for {
			select {
			case job := <-h.JobQueue:
				go func(job RequestProperties) {
					//Get an available job channel from any worker:
					jobChannel := <-h.WorkerPool

					//Give the available worker the job:
					jobChannel <- job
				}(job)

				//Listen for result from any Worker, if a result is recived, then send it to the listener [chan]nel:
			case r := <-result:
				listener <- r
				h.WaitGroup.Done()
			}
		}
	}()
	RequestAmount <- h.appendJobs()
	//Wait until all workers have provided the result for each job given, then send a signal that the core process is done:
	h.WaitGroup.Wait()

	//Send a signal back that all tasks where completed
	done <- true
}

// Add new jobs (tasks) to be performed in the workgroup
// Return the amount of job that was given
func (h *Handler) appendJobs() int {
	h.WaitGroup.Add(1) //<-Adds the core process to avoid the gorutines finish before all jobs is given
	var (
		requestId = 0
		jobAmount = 0
	)
	for _, url := range h.TaskStorage.URLs {
		for _, method := range h.TaskStorage.Methods {
			hash := md5.Sum([]byte(method + url))
			targetHashId := hex.EncodeToString(hash[:])

			for _, tag := range payloads.TAGS {
				//Check if we should adapt to "behavior verification mode":
				if (h.VerifyMode && tag != payloads.TAG_VERIFY) || (!h.VerifyMode && tag == payloads.TAG_VERIFY) {
					continue
				}

				wordlist := h.TaskStorage.Payloads[tag]
				for _, payload := range wordlist {
					requestId++
					//Prepare the request by inserting the current payload into the request:
					//Note : (Some variables given will be modified)
					req := NewInsert(
						&Insert{
							payload: payload,
							keyword: h.TaskStorage.InsertPoint,
						},
						&RequestProperties{
							TargetHashId:    targetHashId,
							RequestId:       requestId,
							Tag:             tag,
							Payload:         payload,
							URL:             url,
							URLOriginal:     url,
							Method:          method,
							HeadersOriginal: h.TaskStorage.Headers,
							PostBody:        h.TaskStorage.PostData,
							RandomUserAgent: h.TaskStorage.RandomUserAgent,
						},
					)

					//Append job to waitgroup
					h.WaitGroup.Add(1)
					h.JobQueue <- *req
					jobAmount++
				}
			}
		}
	}
	//Done with the job task process
	h.WaitGroup.Done()
	return jobAmount
}

// Create a new request worker
func newRequestWorker(client *http.Client, workerPool chan chan RequestProperties, delay int) worker {
	return worker{
		Delay:      delay,
		client:     client,
		workerPool: workerPool,
		jobChannel: make(chan RequestProperties),
	}
}

// start the request worker
func (w worker) spawnRequestWorker(result chan Result) {
	go func() {
		for {
			// Add the current worker into the worker queue:
			w.workerPool <- w.jobChannel

			select {
			case RequestJob := <-w.jobChannel:
				time.Sleep(time.Duration(w.Delay) * time.Second)
				result <- w.request(RequestJob)
			}
		}
	}()
}
