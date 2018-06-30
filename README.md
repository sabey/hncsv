# Golang Hacker News CSV

Prints the first top 20 stories from the front page of https://news.ycombinator.com

Print to stdout:
```
go build -a && ./hncsv
```

Save to path:
```
go build -a && ./hncsv -csv="unittest/result.csv"
```
