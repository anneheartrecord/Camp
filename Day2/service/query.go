package service

import (
	"camp/Day2/dao"
	"errors"
	"fmt"
	"sync"
)

//一个页面的信息 由话题和帖子列表组成
type PageInfo struct {
	Topic    *dao.Topic
	PostList []*dao.Post
}

type QueryPageInfoFlow struct {
	topicId  int64
	pageInfo *PageInfo
	topic    *dao.Topic
	posts    []*dao.Post
}

func QueryPageInfo(topicId int64) (*PageInfo, error) {
	return NewQueryPageInfoFlow(topicId).Do()
}

func NewQueryPageInfoFlow(topId int64) *QueryPageInfoFlow {
	return &QueryPageInfoFlow{
		topicId: topId,
	}
}
func (f *QueryPageInfoFlow) Do() (*PageInfo, error) {
	if err := f.checkParam(); err != nil {
		return nil, err
	}
	if err := f.prepareInfo(); err != nil {
		return nil, err
	}
	if err := f.packPageInfo(); err != nil {
		return nil, err
	}

	return f.pageInfo, nil
}

//验证参数是否正确
func (f *QueryPageInfoFlow) checkParam() error {
	if f.topicId <= 0 {
		return errors.New("topic id must be larger than 0")
	}
	return nil
}

func (f *QueryPageInfoFlow) prepareInfo() error {
	//获取topic信息
	var wg sync.WaitGroup
	//开协程处理
	wg.Add(2)
	//查话题
	go func() {
		defer wg.Done()
		topic := dao.NewTopicDaoInstance().QueryTopicById(f.topicId)

		fmt.Println(topic)
		f.topic = topic
	}()
	//获取post列表
	go func() {
		defer wg.Done()
		posts := dao.NewPostDaoInstance().QueryPostsByParentId(f.topicId)

		fmt.Println(posts)
		f.posts = posts
	}()
	wg.Wait()
	return nil
}

func (f *QueryPageInfoFlow) packPageInfo() error {
	f.pageInfo = &PageInfo{
		Topic:    f.topic,
		PostList: f.posts,
	}
	return nil
}
