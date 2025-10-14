package models

import "errors"

type Status int

const (
    Invalid       Status = 400
    Unauthorized  Status = 401
    Forbidden     Status = 403
    NotFound      Status = 404
    Conflict      Status = 409
    RateLimited   Status = 429
    Precondition  Status = 412
    Timeout       Status = 504
    Unavailable   Status = 503
    Dependency    Status = 502
    Internal      Status = 500
)

type Err struct {
    Op		string	`json:"op"`
    Status	Status	`json:"status"`
    Msg		string	`json:"msg"`
    Err		error	`json:"error"`
}

func (e *Err) Error() string {
	if e == nil {
		return "<nil>"
	}
	base := ""
	if e.Op != "" {
		base = e.Op + ": "
	}
	if e.Msg != "" {
		base += e.Msg
	}
	return base
}

func E(op string, status Status, msg string, err error) error {
    return &Err{Op: op, Status: status, Msg: msg, Err: err}
}

func (e *Err) Unwrap() error { return e.Err }

func StatusOf(err error) Status {
    var ae *Err
    if errors.As(err, &ae) && ae != nil {
        return ae.Status
    }
    return Internal
}

func HTTPStatus(err error) int { return int(StatusOf(err)) }