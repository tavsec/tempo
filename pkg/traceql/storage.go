package traceql

import (
	"context"
)

type Operands []Static

type Condition struct {
	Attribute Attribute
	Op        Operator
	Operands  Operands
}

func SearchMetaConditions() []Condition {
	return []Condition{
		{NewIntrinsic(IntrinsicTraceRootService), OpNone, nil},
		{NewIntrinsic(IntrinsicTraceRootSpan), OpNone, nil},
		{NewIntrinsic(IntrinsicTraceDuration), OpNone, nil},
		{NewIntrinsic(IntrinsicTraceID), OpNone, nil},
		{NewIntrinsic(IntrinsicTraceStartTime), OpNone, nil},
		{NewIntrinsic(IntrinsicSpanID), OpNone, nil},
		{NewIntrinsic(IntrinsicSpanStartTime), OpNone, nil},
		{NewIntrinsic(IntrinsicDuration), OpNone, nil},
	}
}

func SearchMetaConditionsWithoutDuration() []Condition {
	return []Condition{
		{NewIntrinsic(IntrinsicTraceRootService), OpNone, nil},
		{NewIntrinsic(IntrinsicTraceRootSpan), OpNone, nil},
		{NewIntrinsic(IntrinsicTraceDuration), OpNone, nil},
		{NewIntrinsic(IntrinsicTraceID), OpNone, nil},
		{NewIntrinsic(IntrinsicTraceStartTime), OpNone, nil},
		{NewIntrinsic(IntrinsicSpanID), OpNone, nil},
		{NewIntrinsic(IntrinsicSpanStartTime), OpNone, nil},
	}
}

// SecondPassFn is a method that is called in between the first and second
// pass of a fetch spans request. See below.
type SecondPassFn func(*Spanset) ([]*Spanset, error)

type FetchSpansRequest struct {
	StartTimeUnixNanos uint64
	EndTimeUnixNanos   uint64
	Conditions         []Condition

	// Hints

	// By default the storage layer fetches spans meeting any of the criteria.
	// This hint is for common cases like { x && y && z } where the storage layer
	// can make extra optimizations by returning only spansets that meet
	// all criteria.
	AllConditions bool

	// SecondPassFn and Conditions allow a caller to retrieve one set of data
	// in the first pass, filter using the SecondPassFn callback and then
	// request a different set of data in the second pass. This is particularly
	// useful for retrieving data required to resolve a TraceQL query in the first
	// pass and only selecting metadata in the second pass.
	// TODO: extend this to an arbitrary number of passes
	SecondPass           SecondPassFn
	SecondPassConditions []Condition
}

func (f *FetchSpansRequest) appendCondition(c ...Condition) {
	f.Conditions = append(f.Conditions, c...)
}

type Span interface {
	// these are the actual fields used by the engine to evaluate queries
	// if a Filter parameter is passed the spans returned will only have this field populated
	Attributes() map[Attribute]Static

	ID() []byte
	StartTimeUnixNanos() uint64
	DurationNanos() uint64
}

const attributeMatched = "__matched"

type Spanset struct {
	// these fields are actually used by the engine to evaluate queries
	Scalar Static
	Spans  []Span

	TraceID            []byte
	RootSpanName       string
	RootServiceName    string
	StartTimeUnixNanos uint64
	DurationNanos      uint64
	Attributes         map[string]Static
}

func (s *Spanset) AddAttribute(key string, value Static) {
	if s.Attributes == nil {
		s.Attributes = make(map[string]Static)
	}
	s.Attributes[key] = value
}

func (s *Spanset) clone() *Spanset {
	var atts map[string]Static
	if s.Attributes != nil {
		atts = make(map[string]Static, len(s.Attributes))
		for k, v := range s.Attributes {
			atts[k] = v
		}
	}

	return &Spanset{
		TraceID:            s.TraceID,
		Scalar:             s.Scalar,
		RootSpanName:       s.RootSpanName,
		RootServiceName:    s.RootServiceName,
		StartTimeUnixNanos: s.StartTimeUnixNanos,
		DurationNanos:      s.DurationNanos,
		Spans:              s.Spans, // we're not deep cloning into the spans themselves
		Attributes:         atts,
	}
}

type SpansetIterator interface {
	Next(context.Context) (*Spanset, error)
	Close()
}

type FetchSpansResponse struct {
	Results SpansetIterator
	// callback to get the size of data read during Fetch
	Bytes func() uint64
}

type SpansetFetcher interface {
	Fetch(context.Context, FetchSpansRequest) (FetchSpansResponse, error)
}

// MustExtractFetchSpansRequestWithMetadata parses the given traceql query and returns
// the storage layer conditions. Panics if the query fails to parse.
func MustExtractFetchSpansRequestWithMetadata(query string) FetchSpansRequest {
	c, err := ExtractFetchSpansRequest(query)
	if err != nil {
		panic(err)
	}
	c.SecondPass = func(s *Spanset) ([]*Spanset, error) { return []*Spanset{s}, nil }
	c.SecondPassConditions = SearchMetaConditions()
	return c
}

// ExtractFetchSpansRequest parses the given traceql query and returns
// the storage layer conditions. Returns an error if the query fails to parse.
func ExtractFetchSpansRequest(query string) (FetchSpansRequest, error) {
	ast, err := Parse(query)
	if err != nil {
		return FetchSpansRequest{}, err
	}

	req := FetchSpansRequest{
		AllConditions: true,
	}

	ast.Pipeline.extractConditions(&req)
	return req, nil
}

type SpansetFetcherWrapper struct {
	f func(ctx context.Context, req FetchSpansRequest) (FetchSpansResponse, error)
}

var _ = (SpansetFetcher)(&SpansetFetcherWrapper{})

func NewSpansetFetcherWrapper(f func(ctx context.Context, req FetchSpansRequest) (FetchSpansResponse, error)) SpansetFetcher {
	return SpansetFetcherWrapper{f}
}

func (s SpansetFetcherWrapper) Fetch(ctx context.Context, request FetchSpansRequest) (FetchSpansResponse, error) {
	return s.f(ctx, request)
}
