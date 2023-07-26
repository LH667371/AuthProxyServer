.PHONY: help build

NAME:=auth-proxy
IMAGE:=service.golang.auth.proxy
VERSION:=$(shell cat VERSION)

# Terminal coloring
BOLD:=$$(tput bold)
GREEN:=$$(tput setaf 2)
RED:=$$(tput setaf 1)
CLEAR:=$$(tput sgr 0)

## Global commands
help:
		@echo "$(BOLD)Name:$(NAME)"
		@echo
		@echo "$(BOLD)Global commands:$(CLEAR)"
		@echo "  $(BOLD)$(GREEN)help$(CLEAR)    Display this menu."
		@echo
		@echo "$(BOLD)Workstation commands:$(CLEAR)"
		@echo "  $(BOLD)$(GREEN)build$(CLEAR)      Build all tool images."
		@echo

build:
		@echo "$(BOLD)$(GREEN)Building $(IMAGE):$(VERSION) image...$(CLEAR)"
		@docker build -t $(IMAGE):$(VERSION) .
		@echo "$(BOLD)$(GREEN)Image $(IMAGE) was built.$(CLEAR)"

run:
		@echo "$(BOLD)$(GREEN)Runing $(IMAGE):$(VERSION) image...$(CLEAR)"
		@docker run -tid --name=$(NAME) -p 8080:80 $(IMAGE):$(VERSION)

stop:
		@docker rm -f $(NAME)

push:
		@docker push $(IMAGE):$(VERSION)
