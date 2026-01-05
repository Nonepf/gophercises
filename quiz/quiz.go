package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"
)

const path = "problems.csv"
const gameTime = 10

type problem struct {
	q string
	a string
}

func parseLines(lines [][]string) []problem {
	ret := make([]problem, len(lines))
	for i, line := range lines {
		ret[i] = problem{
			q: line[0],
			a: line[1],
		}
	}
	return ret
}

func main() {
	// 1. 打开文件
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("无法打开文件：%s\n", err)
		return
	}
	defer file.Close()

	// 2. 将文件读入内存
	r := csv.NewReader(file)

	lines, err := r.ReadAll()
	if err != nil {
		fmt.Println("读取CSV出错")
		return
	}

	// 3. 解析内存数据
	problems := parseLines(lines)

	// 4. 准备进入
	timer := time.NewTimer(gameTime * time.Second)
	fmt.Println("输入回车，开始游戏！")
	fmt.Scanln()

	// 5. 开始Q-A
	score := 0

	defer func() {
		fmt.Printf("\n最终得分：%d/%d\n", score, len(problems))
	}()

	for i, problem := range problems {
		fmt.Printf("问题%v  %v:", i, problem.q)

		answerCh := make(chan string)

		go func() {
			var answer string
			fmt.Scanf("%s\n", &answer)
			answerCh <- answer
		}()

		select {
		case <-timer.C:
			fmt.Println("\n时间到！")
			return
		case answer := <-answerCh:
			if answer == problem.a {
				score++
			}
		}
	}
}
