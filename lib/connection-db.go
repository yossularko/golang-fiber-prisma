package lib

import (
	"fmt"
	"golang-fiber-prisma/db"
	"log"
)

func ConnectToDatabase(prisma *db.PrismaClient) error {
	err := prisma.Connect()
	if err != nil {
		log.Fatal(err)
	}

	return err
}

func DisconnectFromDatabase(prisma *db.PrismaClient) error {
	fmt.Println("close db")
	err := prisma.Disconnect()
	if err != nil {
		log.Fatal(err)
	}

	return err
}
