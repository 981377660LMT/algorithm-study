# java的坏味道

https://mp.weixin.qq.com/s?__biz=MzkxMjQzMjA0OQ==&mid=2247484455&idx=1&sn=b36839bf02a395b51355f598f0a25742

1. 问题分析

```go
type CourseService interface {
    // 一系列编程课程
    LearnGolang()
    LearnJAVA()
    LearnC()
    // ...


    // 一系列体育课程
    LearnBasketball()
    LearnFootball()
    LearnSki()
    // ...


    // 一系列音乐课程
    LearnPiano()
    LearnHarmonica()
    LearnGuita()
    // ...
}


type courseServiceImpl struct {
}


func NewCourseService() CourseService {
    return &courseServiceImpl{}
}


func (c *courseServiceImpl) LearnGolang() {
    fmt.Println("learn go...")
}


func (c *courseServiceImpl) LearnJAVA() {
    fmt.Println("learn java...")
}


func (c *courseServiceImpl) LearnC() {
    fmt.Println("learn c...")
}


func (c *courseServiceImpl) LearnBasketball() {
    fmt.Println("learn basketball...")
}


func (c *courseServiceImpl) LearnFootball() {
    fmt.Println("learn football...")
}


func (c *courseServiceImpl) LearnSki() {
    fmt.Println("learn ski...")
}


func (c *courseServiceImpl) LearnPiano() {
    fmt.Println("learn piano...")
}


func (c *courseServiceImpl) LearnHarmonica() {
    fmt.Println("learn harmonica...")
}


func (c *courseServiceImpl) LearnGuita() {
    fmt.Println("learn guita...")
}
```

上面这种实现方式中，interface 和实现 class 是紧密绑定的，遵循还是类似于 JAVA 中显式实现的代码风格，然而这种风格在 golang 中是并不值得推崇的，根本原因就在于 Golang 隐式实现与 JAVA 显式实现的差异.
这种实现方式的缺陷在于：

- 对于模块的实现方来说，每次新增或者修改方法，**需要同步变更两处内容**，包括 interface 和具体的 serviceImpl，增加了变更成本

- 对于模块的使用方来说，由于 interface 已经由实现方定好了，这种抽象程度对于使用方来说可能并不合适，导致可能有**使用方不关心的方法也被 interface 暴露出来. 这样一方面会对使用方造成使用上的困扰，另一方面也会增加使用方在对 interface 添加更多实现类时的实现成本**，因为实现时会涉及到对使用方所不关心的一部分抽象方法的实现

`产生这个问题的根本原因在于，构造 interface 的工作不应该由模块的实现方来做。`
定义 interface 本质上是一个对类型的边界和职责进行抽象的过程，作为实现方的角色，它永远无法做到未卜先知地站在未来使用方的视角，来帮助使用方做出”如何使用这个模块“的定义和决策.
使用方有可能只希望通过 CourseService 进行 Golang 课程的学习，那么此时编程课程中的 JAVA 课程和 C 课程对于使用方来说也属于是无须关心的一部分信息.

2. 解决方案

针对于 CourseService 的场景问题，我们进行如下改造：

- 实现方不再进行 interface 的声明，而是`将 CourseService 改为一个具体的实现类型`，将定义 interface 的职责转交给 CourseService 的使用方

- 使用方在使用 CourseService 时，根据使用到的 CourseService 的功能范围，对其身份进行抽象和定义，比如在使用 CourseService 中编程课程有关的方法时，使用方可以定义出一个 CSCourseProxy 的 interface，然后在 interface 中定义好有关编程课程的几个方法，其他无关的方法不再声明，起到屏蔽无关职责的效果

```go
// 实现方
type CourseService struct {
}


func NewCourseService() *CourseService {
    return &CourseService{}
}


func (c *CourseService) LearnGolang() {
    fmt.Println("learn go...")
}


func (c *CourseService) LearnJAVA() {
    fmt.Println("learn java...")
}


func (c *CourseService) LearnC() {
    fmt.Println("learn c...")
}


func (c *CourseService) LearnBasketball() {
    fmt.Println("learn basketball...")
}


func (c *CourseService) LearnFootball() {
    fmt.Println("learn football...")
}


func (c *CourseService) LearnSki() {
    fmt.Println("learn ski...")
}


func (c *CourseService) LearnPiano() {
    fmt.Println("learn piano...")
}


func (c *CourseService) LearnHarmonica() {
    fmt.Println("learn harmonica...")
}


func (c *CourseService) LearnGuita() {
    fmt.Println("learn guita...")
}

// 使用方
type CSCourseProxy interface {
    LearnGolang()
    LearnJAVA()
    LearnC()
}

type PECourseProxy interface {
    LearnBasketball()
    LearnFootball()
    LearnSki()
}

type MusicCourseProxy interface {
    LearnPiano()
    LearnHarmonica()
    LearnSki()
}

// 使用方如果只需要使用到 CourseService 中有关于 Golang 课程的资源，则可以进一步细化。测试的时候也容易mock。
type GolangCourseProxy interface{
    LearnGolang()
}
```

3. 为什么在实际编程中会使用java的这种方式，而不是golang的这种方式？
