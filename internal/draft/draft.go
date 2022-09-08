package draft

import (
	"fmt"
	"os"
	"text/template"
	"time"

	"github.com/gookit/config/v2"
)

type Info struct {
	Title       string
	Filename    string
	Description *string
	Date        *time.Time
	Draft       *bool
	Content     string
}

func (info *Info) fill_defaults() {
	if info.Description == nil {
		empty := ""
		info.Description = &empty
	}

	if info.Date == nil {
		now := time.Now()
		info.Date = &now
	}

	if info.Draft == nil {
		t := true
		info.Draft = &t
	}
}

func GenerateAndSaveDraftFile(info *Info, PWD string) (bool, error) {
	t, err := template.ParseFiles(PWD + "/internal/draft/draft.tmpl")
	if err != nil {
		return false, err
	}

	contentPath := config.String("contentPath")

	f, err := os.Create(fmt.Sprintf("%s/%s", contentPath, info.Filename))
	if err != nil {
		return false, err
	}
	defer f.Close()

	info.fill_defaults()

	err = t.Execute(f, info)
	if err != nil {
		return false, err
	}

	return true, nil
}
