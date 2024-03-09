package request

import (
	"net/http"
	"time"

	"github.com/Brum3ns/firefly/pkg/waitgroup"
)

type Handler struct {
	jobAmount   int
	Worker      worker
	TaskStorage *TaskStorage
	WaitGroup   waitgroup.WaitGroup

	stop        chan bool
	JobReceived chan int
	JobQueue    chan RequestSettings
	WorkerPool  chan chan RequestSettings
	HandlerSettings
}

// HandlerSettings holds all the settings which will be used within the primary Handler structure
type HandlerSettings struct {
	VerifyMode  bool
	Delay       int
	Threads     int
	Client      *http.Client
	RequestBase RequestBase
}

// worker represents the worker that executes the job
type worker struct {
	Delay      int
	client     *http.Client
	jobChannel chan RequestSettings
	workerPool chan chan RequestSettings
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
func NewHandler(settings HandlerSettings) Handler { // httpclient *http.Client, task *TaskStorage, threads int, delay int, verifyMode bool) *Handler {
	return Handler{
		HandlerSettings: settings,
		stop:            make(chan bool),
		JobReceived:     make(chan int),
		JobQueue:        make(chan RequestSettings),
		WorkerPool:      make(chan chan RequestSettings, settings.Threads),
	}
}

// Start all the workers and assign tasks (jobs) to the request workers
// The process will start listen for job and stop once all job sent is done.
// !Note : (To set the job amount to let the process know when to stop use the method "SetJobAmount")
func (h *Handler) Run(listener chan<- Result) {
	var result = make(chan Result)

	//Start the amount of workers related to the amount of given threads:
	for i := 0; i < h.Threads; i++ {
		h.Worker = newRequestWorker(h.Client, h.WorkerPool, h.Delay)
		go h.Worker.spawnRequestWorker(result)
	}

	//Listen for new jobs from the queue and send it to the job channel for the workers to handle it:
	go func() {
		for {
			select {
			case job := <-h.JobQueue:
				go func(job RequestSettings) {
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
	//Wait until all workers have provided the result for each job given, then send a signal that the core process is done:
	<-h.stop
}

// Add a job process to the handler
func (h *Handler) AddJob(job RequestSettings) {
	h.WaitGroup.Add(1)
	h.jobAmount++
	job.RequestId = h.jobAmount
	h.JobQueue <- job
}

// Send a stop signal to the handler
func (h *Handler) Stop() {
	h.stop <- true
}

// Get the amount of active processes that are within the process
func (h *Handler) GetInProcess() int {
	return h.WaitGroup.GetCount()
}

// Get the amount of given jobs
func (h *Handler) GetJobAmount() int {
	return h.jobAmount
}

// Create a new request worker
func newRequestWorker(client *http.Client, workerPool chan chan RequestSettings, delay int) worker {
	return worker{
		Delay:      delay,
		client:     client,
		workerPool: workerPool,
		jobChannel: make(chan RequestSettings),
	}
}

// start the request worker
func (w worker) spawnRequestWorker(result chan Result) {
	for {
		// Add the current worker into the worker queue:
		w.workerPool <- w.jobChannel

		RequestJob := <-w.jobChannel
		time.Sleep(time.Duration(w.Delay) * time.Millisecond)
		result <- Request(w.client, RequestJob)
	}
}
