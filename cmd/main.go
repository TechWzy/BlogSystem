package main

import (
	"Project/BlogSystem/internal/api"
	"Project/BlogSystem/internal/boot"
	"fmt"
)

func AddNum(arr *[]int) {
	*arr = append(*arr, 1)
	fmt.Printf("len:%v,cap:%v\n", len(*arr), cap(*arr))
	fmt.Printf("arrPtr:%p\n", arr)
}

func main() {
	boot.InitProject()
	boot.InitMysql()
	boot.InitRedis()
	api.InitRouter()
	boot.CloseRedis()

}
