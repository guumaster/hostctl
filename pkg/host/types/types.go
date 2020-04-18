package types

const Banner = `
##################################################################
# Content under this line is handled by hostctl. DO NOT EDIT.
##################################################################`

// Content contains complete data of all profiles
type Content struct {
	DefaultProfile DefaultProfile
	ProfileNames   []string
	Profiles       map[string]*Profile
}
