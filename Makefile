.PHONY: all clean clean-all default docker-build local restart start stop

export CGO_ENABLED=0

IMAGE_NAME=gotiny
CONTAINER_NAME=gotiny
API_SPEC_FILE=$(WS_ROOT)/api/gotiny-api.yaml

WS_ROOT:=$(CURDIR)
# Under cygwin, convert path to Windows fortmat
ifneq (,$(findstring /cygdrive,$(WS_ROOT)))
	WS_ROOT:=$(shell cygpath -m "${WS_ROOT}")
endif

default: docker-build
all: docker-build

docker-build:
	docker build -t $(IMAGE_NAME) .

local:
	echo "Building API server"
	cd $(WS_ROOT)/src/cmd && CGO_ENABLED=0 go build -o $(WS_ROOT)/app/apisrv/apisrv$(EXE_EXT) .
	$(COPYCMD) $(WS_ROOT)/conf/config.yaml $(WS_ROOT)/app/apisrv

	echo "Building Web server"
	cd $(WS_ROOT)/tools/tinyweb && CGO_ENABLED=0 go build -o $(WS_ROOT)/app/websrv/tinyweb$(EXE_EXT) .

gen:
	# Build generator
	cd $(WS_ROOT)/tools/rest-gen && go build -o $(WS_ROOT)/app/tools/rest-gen$(EXE_EXT) main.go

	# Generate API models
	$(MKDIRCMD) $(WS_ROOT)/src/apisrv/models
	cd $(WS_ROOT)/src/apisrv/models && $(WS_ROOT)/app/tools/rest-gen model $(API_SPEC_FILE)

	# Generate API markdown
	cd $(WS_ROOT) && $(WS_ROOT)/app/tools/rest-gen markdown $(API_SPEC_FILE) $(WS_ROOT)/api/gotiny-api-spec.md
	$(MKDIRCMD) $(WS_ROOT)/app
	$(COPYCMD) $(WS_ROOT)/api/gotiny_api_spec.md $(WS_ROOT)/app

start:
	-docker stop $(CONTAINER_NAME)
	-docker rm $(IMAGE_NAME)
	docker-compose up -d -f deploy/docker-compose.yaml
	docker run -d -p 80:80 -p 81:81 --name $(CONTAINER_NAME) $(IMAGE_NAME)
	docker ps --latest

stop:
	-docker stop $(CONTAINER_NAME)
	-docker-compose down -f deploy/docker-compose.yaml

clean:
	-rm -rf $(WS_ROOT)/app/*
	-docker rm $(CONTAINER_NAME)
	-docker-compose rm -f deploy/docker-compose.yaml

restart: stop start
