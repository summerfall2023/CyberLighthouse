package src

const (
	ERROR_HEADER_LENGTH      string = "wrong length of header"
	ERROR_HEADER_FLAG_LENGTH string = "wrong length of flag in header"
	ERROR_QNAME_END_MISSING  string = "QNAME parsing error: unexpected end of data"
	ERROR_RECORD_LENGTH      string = "records parsing error: record not ended as expected"
)

// output error
const (
	ERROR_HEADER_FLAG_Z           string = "z must be 0(false)"
	ERROR_QUERY_TYPE_UNSUPPORTED  string = "output error: query type unsupported"
	ERROR_QUERY_CLASS_UNSUPPORTED string = "output error: query class unsupported"
)
