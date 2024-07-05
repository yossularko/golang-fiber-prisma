package inits

import (
	"golang-fiber-prisma/db"
	"golang-fiber-prisma/lib"
)

var Prisma *db.PrismaClient

func PrismaInit() {
	Prisma = db.NewClient()
	err := lib.ConnectToDatabase(Prisma)
	if err != nil {
		panic("failed to connect database")
	}
}
