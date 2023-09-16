package broker

import "time"

type Media struct {
	Url  string `json:"url"`
	Type string `json:"type"`
	Id   string `json:"id"`
}

type PostDocument struct {
	Id           string `json:"id"`
	UserId       string `json:"userId"`
	Text         string `json:"text"`
	Media        Media
	AllowComment bool `json:"allowComment" default:"true"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	Tags         []string `json:"tags"`
	Privacy      string   `json:"privacy" default:"Public"`
}

type LikeDocument struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	PostId    string `json:"postId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CommentDocumment struct {
	Id        string `json:"id"`
	UserId    string `json:"userId"`
	Text      string `json:"text"`
	PostId    string `json:"postId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ReplyDocument struct {
	Id        string `json:"iid"`
	UserId    string `json:"userId"`
	Text      string `json:"text"`
	CommentId string `json:"commentId"`
	CreatedAt time.Time
	UpdatedAt time.Time
}
