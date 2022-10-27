package ConsistentHashing

import (
	"fmt"
	"hash/crc32"
)

func GetHosts() []string {
	return []string{
		"172.16.0.1:3500",
		"172.16.0.2:3500",
		"172.16.0.3:3500",
		"172.16.0.4:3500",
		"172.16.0.5:3500",
	}
}

func GetHostsHashSlice(hosts []string) []uint32 {
	result := []uint32{}
	for _, v := range hosts {
		// 1. hashcode:= 哈希函数（）
		hashcode := crc32.ChecksumIEEE([]byte(v))
		result = append(result, hashcode)
	}

	fmt.Println("Hosts Hash is:")
	fmt.Println(result)
	return result
}

func GetKeysNum(i int) int {
	return i
}

func GetKeysHash(keynum int) []uint32 {
	/* 	result := []uint32{}
	   	for i:=0;i< keynum;i++ {
	   		// 1. hashcode:= 哈希函数（）
	   		hashcode := crc32.ChecksumIEEE([]byte(i))
	   		result = append(result, hashcode)
	   	}

	   	fmt.Println("Keys Hash is:")
	   	fmt.Println(result)
	   	return result */
	return []uint32{}
}

/* func GetMapkeyToHost(keys []int32 ,hosts []int32)map[int32]int32{

	peers:= map[int32]string{
		// 服务器生活在哈希环上的多个位置
		3479877259: "172.16.0.1:3500",
		3778157554: "172.16.0.1:3500",
		2828126987: "172.16.0.2:3500",
		1973046754: "172.16.0.2:3500",
		1973046755: "172.16.0.3:3500",
		1973046756: "172.16.0.3:3500",

	}

	// 3.在perrSlice 中查找第一个大于等于hashcode的元素

	index:= sort.Search(len(perrSlice) func(i int){
		perrSlice[i]>=hashcode
	})

	peer,ok:= peers[perrSlice[index]]
	if ok{
		fmt.Println(peer)
	}


} */
