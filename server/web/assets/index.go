package assets

import "net/http"

func HandleIndex(response http.ResponseWriter, request *http.Request) {
	index := `<!doctype html>
<html>
<head>
    <meta charset="utf-8">
    <title>Websocket example</title>
    <script>
        var exampleSocket = new WebSocket("ws://192.168.56.101:3000/ws");
        exampleSocket.onmessage = function(event) {
            var msg = event.data;

            var existingMessages = document.getElementById("messages").innerHTML;
            var newMessage = "<div>"+msg+"</div>";
            document.getElementById("messages").innerHTML = newMessage + existingMessages;
        };
    </script>
</head>
<body>
<h1>Hello from a static page</h1>
<h2>Messages</h2>
<div id="messages">
</div>
</body>
</html>
`
	response.Write([]byte(index))
}