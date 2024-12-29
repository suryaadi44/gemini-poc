package custom

type ErrorValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type ErrorValues []ErrorValue

func NewErrorValue(key string, value string) *ErrorValue {
	return &ErrorValue{
		Key:   key,
		Value: value,
	}
}

func NewErrorValues(errorValue ...ErrorValue) ErrorValues {
	err := ErrorValues{}
	for _, e := range errorValue {
		err = append(err, e)
	}

	return err
}
