package traceql

import "fmt"

type AttributeScope int

const (
	AttributeScopeNone AttributeScope = iota
	AttributeScopeResource
	AttributeScopeSpan
	AttributeScopeUnknown

	none = "none"
)

func AllAttributeScopes() []AttributeScope {
	return []AttributeScope{AttributeScopeResource, AttributeScopeSpan}
}

func (s AttributeScope) String() string {
	switch s {
	case AttributeScopeNone:
		return none
	case AttributeScopeSpan:
		return "span"
	case AttributeScopeResource:
		return "resource"
	}

	return fmt.Sprintf("att(%d).", s)
}

func AttributeScopeFromString(s string) AttributeScope {
	switch s {
	case "span":
		return AttributeScopeSpan
	case "resource":
		return AttributeScopeResource
	case "":
		fallthrough
	case none:
		return AttributeScopeNone
	}

	return AttributeScopeUnknown
}

type Intrinsic int

const (
	IntrinsicNone Intrinsic = iota
	IntrinsicDuration
	IntrinsicName
	IntrinsicStatus
	IntrinsicKind
	IntrinsicChildCount

	// not yet implemented in traceql but will be
	IntrinsicParent
	IntrinsicTraceRootService
	IntrinsicTraceRootSpan
	IntrinsicTraceDuration

	// not yet implemented in traceql and may never be. these exist so that we can retrieve
	// these fields from the fetch layer
	IntrinsicTraceID
	IntrinsicTraceStartTime
	IntrinsicSpanID
	IntrinsicSpanStartTime
)

func (i Intrinsic) String() string {
	switch i {
	case IntrinsicNone:
		return none
	case IntrinsicDuration:
		return "duration"
	case IntrinsicName:
		return "name"
	case IntrinsicStatus:
		return "status"
	case IntrinsicKind:
		return "kind"
	case IntrinsicChildCount:
		return "childCount"
	// below is unimplemented
	case IntrinsicParent:
		return "parent"
	case IntrinsicTraceRootService:
		return "traceRootService"
	case IntrinsicTraceRootSpan:
		return "traceRootSpan"
	case IntrinsicTraceDuration:
		return "traceDuration"
	case IntrinsicTraceID:
		return "traceID"
	case IntrinsicTraceStartTime:
		return "traceStartTime"
	case IntrinsicSpanID:
		return "spanID"
	case IntrinsicSpanStartTime:
		return "spanStartTime"
	}

	return fmt.Sprintf("intrinsic(%d)", i)
}

// intrinsicFromString returns the matching intrinsic for the given string or -1 if there is none
func intrinsicFromString(s string) Intrinsic {
	switch s {
	case "duration":
		return IntrinsicDuration
	case "name":
		return IntrinsicName
	case "status":
		return IntrinsicStatus
	case "kind":
		return IntrinsicKind
	case "childCount":
		return IntrinsicChildCount
	// unimplemented
	case "parent":
		return IntrinsicParent
	case "traceRootService":
		return IntrinsicTraceRootService
	case "traceRootSpan":
		return IntrinsicTraceRootSpan
	case "traceDuration":
		return IntrinsicTraceDuration
	case "traceID":
		return IntrinsicTraceID
	case "traceStartTime":
		return IntrinsicTraceStartTime
	case "spanID":
		return IntrinsicSpanID
	case "spanStartTime":
		return IntrinsicSpanStartTime
	}

	return IntrinsicNone
}
