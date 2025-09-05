package entities

type UserDetailRequest struct {
	ID           *int64
	ContactValue *string
}

type UserDetailOption interface {
	Apply(*UserDetailRequest)
}

type userDetailOptionFn func(*UserDetailRequest)

func (fn userDetailOptionFn) Apply(r *UserDetailRequest) { fn(r) }

func WithID(id int64) UserDetailOption {
	return userDetailOptionFn(func(r *UserDetailRequest) { r.ID = &id })
}

func WithContactValue(contactValue string) UserDetailOption {
	return userDetailOptionFn(func(r *UserDetailRequest) { r.ContactValue = &contactValue })
}
