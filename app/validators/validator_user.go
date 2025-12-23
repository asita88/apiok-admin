package validators

type UserRegister struct {
	RePassword string `json:"re_password" zh:"确认密码" en:"Confirm Password" binding:"required,eqfield=Password"`
	Password   string `json:"password" zh:"密码" en:"Password" binding:"required,min=8"`
	Email      string `json:"email" zh:"邮箱" en:"Email" binding:"required,email"`
	Name       string `json:"name" zh:"昵称" en:"User name" binding:"required,min=1,max=20"`
}

type UserLogin struct {
	Password string `json:"password" zh:"密码" en:"Password" binding:"required,min=8"`
	Username string `json:"username" zh:"用户名" en:"Username" binding:"required,min=1"`
}

type UserChangePassword struct {
	OldPassword string `json:"old_password" zh:"旧密码" en:"Old Password" binding:"required,min=8"`
	Password    string `json:"password" zh:"新密码" en:"New Password" binding:"required,min=8"`
	RePassword  string `json:"re_password" zh:"确认新密码" en:"Confirm New Password" binding:"required,eqfield=Password"`
}

type UserAddUpdate struct {
	Name     string `json:"name" zh:"昵称" en:"User name" binding:"required,min=1,max=20"`
	Email    string `json:"email" zh:"邮箱" en:"Email" binding:"required,email"`
	Password string `json:"password" zh:"密码" en:"Password" binding:"omitempty,min=8"`
}

type UserList struct {
	Search string `form:"search" json:"search" zh:"搜索内容" en:"Search content" binding:"omitempty"`
	BaseListPage
}
