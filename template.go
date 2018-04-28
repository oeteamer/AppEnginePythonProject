package base

import (
	"html/template"
)

var templateFuncMap = template.FuncMap{
	"getAuthorLink": getAuthorLink,
	"formatTime":    formatTime,
}

var AuthorsTemplate = template.Must(template.New("authors").Funcs(templateFuncMap).Parse(`
<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" type="text/css" href="/stylesheets/dev.css">
</head>
<body>
    <table>
          	<tr>
                <th>Авторы</th>
                <th>Samlib link</th>
          	</tr>
         	{{ range $code, $author := . }}
	        	<tr>
	         	   	<td><a href="/author/{{ $code }}">{{ $author.Name }}</a></td>
	            	<td><a href="{{ getAuthorLink $code }}" target="_blank">Samlib</a></td>
	          	</tr>
          	{{ end }}
    </table>
</body>
</html>
`))

var AuthorTemplate = template.Must(template.New("author").Funcs(templateFuncMap).Parse(`
<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" type="text/css" href="/stylesheets/dev.css">
</head>
<body>
    <h3 align="center">{{ .Name }}</h3>
    <table style="border: 1px;">
        <tr>
            <th>Произведение</th>
            <th>Объем</th>
            <th>Updated</th>
            <th>Updated At</th>
        </tr>
        {{ range $book := .Books }}
	        <tr>
				<td><a href="{{ $book.Href }}" target="_blank">{{ $book.Name }}</a></td>
				<td align="right">{{ $book.Volume }}</td>
				<td style="text-align: right; color: green;">{{ $book.UpdateInfo }}
				<td>{{ formatTime $book.UpdatedAt }}</td>
	        </tr>
       	{{ end }}
    </table>
</body>
</html>
`))

var UpdatedBooksTemplate = template.Must(template.New("updatedBooks").Funcs(templateFuncMap).Parse(`
<!DOCTYPE html>
<html>
<head>
    <link rel="stylesheet" type="text/css" href="/stylesheets/dev.css">
    <script src="/js/zepto.min.js"></script>
    <script>
        $(document).ready(
            getTasksInProgressCount()
        );
        function getTasksInProgressCount() {
            $.ajax({
                type: 'GET',
                url: '/task-stats',
                success: function(data){
                    $('#tasksInProgress').text(data);
                },
                error: function(xhr, type){
                    alert(type);
                }
            });
        }

		function authorsReload() {
            $.ajax({
                type: 'GET',
                url: '/reload-authors',
                success: function(data){
                    alert('authors reloaded');
                }
            });
        }

        function updateList() {
            $.ajax({
                type: 'GET',
                url: '/update-all'
            });
        }
    </script>
</head>
<body>
    <h3 align="center">Обновленные книги</h3>
    <div style="float:left;">
        <button onclick="getTasksInProgressCount()">Количество задач в очереди:</button><span id="tasksInProgress">0</span>
    </div>
    <div style="float:right;">
        <button onclick="updateList()">Обновить</button>
    </div>
    <div style="clear: both"></div>
    <table>
        <tr>
            <th>Автор</th>
            <th>Произведение</th>
            <th>Объем</th>
            <th>Update</th>
            <th>Updated At</th>
        </tr>
        {{ range $author := . }}
			{{ range $book := $author.Books }}
		        <tr>
		            <td><a href="/author/{{ $author.Code }}" target="_blank">{{ $author.Name }}</a></td>
		            <td><a href="{{ $book.Href }}" target="_blank">{{ $book.Name }}</a></td>
		            <td align="right">{{ $book.Volume }}</td>
		            <td align="right">{{ $book.UpdateInfo }}</td>
		            <td align="right">{{ formatTime $book.UpdatedAt }}</td>
		        </tr>
			{{ end }}
        {{ end }}
    </table>
 	<div style="float:left;">
        <button onclick="authorsReload()">обновить авторов в памяти</button>
    </div>
</body>
</html>
`))
