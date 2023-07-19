package global

type IDValidator struct {
	ID string `uri:"id" validate:"required,numeric"`
}
