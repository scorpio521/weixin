package tools

import (
	"math/rand"
	)

//写type的时候第一个字母要大写
type Arrcommonfunc struct {
}

//思想:
//运用map,统计nums1中值出现的次数-map[值]次数
//遍历nums2中的值,查看值是否在map中的出现
//_, ok := demo["a"] 判断map 中是否存在 if ok {}
//if _, ok := myMap[num]; ok {
//方法名和都要大写

//并集
func (c *Arrcommonfunc)Complete_union(nums1 []int64, nums2 []int64) (res []int64) {
	m := make(map[int64]int64)
	res = nums1
	for _,v := range res {
		m[v]++
	}
	for _,v := range nums2 {
		times, _ := m[v]  //v是nums2中的值,m[v]是map中的值.m[v]==times
		if times == 0{
			res = append(res, v)
		}
	}
	return res
}

func (c *Arrcommonfunc)Intersect(nums1 []int64, nums2 []int64) (res []int64) {
	m := make(map[int64]int64)//make 是给map用的
	for _,v := range nums1 {
		m[v]++
	}
	for _,v := range nums2 {
		times, _ := m[v]  //v是nums2中的值,m[v]是map中的值.m[v]==times
		if times != 0{
			res = append(res,v)
			// nums1 = append(nums1, v)
		}
	}
	return res
}

//获取的是nums2 的差集，也就是nums2 数组要比nums1 数据多才行
func (c *Arrcommonfunc)Diff(nums1 []int64, nums2 []int64) (res []int64) {
	m := make(map[int64]int64)//make 是给map用的
	for _,v := range nums1 {
		m[v]++
	}
	for _,v := range nums2 {
		times, _ := m[v]  //v是nums2中的值,m[v]是map中的值.m[v]==times
		if times == 0{
			res = append(res,v)
			// nums1 = append(nums1, v)
		}
	}
	return res
}

func (c *Arrcommonfunc)Random(strings []int64, length int) (res []int64, err string) {
	if len(strings) <= 0 {
		err = "the length of the parameter strings should not be less than 0"
		return
	}

	if length <= 0 || len(strings) <= length {
		err = "the size of the parameter length illegal"
		return
	}

	for i := len(strings) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		strings[i], strings[num] = strings[num], strings[i]
	}

	for i := 0; i < length; i++ {
		res = append(res,strings[i])
	}
	return
}
