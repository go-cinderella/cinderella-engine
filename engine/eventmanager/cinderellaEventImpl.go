package eventmanager

type CinderellaEventImpl struct {
	EventType           CinderellaEventType
	ExecutionId         string
	ProcessInstanceId   string
	ProcessDefinitionId string
}

func (activitiEvent CinderellaEventImpl) GetType() CinderellaEventType {
	return activitiEvent.EventType
}

func (activitiEvent CinderellaEventImpl) SetType(eventType CinderellaEventType) {
	activitiEvent.EventType = eventType
}
func (activitiEvent CinderellaEventImpl) GetProcessDefinitionId() string {
	return activitiEvent.ProcessDefinitionId
}

func (activitiEvent CinderellaEventImpl) SetProcessDefinitionId(processDefinitionId string) {
	activitiEvent.ProcessDefinitionId = processDefinitionId
}

func (activitiEvent CinderellaEventImpl) GetProcessInstanceId() string {
	return activitiEvent.ProcessInstanceId
}

func (activitiEvent CinderellaEventImpl) SetProcessInstanceId(processInstanceId string) {
	activitiEvent.ProcessInstanceId = processInstanceId
}
