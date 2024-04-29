package models

import (
	"fmt"
	"log"
	"my-zhihu/utils"
)

type Comment struct {
	ID            int    `json:"id"`
	Content       string `json:"content"`
	Author        *User  `json:"author"`
	UpVoteCount   uint   `json:"up_vote_count"`
	DownVoteCount uint   `json:"down_vote_count"`
	DateCreated   string `json:"date_created"`
	LikeCount     uint   `json:"like_count"`
	Liked         bool   `json:"is_like"`
}

func InsertQuestionComment(qid, content string, uid uint) (*Comment, error) {
	var err error
	defer func() {
		if err != nil {
			log.Println("*Page.InsertQuestionComment(): ", err)
		}
	}()

	conn := redisPool.Get()
	defer conn.Close()
	comment := new(Comment)
	var author User
	var dateCreated int64
	var urlTokenCode int
	tx, err := db.Begin()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	if _, err = tx.Exec("INSERT question_comments SET question_id=?, user_id=?, content=?", qid, uid, content); err != nil {
		return nil, err
	}
	if _, err = tx.Exec("UPDATE questions SET comment_count=comment_count+1 WHERE id=?", qid); err != nil {
		return nil, err
	}
	if err = tx.QueryRow("SELECT users.id, users.fullName, users.gender, users.headline, "+
		"users.avatar_url, users.url_token, users.url_token_code, users.answer_count, users.follower_count, "+
		"question_comments.id, question_comments.content, unix_timestamp(question_comments.created_at) FROM users, question_comments "+
		"WHERE users.id=question_comments.user_id AND question_comments.id=LAST_INSERT_ID() AND question_comments.user_id=?", uid).Scan(
		&author.ID, &author.Name, &author.Gender, &author.Headline, &author.AvatarURL,
		&author.URLToken, &urlTokenCode, &author.AnswerCount, &author.FollowerCount,
		&comment.ID, &comment.Content, &dateCreated); err != nil {
		return nil, err
	}
	key := fmt.Sprintf("question_comment liked:%d", comment.ID)
	v, err := conn.Do("SCARD", key)
	if err != nil {
		log.Println("*Page.QuestionComments(): ", err)
		return nil, err
	}
	utils.URLToken(&author.URLToken, urlTokenCode)
	comment.LikeCount = uint(v.(int64))
	comment.DateCreated = utils.FormatBeforeUnixTime(dateCreated)
	comment.Author = &author

	if err = tx.Commit(); err != nil {
		return nil, err
	}
	return comment, nil
}

func DeleteQuestionComment(qid, cid string, uid uint) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec("DELETE FROM question_comments WHERE id=? AND user_id=?", cid, uid); err != nil { //deleted by oneself
		return err
	}
	if _, err := tx.Exec("UPDATE questions SET comment_count=comment_count-1 WHERE id=?", qid); err != nil {
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	conn := redisPool.Get()
	defer conn.Close()
	_, _ = conn.Do("DEL", "question_comment liked:"+cid)
	return nil
}

func LikeQuestionComment(cid string, uid uint) error {
	conn := redisPool.Get()
	defer conn.Close()
	v, err := conn.Do(sadd, "question_comment liked:"+cid, uid)
	if err != nil {
		return err
	}
	if v.(int64) == 0 {
		return fmt.Errorf("reply is zero")
	}
	return nil
}

func UnlikeQuestionComment(cid string, uid uint) error {
	conn := redisPool.Get()
	defer conn.Close()
	v, err := conn.Do(srem, "question_comment liked:"+cid, uid)
	if err != nil {
		return err
	}
	if v.(int64) == 0 {
		return fmt.Errorf("reply is zero")
	}
	return nil
}
