package models

import "github.com/gomodule/redigo/redis"

func NewQuestion() *Question {
	question := new(Question)
	question.User = new(User)
	return question
}

func NewAnswer() *Answer {
	answer := new(Answer)
	answer.Question = NewQuestion()
	answer.Author = new(User)
	return answer
}

func (a *Answer) GetAuthorInfo() {

}
func RemoveFromRank(conn redis.Conn, uid string) error {

}
