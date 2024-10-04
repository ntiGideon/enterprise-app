package repository

import (
	"Enterprise/prisma/db"
	"fmt"
	"golang.org/x/net/context"
)

func AuditLogs(ctx context.Context, dbClient *db.PrismaClient, auditId int, action string, details string) error {
	admin, err := dbClient.User.FindUnique(db.User.ID.Equals(auditId)).Exec(ctx)
	lastName, _ := admin.LastName()
	detailMessage := fmt.Sprintf("%v conducted by ==> %v %v", details, admin.FirstName, lastName)
	_, err = dbClient.AuditLog.CreateOne(
		db.AuditLog.User.Link(db.User.ID.Equals(auditId)),
		db.AuditLog.Action.Set(action),
		db.AuditLog.Details.Set(detailMessage),
	).Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}
