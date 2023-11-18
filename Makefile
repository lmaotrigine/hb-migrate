# I hate these

.PHONY: build clean run

build:
	go build ./cmd/migrate_stats.go

clean:
	rm -f migrate_stats

run: build
	./migrate_stats --help
