
# TODO: copy/md for Windows
COPYCMD=cp
MKDIRCMD=mkdir --parent
ifeq ($(OS),Windows_NT)
	EXE_EXT := .exe
endif

API_SPEC_FILE=$(WS_ROOT)/api/gotiny-api.yaml

WS_ROOT:=$(CURDIR)
# Under cygwin, convert path to Windows fortmat
ifneq (,$(findstring /cygdrive,$(WS_ROOT)))
	WS_ROOT:=$(shell cygpath -m "${WS_ROOT}")
endif

default: gen

gen:
	# Generate API models
	$(MKDIRCMD) $(WS_ROOT)/src/apisrv/models
	cd $(WS_ROOT)/src/apisrv/models && $(WS_ROOT)/app/tools/rest-gen model $(API_SPEC_FILE)

	# Generate API markdown
	cd $(WS_ROOT) && $(WS_ROOT)/app/tools/rest-gen markdown $(API_SPEC_FILE) $(WS_ROOT)/api/gotiny-api-spec.md
