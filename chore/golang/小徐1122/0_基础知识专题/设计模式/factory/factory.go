// https://mp.weixin.qq.com/s?__biz=MzkxMjQzMjA0OQ==&mid=2247484392&idx=1&sn=9471c3c01742a3776b6ea9f285dea759
// • 工厂模式的使用背景
// • 简单工厂模式
// • 工厂方法模式
// • 抽象工厂模式
// • 容器工厂模式

package main

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"
)

func main() {
	testSimpleFactory()
	testFactoryMethod()
}

func testSimpleFactory() {
	factory := NewFruitFactory()
	fruit, err := factory.CreateFruit("orange")
	if err != nil {
		fmt.Println(err)
		return
	}

	fruit.Eat()
}

func testFactoryMethod() {
	orangeFactory := NewOrangeFactoryV2()
	orange := orangeFactory.CreateFruit()
	orange.Eat()

	strawberryFactory := NewStrawberryFactoryV2()
	strawberry := strawberryFactory.CreateFruit()
	strawberry.Eat()
}

// #region 简单工厂
type fruitCreator func(name string) Fruit

type FruitFactory struct {
	creators map[string]fruitCreator
}

func NewFruitFactory() *FruitFactory {
	return &FruitFactory{
		creators: map[string]fruitCreator{
			"orange":     NewOrange,
			"strawberry": NewStrawberry,
		},
	}
}

func (f *FruitFactory) CreateFruit(typ string) (Fruit, error) {
	fruitCreator, ok := f.creators[typ]
	if !ok {
		return nil, fmt.Errorf("fruit typ: %s is not supported yet", typ)
	}

	src := rand.NewSource(time.Now().UnixNano())
	rander := rand.New(src)
	name := strconv.Itoa(rander.Int())
	return fruitCreator(name), nil
}

type Fruit interface {
	Eat()
}

type Orange struct {
	name string
}

func NewOrange(name string) Fruit {
	return &Orange{name: name}
}

func (o *Orange) Eat() {
	fmt.Printf("i am orange: %s, i am about to be eaten...", o.name)
}

type Strawberry struct {
	name string
}

func NewStrawberry(name string) Fruit {
	return &Strawberry{name: name}
}

func (s *Strawberry) Eat() {
	fmt.Printf("i am strawberry: %s, i am about to be eaten...", s.name)
}

// #endregion

// #region 工厂方法
type FruitFactoryV2 interface {
	CreateFruit() Fruit
}

type OrangeFactoryV2 struct {
}

func NewOrangeFactoryV2() FruitFactoryV2 {
	return &OrangeFactoryV2{}
}

func (o *OrangeFactoryV2) CreateFruit() Fruit {
	return NewOrange("")
}

type StrawberryFactoryV2 struct {
}

func NewStrawberryFactoryV2() FruitFactoryV2 {
	return &StrawberryFactoryV2{}
}

func (s *StrawberryFactoryV2) CreateFruit() Fruit {
	return NewStrawberry("")
}

// #endregion

// #region 容器工厂

// #endregion
