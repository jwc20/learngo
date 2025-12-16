package poker

//func NewInMemoryPlayerStore() *InMemoryPlayerStore {
//	return &InMemoryPlayerStore{map[string]int{}}
//}
//
//type InMemoryPlayerStore struct {
//	store map[string]int
//}
//
//func (i *InMemoryPlayerStore) GetPlayerScore(name string) int {
//	return i.store[name]
//}
//
//func (i *InMemoryPlayerStore) RecordWin(name string) {
//	i.store[name]++
//}
//
//func (i *InMemoryPlayerStore) GetLeague() League {
//	var League League
//	for name, wins := range i.store {
//		League = append(League, Player{name, wins})
//	}
//	return League
//}
