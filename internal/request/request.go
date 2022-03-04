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
	Message string `json:"message" form:"message" xml:"message" binding:"required"`
}
