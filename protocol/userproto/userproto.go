package userproto

type GetUserDetailReq struct {
	UserId uint64 `form:"user_id" binding:"required"`
}

type CreateUserReq struct {
	Name  string `form:"name" binding:"required,min=1,max=16"`
	Email string `form:"email" binding:"required,email"`
}

type UserInfo struct {
	UserId uint64 `json:"user_id"`
	Email  string `json:"email"`
	Name   string `json:"name"`
}
