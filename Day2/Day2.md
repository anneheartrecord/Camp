多线程编程go的变量是线程安全的 
go的slice和map是非线程安全的 因为底层都有指针
go的channel是线程安全的（加锁了）

Go的依赖管理主要经历了三个阶段 分别是GOPATH Go Vendor Go Module

GOPATH时代GOPATH是Go的一个环境变量目录有以下结构 src 存放Go的源码pkg 存放编译的中间产物 加快编译速度bin 存放项目编译生成的二进制文件
缺点：如果一个pkg有两个版本 都存放在src目录下 而src只能有一个版本存在 这个时候A B项目依赖于不同版本的这个pkg 就会导致A B不能全部编译成功 也就是说GOPATH无法实现pkg的多版本控制
Go Vendor:解决了GOPATH的问题GO Vendor的处理方法：在每个项目目录下都有一个Vendor文件 存放当前项目依赖的副本 在vendor机制下 会优先使用该目录下的依赖 如果依赖不存在 才会去GOPATH下寻找 
缺点：Vendor 无法很好解决依赖包的版本变动问题
为了解决这个问题 mod诞生了
依赖关系A->B->C 我们称A对B是直接依赖 A对C是间接依赖
在mod 中打上//indirect标签的就表明间接依赖在mod 中打上//incompatible标签表明合格仓库已经打上了2或者更高版本的tag 为了兼容这部分仓库 对于没有go mod文件 并且主版本在2以上的依赖 会打上incompatible标签
依赖分发我们每次使用github上的第三方库都不是直接饮用的第三方 原因1.多次拉取会对github造成很大压力2.如果作者改了改代码 可能你的项目直接就动不了了
Go的Proxy就是解决这些问题的答案他是一个服务站点 会缓存站中的软件内容 缓存版本不变 并且在源站删除之后依旧可用 从而实现了依赖分发
GOPROXT=“https://proxy1.cn,direct” //direct就表示源站 这句命令就是说 在找得到的情况下 在proxy1找 找不到 就去源站找

go mod init 初始化mod文件go mod download 下载模块到本地缓存go mod tidy 下载需要的依赖 删去不需要的依赖

Go的测试三种级别的测试从上到下分别是 越下层覆盖率越大 成本越低回归测试:QA手动测试 某些场景是否会有问题集成测试：对系统功能维度做测试验证  单元测试 ：开发者对单独的函数 模块做功能验证

demo 用Go内置的map和文件系统实现了一个查询话题的接口