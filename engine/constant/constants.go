package constant

const (
	ACTIVITI_HANDLER_CODE = "OK"
	TASK_TYPE_CREATE      = "create"
	TASK_TYPE_COMPLETED   = "complete"
	TASK_TYPE_MOVED       = "move"
	TASK_TYPE_ASSIGNED    = "assignment"
)

const DurationUnit = 1_000_000

const HttpRetryCount = 3

var ELEMENT_TASK_LIST = []string{ELEMENT_TASK_USER}

const (
	SERVICE_TASK_HTTP     = "http"
	SERVICE_TASK_PIPELINE = "pipeline"
)
