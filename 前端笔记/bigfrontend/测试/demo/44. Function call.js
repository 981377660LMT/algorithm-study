function a() {
  console.log(1)
  return {
    a: function () {
      console.log(2)
      return a()
    },
  }
}

a().a()
// 1
// 2
// 1
