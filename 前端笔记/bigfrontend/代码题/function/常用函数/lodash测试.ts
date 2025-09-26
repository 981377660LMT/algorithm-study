import _ from 'lodash'

const call = (v?: any) => {
  console.log('call', v)
}

function testDebounce() {
  const debouncedCall = _.debounce(call, 1000)

  debouncedCall()
  debouncedCall()
}

function testThrottle() {
  const throttledCall = _.throttle(call, 1000)
  throttledCall(1)
  throttledCall(2)
  throttledCall(3)
}

testThrottle()
