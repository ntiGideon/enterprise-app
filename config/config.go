package config

import (
	"Enterprise/prisma/db"
	"github.com/rs/zerolog/log"
)

func ConnectDB() (*db.PrismaClient, error) {
	client := db.NewClient()
	if err := client.Connect(); err != nil {
		return nil, err
	}
	log.Info().Msg("Connected to Prisma Database")
	return client, nil
}
