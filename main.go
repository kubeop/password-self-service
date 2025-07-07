/*
Copyright © 2023 Sonic Ma <EMAIL sonic.ma@outlook.com>
*/
package main

import (
	"password-self-service/internal/cmd"
)

// @title						Swagger API
// @version					1.0
// @description				密码自助平台
// @securityDefinitions.apikey	BearerToken
// @in							header
// @name						Authorization
func main() {
	cmd.Execute()
}
