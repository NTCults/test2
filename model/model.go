package model

// Event represents MQ payload
type Event struct {
	Source    string `json:"source"`
	Component string `json:"component"`
	Resource  string `json:"resource"`
	Crit      int    `json:"crit"`
	Message   string `json:"message"`
	Timestamp int    `json:"timestamp"`
}

// constants that represents Alert status
const (
	StatusOngoing  = "ONGOING"
	StatusResolved = "RESOLVED"
)

// Alert is main application entity
type Alert struct {
	Component    string `bson:"component"`
	Resource     string `bson:"resource"`
	Crit         int    `bson:"crit"`
	StartTime    int    `bson:"start_time"`
	LastTime     int    `bson:"last_time"`
	Status       string `bson:"status"`
	LastMessage  string `bson:"last_message"`
	FirstMessage string `bson:"first_message"`
}

func (a *Alert) IsOngoing() bool {
	return a.Status == StatusOngoing
}
