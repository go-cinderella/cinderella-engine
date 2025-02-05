package model

type Task struct {
	Activity
}

func (task Task) GetActivity() Activity {
	return task.Activity
}
