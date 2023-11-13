/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package main

import "password-self-service/internal/cmd"

// @title						Swagger API
// @version					1.0
// @description				AD密码自助平台
// @securityDefinitions.apikey	BearerToken
// @in							header
// @name						Authorization
func main() {
	cmd.Execute()
}
