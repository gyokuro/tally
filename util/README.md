# Utilities

## embedfs 
Generates go source code based on the binary file content and allows the
files to be compiled and embedded into the same executable.

To run from the project directory:

    util/embedfs -match=".*" -generate=true -destDir=resources webapp

This generates `.go` files in the `resources` directory, using file content from
the `webapp` directory.
