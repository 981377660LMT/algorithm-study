var a = 10
;(function () {
  console.log(a)
  a = 5
  console.log(window.a)
  var a = 20
  console.log(a)
})()
// 分别为undefined　10　20，
// 原因是作用域问题，在内部声名var a = 20;
// 相当于先声明var a;然后再执行赋值操作，这是在ＩＩＦＥ内形成的独立作用域，如果把var a=20注释掉，
// 那么a只有在外部有声明，显示的就是外部的Ａ变量的值了。结果Ａ会是10　5　5
