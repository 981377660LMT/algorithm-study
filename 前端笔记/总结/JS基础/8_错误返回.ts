function testFinally() {
  try {
    return '出去玩'
  } catch (error) {
    return '看电视'
  } finally {
    return '做作业'
  }
  // 检测到无法访问的代码。
  return '睡觉'
}
