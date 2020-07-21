##############################################################
#    ____                    __                _      _ __   #
#   / __ \____  __  ______ _/ /_____________  (_)____(_) /_  #
#  / /_/ / __ \/ / / / __ `/ / ___/ ___/ __ \/ / ___/ / __/  #
# / _, _/ /_/ / /_/ / /_/ / (__  |__  ) /_/ / / /  / / /_    #
#/_/ |_|\____/\__, /\__,_/_/____/____/ .___/_/_/  /_/\__/    #
#            /____/                 /_/                      #
##############################################################

.DEFAULT_GOAL := up

GO_VERSION = "1.13"
logfile ?= $(shell pwd)/localfile.log
threshold?= "10"
PROJECT_NAME = "dd-hw"

.PHONY: up test help

base-image:
	@echo "$(BLUE)[$@]$(EOC): Build dd/base-golang (${GO_VERSION})"
	docker build -f build/Golang.base --build-arg GO_VERSION=${GO_VERSION} -t dd/base-golang .

build: base-image
	docker build -t dd-monitoring -f ./build/Dockerfile .

#help up: Run project in container.
up: build
	docker run -e threshold=$(threshold) -v ${logfile}:/app/localfile.log -it dd-monitoring --rm --name monitoring-container || true

#help test: Run unit tests.
test: base-image
	docker run --rm dd/base-golang bash -c 'go test -v ./...'

help: Makefile
	@echo 
	@echo " Choose a command run in ${PROJECT_NAME}"
	@echo
	@sed -n 's/^#help//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

