package workflow

import aw "github.com/deanishe/awgo"

func NewItem(title, subtitle string, vars map[string]string, valid bool) *aw.Item {
	item := new(aw.Item)
	item.Title(title)
	item.Subtitle(subtitle)
	for k, v := range vars {
		item.Var(k, v)
	}
	item.Valid(valid)
	return item
}

func SendItems(items ...*aw.Item) {
	if len(items) == 0 {
		return
	}
	wf := aw.New()
	wf.Feedback.Items = items
	wf.SendFeedback()
}

func HandleError(e error, msg string) {
	SendItems(NewItem(msg, e.Error(), nil, false))
}
