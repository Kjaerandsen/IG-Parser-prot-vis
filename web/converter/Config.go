package converter

import "IG-Parser/exporter"

/*
Contains the default configuration for tabular output.
 */
func SetDefaultConfig() {
	exporter.INCLUDE_SHARED_ELEMENTS_IN_TABULAR_OUTPUT = true
}