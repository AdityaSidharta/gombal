package pkg

type GombalErrors struct {
	message string
}

func (e *GombalErrors) Error() string {
	return e.message
}

var invalidQueryError = &GombalErrors{message: "query does not exist in the dataset"}
var invalidResponseError = &GombalErrors{message: "response does not exist in the dataset"}
var emptyResponseError = &GombalErrors{message: "zero response for the specific query in the dataset"}
var invalidStrategyError = &GombalErrors{message: "invalid strategy"}