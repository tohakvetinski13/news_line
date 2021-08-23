package handler

import "test_task1/store"

type GetLikeResp struct {
	Likes []store.Like `json:"likes"`
}

type GetNewsResp struct {
	News []store.News `json:"news"`
}
