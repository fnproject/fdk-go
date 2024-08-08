RAW_VERSION=$(shell grep -m1 -Eo "[0-9]+\.[0-9]+\.[0-9]+" version.go)
FDK_VERSION=$(RAW_VERSION).$(BLD_NUMBER)$(BLD_BRANCH_SUFFIX)

.PHONY: set_fdk_version set_image_version setup_test_image_build

set_fdk_version:
	echo "version: \"$(FDK_VERSION)\"" > $(BLD_INPUT_DIR)/buildOverrideVariables.conf
	echo "fdk_version=$(FDK_VERSION)" >> $(BLD_INPUT_DIR)/exportVariables.env

set_image_version:
	echo "version: \"$(GO_VERSION)-$(FDK_VERSION)$(if $(VERSION_SUFFIX),-$(VERSION_SUFFIX),)\"" > $(BLD_INPUT_DIR)/buildOverrideVariables.conf

setup_test_image_build:
    # Build tag for the test image
	echo "version: \"go-$(GO_VERSION)-$(FDK_VERSION)\"" >$(BLD_INPUT_DIR)/buildOverrideVariables.conf
	echo "fdkVersion: \"$(FDK_VERSION)\"" >>$(BLD_INPUT_DIR)/buildOverrideVariables.conf
	echo "buildImageVersion: \"$(GO_VERSION)-$(FDK_VERSION)-dev\"" >>$(BLD_INPUT_DIR)/buildOverrideVariables.conf
	echo "runtimeImageVersion: \"$(GO_VERSION)-$(FDK_VERSION)\"" >>$(BLD_INPUT_DIR)/buildOverrideVariables.conf
	cat $(BLD_INPUT_DIR)/buildOverrideVariables.conf
	find $(TEST_DIR) -type d -mindepth 1 -maxdepth 1 -exec "$(BLD_INPUT_DIR)/internal/build-scripts/copy_source_pkg.sh" "{}" \;