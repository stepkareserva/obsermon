package models

type MissingMetricNameError struct {
}

type InvalidMetricNameError struct {
}

type InvalidMetricValueError struct {
}

type ExtraRequestComponentsError struct {
}

func (e *MissingMetricNameError) Error() string {
	return "Missing metric name"
}

func (e *InvalidMetricNameError) Error() string {
	return "Invalid metric name"
}

func (e *InvalidMetricValueError) Error() string {
	return "Invalid metric value"
}

func (e *ExtraRequestComponentsError) Error() string {
	return "Extra request components detected"
}
