loop: for (let i = 0; i < 3; i++) {
  console.log(1)

  for (let j = 0; j < 3; j++) {
    console.log(2)
    continue loop
  }

  console.log(3)

  break loop
}

export {}
