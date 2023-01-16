package ChannelApplicationScenario

import "fmt"

func GoogleSearch(ch chan string) {
	ch <- "golang on Google "
}

func BaiDuSearch(ch chan string) {
	ch <- "golang on BaiDu"
}

func BingSearch(ch chan string) {
	ch <- "golang on Bing"
}

func SearchFind() {
	ch := make(chan string)
	go GoogleSearch(ch)
	go BaiDuSearch(ch)
	go BingSearch(ch)
	t := <-ch
	fmt.Println(t)
}
