package models

import (
	"github.com/gomodule/redigo/redis"
	"html/template"
	"log"
	"strconv"
	"time"
)

type Answer struct {
	ID string `json:"id"`
	*Question
	Author        *User `json:"author"`
	Content       template.HTML
	DateCreated   string
	DateModified  string
	UpVoteCount   uint
	DownVoteCount uint
	MarkedCount   uint
	CommentCount  uint
	Comments      []*AnswerComment
	IsAuthor      bool
	Deleted       bool
	UpVoted       bool
	DownVoted     bool
}

func InsertAnswer(qid string, content string, uid uint) (string, error) {
	var err error
	defer func() {
		if err != nil {
			log.Println("models.InsertAnswer err:", err)
		}
	}()
	tx, err := db.Begin()
	if err != nil {
		return "", err
	}
	defer tx.Rollback()
	now := time.Now()
	res, err := tx.Exec("INSERT answers SET content=?,user_id=?,question_id=?,create_id=?", content, uid, qid, now)
	if err != nil {
		return "", err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return "", err
	}
	aid := strconv.FormatInt(id, 10)
	_, err = tx.Exec("UPDATE questions SET answer_count=answer_count+1 where id=?", aid)
	if err != nil {
		return "", err
	}
	_, err = tx.Exec("UPDATE answers SET answer_count=answer_count+1 where id=?", aid)
	if err != nil {
		return "", err
	}
	conn := redisPool.Get()
	defer conn.Close()
	if err := UpdateRank(conn, aid, now.Unix()); err != nil {
		return "", err
	}
	err = tx.Commit()
	go HandleNewAction(uid, AnswerQuestionAction, aid)
	return aid, err
}

func DeleteAnswer(aid string, uid uint) error {
	var err error
	defer func() {
		if err != nil {
			log.Println("models.DeleteAnswer err:", err)
		}
	}()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var qid string
	err = tx.QueryRow("SELECT question_id FROM answers WHERE id=?", aid).Scan(&qid)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE answers SET is_deleted=1 WHERE id=?", aid)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE questions SET answer_count=answer_count-1 WHERE id=?", aid)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE user SET answer_count=answer_count-1 WHERE id=?", aid)
	if err != nil {
		return err
	}
	conn := redisPool.Get()
	defer conn.Close()
	if err := RemoveFromRank(conn, aid); err != nil {
		return err
	}
	err = tx.Commit()
	go RemoveAction(uid, AnswerQuestionAction, aid)
	return err
}

func RestoreAnswer(aid string, uid uint) error {
	var err error
	defer func() {
		if err != nil {
			log.Println("models.RestoreAnswer err:", err)
		}
	}()
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var qid string
	var time int
	err = tx.QueryRow("SELECT UNIX_TIMESTAMP(created_at),question_id FROM answers WHERE id=?", aid).Scan(&qid, &time)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE answers SET is_deleted=0 WHERE id=?", aid)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE answers SET answer_count=answer_count+1 WHERE id=?", aid)
	if err != nil {
		return err
	}
	_, err = tx.Exec("UPDATE users SET answer_count=answer_count+1 WHERE id=?", aid)
	if err != nil {
		return err
	}
	conn := redisPool.Get()
	defer conn.Close()
	n, err := redis.Int(conn.Do("SCARD", "upvoted:"+aid))
	if err != nil {
		return err
	}
	score := int64(time + n*432)
	if err = UpdateRank(conn, aid, score); err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}
	return err
}
