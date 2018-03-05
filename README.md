# friendly-guacamole :two_men_holding_hands: :avocado: :wine_glass:

This is a sample proxy server. A client sends an RPC message to the server. The server looks at some of the request metadata, and passes the request to an internal server that does the actual work. If the request takes too long, the request handler is killed and an error is returned.

* `go run server.go` starts a GRPC server.

* `go run client.go` sends a request to the server. When you run it without args, it will show you the available options.

Or just `./go -s 20`.
