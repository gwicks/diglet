INSTALL_PATH=$(GOPATH)bin


install: build
	@mkdir -p $(INSTALL_PATH)
	@cp digcli $(INSTALL_PATH)/diglet

build:
	@go build -o diglet ./cli