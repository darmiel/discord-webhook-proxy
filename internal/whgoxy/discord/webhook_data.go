package discord

import (
	"bytes"
	"log"
	"text/template"
)

type WebhookData string

// Exec parses the template
func (d WebhookData) Exec(param ...interface{}) (json string, err error) {
	// replace params in data
	var parse *template.Template
	parse, err = template.New("").Parse(string(d))
	if err != nil {
		return
	}

	// parse data
	var data interface{}
	if param != nil && len(param) >= 1 {
		data = param[0]
	}

	// execute template
	var buffer bytes.Buffer
	if err = parse.Execute(&buffer, data); err != nil {
		log.Println("parsing err:", err, data)
		return
	}

	// read string from buffer
	json = buffer.String()
	return
}
