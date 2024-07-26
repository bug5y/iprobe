# iprobe

## Install

```
▶ go install github.com/bug5y/iprobe@latest
```

## Example Usage

Show the help menu
```
▶ iprobe -h                                                              
Usage of ./iprobe:
  -h    Show help
  -header string
        Custom user agent to add to each request (default "Mozilla/5.0 (Windows NT 10.0; Android 13; Mobile; rv:120.0) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Safari/537.36 OPR/95.0.0.0")
  -i string
        Input (file path or IP:port)
  -p string
        Comma-separated list of protocols (default "http,https")
  -t int
        Max threads (default 20)
  -timeout duration
        Timeout duration (default 10s)
```

Supply single ip port pair using stdin
```
▶ echo 8.8.8.8:443 | ./iprobe
[200] https://8.8.8.8:443
```

Supply a file of ip port pairs using stdin (One pair per line)
```                                                                                                              
▶ cat example.txt | ./iprobe                                                          
[200] https://8.8.8.8:443
[404] https://69.147.65.252:443
[200] http://142.250.72.78:80
[200] https://142.250.72.78:443
```

Supply a file using the input flag
```
▶ iprobe -i example.txt   
[200] https://8.8.8.8:443
[404] https://69.147.65.252:443
[200] http://142.250.72.78:80
[200] https://142.250.72.78:443
```

Only probe for http for all pairs in the supplied file with threads set to 10
```                                                                                                                                                                           
▶ iprobe -i example.txt -p http -t 10
[200] http://142.250.72.78:80
```
