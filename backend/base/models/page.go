package models

type PageSizeBase struct {
	Page int `query:"page" default:"1"`  //页码
	Size int `query:"size" default:"20"` //每页数量
}

func (bps PageSizeBase) Offset() int {
	return (bps.Page - 1) * bps.Size
}
