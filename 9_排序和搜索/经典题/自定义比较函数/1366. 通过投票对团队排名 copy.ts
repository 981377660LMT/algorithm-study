// # 每个投票者都需要按从高到低的顺序对参与排名的所有团队进行排位。
// # 参赛团队的排名次序依照其所获「排位第一」的票的多少决定。如果存在多个团队并列的情况，将继续考虑其「排位第二」的票的数量。以此类推，直到不再存在并列的情况。
// # 如果在考虑完所有投票情况后仍然出现并列现象，则根据团队字母的字母顺序进行排名。

function rankTeams(votes: string[]): string {
  const level = votes[0].length
  const levelCounter = Array.from<unknown, Map<string, number>>({ length: level }, () => new Map())
  const team = new Set<string>()
  for (const vote of votes) {
    for (let i = 0; i < level; i++) {
      const teamName = vote.charAt(i)
      levelCounter[i].set(teamName, (levelCounter[i].get(teamName) || 0) + 1)
      team.add(teamName)
    }
  }

  return [...team].sort((t1, t2) => compare(t1, t2)).join('')

  function compare(team1: string, team2: string) {
    for (let i = 0; i < level; i++) {
      const count1 = levelCounter[i].get(team1) || 0
      const count2 = levelCounter[i].get(team2) || 0
      if (count1 !== count2) return count2 - count1
    }

    return team1.localeCompare(team2)
  }
}

console.log(rankTeams(['ABC', 'ACB', 'ABC', 'ACB', 'ACB']))
