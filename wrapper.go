package rodutils

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// RodOperationWrapperOptions encapsulates the optional parameters for RodOperationWrapper.
type RodOperationWrapperOptions struct {
	TimeoutDuration *time.Duration
	Path            *string
	Name            *string
}

// RodOperationWrapper wraps a rod operation with error handling and screenshot capture.
// It executes the given operation and captures a screenshot if an error occurs.
// It returns an error if the operation fails.
func RodOperationWrapper(page *rod.Page, operation func() error, opts *RodOperationWrapperOptions) error {
	var timeout *time.Duration
	var path *string
	var name *string

	if opts != nil {
		timeout = opts.TimeoutDuration
		path = opts.Path
		name = opts.Name
	}

	if path == nil {
		defaultValue := DefaultScreenshotPath
		path = &defaultValue
	}
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	if name == nil {
		name = &timestamp
	} else {
		tmp := fmt.Sprintf("%s_%s", *name, timestamp)
		name = &tmp
	}
	if timeout == nil {
		defaultValue := DefaultTimeoutDuration
		timeout = &defaultValue
	}
	limit := time.Duration(*timeout) * time.Second
	err := timeLimit(limit, func() error {
		return operation()
	})
	if err != nil {
		screenshotName := fmt.Sprintf("%s/%s.png", *path, *name)
		dir := filepath.Dir(screenshotName)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create screenshot directory: %w", err)
		}
		screenshotData, screenErr := page.Screenshot(true, &proto.PageCaptureScreenshot{Format: proto.PageCaptureScreenshotFormatPng})
		if screenErr != nil {
			return fmt.Errorf("failed to capture screenshot: %w", screenErr)
		}
		// Save the screenshot data to a file
		if fileErr := os.WriteFile(screenshotName, screenshotData, 0644); fileErr != nil {
			return fmt.Errorf("failed to save screenshot: %v", fileErr)
		}
	}
	return err
}

// timeLimit executes the given function with a time limit.
// It returns an error if the operation takes longer than the given timeout.
func timeLimit(timeout time.Duration, f func() error) error {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Receive results using channels
	errorChan := make(chan error)

	go func() {
		err := f()
		if err != nil {
			errorChan <- err
			return
		}
	}()

	select {
	case err := <-errorChan:
		return err
	case <-ctx.Done():
		return fmt.Errorf("operation timed out: %v", ctx.Err())
	}
}
