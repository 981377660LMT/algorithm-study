if (require.main === module) {
  const getBig = (len = 1e5): number[] => {
    const big = Array(len).fill(0)
    for (let i = 0; i < len; i++) {
      big[i] = ~~(Math.random() * 2e9) - 1e9
    }
    return big
  }
}

export {}
