module comparison

go 1.24.1

require (
	github.com/go-chi/chi/v5 v5.2.2
	github.com/pal-paul/go-libraries v0.0.0
	github.com/valyala/fasthttp v1.64.0
)

replace github.com/pal-paul/go-libraries => ../..

require (
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/klauspost/compress v1.18.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
)
