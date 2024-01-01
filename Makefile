
ALL: abs-importer

.PHONY: abs-importer


abs-importer:
	go mod tidy
	go build


