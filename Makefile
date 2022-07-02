.PHONY: all clean clean-all default docker-build local restart start stop

export CGO_ENABLED=0

IMAGE_NAME:=gotiny
CONTAINER_NAME:=gotiny

WS_ROOT:=$(CURDIR)
# Under cygwin, convert path to Windows fortmat
ifneq (,$(findstring /cygdrive,$(WS_ROOT)))
	WS_ROOT:=$(shell cygpath -m "${WS_ROOT}")
endif

API_SPEC_FILE:=$(WS_ROOT)/api/gotiny-api.yaml
DOCKER_COMPOSE:=$(WS_ROOT)/deploy/docker-compose.yaml
APP_DIR:=$(WS_ROOT)/app

# TODO: copy/md for Windows
COPYCMD=cp
MKDIRCMD=mkdir --parent
ifeq ($(OS),Windows_NT)
	EXE_EXT := .exe
endif


default: local
all: gen local docker-build

docker-build:
	docker build -t $(IMAGE_NAME) .

local:
	echo "Building API server"
	cd $(WS_ROOT)/src/cmd && go build -o $(APP_DIR)/apisrv/apisrv$(EXE_EXT) .
	$(COPYCMD) $(WS_ROOT)/conf/config.yaml $(APP_DIR)/apisrv

	echo "Building Web server"
	cd $(WS_ROOT)/tools/tinyweb && go build -o $(APP_DIR)/websrv/tinyweb$(EXE_EXT) .

gen:
	# Build generator
	cd $(WS_ROOT)/tools/rest-gen && go build -o $(APP_DIR)/tools/rest-gen$(EXE_EXT) main.go

	# Generate API models
	$(MKDIRCMD) $(WS_ROOT)/src/apisrv/models
	cd $(WS_ROOT)/src/apisrv/models && $(APP_DIR)/tools/rest-gen model $(API_SPEC_FILE)

	# Generate API markdown
	cd $(WS_ROOT) && $(APP_DIR)/tools/rest-gen markdown $(API_SPEC_FILE) $(WS_ROOT)/api/gotiny-api-spec.md
	$(MKDIRCMD) $(APP_DIR)
	$(COPYCMD) $(WS_ROOT)/api/gotiny_api_spec.md $(APP_DIR)

start: local
	-docker rm $(IMAGE_NAME)  2>&1 > /dev/null
	docker-compose -f $(DOCKER_COMPOSE) up -d
	$(APP_DIR)/apisrv/apisrv

stop:
	-docker stop $(CONTAINER_NAME)
	docker-compose -f $(DOCKER_COMPOSE) down

clean:
	-rm -rf $(APP_DIR)/*
	-docker rm $(CONTAINER_NAME) 2>&1 > /dev/null
	docker-compose -f $(DOCKER_COMPOSE) rm

