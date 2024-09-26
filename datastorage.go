/*
实现一个简易的键对值数据存储系统，使用命令行接收一些命令并将其对应的数据存储在JSON文件中

	为 SET 命令和 SETNX 命令 新增⼀个 过期时间 的参数 并在 GET 时判断 如果已经过期 则删除这⼀

条数据并提⽰给⽤⼾。
*/
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	// "reflect"
)

var date map[string]time.Time = make(map[string]time.Time)

func showmenu() {
	fmt.Println("----------使用说明----------")
	fmt.Println("请输入对应序号选择所要实现的功能：")
	fmt.Println("1.SET key value - 设置键值对")
	fmt.Println("2.SETNX key value - 设置键值对，如果键不存在则设置")
	fmt.Println("3.GET key - 获取键对应的值")
	fmt.Println("4.DEL key - 删除键值对")
	fmt.Println("5.SADD SetName value - 添加元素到集合")
	fmt.Println("6.SMEMBER SetName - 获取集合中的所有元素")
	fmt.Println("7.LPUSH ListName value - 左侧添加元素到列表")
	fmt.Println("8.LRANGE ListName start end - 获取列表中指定范围的元素")
	fmt.Println("0.退出程序")
	fmt.Println("------------------------")

}
func set() {
	var m map[string]interface{} = make(map[string]interface{})
	fmt.Println("请输入命令：")
	var commend, key, value string
	fmt.Scanln(&commend, &key, &value)

	if commend == "SET" || commend == "set" || commend == "Set" {
		var Duration time.Duration
		fmt.Println("请设置过期时间（单位：秒）：")
		fmt.Scanln(&Duration)
		currenttime := time.Now()
		expirationDuration := time.Duration(Duration * time.Second) //注意函数参数要求的类型与变量类型要一致
		expirationDate := currenttime.Add(expirationDuration)
		m[key] = value
		// m[key+"_expirationDate"] = expirationDate--不将过期时间也存入JSON文件中是因为JSON反序列化后得到的expirationDate为string类型而非time.Time
		date[key] = expirationDate
		// fmt.Println(reflect.TypeOf(expirationDate))
		marshal, err := json.Marshal(m)
		if err != nil {
			fmt.Println("json errors", err)
			return
		}
		open, err := os.OpenFile("server.json", os.O_RDWR|os.O_CREATE, 0666) //打开或创建文件,o666表示文件权限--可读写
		if err != nil {
			fmt.Println("open error", err)
			return
		}

		err = os.WriteFile("server.json", marshal, 0666)
		if err != nil {
			fmt.Println("write error", err)
			return
		} else {
			fmt.Println("设置键值对成功")
		}
		defer open.Close()

	} else {
		fmt.Println("命令格式不合理")
	}

}
func setnx() {
	var m map[string]interface{} = make(map[string]interface{})
	fmt.Println("请输入命令：")
	var commend, key, value string
	fmt.Scanln(&commend, &key, &value)
	if commend == "SETNX" || commend == "setnx" || commend == "Setnx" {
		open, err := os.OpenFile("server.json", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("open error", err)
			return
		}
		content, err := os.ReadFile("server.json")
		if err != nil {
			fmt.Println("read error", err)
			return
		}

		json.Unmarshal(content, &m)
		if _, ok := m[key]; !ok { //判断key是否存在于map中，存在则返回对应的value和true，不存在则返回nil和false

			m[key] = value
			marshal, err := json.Marshal(m)
			if err != nil {
				fmt.Println("json errors", err)
				return
			}
			err = os.WriteFile("server.json", marshal, 0666)
			if err != nil {
				fmt.Println("write error", err)
				return
			} else {
				fmt.Println("设置键值对成功")
			}
		} else {
			fmt.Println("键已存在，设置失败")
		}
		defer open.Close()

	} else {
		fmt.Println("命令格式不合理")
	}

}

// 获取元素
func get() {
	var m map[string]interface{} = make(map[string]interface{})
	fmt.Println("请输入命令：")
	var commend, key string
	fmt.Scanln(&commend, &key)
	if commend == "GET" || commend == "get" || commend == "Get" {

		open, err := os.OpenFile("server.json", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("open error", err)
			return
		}
		content, err := os.ReadFile("server.json")
		if err != nil {
			fmt.Println("read error", err)
			return
		}

		json.Unmarshal(content, &m)
		// currenttime := time.Now()
		// expirationDate := m[key+"_expirationDate"]
		// fmt.Println(reflect.TypeOf(expirationDate))---得string,变量不可类型转换，也不可类型断言为time.Time
		if _, ok := m[key]; ok {
			currenttime := time.Now()
			expirationDate := date[key]
			if currenttime.After(expirationDate) {
				delete(m, key)
				marshal, err := json.Marshal(m)
				if err != nil {
					fmt.Println("json errors", err)
					return
				}
				err = os.WriteFile("server.json", marshal, 0666)
				if err != nil {
					fmt.Println("write error", err)
					return
				}
				fmt.Println("键已过期，已被删除")
				return
			}
			fmt.Println("该键对应的值为:", m[key])
		} else {
			fmt.Println("不存在该键")
		}
		defer open.Close()
	} else {
		fmt.Println("命令格式不合理")
	}

}

// 删除元素
func del() {
	var m map[string]interface{} = make(map[string]interface{})
	fmt.Println("请输入命令：")
	var commend, key string
	fmt.Scanln(&commend, &key)
	if commend == "DEL" || commend == "del" || commend == "Del" {
		open, err := os.OpenFile("server.json", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("open error", err)
			return
		}
		content, err := os.ReadFile("server.json")
		if err != nil {
			fmt.Println("read error", err)
			return
		}

		json.Unmarshal(content, &m)

		if _, ok := m[key]; ok {
			delete(m, key)
			marshal, err := json.Marshal(m)
			if err != nil {
				fmt.Println("json errors", err)
				return
			}
			err = os.WriteFile("server.json", marshal, 0666)
			if err != nil {
				fmt.Println("write error", err)
				return
			} else {
				fmt.Println("删除键值对成功")
			}

		} else {
			fmt.Println("不存在该键")
		}
		defer open.Close()

	} else {
		fmt.Println("命令格式不合理")
	}
}

