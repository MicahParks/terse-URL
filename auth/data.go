package auth

// ShortenedData is a mapping of users to Authorization data.
type ShortenedData map[string]Authorization

// UserData is a mapping of shortened URLs to Authorization data.
type UserData map[string]Authorization

// TODO Remove all these methods?

// Delete deletes the Authorization data for the given shortened URL.
func (s ShortenedData) Delete(shortened string) {
	del(s, shortened)
}

// Get retrieves the Authorization data for the given URL.
func (s ShortenedData) Get(shortened string) (a Authorization, ok bool) {
	return get(s, shortened)
}

// Set sets the shortened URL's value to the given Authorization data.
func (s ShortenedData) Set(shortened string, authorization Authorization) {
	set(s, shortened, authorization)
}

// Delete deletes the Authorization data for the given shortened URL.
func (u UserData) Delete(shortened string) {
	del(u, shortened)
}

// Get retrieves the Authorization data for the given URL.
func (u UserData) Get(shortened string) {
	get(u, shortened)
}

// Set sets the shortened URL's value to the given Authorization data.
func (u UserData) Set(shortened string, a Authorization) {
	set(u, shortened, a)
}

func del(m map[string]Authorization, shortened string) {
	delete(m, shortened)
}

func get(m map[string]Authorization, shortened string) (a Authorization, ok bool) {
	a, ok = m[shortened]
	return a, ok
}

func set(m map[string]Authorization, shortened string, a Authorization) {
	m[shortened] = a
}
