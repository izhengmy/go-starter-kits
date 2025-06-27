package auth

type options struct {
	ctxKey              string
	retrieveTokenFunc   RetrieveTokenFunc
	unauthenticatedFunc UnauthenticatedFunc
}

type Option interface {
	apply(opts *options)
}

type optionFunc func(opts *options)

func (f optionFunc) apply(opts *options) {
	f(opts)
}

func WithContextKey(key string) Option {
	return optionFunc(func(opts *options) {
		opts.ctxKey = key
	})
}

func WithRetrieveTokenFunc(retrieveTokenFunc RetrieveTokenFunc) Option {
	return optionFunc(func(opts *options) {
		opts.retrieveTokenFunc = retrieveTokenFunc
	})
}

func WithUnauthenticatedFunc(unauthenticatedFunc UnauthenticatedFunc) Option {
	return optionFunc(func(opts *options) {
		opts.unauthenticatedFunc = unauthenticatedFunc
	})
}
