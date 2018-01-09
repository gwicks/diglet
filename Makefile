INSTALL_PATH=$(GOPATH)bin


install: build
	@mkdir -p $(INSTALL_PATH)
	@cp lincoln $(INSTALL_PATH)/lincoln

build:
	@go build