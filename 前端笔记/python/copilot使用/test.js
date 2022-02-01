// fetch an image
// http://www.pythonchallenge.com/pc/def/oxygen.html
// http://www.pythonchallenge.com/pc/def/integrity.html
// http://www.pythonchallenge.com/pc/def/banner.p

// randint function
function randint(a, b) {
  return Math.floor(Math.random() * (b - a + 1) + a)
}

// koa server
const Koa = require('koa')
const app = new Koa()
