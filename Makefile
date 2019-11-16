DST_BIN_DIR=${HOME}/bin/workflows
.PHONY: compile
compile:
	go build

.PHONY: install
install:
	go build -o $(DST_BIN_DIR)