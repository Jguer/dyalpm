//go:build integration

package dyalpm

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"
)

func TestHandleCallbacks_RegistrationAndClear(t *testing.T) {
	h := mustInitializeTestHandle(t)

	if err := h.SetDownloadCallbackFunc(func(ev DownloadEvent) {
		// no-op
	}); err != nil {
		t.Fatalf("failed to set download callback: %v", err)
	}
	if cb, ctx := h.DownloadCallback(); cb == 0 || ctx == 0 {
		t.Fatalf("download callback did not register: cb=%d ctx=%d", cb, ctx)
	}
	if err := h.SetDownloadCallbackFunc(nil); err != nil {
		t.Fatalf("failed to clear download callback: %v", err)
	}
	if cb, ctx := h.DownloadCallback(); cb != 0 || ctx != 0 {
		t.Fatalf("expected download callback to clear: cb=%d ctx=%d", cb, ctx)
	}

	if err := h.SetFetchCallbackFunc(func(url string, localpath string, force bool) int {
		return 0
	}); err != nil {
		t.Fatalf("failed to set fetch callback: %v", err)
	}
	if cb, ctx := h.FetchCallback(); cb == 0 || ctx == 0 {
		t.Fatalf("fetch callback did not register: cb=%d ctx=%d", cb, ctx)
	}
	if err := h.SetFetchCallbackFunc(nil); err != nil {
		t.Fatalf("failed to clear fetch callback: %v", err)
	}
	if cb, ctx := h.FetchCallback(); cb != 0 || ctx != 0 {
		t.Fatalf("expected fetch callback to clear: cb=%d ctx=%d", cb, ctx)
	}

	if err := h.SetEventCallbackFunc(func(ev Event) {
		// no-op
	}); err != nil {
		t.Fatalf("failed to set event callback: %v", err)
	}
	if cb, ctx := h.EventCallback(); cb == 0 || ctx == 0 {
		t.Fatalf("event callback did not register: cb=%d ctx=%d", cb, ctx)
	}
	if err := h.SetEventCallbackFunc(nil); err != nil {
		t.Fatalf("failed to clear event callback: %v", err)
	}
	if cb, ctx := h.EventCallback(); cb != 0 || ctx != 0 {
		t.Fatalf("expected event callback to clear: cb=%d ctx=%d", cb, ctx)
	}

	if err := h.SetQuestionCallbackFunc(func(q Question) {
		// no-op
	}); err != nil {
		t.Fatalf("failed to set question callback: %v", err)
	}
	if cb, ctx := h.QuestionCallback(); cb == 0 || ctx == 0 {
		t.Fatalf("question callback did not register: cb=%d ctx=%d", cb, ctx)
	}
	if err := h.SetQuestionCallbackFunc(nil); err != nil {
		t.Fatalf("failed to clear question callback: %v", err)
	}
	if cb, ctx := h.QuestionCallback(); cb != 0 || ctx != 0 {
		t.Fatalf("expected question callback to clear: cb=%d ctx=%d", cb, ctx)
	}

	if err := h.SetProgressCallbackFunc(func(progress int32, pkg string, percent int, howmany uint64, current uint64) {
		// no-op
	}); err != nil {
		t.Fatalf("failed to set progress callback: %v", err)
	}
	if cb, ctx := h.ProgressCallback(); cb == 0 || ctx == 0 {
		t.Fatalf("progress callback did not register: cb=%d ctx=%d", cb, ctx)
	}
	if err := h.SetProgressCallbackFunc(nil); err != nil {
		t.Fatalf("failed to clear progress callback: %v", err)
	}
	if cb, ctx := h.ProgressCallback(); cb != 0 || ctx != 0 {
		t.Fatalf("expected progress callback to clear: cb=%d ctx=%d", cb, ctx)
	}
}

func TestHandleCallbacks_FetchPkgURL_HttpServer(t *testing.T) {
	requireDownloaderCapability(t)

	h := mustInitializeTestHandle(t)
	if err := h.SetCacheDirs([]string{t.TempDir()}); err != nil {
		t.Fatalf("failed to set cache dirs: %v", err)
	}

	var downloadCalls int32
	var fetchCalls int32
	var progressCalls int32

	if err := h.SetDownloadCallbackFunc(func(ev DownloadEvent) {
		atomic.AddInt32(&downloadCalls, 1)
	}); err != nil {
		t.Fatalf("failed to set download callback: %v", err)
	}
	defer func() {
		_ = h.SetDownloadCallbackFunc(nil)
	}()

	if err := h.SetFetchCallbackFunc(func(url string, localpath string, force bool) int {
		atomic.AddInt32(&fetchCalls, 1)
		resp, err := http.Get(url)
		if err != nil {
			return -1
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return -1
		}

		payload, err := io.ReadAll(resp.Body)
		if err != nil {
			return -1
		}

		destination := localpath
		if info, err := os.Stat(localpath); err == nil && info.IsDir() {
			destination = filepath.Join(localpath, filepath.Base(url))
		}

		if err := os.WriteFile(destination, payload, 0o644); err != nil {
			return -1
		}
		return 0
	}); err != nil {
		t.Fatalf("failed to set fetch callback: %v", err)
	}
	defer func() {
		_ = h.SetFetchCallbackFunc(nil)
	}()

	if err := h.SetProgressCallbackFunc(func(progress int32, pkgName string, percent int, howmany uint64, current uint64) {
		atomic.AddInt32(&progressCalls, 1)
	}); err != nil {
		t.Fatalf("failed to set progress callback: %v", err)
	}
	defer func() {
		_ = h.SetProgressCallbackFunc(nil)
	}()

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("integration package payload"))
	}))
	defer server.Close()

	fetchedPath, err := h.FetchPkgURL(server.URL + "/pkg-0-any.pkg.tar.zst")
	if err != nil {
		t.Fatalf("FetchPkgURL failed: %v", err)
	}
	if fetchedPath == "" {
		t.Fatalf("FetchPkgURL returned empty path")
	}
	if _, err := os.Stat(fetchedPath); err != nil {
		t.Fatalf("fetched file path does not exist: %s (%v)", fetchedPath, err)
	}

	if fetchCalls == 0 {
		t.Fatalf("expected fetch callback to be invoked")
	}
	if downloadCalls == 0 {
		t.Logf("download callback did not fire when fetch callback handled download")
	}
	if progressCalls == 0 {
		t.Logf("progress callback did not fire for this invocation")
	}
}
