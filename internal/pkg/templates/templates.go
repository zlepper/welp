// Code Generated by go run scripts/embedTemplates DO NOT EDIT.

/*
 * Copyright © 2018 Rasmus Hansen
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
 * THE SOFTWARE.
 */

package templates

import "html/template"

type templateContent struct {
	Filename string
	Content  string
}

var contents = []templateContent{

	templateContent{
		Filename: "create-new-user",
		Content:  "<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <title>Create new user</title>\r\n</head>\r\n<body>\r\n{{template \"header\" .AuthState}}\r\n\r\n<style>\r\n    .error {\r\n        color: red;\r\n    }\r\n\r\n    .error.hidden {\r\n        display: none;\r\n    }\r\n</style>\r\n\r\n<form id=\"create-user-form\" action=\"/users/new\" method=\"post\" onsubmit=\"return createUser(event)\">\r\n\r\n    <div>\r\n        <label>\r\n            Name\r\n            <input type=\"text\" name=\"name\" required>\r\n        </label>\r\n    </div>\r\n\r\n    <div>\r\n        <label>\r\n            Email\r\n            <input type=\"email\" name=\"email\" required>\r\n        </label>\r\n    </div>\r\n\r\n    <div>\r\n        <label>\r\n            Password\r\n            <input type=\"password\" name=\"password\" required>\r\n        </label>\r\n    </div>\r\n\r\n    <div>\r\n        <label>\r\n            Repeat Password\r\n            <input type=\"password\" name=\"repeatPassword\" required>\r\n        </label>\r\n    </div>\r\n\r\n    <div id=\"password-no-match-error\" class=\"error hidden\">\r\n        Passwords do not match\r\n    </div>\r\n\r\n    <div>\r\n        <span>Available roles</span>\r\n    {{range .AvailableRoles}}\r\n        <label>\r\n            <input type=\"checkbox\" value=\"{{.Key}}\" name=\"roles\">\r\n        {{.Name}}\r\n        </label>\r\n    {{end}}\r\n    </div>\r\n\r\n    <button type=\"submit\">\r\n        Create\r\n    </button>\r\n</form>\r\n\r\n<script>\r\n    function createUser(event) {\r\n\r\n        event.preventDefault();\r\n\r\n        var form = document.getElementById('create-user-form');\r\n\r\n        var password = form.password.value;\r\n        var repeatPassword = form.repeatPassword.value;\r\n\r\n        var passwordMatchError = document.getElementById('password-no-match-error');\r\n        if (password !== repeatPassword) {\r\n            passwordMatchError.classList.remove('hidden');\r\n            return false;\r\n        } else {\r\n            passwordMatchError.classList.add('hidden');\r\n        }\r\n\r\n        console.log(form);\r\n\r\n        var fd = new FormData(form);\r\n\r\n        var xhr = new XMLHttpRequest();\r\n\r\n        xhr.addEventListener('load', function (event) {\r\n            console.log('response', xhr.responseText);\r\n        });\r\n\r\n        xhr.addEventListener('error', function (event) {\r\n            console.error('Request failed', event);\r\n        });\r\n\r\n        xhr.open('POST', '/users');\r\n\r\n        xhr.setRequestHeader('Accept', 'application/json');\r\n\r\n        xhr.send(fd);\r\n\r\n        return false;\r\n    }\r\n</script>\r\n</body>\r\n</html>",
	},

	templateContent{
		Filename: "edit-user",
		Content:  "<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <title>Edit user</title>\r\n</head>\r\n<body>\r\n\r\n{{template \"header\" .AuthState}}\r\n\r\n<style>\r\n    .error {\r\n        color: red;\r\n    }\r\n\r\n    .hidden {\r\n        display: none;\r\n    }\r\n</style>\r\n\r\n<form id=\"create-user-form\" action=\"/users/{{.User.Email}}/update\" method=\"post\">\r\n\r\n    <div>\r\n        <label>\r\n            Name\r\n            <input type=\"text\" name=\"name\" required value=\"{{.User.Name}}\">\r\n        </label>\r\n    </div>\r\n\r\n    <div>\r\n        <label>\r\n            Email\r\n            <input type=\"email\" name=\"email\" required value=\"{{.User.Email}}\">\r\n        </label>\r\n    </div>\r\n\r\n    <div>\r\n        <span>Available roles</span>\r\n    {{range .AvailableRoles}}\r\n        <label>\r\n            <input type=\"checkbox\" value=\"{{.Key}}\" name=\"roles\" {{if $.User.HasRole .Key }}checked{{end}}>\r\n        {{.Name}}\r\n        </label>\r\n    {{end}}\r\n    </div>\r\n\r\n    <div>\r\n        Email notification schedule\r\n        <div>\r\n            <label>\r\n                <input name=\"emailUpdate\" type=\"radio\" value=\"never\" {{if eq .User.EmailUpdate \"never\"}}checked{{end}}/>\r\n                Never\r\n                <small>(The user will never receive any updates with new feedback, they will have to check the system\r\n                    manually)\r\n                </small>\r\n            </label>\r\n        </div>\r\n        <div>\r\n            <label>\r\n                <input name=\"emailUpdate\" type=\"radio\" value=\"daily\" {{if eq .User.EmailUpdate \"daily\"}}checked{{end}}/>\r\n                Daily\r\n                <small>(The user will a daily digest of all feedback that came in that day)</small>\r\n            </label>\r\n        </div>\r\n        <div>\r\n            <label>\r\n                <input name=\"emailUpdate\" type=\"radio\" value=\"immediately\"\r\n                       {{if eq .User.EmailUpdate \"immediately\"}}checked{{end}}/>\r\n                Immediately\r\n                <small>(The user will receive an email with the feedback the moment it comes in)</small>\r\n            </label>\r\n        </div>\r\n    </div>\r\n\r\n    <button type=\"submit\">\r\n        Update\r\n    </button>\r\n</form>\r\n\r\n<form action=\"/users/{{.User.Email}}/change-password\" method=\"post\">\r\n    <div>\r\n        <label>\r\n            Password\r\n            <input type=\"password\" name=\"password\" required>\r\n        </label>\r\n    </div>\r\n\r\n    <div>\r\n        <label>\r\n            Repeat Password\r\n            <input type=\"password\" name=\"repeatPassword\" required>\r\n        </label>\r\n    </div>\r\n\r\n    <div id=\"password-no-match-error\" class=\"error hidden\">\r\n        Passwords do not match\r\n    </div>\r\n\r\n\r\n    <button type=\"submit\">\r\n        Change password\r\n    </button>\r\n</form>\r\n\r\n\r\n</body>\r\n</html>",
	},

	templateContent{
		Filename: "embed",
		Content:  "<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <title>Feedback</title>\r\n</head>\r\n<body>\r\n\r\n<form method=\"post\" action=\"/\" onsubmit=\"return handleSubmit(event)\" id=\"form\">\r\n    <label>\r\n        Message\r\n        <textarea id=\"message\" name=\"message\" required></textarea>\r\n    </label>\r\n\r\n    <label>\r\n        Contact address\r\n        <input type=\"email\" name=\"contactAddress\" id=\"email\" />\r\n    </label>\r\n\r\n    <label>\r\n        Attach files\r\n        <input type=\"file\" name=\"files\" id=\"files\" multiple>\r\n    </label>\r\n\r\n    <button type=\"submit\">\r\n        Send feedback\r\n    </button>\r\n</form>\r\n\r\n\r\n<script>\r\n\r\n    function handleSubmit(event) {\r\n\r\n        console.log(event);\r\n\r\n        var form = document.getElementById('form');\r\n        var fd = new FormData(form);\r\n\r\n        var xhr = new XMLHttpRequest();\r\n\r\n        xhr.addEventListener('load', function(event) {\r\n           console.log('submitted without issues', event);\r\n        });\r\n\r\n        xhr.addEventListener('error', function(event) {\r\n            console.error('something went wrong when submitting feedback', event);\r\n        });\r\n\r\n        xhr.open('POST', form.action);\r\n\r\n        xhr.setRequestHeader('Accept', 'application/json');\r\n\r\n        xhr.send(fd);\r\n\r\n        return false;\r\n    }\r\n\r\n</script>\r\n\r\n</body>\r\n</html>",
	},

	templateContent{
		Filename: "feedback-created",
		Content:  "<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <title>Feedback created</title>\r\n</head>\r\n<body>\r\nHello WelD!!\r\n\r\n\r\n<script>\r\n    if(window) {\r\n        console.log(\"I'm alive mutha\")\r\n    }\r\n</script>\r\n</body>\r\n</html>",
	},

	templateContent{
		Filename: "feedback-list",
		Content:  "<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <title>Feedback list</title>\r\n</head>\r\n<body>\r\n\r\n{{template \"header\" .AuthState}}\r\n\r\n<main>\r\n{{if .Feedback}}\r\n\r\n    <style type=\"text/css\">\r\n        .flex {\r\n            display: flex;\r\n        }\r\n\r\n        .flex.column {\r\n            flex-direction: column;\r\n        }\r\n\r\n        .flex.row {\r\n            flex-direction: row;\r\n        }\r\n\r\n        .flex.wrap {\r\n            flex-wrap: wrap;\r\n        }\r\n\r\n        .feedback-list {\r\n            margin: 0 1rem;\r\n        }\r\n\r\n        .feedback-item {\r\n            flex: 0 0 auto;\r\n            min-height: 0;\r\n            box-shadow: 0 0 15px rgba(0, 0, 0, .15);\r\n            margin-bottom: 1rem;\r\n            border-radius: 0.5rem;\r\n            padding: 1rem;\r\n            box-sizing: border-box;\r\n        }\r\n\r\n        .feedback-item-header {\r\n            flex: 0 0 3rem;\r\n            align-items: center;\r\n        }\r\n\r\n        .feedback-item-filler {\r\n            flex: 1;\r\n        }\r\n\r\n        .feedback-item-header-button {\r\n            text-decoration: none;\r\n            margin-left: 0.3rem;\r\n            height: 3rem;\r\n            color: white;\r\n            background-color: #B63332;\r\n            border: none;\r\n        }\r\n\r\n        .feedback-item-attachment {\r\n            background: no-repeat center;\r\n            background-size: contain;\r\n            position: relative;\r\n            height: 10rem;\r\n            width: 10rem;\r\n            margin: 1rem;\r\n        }\r\n\r\n        .feedback-item-attachment-button {\r\n            text-decoration: none;\r\n            border: none;\r\n            background-color: #B63332;\r\n            color: white;\r\n            padding: 0.5rem;\r\n            line-height: 1rem;\r\n            box-sizing: border-box;\r\n            font-size: 1rem;\r\n            display: inline-block;\r\n        }\r\n\r\n        .feedback-item-message {\r\n            width: 100%;\r\n            padding: 1rem;\r\n            margin-top: .5rem;\r\n            box-sizing: border-box;\r\n            color: black;\r\n        }\r\n\r\n        .feedback-item-attachments {\r\n            max-width: 100%;\r\n            padding: 1rem;\r\n        }\r\n\r\n        .feedback-item-body-item {\r\n            border: 1px solid #ebebeb;\r\n        }\r\n\r\n        .feedback-item-body-item:first-child {\r\n            border-top-left-radius: 3px;\r\n            border-top-right-radius: 3px;\r\n        }\r\n\r\n        .feedback-item-body-item:last-child {\r\n            border-bottom-left-radius: 3px;\r\n            border-bottom-right-radius: 3px;\r\n        }\r\n\r\n        .contact-address {\r\n            color: black;\r\n        }\r\n    </style>\r\n\r\n    <div class=\"feedback-list flex column\">\r\n    {{range .Feedback}}\r\n        <div class=\"feedback-item flex column\">\r\n            <div class=\"feedback-item-header flex row\">\r\n\r\n            {{if .ContactAddress}}\r\n                <span>From <a class=\"contact-address\" href=\"mailto:{{.ContactAddress}}\">{{.ContactAddress}}</a></span>\r\n            {{else}}\r\n                <span>No contact address provided</span>\r\n            {{end}}\r\n\r\n                <span class=\"feedback-item-filler\"></span>\r\n\r\n                <button class=\"feedback-item-header-button\">\r\n                    Mark as read\r\n                </button>\r\n\r\n                <button class=\"feedback-item-header-button\">\r\n                    Reply\r\n                </button>\r\n\r\n            </div>\r\n\r\n            <div class=\"feedback-item-body flex column\">\r\n                <div class=\"feedback-item-message feedback-item-body-item\">\r\n                {{.Message}}\r\n                </div>\r\n\r\n            {{if .Files}}\r\n                <div class=\"feedback-item-attachments feedback-item-body-item flex row wrap\">\r\n                {{range .Files}}\r\n                {{if .IsImage}}\r\n                    <div class=\"feedback-item-attachment\"\r\n                         style=\"background-image: url(/files/{{.Id}})\">\r\n                        <a href=\"/files/{{.Id}}\" download class=\"feedback-item-attachment-button\">\r\n                            Download\r\n                        </a>\r\n                        <button class=\"feedback-item-attachment-button\">\r\n                            Preview\r\n                        </button>\r\n                    </div>\r\n                {{else}}\r\n                    <div class=\"feedback-item-attachment\">\r\n                        <a href=\"/files/{{.Id}}\" download class=\"feedback-item-attachment-button\">\r\n                            Download\r\n                        </a>\r\n                    </div>\r\n                {{end}}\r\n                {{end}}\r\n                </div>\r\n            {{end}}\r\n            </div>\r\n        </div>\r\n    {{end}}\r\n    </div>\r\n{{else}}\r\n\r\n    <style>\r\n        .no-feedback {\r\n            margin: auto;\r\n        }\r\n    </style>\r\n\r\n    <div class=\"no-feedback\">\r\n        No feedback has been sent so far.\r\n    </div>\r\n\r\n{{end}}\r\n</main>\r\n\r\n</body>\r\n</html>",
	},

	templateContent{
		Filename: "file-not-found",
		Content:  "<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <title>Title</title>\r\n</head>\r\n<body>\r\n\r\n</body>\r\n</html>",
	},

	templateContent{
		Filename: "header",
		Content:  "<div class=\"header\">\r\n\r\n{{if .Authenticated}}\r\n    <a href=\"/\" class=\"header-button\">\r\n        Feedback list\r\n    </a>\r\n\r\n{{if .User.HasRole \"admin\"}}\r\n    <a href=\"/users\" class=\"header-button\">\r\n        Users\r\n    </a>\r\n{{end}}\r\n{{end}}\r\n\r\n    <span class=\"filler\"></span>\r\n\r\n{{if .Authenticated}}\r\n    <a href=\"/logout\" class=\"header-button\" onclick=\"return logout()\">\r\n        Logout\r\n    </a>\r\n{{else}}\r\n    <a href=\"/login\" class=\"header-button\">\r\n        Login\r\n    </a>\r\n{{end}}\r\n</div>\r\n<style>\r\n    body {\r\n        margin: 0;\r\n    }\r\n\r\n    .header {\r\n        display: flex;\r\n        flex-direction: row;\r\n        align-items: center;\r\n        height: 3rem;\r\n        box-sizing: border-box;\r\n    }\r\n\r\n    .filler {\r\n        flex: 1;\r\n    }\r\n\r\n    .header-button {\r\n        text-decoration: none;\r\n        border: none;\r\n        background-color: #B63332;\r\n        color: white;\r\n        padding: 0.5rem;\r\n        line-height: 1rem;\r\n        box-sizing: border-box;\r\n        font-size: 1rem;\r\n        display: flex;\r\n        align-items: center;\r\n        justify-content: center;\r\n        margin-right: 1rem;\r\n    }\r\n\r\n</style>\r\n\r\n<script>\r\n    function logout() {\r\n        localStorage.removeItem('token');\r\n        return true;\r\n    }\r\n</script>",
	},

	templateContent{
		Filename: "login",
		Content:  "<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <title>Title</title>\r\n</head>\r\n<body>\r\n\r\n{{template \"header\" .AuthState}}\r\n\r\n<div>\r\n    <form method=\"post\" action=\"/login\" onsubmit=\"return login(event)\" id=\"login-form\">\r\n\r\n        <label>\r\n            Email\r\n            <input type=\"email\" name=\"email\" id=\"email\">\r\n        </label>\r\n\r\n        <label>\r\n            Password\r\n            <input type=\"password\" name=\"password\" id=\"password\">\r\n        </label>\r\n\r\n        <button type=\"submit\">\r\n            Login\r\n        </button>\r\n\r\n    </form>\r\n</div>\r\n\r\n<script>\r\n    var urlParams = {};\r\n    (function init() {\r\n        var match,\r\n                pl = /\\+/g,  // Regex for replacing addition symbol with a space\r\n                search = /([^&=]+)=?([^&]*)/g,\r\n                decode = function (s) {\r\n                    return decodeURIComponent(s.replace(pl, \" \"));\r\n                },\r\n                query = window.location.search.substring(1);\r\n\r\n        while (match = search.exec(query))\r\n            urlParams[decode(match[1])] = decode(match[2]);\r\n    })();\r\n\r\n    function login(event) {\r\n\r\n        console.log(event);\r\n\r\n        var form = document.getElementById('login-form');\r\n        var fd = new FormData(form);\r\n\r\n        var xhr = new XMLHttpRequest();\r\n\r\n        xhr.addEventListener('load', function (event) {\r\n            var response = JSON.parse(xhr.responseText);\r\n            if (response.token) {\r\n                console.log('login was successful');\r\n                localStorage.setItem('token', response.token);\r\n\r\n                if (urlParams.returnUrl) {\r\n                    window.location.replace(urlParams.returnUrl);\r\n                } else {\r\n                    window.location.replace('/');\r\n                }\r\n\r\n            } else {\r\n                console.log('Got response', event);\r\n            }\r\n        });\r\n\r\n        xhr.addEventListener('error', function (event) {\r\n            console.error('Request failed', event);\r\n        });\r\n\r\n        xhr.open('POST', '/login');\r\n\r\n        xhr.setRequestHeader('Accept', 'application/json');\r\n\r\n        xhr.send(fd);\r\n\r\n        return false;\r\n    }\r\n</script>\r\n\r\n</body>\r\n</html>",
	},

	templateContent{
		Filename: "user-list",
		Content:  "<!DOCTYPE html>\r\n<html lang=\"en\">\r\n<head>\r\n    <meta charset=\"UTF-8\">\r\n    <title>User list</title>\r\n</head>\r\n<body>\r\n\r\n{{template \"header\" .AuthState}}\r\n\r\n<style>\r\n    .user-table {\r\n        width: 100%;\r\n    }\r\n\r\n    .user-table td {\r\n        text-align: center;\r\n    }\r\n\r\n    .user-table .options {\r\n        display: flex;\r\n        flex-direction: row;\r\n    }\r\n\r\n    .option-button {\r\n        text-decoration: none;\r\n        border: none;\r\n        background-color: #B63332;\r\n        color: white;\r\n        padding: 0.5rem;\r\n        line-height: 1rem;\r\n        box-sizing: border-box;\r\n        font-size: 1rem;\r\n        align-items: center;\r\n        justify-content: center;\r\n        margin-right: 0.5rem;\r\n    }\r\n</style>\r\n\r\n<table class=\"user-table\">\r\n    <tr>\r\n        <th>Name</th>\r\n        <th>Email</th>\r\n        <th>Roles</th>\r\n        <th>Options</th>\r\n    </tr>\r\n{{range .Users}}\r\n    <tr>\r\n        <td>{{.Name}}</td>\r\n        <td>{{.Email}}</td>\r\n        <td class=\"role-list\">\r\n        {{range .Roles}}\r\n        {{with $.GetRole .}}\r\n            <span data-key=\"{{.Key}}\">{{.Name}}</span>\r\n        {{end}}\r\n        {{end}}\r\n        </td>\r\n        <td class=\"options\">\r\n        {{if $.IsLastAdmin . | not}}\r\n            <form class=\"admin-danger\" action=\"/users/deleteUser/{{.Email}}\" method=\"post\"\r\n                  onsubmit=\"return deleteUser(event, {{.Email}})\">\r\n                <button type=\"submit\" class=\"option-button\">\r\n                    Delete\r\n                </button>\r\n            </form>\r\n        {{end}}\r\n\r\n            <a href=\"/users/{{.Email}}\" class=\"option-button\">\r\n                Edit\r\n            </a>\r\n        </td>\r\n    </tr>\r\n{{end}}\r\n</table>\r\n<a href=\"/users/new\" class=\"option-button\">\r\n    Create new user\r\n</a>\r\n\r\n<script>\r\n    function deleteUser(event, email) {\r\n        var xhr = new XMLHttpRequest();\r\n\r\n        xhr.addEventListener('load', function (event) {\r\n            console.log('response', event);\r\n\r\n            deleteUserRow(event);\r\n        });\r\n\r\n        xhr.addEventListener('error', function (event) {\r\n            console.error('Request failed', event);\r\n        });\r\n\r\n        xhr.open('POST', '/login');\r\n\r\n        xhr.setRequestHeader('Accept', 'application/json');\r\n        xhr.setRequestHeader(\"Content-Type\", \"application/json;charset=UTF-8\");\r\n\r\n        xhr.send(JSON.stringify({email: email}));\r\n\r\n        return false;\r\n    }\r\n\r\n    function deleteUserRow(event) {\r\n        var target = event.currentTarget;\r\n\r\n        while (target && target.tagName !== 'TR') {\r\n            target = target.parentElement;\r\n        }\r\n\r\n        if (target) {\r\n            target.parentNode.removeChild(target);\r\n        }\r\n        console.log('removed user row');\r\n\r\n        var adminRows = [];\r\n        var roleListElements = document.querySelectorAll('.role-list');\r\n        for (var i = 0; i < roleListElements.length; i++) {\r\n            var rle = roleListElements[i];\r\n            var roleElements = rle.querySelectorAll('span');\r\n\r\n            for (var j = 0; j < roleElements.length; j++) {\r\n                var el = roleElements[i];\r\n                var key = el.dataset.key;\r\n                if (key === 'admin') {\r\n                    adminRows.push(rle.parentNode);\r\n                    break;\r\n                }\r\n            }\r\n        }\r\n\r\n        console.log('admin rows', adminRows);\r\n\r\n        // Remove the last delete forms\r\n        if (adminRows.length === 1) {\r\n            var forms = adminRows[0].querySelectorAll('admin-danger');\r\n            for (i = 0; i < forms.length; i++) {\r\n                var form = forms[i];\r\n                form.parentNode.removeChild(form);\r\n            }\r\n        }\r\n    }\r\n</script>\r\n\r\n</body>\r\n</html>",
	},
}

func getTemplates() (*template.Template, error) {
	var t *template.Template

	for _, file := range contents {

		filename := file.Filename
		var tmpl *template.Template
		if t == nil {
			t = template.New(filename)
		}

		if filename == t.Name() {
			tmpl = t
		} else {
			tmpl = t.New(filename)
		}
		_, err := tmpl.Parse(file.Content)
		if err != nil {
			return nil, err
		}
	}

	return t, nil
}
