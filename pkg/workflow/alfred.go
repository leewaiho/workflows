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

type Workflow interface {
	AddItem(items ...*aw.Item)
	AddVar(k, v string)
	Send()
	Error(e error, msg string)
}

type workflow struct {
	wf    *aw.Workflow
	items []*aw.Item
	done  bool
}

func New() Workflow {
	return &workflow{
		wf:    aw.New(),
		items: []*aw.Item{},
		done:  false,
	}
}

func (w *workflow) AddVar(k, v string) {
	w.wf.Feedback.Var(k, v)
}

func (w *workflow) AddItem(items ...*aw.Item) {
	w.items = append(w.items, items...)
}

func (w *workflow) Send() {
	if !w.done {
		w.send(w.items)
	}
}

func (w *workflow) Error(e error, msg string) {
	if !w.done {
		w.send([]*aw.Item{NewItem(msg, e.Error(), nil, false)})
	}
}

func (w *workflow) send(items []*aw.Item) {
	w.wf.Feedback.Items = items
	w.wf.SendFeedback()
}
