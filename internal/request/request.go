package request

// Topic topic api request
type Topic struct {
	Name string `json:"name" form:"name" binding:"required,max=64"`
}

// Client clinet api request
type Client struct {
	Name string `json:"name" form:"name" binding:"required,max=64"`
}

// Event publish api request
type Event struct {
	Topic   string `json:"topic" form:"topic" xml:"topic" binding:"required,max=64"`
	Payload string `json:"payload" form:"payload" xml:"payload" binding:"required,max=2048"`
}

// MPubReq multi publish event with multi payload api request
type MPubReq struct {
	Topic   string `json:"topic" form:"topic" xml:"topic" binding:"required,max=64"`
	Payload string `json:"payload" form:"payload" xml:"payload" binding:"required"`
}

// MPubWithTopicReq publish event with multi topic api request
type MPubWithTopicReq struct {
	Topic   string `json:"topic" form:"topic" xml:"topic" binding:"required"`
	Payload string `json:"payload" form:"payload" xml:"payload" binding:"required,max=2048"`
}
