package downloader

import (
	"github.com/hegeng1212/pholcus/app/downloader/request"
	"github.com/hegeng1212/pholcus/app/spider"
)

// The Downloader interface.
// You can implement the interface by implement function Download.
// Function Download need to return Page instance pointer that has request result downloaded from Request.
type Downloader interface {
	Download(*spider.Spider, *request.Request) *spider.Context
}
