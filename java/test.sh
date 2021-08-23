#!/bin/bash

go generate ../main.go
javac com/cos/*.java
javac com/ds/*.java
javac Test.java
LD_LIBRARY_PATH=. java Test
