package events

const (
	TopicTaskAssigned = "task_assigned_stream"
	TopicTaskDone     = "task_done_stream"
)

type TaskAssignedEvent struct {
	TaskID    string `json:"task_id"`
	TaskTitle string `json:"task_title"`
	UserID    string `json:"user_id"`
	UserEmail string `json:"user_email"`
	UserName  string `json:"user_name"`
}
