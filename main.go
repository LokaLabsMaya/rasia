package main

import "github.com/lokalabsmaya/rasia/cmd"

// @title rasia Service API
// @description This is a documentation for rasia Service RESTful APIs. <br>

// @securityDefinitions.basic BasicAuth
// @in header
// @name Authorization

// @securityDefinitions.apikey BearerToken
// @in header
// @name Authorization

func main() {

	cmd.Run()

}