//向一个集合中添加元素value

func sadd() {
	// var m map[string][]interface{} = make(map[string][]interface{})--此种写法没有给map中的切片分配空间，导致后面添加元素时会panic
	var m map[string]interface{} = make(map[string]interface{})
	fmt.Println("请输入命令：")
	var commend, setname, value string
	fmt.Scanln(&commend, &setname, &value)
	if commend == "SADD" || commend == "sadd" || commend == "SAdd" {
		open, err := os.OpenFile("server.json", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("open error", err)
			return
		}
		content, err := os.ReadFile("server.json")
		if err != nil {
			fmt.Println("read error", err)
			return
		}

		json.Unmarshal(content, &m)
		if _, ok := m[setname]; !ok {
			m[setname] = []interface{}{value}
			marshal, err := json.Marshal(m)
			if err != nil {
				fmt.Println("json errors", err)
				return
			}
			err = os.WriteFile("server.json", marshal, 0666)
			if err != nil {
				fmt.Println("write error", err)
				return
			} else {
				fmt.Println("添加元素到集合成功")
			}
		} else {
			m[setname] = append(m[setname].([]interface{}), value)
			marshal, err := json.Marshal(m)
			if err != nil {
				fmt.Println("json errors", err)
				return
			}
			err = os.WriteFile("server.json", marshal, 0666)
			if err != nil {
				fmt.Println("write error", err)
				return
			} else {
				fmt.Println("添加元素到集合成功")
			}
		}
		defer open.Close()
	} else {
		fmt.Println("命令格式不合理")
	}

}

//获取集合中的所有元素

func smember() {
	var m map[string]interface{} = make(map[string]interface{})
	fmt.Println("请输入命令：")
	var commend, setname string
	fmt.Scanln(&commend, &setname)
	if commend == "SMEMBER" || commend == "smember" || commend == "SMember" {
		open, err := os.OpenFile("server.json", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("open error", err)
			return
		}
		content, err := os.ReadFile("server.json")
		if err != nil {
			fmt.Println("read error", err)
			return
		}

		json.Unmarshal(content, &m)
		if _, ok := m[setname]; ok {
			fmt.Println("集合中的所有元素为：", m[setname])
		} else {
			fmt.Println("集合不存在")
		}
		defer open.Close()
	} else {
		fmt.Println("命令格式不合理")
	}
}

//向一个列表的左侧添加元素value

func lpush() {
	var m map[string]interface{} = make(map[string]interface{})
	fmt.Println("请输入命令：")
	var commend, listname, value string
	fmt.Scanln(&commend, &listname, &value)
	if commend == "LPUSH" || commend == "lpush" || commend == "Lpush" {
		open, err := os.OpenFile("server.json", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("open error", err)
			return
		}
		content, err := os.ReadFile("server.json")
		if err != nil {
			fmt.Println("read error", err)
			return
		}

		json.Unmarshal(content, &m)
		if _, ok := m[listname]; !ok {
			m[listname] = []interface{}{value}
			marshal, err := json.Marshal(m)
			if err != nil {
				fmt.Println("json errors", err)
				return
			}
			err = os.WriteFile("server.json", marshal, 0666)
			if err != nil {
				fmt.Println("write error", err)
				return
			} else {
				fmt.Println("向列表的左侧成功添加元素")
			}
		} else {
			m[listname] = append([]interface{}{value}, m[listname].([]interface{})...)
			marshal, err := json.Marshal(m)
			if err != nil {
				fmt.Println("json errors", err)
				return
			}
			err = os.WriteFile("server.json", marshal, 0666)
			if err != nil {
				fmt.Println("write error", err)
				return
			} else {
				fmt.Println("向列表的左侧成功添加元素")
			}
		}
		defer open.Close()
	} else {
		fmt.Println("命令格式不合理")
	}
}

//获取列表中指定范围的元素

func lrange() {
	var m map[string]interface{} = make(map[string]interface{})
	fmt.Println("请输入命令：")
	var commend, listname string
	var start, end int
	fmt.Scanln(&commend, &listname, &start, &end)
	if commend == "LRANGE" || commend == "lrange" || commend == "LRange" {
		open, err := os.OpenFile("server.json", os.O_RDWR|os.O_CREATE, 0666)
		if err != nil {
			fmt.Println("open error", err)
			return
		}
		content, err := os.ReadFile("server.json")
		if err != nil {
			fmt.Println("read error", err)
			return
		}

		json.Unmarshal(content, &m)
		if _, ok := m[listname]; ok {
			fmt.Println("列表中指定范围的元素为：", m[listname].([]interface{})[start:end+1]) //因为map中存储的元素为interface{}类型，所以需要转换为[]interface{}类型
		} else {
			fmt.Println("列表不存在")
		}
		defer open.Close()
	} else {
		fmt.Println("命令格式不合理")
	}
}

func main() {

	var choice int
	for {
		showmenu()
		fmt.Scanln(&choice)
		switch choice {
		case 1:
			set()
		case 2:
			setnx()
		case 3:
			get()
		case 4:
			del()
		case 5:
			sadd()
		case 6:
			smember()
		case 7:
			lpush()
		case 8:
			lrange()
		case 0:
			fmt.Println("欢迎下次使用")
			return
		default:
			fmt.Println("输入错误，请重新输入")

		}
	}
}
