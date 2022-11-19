package userproto

type GetUserDetailReq struct {
	UserId uint64 `form:"user_id" binding:"required"`
}
