package models

import (
	"github.com/gin-gonic/gin"
	"log"
	"strconv"
)

type Question struct {
	ID           string `json:"id"`
	User         *User
	Title        string `json:"title"`
	Detail       string `json:"detail"`
	DateCreated  string `json:"date_created"`
	DateModified string `json:"date_modified"`

	VisitCount    uint `json:"visit_count"`
	AnswerCount   uint `json:"answer_count"`
	CommentCount  uint `json:"comment_count"`
	FollowerCount uint `json:"follower_count"`

	Answers        []*Answer `json:"answers"`
	Topics         []*Topic  `json:"topics"`
	TopicURLTokens []*Topic  `json:"topic_url_tokens"`

	Followed             bool    `json:"is_followed"`
	Answered             bool    `json:"is_answered"`
	VisitorAnswerID      uint    `json:"visitor_answer_id"`
	VisitorAnswerDeleted *Answer `json:"visitor_answer_deleted"`
	Anonymous            bool    `json:"is_anonymous"`
}

func (q *Question) GetTopics(c *gin.Context) {

}

func (q *Question) GetAnswers(c *gin.Context) {

}

func (q *Question) UpdateVisitCount(c *gin.Context) {

}

func InsertQuestion(question *Question, uid uint) error {
	var err error
	defer func() {
		if err != nil {
			log.Println("models.InsertQuestion(): ", err)
		}
	}()

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	row, err := tx.Exec("INSERT questions SET user_id=?, title=?, detail=?",
		uid, question.Title, question.Detail)
	if err != nil {
		return err
	}
	id, err := row.LastInsertId()
	if err != nil {
		return err
	}
	qid := strconv.FormatInt(id, 10)
	question.ID = qid
	for _, topic := range question.TopicURLTokens {

		if _, err = tx.Exec("INSERT question_topics SET question_id=?, topic_id=?",
			qid, topic); err != nil {
			return err
		}
	}
	if _, err = tx.Exec("UPDATE users SET question_count=question_count+1 WHERE id=?",
		uid); err != nil {
		return err
	}
	if err = tx.Commit(); err != nil {
		return err
	}
	go HandleNewAction(uid, AskQuestionAction, qid)
	return nil
}

func FollowQuestion(qid string, uid uint) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("INSERT question_followers SET question_id=?, user_id=?",
		qid, uid); err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE questions SET follower_count=follower_count+1 WHERE id=?",
		qid); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	go HandleNewAction(uid, FollowQuestionAction, qid)
	return nil
}

func UnfollowQuestion(qid string, uid uint) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM question_followers WHERE question_id=? AND user_id=?",
		qid, uid); err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE questions SET follower_count=follower_count-1 WHERE id=?",
		qid); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	go RemoveAction(uid, FollowQuestionAction, qid)
	return nil
}
