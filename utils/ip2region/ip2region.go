package ip2region

import (
	"github.com/lionsoul2014/ip2region/binding/golang/xdb"
)

var regionSearcher *xdb.Searcher

func init() {
	var err error
	// 将文件读取加载到内存
	file, err := xdb.LoadContentFromFile("./data/ip2region.xdb")
	if err != nil {
		panic(err)
	}
	regionSearcher, err = xdb.NewWithBuffer(file)
	if err != nil {
		panic(err)
	}
}

func GetRegionByIpList(ipList []string) (map[string]string, error) {
	regionMap := make(map[string]string, len(ipList))
	for _, ip := range ipList {
		res, err := regionSearcher.SearchByStr(ip)
		if err != nil {
			return regionMap, err
		}
		regionMap[ip] = res
	}
	return regionMap, nil
}
