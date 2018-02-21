# Tagged article in pocket, send to kindle
If an article (PDF) is tagged #kindle in Pocket, reflow the PDF file using [k2pdfopt](http://www.willus.com/k2pdfopt/) for viewing on the Kindle Paperwhite and send it to your Kindle device.

# Usage
### Install

```
go get github.com/c633/if-pocket-then-kindle/...
```

### Docker deploy
```
docker-compose build
docker rmi $(docker images --filter "dangling=true" -q --no-trunc)
docker-compose up -d
```

# TODO
- Remove PDF files after sending.