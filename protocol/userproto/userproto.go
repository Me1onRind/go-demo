package userproto

type GetUserDetailReq struct {
	UserId uint64 `form:"user_id" binding:"required"`
}

type CreateUserReq struct {
	Name string `form:"user_id" binding:"required,min=1,max=16"`
}
