package test

import (
	"fmt"
	"github.com/go-zookeeper/zk"
	"testing"
	"time"
)

func TestZookeeper(t *testing.T) {
	conn,_,err:=zk.Connect([]string{"127.0.0.1"},5*time.Second)
	if err != nil {
		panic(err)
	}
	fmt.Println("connect success")
	if _,err:=conn.Create("/test_tree2",[]byte("tree_content"),0,zk.WorldACL(zk.PermAll));err!=nil{
		fmt.Println("create err :",err)
	}
	nodeValue,dStat,err := conn.Get("/test_tree2")
	if err != nil {
		fmt.Println("get err:",err)
		return
	}
	fmt.Println("nodeValue:",string(nodeValue))

	if _,err := conn.Set("/test_tree2",[]byte("new_content"),dStat.Version); err != nil {
		fmt.Println("update err :",err)
	}

	_,dStat,_ =conn.Get("/test_tree2")
	if err:=conn.Delete("/test_tree2",dStat.Version) ;err != nil {
		fmt.Println("delete err:",err)
	}

	hasNode ,_,err:= conn.Exists("/test_tree2")
	if err != nil {
		fmt.Println("Exists err: ",err)
	}
	fmt.Println("node exists :",hasNode)

}