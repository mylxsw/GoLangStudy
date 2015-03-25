<html>
<head>
<title> Test Page </title>
</head>

<body>
    <form action="/login" method="post">
        Username: <input type="text" name="username" /><br />
        Password: <input type="password" name="password" /><br />
        <input type="submit" value="Login" />
        <input type="hidden" name="token" value="{{.}}" />
    </form>
</body>
</html>
