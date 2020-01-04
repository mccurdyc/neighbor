package github

func searchRepository() {
	panic("not implemented")
}

func searchCode() {
	panic("not implemented")
}

func searchMeta(entity searchMethodEntity) {
	switch entity {
	case topic:
		searchTopic()
	case textMatch:
		searchTextMatch()
	case label:
		searchLabel()
	}
}

func searchTopic() {
	panic("not implemented")
}

func searchLabel() {
	panic("not implemented")
}

func searchTextMatch() {
	panic("not implemented")
}
