var i = 0
function inc() {
  i++
  inc()
}

try {
  inc()
} catch (e) {
  // The StackOverflow sandbox adds one frame that is not being counted by this code
  // Incrementing once manually
  i++
  console.log('Maximum stack size is', i, 'in your current browser')
}
