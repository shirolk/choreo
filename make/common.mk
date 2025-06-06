# This makefile contains various helper functions and variables used across other makefiles.

# Setting SHELL to bash allows bash commands to be executed by recipes.
# Options are set to exit when a recipe line exits non-zero or a piped command fails.
SHELL = /usr/bin/env bash -o pipefail
.SHELLFLAGS = -ec

PROJECT_BIN_DIR := $(PROJECT_DIR)/bin

# Read the version from the VERSION file
RELEASE_VERSION ?= $(shell cat VERSION)

# Store the short git sha of latest commit to be used in the make targets
GIT_REV := $(shell git rev-parse --short HEAD)

# Helper functions for logging
define log_info
echo -e "\033[36m--->$1\033[0m"
endef

define log_error
echo -e "\033[0;31m--->$1\033[0m"
endef

# Helper functions to get the OS and architecture from the platform string
# Format: <os>/<arch>
get_platform_os = $(word 1, $(subst /, ,$(1)))
get_platform_arch = $(word 2, $(subst /, ,$(1)))

##@ General

# The help target prints out all targets with their descriptions organized
# beneath their categories. The categories are represented by '##@' and the
# target descriptions by '##'. The awk command is responsible for reading the
# entire set of makefiles included in this invocation, looking for lines of the
# file as xyz: ## something, and then pretty-format the target and help. Then,
# if there's a line with ##@ something, that gets pretty-printed as a category.
# More info on the usage of ANSI control characters for terminal formatting:
# https://en.wikipedia.org/wiki/ANSI_escape_code#SGR_parameters
# More info on the awk command:
# http://linuxcommand.org/lc3_adv_awk.php

.PHONY: help
help: ## Display this help.
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9.%-]+:.*?##/ { printf "  \033[36m%-24s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
