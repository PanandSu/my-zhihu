package models

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"github.com/henrylee2cn/pholcus/common/pinyin"
	"log"
	"my-zhihu/utils"
	"strings"
	"unicode"
)

type User struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	Email     string `json:"-"`
	Password  string `json:"-"`
	Gender    string `json:"gender"`
	URLToken  string `json:"url_token"`
	Headline  string `json:"headline"` //大字标题
	AvatarURL string `json:"avatar_url"`

	Posts         []*Post `json:"posts"`
	AnswerCount   uint    `json:"answer_count"`
	FollowerCount uint    `json:"follower_count"` //恒正,使用uint即可

	Followed  bool `json:"is_followed"`
	Following bool `json:"is_following"`
	Anonymous bool `json:"is_anonymous"`
}

func (u *User) QueryRelationWithVisitor() {

}

func (u *User) VoteUp(aid string) bool {
	conn := redisPool.Get()
	err := conn.Send("SADD", "upvoted:"+aid, u.ID)
	if err != nil {
		return false
	}
	err = conn.Send("SREM", "downvoted:"+aid, u.ID)
	if err != nil {
		return false
	}
	conn.Flush()
	upvoteAddedCount, err := conn.Receive()
	if err != nil {
		return false
	}
	if _, err := conn.Receive(); err != nil {
		return false
	}

	UpdateRank(conn, aid, 432)
	if upvoteAddedCount.(int64) == 1 {
		go func() {
			_, err := db.Exec("INSERT answer_voters SET answer_id=?, user_id=?", aid, u.ID)
			if err != nil {
				log.Println("*models.User.UpVote: ", err)
			}
			HandleNewAction(u.ID, VoteUpAnswerAction, aid)
		}()
	}

	return true
}

func (u *User) VoteDown(aid string) bool {
	conn := redisPool.Get()
	conn.Send("SADD", "downvoted:"+aid, u.ID)
	conn.Send("SREM", "upvoted:"+aid, u.ID)
	conn.Flush()
	if v, err := conn.Receive(); err != nil || v == 0 {
		log.Println(err, v.(int64))
		return false
	}
	upvoteRemovedCount, err := conn.Receive()
	if err != nil {
		return false
	}
	UpdateRank(conn, aid, -432)

	if upvoteRemovedCount.(int64) == 1 {
		go func() {
			_, err := db.Exec("DELETE FROM answer_voters WHERE answer_id=? AND user_id=?", aid, user.ID)
			if err != nil {
				log.Println("*User.DownVote: ", err)
			}
		}()
	}
	return true
}

func (u *User) Neutral(aid string) bool {
	conn := redisPool.Get()
	conn.Send("SREM", "upvoted:"+aid, u.ID)
	conn.Send("SREM", "downvoted:"+aid, u.ID)
	conn.Flush()
	upvoteRemovedCount, err := conn.Receive()
	if err != nil {
		return false
	}
	if _, err := conn.Receive(); err != nil {
		return false
	}
	UpdateRank(conn, aid, -432)

	if upvoteRemovedCount == 1 {
		go func() {
			_, err := db.Exec("DELETE FROM answer_voters WHERE answer_id=? AND user_id=?", aid, user.ID)
			if err != nil {
				log.Println("*User.DownVote: ", err)
			}
		}()
	}
	return true
}

func UpdateRank(conn redis.Conn, aid string) error {
	_, err := conn.Do("ZREM", "rank", aid)
	return err
}

func InsertUser(user *User) (uid uint, err error) {
	defer func() {
		if err != nil {
			log.Println("models.InsertUser: ", err)
		}
	}()
	tx, err := db.Begin()
	if err != nil {
		return 0, err
	}
	defer tx.Rollback()

	urlToken, urlTokenCode, err := CreateURLToken(tx, user.Name)
	if err != nil {
		return 0, err
	}
	res, err := tx.Exec("INSERT users SET email=?, fullname=?, password=?, url_token=?, url_token_code=?",
		user.Email,
		user.Name,
		user.Password,
		urlToken,
		urlTokenCode,
	)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	uid = uint(id)
	err = tx.Commit()
	if err != nil {
		return 0, err
	}
	return uid, nil
}

