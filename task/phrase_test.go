package task

import (
	"testing"

	"github.com/blendlabs/go-util/uuid"
)

func TestAddPhrase(t *testing.T) {
	p := Phrase{}
	p.Phrase = "主人我错了啦,再也不凶你了,请去工作QwQ"
	p.UUID = uuid.V4().String()
	p.Show = "on"
	p.GroupUUID = uuid.V4().String()
	InsertPhrase(DB, p)
}
