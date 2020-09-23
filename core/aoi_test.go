package core

import (
	"fmt"
	"testing"
)

func TestNewAOIManager(t *testing.T) {
	// 初始化 AOIManager
	aoiMgr := NewAOIManager(0, 250, 5, 0, 250, 5)

	// 打印 AOIManager
	fmt.Println(aoiMgr)
}

func TestGetSurroundGridsByGid(t *testing.T) {
	// 初始化 AOIManager
	aoiMgr := NewAOIManager(0, 250, 5, 0, 100, 2)

	// 打印
	fmt.Println(aoiMgr.GetSurroundGridsByGid(4))
}
