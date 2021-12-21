package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sync"
	"time"

	"github.com/alexanderritola/nntp"
)

var total int
var completed int

func archiveGroup(name string, wg *sync.WaitGroup) {
	c, err := nntp.Dial("tcp", "news.cri.epita.fr:119")
	if err != nil {
		panic(err)
	}

	defer c.Quit()

	g, err := c.Group(name)
	if err != nil {
		panic(err)
	}

	overview, err := c.Overview(g.Low, g.High)
	if err != nil {
		fmt.Println(g.Name, "Error:", err)
		wg.Done()
		return
	}
	total += len(overview)
	for _, a := range overview {
		p := path.Join("output", name, a.MessageId+".news")
		_, err = os.Stat(p)
		if err == nil {
			fmt.Println("Already downloaded:", a.MessageId)
			completed++
			continue
		}

		reader, err := c.ArticleText(a.MessageId)
		if err != nil {
			println(err.Error())
		}

		data, err := ioutil.ReadAll(reader)
		if err != nil {
			println(err.Error())
		}

		os.WriteFile(p, data, 0644)

		fmt.Println("Downloaded:", a.MessageId)
		completed++
	}

	wg.Done()
}

func main() {
	c, err := nntp.Dial("tcp", "news.cri.epita.fr:119")
	if err != nil {
		panic(err)
	}
	defer c.Quit()

	groups, err := c.List()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, group := range groups {
		os.MkdirAll(path.Join("output", group.Name), 0755)

		wg.Add(1)
		fmt.Println("Starting archiving:", group.Name)
		go archiveGroup(group.Name, &wg)
	}

	go func() {
		for {
			time.Sleep(time.Second * 1)
			fmt.Println(completed, "/", total, "news downloaded")
		}
	}()

	wg.Wait()

}
