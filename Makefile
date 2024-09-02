migrateup:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable" -verbose down

createnewmigration/%:# you can run createnewmigration/{new_name_schema}
	migrate create -ext sql -dir db/migrations -seq $(shell echo $@ | cut -d '/' -f2-)

.PHONY: migrateup migratedown createnewmigration/%