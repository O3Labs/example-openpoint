package helper

import (
	"html/template"
	_ "log"
	"runtime"
	"time"

	"github.com/o3labs/openpoint/platform/log"
)

func TrackFunctionTime(start time.Time) {
	pc := make([]uintptr, 10) // at least 1 entry needed
	runtime.Callers(2, pc)
	f := runtime.FuncForPC(pc[0])
	file, _ := f.FileLine(pc[0])

	elapsed := time.Since(start)
	log.Printf("%s in %s took %s", f.Name(), file, elapsed)
}

func LoadTemplateWithContentFile(files ...string) *template.Template {

	t := template.New("_all")

	var err error
	for _, v := range files {
		fileData, err := Asset("templates/" + v)
		if err != nil {
			log.Errorf("error Asset templates %+v", err)
			return nil
		}
		t, err = t.Parse(string(fileData))
	}

	mustTeamplate := template.Must(t, err)
	if err != nil {
		log.Errorf("error must templates %+v", err)
		return nil
	}
	return mustTeamplate
}
