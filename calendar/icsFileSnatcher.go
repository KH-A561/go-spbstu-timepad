package calendar

import (
	"bytes"
	"context"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"os"

	"time"
)

const (
	spbstuGetIcalAddress = "https://ruz.spbstu.ru/faculty/%d/groups/%d/ical?date=%s"
	DateFormat           = "2006-1-02"
)

type FileSnatcher interface {
	SnatchCalendar(ctx context.Context, facultyId int, groupId int, date time.Time)
}

type FileSnatcherImpl struct {
	Client *fasthttp.Client
}

var DefaultFileSnatcher = &FileSnatcherImpl{
	Client: &fasthttp.Client{},
}

type Response struct {
	Ok bool `json:"ok"`
}

func (a FileSnatcherImpl) SnatchCalendar(ctx context.Context, facultyId int, groupId int, date time.Time) (
	*Response, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// Continue
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	var formattedDate = date.Format(DateFormat)
	var url = fmt.Sprintf(spbstuGetIcalAddress, facultyId, groupId, formattedDate)

	req.SetRequestURI(url)
	req.Header.SetMethod(fasthttp.MethodGet)

	resp := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(resp)

	var err error
	deadline, ok := ctx.Deadline()
	if ok {
		err = a.Client.DoDeadline(req, resp, deadline)
	} else {
		err = a.Client.Do(req, resp)
	}
	if err != nil {
		return nil, fmt.Errorf("fasthttp do request: %w", err)
	}

	if statusCode := resp.StatusCode(); statusCode >= fasthttp.StatusInternalServerError {
		return nil, fmt.Errorf("internal server error: %d", statusCode)
	}

	outFile, err := os.Create("calendar.ics")
	if err != nil {
		return nil, fmt.Errorf("create calendar.ics file: %w", err)
	}
	defer outFile.Close()

	_, err = io.Copy(outFile, bytes.NewReader(resp.Body()))
	if err != nil {
		return nil, fmt.Errorf("write calendar.ics file: %w", err)
	}

	apiResp := &Response{}
	//err = json.Unmarshal(resp.Body(), apiResp)
	if err != nil {
		return nil, fmt.Errorf("decode json: %w", err)
	}

	return apiResp, nil
}
