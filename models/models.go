package models

type Member struct {
	User
	MarkedCount         uint   `json:"marked_count"`
	FollowingCount      uint   `json:"following_count"`
	PrivacyProtected    bool   `json:"is_privacy_protected"`
	VoteUpCount         uint   `json:"vote_up_count"`
	Description         string `json:"description"`
	FollowingTopicCount uint   `json:"following_topic_count"`
	ThankedCount        uint   `json:"thanked_count"`
}

type Post struct{}

const (
	AskQuestionAction = iota
	AnswerQuestionAction
	FollowQuestionAction
	VoteUpAnswerAction
	OtherAction
)

type Topic struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type AnswerComment struct {
	*Answer
	Comment
}

type QuestionComment struct {
	*Question
	Comment
}
