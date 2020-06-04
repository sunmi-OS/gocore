package istio

import (
	"context"
	"google.golang.org/grpc/metadata"
	"net/http"
)

const (
	X_REQUEST_ID      = "x-request-id"
	X_B3_TRACEID      = "x-b3-traceid"
	X_B3_SPANID       = "x-b3-spanid"
	X_B3_PARENTSPANID = "x-b3-parentspanid"
	X_B3_SAMPLED      = "x-b3-sampled"
	X_B3_FLAGS        = "x-b3-flags"
	X_OT_SPAN_CONTEXT = "x-ot-span-context"
)

type TraceHeader struct {
	Http_Header http.Header
	Grpc_MD     metadata.MD
}

func SetHttp(header http.Header) TraceHeader {

	return TraceHeader{
		Http_Header: header,
		Grpc_MD:     httpTogrpc(header),
	}
}

func SetGrpc(ctx context.Context) TraceHeader {
	headersIn, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return TraceHeader{}
	}

	return TraceHeader{
		Grpc_MD:     headersIn,
		Http_Header: grpcTohttp(headersIn),
	}
}

func grpcTohttp(headersIn metadata.MD) http.Header {
	httpHeader := http.Header{}

	requestId := headersIn.Get(X_REQUEST_ID)
	traceId := headersIn.Get(X_B3_TRACEID)
	spanId := headersIn.Get(X_B3_SPANID)
	panrentSpanId := headersIn.Get(X_B3_PARENTSPANID)
	sampled := headersIn.Get(X_B3_SAMPLED)
	flags := headersIn.Get(X_B3_FLAGS)
	spanContext := headersIn.Get(X_OT_SPAN_CONTEXT)

	if len(requestId) > 0 {
		httpHeader.Add(X_REQUEST_ID, requestId[0])
	}
	if len(traceId) > 0 {
		httpHeader.Add(X_B3_TRACEID, requestId[0])
	}
	if len(spanId) > 0 {
		httpHeader.Add(X_B3_SPANID, requestId[0])
	}
	if len(panrentSpanId) > 0 {
		httpHeader.Add(X_B3_PARENTSPANID, requestId[0])
	}
	if len(sampled) > 0 {
		httpHeader.Add(X_B3_SAMPLED, requestId[0])
	}
	if len(flags) > 0 {
		httpHeader.Add(X_B3_FLAGS, requestId[0])
	}
	if len(spanContext) > 0 {
		httpHeader.Add(X_OT_SPAN_CONTEXT, requestId[0])
	}

	return httpHeader
}

func httpTogrpc(header http.Header) metadata.MD {

	mddata := map[string]string{}
	mddata[X_REQUEST_ID] = header.Get(X_REQUEST_ID)
	mddata[X_B3_TRACEID] = header.Get(X_B3_TRACEID)
	mddata[X_B3_SPANID] = header.Get(X_B3_SPANID)
	mddata[X_B3_PARENTSPANID] = header.Get(X_B3_PARENTSPANID)
	mddata[X_B3_SAMPLED] = header.Get(X_B3_SAMPLED)
	mddata[X_B3_FLAGS] = header.Get(X_B3_FLAGS)
	mddata[X_OT_SPAN_CONTEXT] = header.Get(X_OT_SPAN_CONTEXT)

	return metadata.New(mddata)
}
