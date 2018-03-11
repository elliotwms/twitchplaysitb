package main

import "fmt"

type VoteResult struct {
	Command *Command
	Votes   int
	Users   []string
}

type VoteResultList []VoteResult

func (v VoteResultList) Len() int {
	return len(v)
}

func (v VoteResultList) Less(i, j int) bool {
	return v[i].Votes < v[j].Votes
}

func (v VoteResultList) Swap(i, j int) {
	v[i], v[j] = v[j], v[i]
}

func tallyVotes(cq map[string]*Command) VoteResultList {
	fmt.Printf("Processing %d commands...\n", len(cq))

	results := make(map[string]VoteResult)

	for k, c := range cq {
		hash := c.GetHash()

		r := VoteResult{
			Command: c,
		}

		if _, found := results[hash]; found {
			// If we already have votes for this command then bring those results up
			r = results[hash]
		}

		r.Votes++
		r.Users = append(r.Users, k)

		results[hash] = r

		delete(cq, k)
	}

	// Group the results into a VoteResultList so they can be sorted

	resultList := make(VoteResultList, len(results))

	i := 0
	for _, r := range results {
		resultList[i] = r
		i++
	}

	return resultList
}

func getWinningVote(list VoteResultList) VoteResult {
	return list[0] //todo
}
