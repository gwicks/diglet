INSTALL_PATH=$(GOPATH)bin


install: build
	@mkdir -p $(INSTALL_PATH)
	@cp diglet $(INSTALL_PATH)/diglet

build:
	@go build -o diglet ./cli
