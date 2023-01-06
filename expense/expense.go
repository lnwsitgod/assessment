package expense

import "errors"

type Expense struct {
	ID     int      `json:"id"`
	Title  string   `json:"title"`
	Amount float64  `json:"amount"`
	Note   string   `json:"note"`
	Tags   []string `json:"tags"`
}

type Err struct {
	Message string `json:"message"`
}

func (e *Expense) Validate() error {
	if e.Title == "" {
		return errors.New("title is required")
	}
	if e.Amount <= 0 {
		return errors.New("amount is required and must be greater than 0")
	}
	if e.Note == "" {
		return errors.New("note is required")
	}
	if len(e.Tags) == 0 {
		return errors.New("at least one tag is required")
	}
	return nil
}
