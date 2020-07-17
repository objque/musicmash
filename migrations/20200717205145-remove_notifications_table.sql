-- +migrate Up
ALTER TABLE "notifications" RENAME TO "_notifications";
ALTER TABLE "notification_settings" RENAME TO "_notification_settings";
ALTER TABLE "notification_services" RENAME TO "_notification_services";

-- +migrate Down
ALTER TABLE "_notifications" RENAME TO "notifications";
ALTER TABLE "_notification_settings" RENAME TO "notification_settings";
ALTER TABLE "_notification_services" RENAME TO "notification_services";
