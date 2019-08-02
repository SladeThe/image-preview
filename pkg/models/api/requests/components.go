package requests

type ImageCollection struct {
	Images []string `json:"images" form:"image" query:"image" validate:"required,min=1,max=100,dive,required,min=1,max=100000000"`
}
