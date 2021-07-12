package main

import (
	pgs "github.com/lyft/protoc-gen-star"
	// pgsgo "github.com/lyft/protoc-gen-star/lang/go"
)

func main() {
	g := pgs.Init()
	g.RegisterModule(NewRpc())
	g.Render()
}
