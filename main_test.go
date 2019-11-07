package main

import (
	"fmt"
	"testing"
)

func TestIPParse(t *testing.T){
	//
	//for _,x :=range ListReader("/Users/yunxu/iplist.txt","ip"){
	//
	//	fmt.Println(x)
	//}


	aa,err:=PortParse("3389,80,1-33")
	fmt.Println(aa)
	fmt.Println(err)

}
func RemoveRep(s []string) []string {
	result := []string{}
	m := make(map[string]bool) //map的值不重要
	for _, v := range s {
		if _, ok := m[v]; !ok {
			result = append(result, v)
			m[v] = true
		}
	}
	return result
}

func TestRemoveRep(t *testing.T){
	fmt.Println(RemoveRep([]string{"1","1","3","6"}))
}