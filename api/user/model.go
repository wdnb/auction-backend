package user

type User struct {
	ID       uint32 `db:"id" json:"id"`
	Username string `db:"username" json:"username" validate:"min=6,max=30"`
	Password string `json:"password" validate:"required,min=6,max=20,containsany=!@#$%^&*()_+"`
	Phone    string `db:"phone" json:"phone" validate:"e164"`
	Email    string `db:"email" json:"email" validate:"email"`
}

type Login struct {
	Username string `db:"username" json:"username" validate:"required,min=6,max=30"`
	Password string `json:"password" validate:"required,min=6,max=20,containsany=!@#$%^&*()_+"`
}

type UpdateUser struct {
	ID       uint32 `db:"id" json:"-" validate:"numeric"`
	Username string `db:"username" json:"username" validate:"min=6,max=30"`
	Password string `json:"password" validate:"min=6,max=20,containsany=!@#$%^&*()_+"`
	Email    string `json:"email" validate:"email"`
}

type CreateUser struct {
	ID       uint32  `db:"id" json:"id"`
	Username *string `json:"username" validate:"required"`
	Password string  `json:"password" validate:"required,min=6,max=20,containsany=!@#$%^&*()_+"`
	Phone    *string `db:"phone" json:"phone" validate:"e164"`
	Email    *string `db:"email" json:"email" validate:"email"`
}

type VerificationCode struct {
	ID    uint32 `db:"id" json:"-"`
	Phone string `db:"phone" json:"phone" validate:"e164"`
	Code  string `json:"code" validate:"required"`
	Kind  string `json:"kind" validate:"oneof=login register reset_password"`
}

type KindValidator struct {
	Kind string `uri:"kind" validate:"oneof=login register reset_password"`
}

type CodeSendUser struct {
	Phone string `db:"phone" json:"phone" validate:"e164"`
	Code  string `db:"code" json:"-"`
}

type RoleUser struct {
	User User
	Role Role
}

type Role struct {
	//ID       uint32 `db:"id" json:"id"`
	PType    string `db:"p_type" json:"p_type"`
	RoleName string `db:"role_name" json:"role_name"`
	V0       string `db:"v0" json:"v0"`
	V1       string `db:"v1" json:"v1"`
	V2       string `db:"v2" json:"v2"`
}

type CreateShippingAddress struct {
	UserID        uint32 `db:"user_id" json:"-" validate:"number"`
	IsActive      bool   `db:"is_active" json:"is_active" validate:"boolean"`
	Phone         string `db:"phone" json:"phone" validate:"required,e164"`
	RecipientName string `db:"recipient_name" json:"recipient_name" validate:"required,min=2,max=120"`
	Region        string `db:"region" json:"region" validate:"required,min=2,max=120"`
	Address       string `db:"address" json:"address" validate:"required,min=2,max=120"`
}

type ShippingAddress struct {
	ID            uint32 `db:"id" json:"id"`
	UserID        uint32 `db:"user_id" json:"-" validate:"number"`
	IsActive      bool   `db:"is_active" json:"is_active" validate:"boolean"`
	Phone         string `db:"phone" json:"phone" validate:"required,e164"`
	RecipientName string `db:"recipient_name" json:"recipient_name" validate:"required,min=2,max=120"`
	Region        string `db:"region" json:"region" validate:"required,min=2,max=120"`
	Address       string `db:"address" json:"address" validate:"required,min=2,max=120"`
}
