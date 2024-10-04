package repository

import (
	"Enterprise/prisma/db"
	"golang.org/x/net/context"
)

func ExistingProductByName(ctx context.Context, dbClient *db.PrismaClient, name string) bool {
	existingProduct, _ := dbClient.Product.FindFirst(db.Product.Name.Equals(name)).Exec(ctx)
	return existingProduct != nil
}
