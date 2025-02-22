package rodutils

import (
	"errors"
	"fmt"
	"time"

	"github.com/go-rod/rod"
	"github.com/go-rod/rod/lib/proto"
)

// Element finds the first element matching the selector.
// It returns the element and an error, if any.
// If no element is found, it returns an error.
func Element(e *rod.Element, selector string) (*rod.Element, error) {
	if e == nil {
		return nil, errors.New("rod.Element is nil")
	}
	has, elem, err := e.Has(selector)
	if err != nil {
		return nil, fmt.Errorf("failed to check element existence: %s\n%v", selector, err)
	}
	if !has {
		return nil, fmt.Errorf("element not found: %s", selector)
	}
	return elem, nil
}

// ElementVisible finds the first element matching the selector and waits for it to be visible.
// It returns the element and an error, if any.
// If the element is not found or not visible, it returns an error.
func ElementVisible(e *rod.Element, selector string) (*rod.Element, error) {
	if e == nil {
		return nil, errors.New("rod.Element is nil")
	}
	has, elem, err := e.Has(selector)
	if err != nil {
		return nil, fmt.Errorf("failed to check element existence: %s\n%v", selector, err)
	}
	if !has {
		return nil, fmt.Errorf("element not found: %s", selector)
	}
	err = elem.WaitVisible()
	if err != nil {
		return nil, fmt.Errorf("failed to wait for element to be visible: %s\n%v", selector, err)
	}
	return elem, nil
}

// ElementVisible finds the first element matching the selector and waits for it to be visible.
// It returns the element and an error, if any.
// If the element is not found or not visible, it returns an error.
func ElementStable(e *rod.Element, selector string, duration *time.Duration) (*rod.Element, error) {
	if e == nil {
		return nil, errors.New("rod.Element is nil")
	}
	has, elem, err := e.Has(selector)
	if err != nil {
		return nil, fmt.Errorf("failed to check element existence: %s\n%v", selector, err)
	}
	if !has {
		return nil, fmt.Errorf("element not found: %s", selector)
	}
	if duration == nil {
		tmp := DefaultStableDuration
		duration = &tmp
	}
	err = elem.WaitStable(*duration)
	if err != nil {
		return nil, fmt.Errorf("failed to wait for element to be visible: %s\n%v", selector, err)
	}
	return elem, nil
}

// ElementText finds the first element matching the selector and returns its text content.
// It returns the text content and an error, if any.
// If the element is not found or the text cannot be retrieved, it returns an error.
func ElementText(e *rod.Element, selector string) (*string, error) {
	if e == nil {
		return nil, errors.New("rod.Element is nil")
	}
	has, elem, err := e.Has(selector)
	if err != nil {
		return nil, fmt.Errorf("failed to check element existence: %s\n%v", selector, err)
	}
	if !has {
		return nil, fmt.Errorf("element not found: %s", selector)
	}
	text, err := elem.Text()
	if err != nil {
		return nil, fmt.Errorf("failed to get text: %s\n%v", selector, err)
	}
	return &text, nil
}

// Elements finds all elements matching the selector.
// It returns a slice of elements and an error, if any.
// If no elements are found, it returns an error.
func Elements(e *rod.Element, selector string) (rod.Elements, error) {
	if e == nil {
		return nil, errors.New("rod.Element is nil")
	}
	elems, err := e.Elements(selector)
	if err != nil {
		return nil, fmt.Errorf("failed to get element: %s\n%v", selector, err)
	}
	if len(elems) == 0 {
		return nil, fmt.Errorf("the number of acquired elements was 0: %s", selector)
	}
	return elems, nil
}

// Click clicks the element.
// It returns an error if the element is not enabled or if the click fails.
func Click(e *rod.Element) error {
	if e == nil {
		return errors.New("rod.Element is nil")
	}
	err := e.WaitEnabled()
	if err != nil {
		return errors.New("failed to wait for element to be enabled")
	}

	err = e.Click(proto.InputMouseButtonLeft, 1)
	if err != nil {
		return fmt.Errorf("failed to click: %v", err)
	}

	return nil
}

// ClickAndLoad clicks the element and waits for the page to load.
// It returns an error if the click fails or if the page fails to load.
func ClickAndLoad(e *rod.Element) error {
	if e == nil {
		return errors.New("rod.Element is nil")
	}
	err := Click(e)
	if err != nil {
		return err
	}
	err = e.Page().WaitLoad()
	if err != nil {
		return fmt.Errorf("error waiting for page load to complete: %v", err)
	}

	return nil
}

// Input inputs text into the element.
// It returns an error if the element is not writable or if the input fails.
func Input(e *rod.Element, txt string) error {
	if e == nil {
		return errors.New("rod.Element is nil")
	}
	err := e.WaitWritable()
	if err != nil {
		return errors.New("failed to wait for element to be writable")
	}

	err = e.Input(txt)
	if err != nil {
		return fmt.Errorf("failed to input text: %v", err)
	}

	return nil
}

// Attribute returns the value of the attribute with the given name.
// It returns the attribute value and an error, if any.
// If the attribute is not found, it returns an error.
func Attribute(e *rod.Element, name string) (*string, error) {
	if e == nil {
		return nil, errors.New("rod.Element is nil")
	}
	attr, err := e.Attribute(name)
	if err != nil {
		return nil, fmt.Errorf("failed to get attribute: %s\n%v", name, err)
	}
	return attr, nil
}
