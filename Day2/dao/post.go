package dao

import "sync"

//每个帖子都属于某个话题 所以应该有ParentId 字段
type Post struct {
	Id         int64  `json:"id"`
	ParentId   int64  `json:"parent_id"`
	Content    string `json:"content"`
	CreateTime int64  `json:"create_time"`
}
type PostDao struct {
}

//Once 只执行一次 就是我们平常所说的单例模式
var (
	postDao  *PostDao
	postOnce sync.Once
)

//实例化一个空结构体
func NewPostDaoInstance() *PostDao {
	postOnce.Do(
		func() {
			postDao = &PostDao{}
		})
	return postDao
}

func (*PostDao) QueryPostsByParentId(parentId int64) []*Post {
	return postIndexMap[parentId]
}
