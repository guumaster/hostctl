package host

// Content contains complete data of all profiles
type Content struct {
	DefaultProfile DefaultProfile
	ProfileNames   []string
	Profiles       map[string]*Profile
}
