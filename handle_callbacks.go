package dyalpm

import (
	stderrors "errors"
	"runtime"
	"sync"
	"unsafe"

	"github.com/ebitengine/purego"

	alpmerrors "github.com/Jguer/dyalpm/internal/dyerrors"
	"github.com/Jguer/dyalpm/internal/lib"
)

// DownloadEventType corresponds to alpm_download_event_type_t.
type DownloadEventType int32

const (
	DownloadInit      DownloadEventType = 0 // ALPM_DOWNLOAD_INIT
	DownloadProgress  DownloadEventType = 1 // ALPM_DOWNLOAD_PROGRESS
	DownloadRetry     DownloadEventType = 2 // ALPM_DOWNLOAD_RETRY
	DownloadCompleted DownloadEventType = 3 // ALPM_DOWNLOAD_COMPLETED
)

// DownloadEventData is an optional decoded view of the download event payload.
type DownloadEventData interface {
	isDownloadEventData()
}

type DownloadInitData struct {
	Optional bool
}

func (DownloadInitData) isDownloadEventData() {}

type DownloadProgressData struct {
	Downloaded int64
	Total      int64
}

func (DownloadProgressData) isDownloadEventData() {}

type DownloadRetryData struct {
	Resume bool
}

func (DownloadRetryData) isDownloadEventData() {}

type DownloadCompletedData struct {
	Total  int64
	Result int32
}

func (DownloadCompletedData) isDownloadEventData() {}

// DownloadEvent is passed to DownloadCallback.
type DownloadEvent struct {
	Filename string
	Type     DownloadEventType
	Data     DownloadEventData
	RawData  uintptr
}

// DownloadCallback corresponds to alpm_cb_download, with decoded data when possible.
type DownloadCallback func(ev DownloadEvent)

// FetchCallback corresponds to alpm_cb_fetch.
// Return values: 0 success, 1 up-to-date, -1 error.
type FetchCallback func(url string, localpath string, force bool) int

// Event is a minimal wrapper around alpm_event_t*.
// For advanced use, consult alpm.h and decode based on Type and Ptr.
type Event struct {
	Type int32
	Ptr  uintptr
}

// EventCallback corresponds to alpm_cb_event.
type EventCallback func(ev Event)

// Question is a minimal wrapper around alpm_question_t*.
// You can generally answer by setting the second field of the union (an int)
// using SetAnswerInt, since all question structs begin with (type, answer/int).
type Question struct {
	Type int32
	Ptr  uintptr
}

type cQuestionAny struct {
	Type   int32
	Answer int32
}

const (
	maxInt32 = int(^uint32(0) >> 1)
	minInt32 = -maxInt32 - 1
)

func clampIntToInt32(value int) int32 {
	return int32(max(minInt32, min(maxInt32, value)))
}

// SetAnswerInt sets the question's "answer" field. The meaning depends on q.Type.
func (q Question) SetAnswerInt(answer int) {
	if q.Ptr == 0 {
		return
	}
	qa := (*cQuestionAny)(unsafe.Pointer(q.Ptr))
	qa.Answer = clampIntToInt32(answer)
}

// QuestionCallback corresponds to alpm_cb_question.
type QuestionCallback func(q Question)

// ProgressCallback corresponds to alpm_cb_progress.
type ProgressCallback func(progress int32, pkg string, percent int, howmany uint64, current uint64)

// LogCallback corresponds to alpm_cb_log. The libalpm callback's final
// argument is a va_list; it is formatted into msg with vsnprintf before being
// handed to Go, so the callback receives a plain, ready-to-print string.
type LogCallback func(level LogLevel, msg string)

type handleCallbackSet struct {
	mu sync.RWMutex

	download DownloadCallback
	fetch    FetchCallback
	event    EventCallback
	question QuestionCallback
	progress ProgressCallback
	log      LogCallback
}

var (
	callbackSetsMu sync.RWMutex
	callbackSets   = map[uintptr]*handleCallbackSet{} // key: handle pointer or ctx value
)

func getOrCreateCallbackSet(key uintptr) *handleCallbackSet {
	callbackSetsMu.RLock()
	if s, ok := callbackSets[key]; ok {
		callbackSetsMu.RUnlock()
		return s
	}
	callbackSetsMu.RUnlock()

	callbackSetsMu.Lock()
	defer callbackSetsMu.Unlock()
	if s, ok := callbackSets[key]; ok {
		return s
	}
	s := &handleCallbackSet{}
	callbackSets[key] = s
	return s
}

