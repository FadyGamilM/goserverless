package golambda

type Note struct {
	Id        string `json:"id"`
	CreatedAt string `json:"createdAt"`
	Content   string `json:"content"`
	Password  string `json:"password"`
	URL       string `json:"url"` // maybe later i will add the note at s3 and get presigned url to fetch it too (not sure what is this useful for, but its just for having fun with aws)
}
