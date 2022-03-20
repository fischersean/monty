package main

import (
	"os"
	"strconv"
)

type dbConnInfo struct {
	Password   string `json:"password"`
	Port       int    `json:"port"`
	Host       string `json:"Host"`
	Username   string `json:"username"`
	Identifier string `json:"dbInstanceIdentifier"`
}

func getDbConn() (info dbConnInfo, err error) {

	port, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		return info, err
	}
	info.Password = os.Getenv("DB_PASSWORD")
	info.Host = os.Getenv("DB_HOST")
	info.Port = port
	info.Username = os.Getenv("DB_USER")

	return info, err
}
