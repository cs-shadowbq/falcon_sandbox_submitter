#!make
.PHONY: reportVersion depend clean build all
.DEFAULT: all
.EXPORT_ALL_VARIABLES:

-include .env

cat := $(if $(filter $(OS),Windows_NT),type,cat)
PACKAGE_VERSION := $(shell $(cat) VERSION)

GOLD_FLAGS=-X github.com/cs-shadowbq/falcon_sandbox/cmd.Version=$(PACKAGE_VERSION)

SILVER_FALGS := 

# Combine SILVER_FLAGS if Variable not empty
ifdef FALCON_CLIENT_ID
SILVER_FALGS+=-X github.com/cs-shadowbq/falcon_sandbox/cmd.buildClientId=${FALCON_CLIENT_ID}
endif

ifdef FALCON_CLIENT_SECRET
SILVER_FALGS+= -X github.com/cs-shadowbq/falcon_sandbox/cmd.buildClientSecret=${FALCON_CLIENT_SECRET}
endif

ifdef FALCON_API_BASE_URL
SILVER_FALGS+= -X github.com/cs-shadowbq/falcon_sandbox/cmd.buildApiBaseUrl=${FALCON_API_BASE_URL}
endif

BUILD_DOCS := README.md LICENSE example_config.yml

OS := $(shell uname)
ifeq ($(OS), Darwin)
  certname := $(shell security find-identity -v -p codesigning | head -1 | tr -d '"' | awk '{ print $$3 }')
else
  // assign certname to UNDEFINED
  certname := undefined
endif

all: reportVersion depend clean build

reportVersion: 
	@echo "\033[32mProduct Version $(PACKAGE_VERSION)"

build:
	@echo
	@echo "\033[32mBuilding ----> \033[m"
	
	env GOOS=linux GOARCH=amd64 go build -ldflags "$(GOLD_FLAGS) ${SILVER_FALGS}" -o build/falcon_sandbox_linux_amd64 main.go
	env GOOS=windows GOARCH=amd64 go build -ldflags "$(GOLD_FLAGS) ${SILVER_FALGS}" -o build/falcon_sandbox.exe main.go
	env GOOS=darwin GOARCH=amd64 go build -ldflags "$(GOLD_FLAGS) ${SILVER_FALGS}" -o build/falcon_sandbox_darwin_amd64 main.go

clean:
	@echo
	@echo "\033[32mCleaning Build ----> \033[m"
	$(RM) -rf pkg/*
	$(RM) -rf build/*
	$(RM) -rf tmp/*

depend:
	@echo
	@echo "\033[32mChecking Build Dependencies ----> \033[m"

ifndef PACKAGE_VERSION
	@echo "\033[1;33mPACKAGE_VERSION is not set. In order to build a package I need PACKAGE_VERSION=n\033[m"
	exit 1;
endif

ifndef GOPATH
	@echo "\033[1;33mGOPATH is not set. This means that you do not have go setup properly on this machine\033[m"
	@echo "$$ mkdir ~/gocode";
	@echo "$$ echo 'export GOPATH=~/gocode' >> ~/.bash_profile";
	@echo "$$ echo 'export PATH=\"\$$GOPATH/bin:\$$PATH\"' >> ~/.bash_profile";
	@echo "$$ source ~/.bash_profile";
	exit 1;
endif

	@type go >/dev/null 2>&1|| { \
	  echo "\033[1;33mGo is required to build this application\033[m"; \
	  echo "\033[1;33mIf you are using homebrew on OSX, run\033[m"; \
	  echo "Recommend: $$ brew install go --cross-compile-all"; \
	  exit 1; \
	}

codesign:
	
	@echo
	@echo "\033[32mCodesigning ----> \033[m"
	@echo "\033[32mCodesigning with certificate: $(certname)\033[m"
	@echo "\033[32mCodesigning build/falcon_sandbox_darwin_amd64\033[m"
	@codesign -fs $(certname) build/falcon_sandbox_darwin_amd64
	@echo "\033[32mCodesigning build/falcon_sandbox_linux_amd64\033[m"
	@codesign -fs $(certname) build/falcon_sandbox_linux_amd64
	@echo "\033[32mCodesigning build/falcon_sandbox.exe\033[m"
	@codesign -fs $(certname) build/falcon_sandbox.exe