func unregisterCallbackSet(key uintptr) {
	if key == 0 {
		return
	}
	callbackSetsMu.Lock()
	delete(callbackSets, key)
	callbackSetsMu.Unlock()
}

// --- Raw callback accessors (alpm_option_get_*cb / set_*cb) ---

func (h *handle) getCallbackPair(getFn, getCtxFn func(uintptr) uintptr) (cb uintptr, ctx uintptr) {
	if h.ptr == 0 || getFn == nil || getCtxFn == nil {
		return 0, 0
	}
	return getFn(h.ptr), getCtxFn(h.ptr)
}

func (h *handle) setCallback(setFn func(uintptr, uintptr, uintptr) int32, name string, cb uintptr, ctx uintptr) error {
	if h.ptr == 0 {
		return alpmerrors.ErrHandleNull
	}
	if setFn == nil {
		return stderrors.New("missing function: " + name)
	}
	if setFn(h.ptr, cb, ctx) != 0 {
		return alpmerrors.NewError(h.Errno(), "failed to set callback")
	}
	return nil
}

func (h *handle) LogCallback() (cb uintptr, ctx uintptr) {
	return h.getCallbackPair(lib.AlpmOptionGetLogcb, lib.AlpmOptionGetLogcbCtx)
}

func (h *handle) SetLogCallback(cb uintptr, ctx uintptr) error {
	return h.setCallback(lib.AlpmOptionSetLogcb, "alpm_option_set_logcb", cb, ctx)
}

func (h *handle) DownloadCallback() (cb uintptr, ctx uintptr) {
	return h.getCallbackPair(lib.AlpmOptionGetDlcb, lib.AlpmOptionGetDlcbCtx)
}

func (h *handle) SetDownloadCallback(cb uintptr, ctx uintptr) error {
	return h.setCallback(lib.AlpmOptionSetDlcb, "alpm_option_set_dlcb", cb, ctx)
}

func (h *handle) FetchCallback() (cb uintptr, ctx uintptr) {
	return h.getCallbackPair(lib.AlpmOptionGetFetchcb, lib.AlpmOptionGetFetchcbCtx)
}

func (h *handle) SetFetchCallback(cb uintptr, ctx uintptr) error {
	return h.setCallback(lib.AlpmOptionSetFetchcb, "alpm_option_set_fetchcb", cb, ctx)
}

func (h *handle) EventCallback() (cb uintptr, ctx uintptr) {
	return h.getCallbackPair(lib.AlpmOptionGetEventcb, lib.AlpmOptionGetEventcbCtx)
}

func (h *handle) SetEventCallback(cb uintptr, ctx uintptr) error {
	return h.setCallback(lib.AlpmOptionSetEventcb, "alpm_option_set_eventcb", cb, ctx)
}

func (h *handle) QuestionCallback() (cb uintptr, ctx uintptr) {
	return h.getCallbackPair(lib.AlpmOptionGetQuestioncb, lib.AlpmOptionGetQuestioncbCtx)
}

func (h *handle) SetQuestionCallback(cb uintptr, ctx uintptr) error {
	return h.setCallback(lib.AlpmOptionSetQuestioncb, "alpm_option_set_questioncb", cb, ctx)
}

func (h *handle) ProgressCallback() (cb uintptr, ctx uintptr) {
	return h.getCallbackPair(lib.AlpmOptionGetProgresscb, lib.AlpmOptionGetProgresscbCtx)
}

func (h *handle) SetProgressCallback(cb uintptr, ctx uintptr) error {
	return h.setCallback(lib.AlpmOptionSetProgresscb, "alpm_option_set_progresscb", cb, ctx)
}

// --- Go callback helpers (purego.NewCallback + ctx keyed by handle ptr) ---

var (
	getDlCbPtr       = sync.OnceValue(func() uintptr { return purego.NewCallback(dlcbTrampoline) })
	getFetchCbPtr    = sync.OnceValue(func() uintptr { return purego.NewCallback(fetchcbTrampoline) })
	getEventCbPtr    = sync.OnceValue(func() uintptr { return purego.NewCallback(eventcbTrampoline) })
	getQuestionCbPtr = sync.OnceValue(func() uintptr { return purego.NewCallback(questioncbTrampoline) })
	getProgressCbPtr = sync.OnceValue(func() uintptr { return purego.NewCallback(progresscbTrampoline) })
	getLogCbPtr      = sync.OnceValue(func() uintptr { return purego.NewCallback(logcbTrampoline) })
)

