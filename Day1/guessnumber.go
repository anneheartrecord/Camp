package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	maxNumber := 100
	//种下时间种子 时间戳
	rand.Seed(time.Now().UnixNano())
	secretNumber := rand.Intn(maxNumber)
	for {
		//从标准输入读入
		reader := bufio.NewReader(os.Stdin)
		//读到\n结束 但是会读进\n字符
		read, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("read failed", err)
			break
		}
		//删掉 \n字符
		data := strings.TrimSuffix(read, "\n")
		//将字符串转成 int 类型
		yourGuess, err := strconv.Atoi(data)
		if err != nil {
			fmt.Println("atoi failed", err)
			break
		}
		if yourGuess > secretNumber {
			fmt.Println("your guess is bigger than secretNumber")
		} else if yourGuess < secretNumber {
			fmt.Println("your guess is smaller than secretNumber")
		} else {
			fmt.Println("you win,the secretNumber is", yourGuess)
			break
		}
	}
}
