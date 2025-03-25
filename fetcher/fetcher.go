package fetcher

import (
	"bytes"
	"context"
	"fmt"
	"github.com/valyala/fasthttp"
	"io"
	"os"
	"universityTimepad/model"
)

const (
	spbstuGetIcalAddress = "https://ruz.spbstu.ru/faculty/%d/groups/%d/ical?date=%s"
	DateFormat           = "2006-1-02"
)

type DefaultFetcherImpl[E any] struct {
	Client *fasthttp.Client
}

type FacultyFetcher[E model.Faculty] struct {
	DefaultFetcherImpl[E]
}

func (a DefaultFetcherImpl[E]) Fetch(ctx context.Context, reqUri string, respFilename string) (*os.File, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		// Continue
	}

	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(reqUri)
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

	outFile, err := os.Create(respFilename)
	if err != nil {
		return nil, fmt.Errorf("create calendar.ics file: %w", err)
	}
	defer closeFile(outFile)
	if err != nil {
		panic(err)
	}

	_, err = io.Copy(outFile, bytes.NewReader(resp.Body()))
	if err != nil {
		return nil, fmt.Errorf("write %s.ics file: %w", respFilename, err)
	}

	return outFile, nil
}

func closeFile(f *os.File) {
	fmt.Println("closing")
	err := f.Close()
	if err != nil {
		panic(err)
	}
}
