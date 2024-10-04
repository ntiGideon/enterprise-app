package repository

import (
	"Enterprise/prisma/db"
	"golang.org/x/net/context"
)

func ExistingCategoryByName(ctx context.Context, dbClient *db.PrismaClient, name string) (bool, error) {
	exitingCategory, _ := dbClient.Category.FindFirst(db.Category.Name.Equals(name)).Exec(ctx)
	if exitingCategory != nil {
		return true, nil
	}
	return false, nil
}
