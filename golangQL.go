package golangQL


func Filter(v interface{}, query string) (interface{}, error) {
	if len(query) == 0 {
		return v, nil
	}

	tree, err := parse(query)
	if err != nil {
		return nil, err
	}

	return filter(v, tree)
}