// C data structures for download events (matching alpm.h)
type cDownloadInit struct {
	Optional int32
}

type cDownloadProgress struct {
	Downloaded int64
	Total      int64
}

type cDownloadRetry struct {
	Resume int32
}

type cDownloadCompleted struct {
	Total  int64
	Result int32
	_      [4]byte // padding for 8-byte alignment
}

func dlcbTrampoline(_ purego.CDecl, ctx uintptr, filename uintptr, event int32, data uintptr) {
	callbackSetsMu.RLock()
	set := callbackSets[ctx]
	callbackSetsMu.RUnlock()
	if set == nil {
		return
	}

	set.mu.RLock()
	cb := set.download
	set.mu.RUnlock()
	if cb == nil {
		return
	}

	evType := DownloadEventType(event)
	decoded, raw := decodeDownloadEventData(evType, data)
	cb(DownloadEvent{
		Filename: lib.PtrToString(filename),
		Type:     evType,
		Data:     decoded,
		RawData:  raw,
	})
}

func decodeDownloadEventData(t DownloadEventType, data uintptr) (DownloadEventData, uintptr) {
	if data == 0 {
		return nil, 0
	}
	switch t {
	case DownloadInit:
		v := (*cDownloadInit)(unsafe.Pointer(data))
		return DownloadInitData{Optional: v.Optional != 0}, data
	case DownloadProgress:
		v := (*cDownloadProgress)(unsafe.Pointer(data))
		return DownloadProgressData{Downloaded: v.Downloaded, Total: v.Total}, data
	case DownloadRetry:
		v := (*cDownloadRetry)(unsafe.Pointer(data))
		return DownloadRetryData{Resume: v.Resume != 0}, data
	case DownloadCompleted:
		v := (*cDownloadCompleted)(unsafe.Pointer(data))
		return DownloadCompletedData{Total: v.Total, Result: v.Result}, data
	default:
		return nil, data
	}
}

func fetchcbTrampoline(_ purego.CDecl, ctx uintptr, url uintptr, localpath uintptr, force int32) int32 {
	callbackSetsMu.RLock()
	set := callbackSets[ctx]
	callbackSetsMu.RUnlock()
	if set == nil {
		return -1
	}

	set.mu.RLock()
	cb := set.fetch
	set.mu.RUnlock()
	if cb == nil {
		return -1
	}

	res := cb(lib.PtrToString(url), lib.PtrToString(localpath), force != 0)
	return clampIntToInt32(res)
}

func eventcbTrampoline(_ purego.CDecl, ctx uintptr, eventPtr uintptr) {
	callbackSetsMu.RLock()
	set := callbackSets[ctx]
	callbackSetsMu.RUnlock()
	if set == nil {
		return
	}

	set.mu.RLock()
	cb := set.event
	set.mu.RUnlock()
	if cb == nil || eventPtr == 0 {
		return
	}

	typ := *(*int32)(unsafe.Pointer(eventPtr))
	cb(Event{Type: typ, Ptr: eventPtr})
}

func questioncbTrampoline(_ purego.CDecl, ctx uintptr, questionPtr uintptr) {
	callbackSetsMu.RLock()
	set := callbackSets[ctx]
	callbackSetsMu.RUnlock()
	if set == nil {
		return
	}

	set.mu.RLock()
	cb := set.question
	set.mu.RUnlock()
	if cb == nil || questionPtr == 0 {
		return
	}

	typ := *(*int32)(unsafe.Pointer(questionPtr))
	cb(Question{Type: typ, Ptr: questionPtr})
}

func progresscbTrampoline(_ purego.CDecl, ctx uintptr, progress int32, pkg uintptr, percent int32, howmany uintptr, current uintptr) {
	callbackSetsMu.RLock()
	set := callbackSets[ctx]
	callbackSetsMu.RUnlock()
	if set == nil {
		return
	}

	set.mu.RLock()
	cb := set.progress
	set.mu.RUnlock()
	if cb == nil {
		return
	}

	cb(progress, lib.PtrToString(pkg), int(percent), uint64(howmany), uint64(current))
}

