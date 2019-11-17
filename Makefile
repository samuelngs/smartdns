NAME = smartdns

APP_SRC_DIR = .
APP_CMD_DIR = $(APP_SRC_DIR)/cmd
APP_BUILD_DIR = $(APP_SRC_DIR)/build

PROGRAMS = $(shell ls $(APP_CMD_DIR) | xargs -I* echo build-*)

$(NAME): install-deps test $(PROGRAMS)

install-deps:
	go mod download

build-%:
	go build -o $(APP_BUILD_DIR)/$* -v ./cmd/$*

test:
	go test -v ./...

clean:
	go clean
	rm -rf $(APP_BUILD_DIR)
