-- CreateEnum
CREATE TYPE "StateEnum" AS ENUM ('FRESH', 'VERIFIED', 'DISABLED', 'DELETED');

-- AlterTable
ALTER TABLE "User" ADD COLUMN     "state" "StateEnum" NOT NULL DEFAULT 'FRESH';