// logcbTrampoline bridges libalpm's alpm_cb_log. Its final C argument is a
// va_list, which on every supported platform is passed to the callee as a
// pointer (an array that decays to a pointer on x86-64 SysV, or a >16-byte
// struct passed indirectly on arm64 AAPCS64). That same pointer is forwarded
// verbatim to vsnprintf to format the message, mirroring the C idiom of
// va_start followed by a v*printf call.
func logcbTrampoline(_ purego.CDecl, ctx uintptr, level int32, format uintptr, ap uintptr) {
	callbackSetsMu.RLock()
	set := callbackSets[ctx]
	callbackSetsMu.RUnlock()
	if set == nil {
		return
	}

	set.mu.RLock()
	cb := set.log
	set.mu.RUnlock()
	if cb == nil || format == 0 || lib.LibcVsnprintf == nil {
		return
	}

	buf := make([]byte, 4096)
	n := lib.LibcVsnprintf(uintptr(unsafe.Pointer(&buf[0])), uintptr(len(buf)), format, ap)
	runtime.KeepAlive(buf)
	if n <= 0 {
		return
	}

	length := min(int(n),
		// vsnprintf returns the would-be length on truncation
		len(buf)-1)

	if level < 0 {
		return
	}
	cb(LogLevel(uint(level)), string(buf[:length]))
}

func (h *handle) setCallbackFunc(
	cbPtr uintptr,
	shouldClear bool,
	update func(*handleCallbackSet),
	setFn func(uintptr, uintptr) error,
) error {
	if h.ptr == 0 {
		return alpmerrors.ErrHandleNull
	}

	set := getOrCreateCallbackSet(h.ptr)
	set.mu.Lock()
	update(set)
	set.mu.Unlock()

	if shouldClear {
		return setFn(0, 0)
	}
	return setFn(cbPtr, h.ptr)
}

func (h *handle) SetDownloadCallbackFunc(cb DownloadCallback) error {
	if cb == nil {
		return h.setCallbackFunc(0, true, func(set *handleCallbackSet) {
			set.download = nil
		}, h.SetDownloadCallback)
	}
	return h.setCallbackFunc(getDlCbPtr(), false, func(set *handleCallbackSet) {
		set.download = cb
	}, h.SetDownloadCallback)
}

func (h *handle) SetFetchCallbackFunc(cb FetchCallback) error {
	if cb == nil {
		return h.setCallbackFunc(0, true, func(set *handleCallbackSet) {
			set.fetch = nil
		}, h.SetFetchCallback)
	}
	return h.setCallbackFunc(getFetchCbPtr(), false, func(set *handleCallbackSet) {
		set.fetch = cb
	}, h.SetFetchCallback)
}

func (h *handle) SetEventCallbackFunc(cb EventCallback) error {
	if cb == nil {
		return h.setCallbackFunc(0, true, func(set *handleCallbackSet) {
			set.event = nil
		}, h.SetEventCallback)
	}
	return h.setCallbackFunc(getEventCbPtr(), false, func(set *handleCallbackSet) {
		set.event = cb
	}, h.SetEventCallback)
}

func (h *handle) SetQuestionCallbackFunc(cb QuestionCallback) error {
	if cb == nil {
		return h.setCallbackFunc(0, true, func(set *handleCallbackSet) {
			set.question = nil
		}, h.SetQuestionCallback)
	}
	return h.setCallbackFunc(getQuestionCbPtr(), false, func(set *handleCallbackSet) {
		set.question = cb
	}, h.SetQuestionCallback)
}

func (h *handle) SetProgressCallbackFunc(cb ProgressCallback) error {
	if cb == nil {
		return h.setCallbackFunc(0, true, func(set *handleCallbackSet) {
			set.progress = nil
		}, h.SetProgressCallback)
	}
	return h.setCallbackFunc(getProgressCbPtr(), false, func(set *handleCallbackSet) {
		set.progress = cb
	}, h.SetProgressCallback)
}

func (h *handle) SetLogCallbackFunc(cb LogCallback) error {
	if cb == nil {
		return h.setCallbackFunc(0, true, func(set *handleCallbackSet) {
			set.log = nil
		}, h.SetLogCallback)
	}
	return h.setCallbackFunc(getLogCbPtr(), false, func(set *handleCallbackSet) {
		set.log = cb
	}, h.SetLogCallback)
}
