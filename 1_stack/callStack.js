const func1 = () => {
  func2()
}
const func2 = () => {
  func3()
}
const func3 = () => {}

func1()
