package message

type Document struct {
	Id string `json:"id"`
}

func NewDocument(id string) *Document {
	return &Document{
		Id: id,
	}
}

type Audio struct {
	Id string `json:"id"`
}

func NewAudio(id string) *Audio {
	return &Audio{
		Id: id,
	}
}

type Image struct {
	Id string `json:"id"`
}

func NewImage(id string) *Image {
	return &Image{
		Id: id,
	}
}

type Video struct {
	Id string `json:"id"`
}

func NewVideo(id string) *Video {
	return &Video{
		Id: id,
	}
}

type Sticker struct {
	Id string `json:"id"`
}

func NewSticker(id string) *Sticker {
	return &Sticker{
		Id: id,
	}
}
