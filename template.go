package base

import (
	"html/template"
)

var guestbookTemplate = template.Must(template.New("book").Parse(`
<html>
  <head>
    <title>Authors</title>
  </head>
  <body>
    {{ . }}
  </body>
</html>
`))
