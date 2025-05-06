package persistence

type StateStorage interface {
	LoadState() (*State, error)
	StoreState(State) error
}
