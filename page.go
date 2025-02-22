package rodutils

import (
	"context"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// Navigate navigates the page to the given URL.
// It returns the page and an error, if any.
// If the navigation fails, it returns an error.
func Navigate(p *rod.Page, url string) (*rod.Page, error) {
	if p == nil {
		return nil, fmt.Errorf("rod.Page is nil")
	}
	err := p.Navigate(url)
	if err != nil {
		return nil, fmt.Errorf("failed to navigate to page: %s\n%v", url, err)
	}
	return p, nil
}

// PageElement finds the first element matching the selector in the page.
// It returns the element and an error, if any.
// If the element is not found, it returns an error.
func PageElement(p *rod.Page, selector string) (*rod.Element, error) {
	if p == nil {
		return nil, fmt.Errorf("rod.Page is nil")
	}
	elem, err := p.Element(selector)
	if err != nil {
		return nil, fmt.Errorf("failed to get element: %s\n%v", selector, err)
	}
	return elem, nil
}

// PageElementVisible finds the first element matching the selector in the page and waits for it to be visible.
// It returns the element and an error, if any.
// If the element is not found or not visible, it returns an error.
func PageElementVisible(p *rod.Page, selector string) (*rod.Element, error) {
	if p == nil {
		return nil, fmt.Errorf("rod.Page is nil")
	}
	ok, _, err := p.Has(selector)
	if !ok || err != nil {
		return nil, fmt.Errorf("element not found: %s\n%v", selector, err)
	}
	elem, err := p.Element(selector)
	if err != nil {
		return nil, fmt.Errorf("failed to get element: %s\n%v", selector, err)
	}
	err = elem.WaitVisible()
	if err != nil {
		return nil, fmt.Errorf("failed to wait for element to be visible: %s\n%v", selector, err)
	}
	return elem, nil
}

// PageElements finds all elements matching the selector in the page.
// It returns a slice of elements and an error, if any.
// If no elements are found, it returns an error.
func PageElements(p *rod.Page, selector string) (rod.Elements, error) {
	if p == nil {
		return nil, fmt.Errorf("rod.Page is nil")
	}
	ok, _, err := p.Has(selector)
	if !ok || err != nil {
		return nil, fmt.Errorf("element not found: %s\n%v", selector, err)
	}
	elems, err := p.Elements(selector)
	if err != nil {
		return nil, fmt.Errorf("failed to get elements: %s\n%v", selector, err)
	}
	if len(elems) == 0 {
		return nil, fmt.Errorf("the number of acquired elements was 0: %s", selector)
	}
	return elems, nil
}

// ScrollToBottom scrolls the page to the bottom.
// It returns an error if scrolling fails.
func ScrollToBottom(page *rod.Page) error {
	if page == nil {
		return fmt.Errorf("rod.Page is nil")
	}
	err := page.Mouse.Scroll(0, 0, 0) // Reset scroll position
	if err != nil {
		return fmt.Errorf("failed to reset scroll position: %v", err)
	}
	var lastHeight float64
	r, err := page.Eval(`() => document.body.scrollHeight`)
	if err != nil {
		return err
	}

	lastHeight = r.Value.Num()

	err = page.Mouse.Scroll(0, lastHeight, 0)
	if err != nil {
		return err
	}

	return nil
}

type RodOptions struct {
	Timeout        time.Duration // Overall timeout
	StableDuration time.Duration // Time the element needs to be stable
	RetryCount     int           // Number of retries
	RetryDelay     time.Duration // Wait time between retries
	MustVisible    bool          // Whether the element needs to be visible
	MustStable     bool          // Whether the element needs to be stable
	MustWaitLoad   bool          // Whether the page load needs to be complete
}

// DefaultRodOptions returns the default options.
func DefaultRodOptions() *RodOptions {
	return &RodOptions{
		Timeout:        10 * time.Second,
		StableDuration: 200 * time.Millisecond,
		RetryCount:     3,
		RetryDelay:     500 * time.Millisecond,
		MustVisible:    true,
		MustWaitLoad:   true,
		MustStable:     true,
	}
}

// SafeClick executes a click after waiting for the element to stabilize.
// It returns an error if the click fails.
func SafeClick(page *rod.Page, selector string, opts *RodOptions) error {
	if opts == nil {
		opts = DefaultRodOptions()
	}
	var lastErr error

	for i := 0; i <= opts.RetryCount; i++ {

		// If an error occurs, wait a bit and then retry
		if i > 0 {
			time.Sleep(opts.RetryDelay)
		}

		// Timeout context
		ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
		defer cancel()

		// Wait for element
		el, err := page.Context(ctx).Element(selector)
		if err != nil {
			lastErr = fmt.Errorf("element not found: %w", err)
			continue
		}

		// // Check visibility
		// if err := el.WaitVisible(); err != nil {
		// 	lastErr = fmt.Errorf("element not visible: %w", err)
		// 	continue
		// }

		// Wait for element to stabilize
		if err := el.WaitStable(opts.StableDuration); err != nil {
			lastErr = fmt.Errorf("element not stable: %w", err)
			continue
		}
		// Execute click
		if err := el.Click(proto.InputMouseButtonLeft, 1); err != nil {
			lastErr = fmt.Errorf("click failed: %w", err)
			continue
		}

		// Exit the loop if successful
		if lastErr == nil {
			return nil
		}
	}

	return fmt.Errorf("all click attempts failed: %w", lastErr)
}

// SafeElement retrieves an element safely.
// It returns the element and an error, if any.
// If the element is not found, it returns an error.
func SafeElement(p *rod.Page, selector string, opts *RodOptions) (*rod.Element, error) {
	if p == nil {
		return nil, fmt.Errorf("rod.Page is nil")
	}
	if opts == nil {
		opts = DefaultRodOptions()
	}
	var lastErr error
	var element *rod.Element

	for i := 0; i <= opts.RetryCount; i++ {
		if i > 0 {
			time.Sleep(opts.RetryDelay)
		}

		// Timeout context
		ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
		defer cancel()

		// Wait for element
		el, err := p.Context(ctx).Element(selector)
		if err != nil {
			lastErr = fmt.Errorf("element not found: %w", err)
			continue
		}

		// Visibility check (optional)
		if opts.MustVisible {
			if err := el.WaitVisible(); err != nil {
				lastErr = fmt.Errorf("element not visible: %w", err)
				continue
			}
		}

		// Stability check (optional)
		if opts.MustStable {
			if err := el.WaitStable(opts.StableDuration); err != nil {
				lastErr = fmt.Errorf("element not stable: %w", err)
				continue
			}
		}

		// If all checks pass
		element = el

		return element, lastErr
	}

	return nil, fmt.Errorf("all attempts to get element failed: %w", lastErr)
}
