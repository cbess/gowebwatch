package web

import (
	"net/http"
	"time"
)

// Returns a new Watcher
func New(url string, checkNow bool) *WebWatcher {
	watcher := &WebWatcher{URL: url}
	watcher.active = false
	watcher.StatusCode = -1
	watcher.Interval = 1 // 1 second default
	watcher.watchChan = make(chan StatusItem)
	watcher.checkNow = checkNow
	return watcher
}

// The struct that is sent over the channel to report
// the status of the watcher
type StatusItem struct {
	StatusCode int
	URL        string
}

func (s StatusItem) IsOk() bool {
	return s.StatusCode == 200
}

type Watcher interface {
	Watch()
	Stop()
	IsOk() bool
	IsActive() bool
}

type WebWatcher struct {
	URL    string
	active bool
	// The status code from the last check
	// Value -1, indicates the url has not been watched
	StatusCode int
	watchChan  chan StatusItem
	// The interval between url checks. Defaults to 1.
	Interval int
	// true if checking the web as soon as possible
	checkNow bool
}

func (w *WebWatcher) Watch() chan StatusItem {
	go w.startWatch()
	return w.GetChan()
}

func (w *WebWatcher) GetChan() chan StatusItem {
	if w.watchChan == nil {
		panic("No watcher chan.")
	}

	return w.watchChan
}

func (w *WebWatcher) SetChan(wChan chan StatusItem) {
	w.watchChan = wChan
}

// Starts watching the url on an interval, should be used as
// a go routine
func (w *WebWatcher) startWatch() {
	w.active = true

	// do the initial check?
	if w.checkNow {
		w.performWatchCheck()
	}

	// now check every interval
	duration := time.Duration(w.Interval) * time.Second
	// tick := time.Tick(duration)
	for {
		time.Sleep(duration)

		w.performWatchCheck()
	}
}

// Performs the actual watch url check, then sends on the chan
func (w *WebWatcher) performWatchCheck() {
	w.StatusCode = w.checkURL()
	sItem := w.getStatusItem()

	// send on chan
	w.watchChan <- sItem
}

// Checks the url, returning the status code from it
func (w *WebWatcher) checkURL() int {
	resp, err := http.Get(w.URL)
	if err != nil {
		return 0
	}

	return resp.StatusCode
}

func (w WebWatcher) getStatusItem() StatusItem {
	return StatusItem{StatusCode: w.StatusCode, URL: w.URL}
}

func (w *WebWatcher) Stop() {
	w.active = false
}

// Returns true while the Watcher is watching the url
func (w *WebWatcher) IsActive() bool {
	return w.active
}

// Returns true if the watched resource is available/ok
func (w *WebWatcher) IsOk() bool {
	return w.getStatusItem().IsOk()
}
