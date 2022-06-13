package dao

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

//map作为内存索引
var (
	topicIndexMap map[int64]*Topic
	postIndexMap  map[int64][]*Post
)

func Init(filePath string) error {
	if err := initTopicIndexMap(filePath); err != nil {
		return err
	}
	if err := initPostIndexMap(filePath); err != nil {
		return err
	}
	return nil
}

func initTopicIndexMap(filePath string) error {
	//Open 只能打开已经存在的文件 也就是只有R操作
	open, err := os.Open(filePath + "topic")
	if err != nil {
		fmt.Println("os Open failed", err)
		return err
	}
	//scanner 读取reader
	scanner := bufio.NewScanner(open)
	topicTmpMap := make(map[int64]*Topic)
	for scanner.Scan() {
		//读出每一行的数据 反序列化到结构体中
		text := scanner.Text()
		var topic Topic
		if err := json.Unmarshal([]byte(text), &topic); err != nil {
			fmt.Println("json unmarshal failed", err)
			return err
		}
		//添加到索引map里
		topicTmpMap[topic.Id] = &topic
	}
	topicIndexMap = topicTmpMap
	return nil
}

func initPostIndexMap(filepath string) error {
	open, err := os.Open(filepath + "post")
	if err != nil {
		fmt.Println("os Open failed", err)
		return err
	}
	scanner := bufio.NewScanner(open)
	postTmpMap := make(map[int64][]*Post)
	for scanner.Scan() {
		text := scanner.Text()
		var post Post
		if err := json.Unmarshal([]byte(text), &post); err != nil {
			fmt.Println("json unmarshal failed", err)
			return err
		}
		//当没有这个Id
		posts, ok := postTmpMap[post.ParentId]
		if !ok {
			postTmpMap[post.ParentId] = []*Post{&post}
			continue
		}
		posts = append(posts, &post)
		postTmpMap[post.ParentId] = posts
	}

	postIndexMap = postTmpMap
	return nil
}
