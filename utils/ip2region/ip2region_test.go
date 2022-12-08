package ip2region

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetRegionByIpList(t *testing.T) {
	regionMap, err := GetRegionByIpList([]string{"171.34.169.122", "101.227.131.220"})
	fmt.Printf("%+v,err:%v\n", regionMap, err)
	if !strings.Contains(regionMap["171.34.169.122"], "南昌") || strings.Contains(regionMap["101.227.131.220"], "上海") {
		assert.Fail(t, err.Error())
		return
	}
	assert.Equal(t, err, nil)
}
