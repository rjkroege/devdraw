/*
	Generate a JavaScript file for output that is easier to
	read.
*/

package main

const (

visualizer_prefx = 
`<html>
<head>
	<style type="text/css">
pre {outline: 1px solid #ccc; padding: 5px; margin: 5px; }
.string { color: green; }
.number { color: darkorange; }
.boolean { color: blue; }
.null { color: magenta; }
.key { color: red; }
	 </style>	

	<script language="JavaScript" type="text/javascript">

function output(inp) {
	var v = document.createElement('pre')
	 document.body.appendChild(v).innerHTML = inp;
}

function syntaxHighlight(json) {
    json = json.replace(/&/g, '&amp;').replace(/</g, '&lt;').replace(/>/g, '&gt;');
    return json.replace(/("(\\u[a-zA-Z0-9]{4}|\\[^u]|[^\\"])*"(\s*:)?|\b(true|false|null)\b|-?\d+(?:\.\d*)?(?:[eE][+\-]?\d+)?)/g, function (match) {
        var cls = 'number';
        if (/^"/.test(match)) {
            if (/:$/.test(match)) {
                cls = 'key';
            } else {
                cls = 'string';
            }
        } else if (/true|false/.test(match)) {
            cls = 'boolean';
        } else if (/null/.test(match)) {
            cls = 'null';
        }
        return '<span class="' + cls + '">' + match + '</span>';
    });
}
	</script>


	<script language="JavaScript" type="text/javascript">
`

visualizer_suffix =
`
</script>
</head>
<body onload="output(syntaxHighlight(JSON.stringify(obj, undefined, 4)));"></body>
</html>
`

)