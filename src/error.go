package src

const (
	ERROR_HEADER_LENGTH      string = "wrong length of header"
	ERROR_HEADER_FLAG_LENGTH string = "wrong length of flag in header"
	ERROR_QNAME_END_MISSING  string = "QNAME parsing error: unexpected end of data"
	ERROR_RECORD_LENGTH      string = "records parsing error: record not ended as expected"
)