func CreateURLToken(tx *sql.Tx, name string) (string, int, error) {
	s := []rune(name)
	var res []string
	for i := len(s) - 1; i >= 0; i-- {
		r := s[i]
		if unicode.IsDigit(r) || unicode.IsLower(r) || unicode.IsUpper(r) {
			c := fmt.Sprintf("%c", r)
			if res == nil {
				res = append(res, c)
			} else {
				res[len(res)-1] = c + res[len(res)-1]
			}
		} else {
			res = append(res, pinyin.SinglePinyin(r, pinyin.NewArgs())[0])
		}
	}
	for to, from := 0, len(res)-1; to < from; to, from = to+1, from-1 {
		res[to], res[from] = res[from], res[to]
	}
	urlToken := strings.Join(res, "-")

	urlTokenCode := 0
	if err := tx.QueryRow("SELECT url_token_code FROM users WHERE url_token=? ORDER BY id DESC limit 1",
		urlToken).Scan(&urlTokenCode); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return urlToken, urlTokenCode, nil
		}
		return "", 0, err
	}
	return urlToken, urlTokenCode + 1, nil
}

func GetUserByUsername(username string) *User {
	stmt, err := db.Prepare("SELECT id, email, password FROM users WHERE email=?")
	if err != nil {
		log.Println(err)
		return nil
	}
	defer stmt.Close()

	user := new(User)
	if err := stmt.QueryRow(username).Scan(&user.ID, &user.Email, &user.Password); err != nil {
		log.Printf("user %s: %v", username, err)
		return nil
	}
	return user
}

func GetUserByID(uid int64) *User {
	var err error
	defer func() {
		if err != nil {
			log.Println("models.GetUserByID(): uid =", uid, err)
		}
	}()
	user := new(User)
	stmt, err := db.Prepare("SELECT id, fullname, gender, headline, url_token, " +
		"url_token_code, avatar_url, answer_count, follower_count FROM users WHERE id=?")
	if err != nil {
		return nil
	}
	defer stmt.Close()

	urlTokenCode := 0
	err = stmt.QueryRow(uid).Scan(
		&user.ID,
		&user.Name,
		&user.Gender,
		&user.Headline,
		&user.URLToken,
		&urlTokenCode,
		&user.AvatarURL,
		&user.AnswerCount,
		&user.FollowerCount,
	)
	if err != nil {
		return nil
	}
	utils.URLToken(&user.URLToken, urlTokenCode)

	return user
}

func GetUserByURLToken(urlToken string, uid uint) *User {
	user := new(User)
	if err := db.QueryRow("SELECT id, fullname, gender, headline, avatar_url, "+
		"answer_count, follower_count FROM users WHERE url_token=? AND url_token_code=0", urlToken).Scan(
		&user.ID,
		&user.Name,
		&user.Gender,
		&user.Headline,
		&user.AvatarURL,
		&user.AnswerCount,
		&user.FollowerCount,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			s := strings.Split(urlToken, "-")
			if len(s) == 1 {
				return nil
			}
			code := s[len(s)-1]
			pre := strings.TrimSuffix(urlToken, "-"+code)
			if err := db.QueryRow(
				"SELECT id, fullname, gender, headline, avatar_url, "+
					"answer_count, follower_count FROM users WHERE url_token=? AND url_token_code=?", pre, code).Scan(
				&user.ID,
				&user.Name,
				&user.Gender,
				&user.Headline,
				&user.AvatarURL,
				&user.AnswerCount,
				&user.FollowerCount,
			); err != nil {
				return nil
			}
		} else {
			return nil
		}
	}
	user.URLToken = urlToken
	_ = user.QueryRelationWithVisitor(uid)
	return user
}

func FollowMember(urlToken string, uid uint) error {
	var memberID uint
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.QueryRow("SELECT id FROM users WHERE url_token=?", urlToken).Scan(&memberID); err != nil {
		return err
	}
	if _, err := tx.Exec("INSERT member_followers SET member_id=?, follower_id=?", memberID, uid); err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE users SET follower_count=follower_count+1 WHERE id=?", memberID); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func UnfollowMember(urlToken string, uid uint) error {
	var memberID uint
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := tx.QueryRow("SELECT id FROM users WHERE url_token=?", urlToken).Scan(&memberID); err != nil {
		return err
	}
	if _, err := tx.Exec("DELETE FROM member_followers WHERE member_id=? AND follower_id=?", memberID, uid); err != nil {
		return err
	}
	if _, err := tx.Exec("UPDATE users SET follower_count=follower_count-1 WHERE id=?", memberID); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}
