package store

import (
	"database/sql"
	"encoding/json"
	"fmt"
)

type Service struct {
	tx *sql.Tx
}

func New(tx *sql.Tx) *Service {
	return &Service{
		tx: tx,
	}
}

func (s *Service) GetCursor(UserID string) error {

	filter, err := s.GetFilterSub(UserID)
	if err != nil {
		return err
	}

	query := fmt.Sprintf("DECLARE %s CURSOR FOR select *from news %s order by data_create desc", "cursor"+UserID, filter)

	_, err = s.tx.Exec(query)
	if err != nil {

		return err
	}

	return nil
}

func (s *Service) CloseCursor(UserID string) error {

	query := fmt.Sprintf("ClOSE %s", "cursor"+UserID)

	_, err := s.tx.Exec(query)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetLikesByNewsID(newsID string) ([]Like, error) {

	row := s.tx.QueryRow("select likes from news where id=$1", newsID)

	var likesJSON string
	err := row.Scan(
		&likesJSON,
	)
	if err != nil {
		return nil, err
	}

	u := &Users{}
	err = json.Unmarshal([]byte(likesJSON), u)
	if err != nil {
		return nil, err
	}

	filter := ""
	for i, v := range u.Users {
		if i == 0 {
			filter = fmt.Sprintf("where user_id='%v'", v)
			continue
		}

		filter = filter + fmt.Sprintf(" or user_id='%v'", v)
	}
	fmt.Println("filter", filter)
	likes := make([]Like, 0)
	if filter != "" {

		rows, err := s.tx.Query("select id, first_name, last_name from users " + filter)
		if err != nil {
			return nil, err
		}

		for rows.Next() {
			var like Like
			err := rows.Scan(
				&like.UserID,
				&like.FirstName,
				&like.LastName,
			)

			if err != nil {
				return nil, err
			}

			likes = append(likes, like)
		}
	}

	return likes, nil
}

func (s *Service) GetFetch(UserID, count string) ([]News, error) {

	query := fmt.Sprintf("FETCH %s FROM %s", count, "cursor"+UserID)

	rows, err := s.tx.Query(query)
	if err != nil {
		return nil, err
	}

	newsSlice := make([]News, 0)
	for rows.Next() {
		var news News
		err = rows.Scan(
			&news.ID,
			&news.Title,
			&news.Text,
			&news.Date,
			&news.Likes,
			&news.UserID,
		)
		if err != nil {
			return nil, err
		}
		newsSlice = append(newsSlice, news)
	}

	return newsSlice, nil
}

// get filter with user subscriptions
func (s *Service) GetFilterSub(UserID string) (string, error) {

	row := s.tx.QueryRow("select subscriptions from users where id=$1", UserID)

	var sub string
	err := row.Scan(
		&sub,
	)
	if err != nil {
		return "", err
	}

	st := &Subscriptions{}
	err = json.Unmarshal([]byte(sub), st)
	if err != nil {
		return "", err
	}

	filter := ""
	for i, v := range st.Subscriptions {
		if i == 0 {
			filter = fmt.Sprintf("where user_id='%v'", v)
			continue
		}

		filter = filter + fmt.Sprintf(" or user_id='%v'", v)
	}

	return filter, nil
}

// query := `
// 	INSERT INTO "news" (
// 		 "title", "text_news", "data_create", "likes", "user_id")
// 	VALUES ($1, $2, $3, $4, $5)
// `
// s.db.QueryRow(query, "title", "text", time.Now(), "{}", UserID)

// query := `
// 	INSERT INTO "users" (
// 		 "first_name", "last_name", "password", "subscriptions", "data_create")
// 	VALUES ($1, $2, $3, $4, $5)`

// sq := s.db.QueryRow(query, "anton", "kvetinski", "testuser", `{"members":["1","2"]}`, time.Now())
// errsq := sq.Err()
// if errsq != nil {
// 	log.Fatal(errsq)
// }
