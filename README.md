
![11](image/demo.jpg)
# WpGo
## Wordpress Batch Brute Force by Go

```
Usage of wpgo:
  -a string
        attack type  login / xmlrpc (default "login")
  -c int
    	max auto get user count (default 5)
  -o string
    	out filepath (default "result.txt")
  -p string
    	password list filepath
  -t int
    	max thread (default 20)
  -u string
    	username list filepath
  -w string
    	website list filepath

```

auto get author and try login


```
go run main.go -p pass.txt -w site.txt
WpGo.exe -p pass.txt -w site.txt

WpGo.exe -a xmlrpc -p ./dict/p.txt -u ./dict/u.txt -w ./dict/site.txt

```

auto get author and custom author list

```
WpGo.exe -p pass.txt -w site.txt -u user.txt
```
