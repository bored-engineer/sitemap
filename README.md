# sitemap
A Golang parser and client for the [Sitemap XML format](https://www.sitemaps.org/protocol.html):
```go
urls, err := sitemap.Fetch(context.TODO(), "https://sitemaps.org/sitemap.xml")
if err != nil {
    panic(err)
}
for _, url := range urls {
    log.Println(url.LastModification, url.Location)
}
```
