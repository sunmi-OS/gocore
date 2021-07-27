package utils

import (
	"context"
	"net/http"

	"google.golang.org/grpc/metadata"
)

const (
	X_REQUEST_ID      = "x-request-id"
	X_B3_TRACEID      = "x-b3-traceid"
	X_B3_SPANID       = "x-b3-spanid"
	X_B3_PARENTSPANID = "x-b3-parentspanid"
	X_B3_SAMPLED      = "x-b3-sampled"
	X_B3_FLAGS        = "x-b3-flags"
	B3                = "b3"
	X_OT_SPAN_CONTEXT = "x-ot-span-context"
)

type TraceHeader struct {
	HttpHeader http.Header
	GrpcMd     metadata.MD
}

func SetHttp(header http.Header) TraceHeader {

	return TraceHeader{
		HttpHeader: header,
		GrpcMd:     httpToGrpc(header),
	}
}

func SetGrpc(ctx context.Context) TraceHeader {
	headersIn, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return TraceHeader{}
	}

	return TraceHeader{
		GrpcMd:     headersIn,
		HttpHeader: grpcToHttp(headersIn),
	}
}

func SetHeader(header interface{}) TraceHeader {

	switch header := header.(type) {
	case http.Header:
		return SetHttp(header)
	case context.Context:
		return SetGrpc(header)
	default:
		return TraceHeader{}
	}
}

func grpcToHttp(headersIn metadata.MD) http.Header {
	httpHeader := http.Header{}

	requestId := headersIn.Get(X_REQUEST_ID)
	traceId := headersIn.Get(X_B3_TRACEID)
	spanId := headersIn.Get(X_B3_SPANID)
	panrentSpanId := headersIn.Get(X_B3_PARENTSPANID)
	sampled := headersIn.Get(X_B3_SAMPLED)
	flags := headersIn.Get(X_B3_FLAGS)
	spanContext := headersIn.Get(X_OT_SPAN_CONTEXT)
	b3 := headersIn.Get(B3)

	if len(requestId) > 0 {
		httpHeader.Add(X_REQUEST_ID, requestId[0])
	}
	if len(traceId) > 0 {
		httpHeader.Add(X_B3_TRACEID, traceId[0])
	}
	if len(spanId) > 0 {
		httpHeader.Add(X_B3_SPANID, spanId[0])
	}
	if len(panrentSpanId) > 0 {
		httpHeader.Add(X_B3_PARENTSPANID, panrentSpanId[0])
	}
	if len(sampled) > 0 {
		httpHeader.Add(X_B3_SAMPLED, sampled[0])
	}
	if len(flags) > 0 {
		httpHeader.Add(X_B3_FLAGS, flags[0])
	}
	if len(spanContext) > 0 {
		httpHeader.Add(X_OT_SPAN_CONTEXT, spanContext[0])
	}
	if len(b3) > 0 {
		httpHeader.Add(B3, b3[0])
	}

	return httpHeader
}

func httpToGrpc(header http.Header) metadata.MD {

	medata := map[string]string{}

	if header.Get(X_REQUEST_ID) != "" {
		medata[X_REQUEST_ID] = header.Get(X_REQUEST_ID)
	}
	if header.Get(X_B3_TRACEID) != "" {
		medata[X_B3_TRACEID] = header.Get(X_B3_TRACEID)
	}
	if header.Get(X_B3_SPANID) != "" {
		medata[X_B3_SPANID] = header.Get(X_B3_SPANID)
	}
	if header.Get(X_B3_PARENTSPANID) != "" {
		medata[X_B3_PARENTSPANID] = header.Get(X_B3_PARENTSPANID)
	}
	if header.Get(X_B3_SAMPLED) != "" {
		medata[X_B3_SAMPLED] = header.Get(X_B3_SAMPLED)
	}
	if header.Get(X_B3_FLAGS) != "" {
		medata[X_B3_FLAGS] = header.Get(X_B3_FLAGS)
	}
	if header.Get(X_OT_SPAN_CONTEXT) != "" {
		medata[X_OT_SPAN_CONTEXT] = header.Get(X_OT_SPAN_CONTEXT)
	}
	if header.Get(B3) != "" {
		medata[B3] = header.Get(B3)
	}

	return metadata.New(medata)
}
