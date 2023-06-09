#!/bin/bash

function assert() {
  SCRIPT=$1
  EXPECT_EXIT_CODE=$2
  EXPECT_OUT=$3

  OUT=`$SCRIPT 2>/dev/null`
  EXIT_CODE=$?
  if [[ $EXIT_CODE -ne $EXPECT_EXIT_CODE ]]; then
    echo "run '$SCRIPT'. exit status want '$EXPECT_EXIT_CODE', but got '$EXIT_CODE'"
  fi
  if [[ $OUT != "$EXPECT_OUT" ]]; then
    echo "run '$SCRIPT'. want '$EXPECT_OUT', but got '$OUT'"
  fi
}

# parameter error
assert "./xlgrep" 2 ""
assert "./xlgrep hoge" 2 ""

# no support file
assert "./xlgrep hoge main.go" 1 ""

# not match
assert "./xlgrep foobar test/testdata.xlsx" 1 ""

# single match
assert "./xlgrep xxx test/testdata.xlsx" 0 "test/testdata.xlsx:Sheet1:B9:xxx"

# multiple match in single file
assert "./xlgrep xx test/testdata.xlsx" 0 "test/testdata.xlsx:Sheet1:B9:xxx
test/testdata.xlsx:Sheet1:C15:xx"

# multiple match in multiple file
assert "./xlgrep test test/testdata.xlsx test/testdata2.xlsx" 0 "test/testdata.xlsx:Sheet1:B2:test
test/testdata.xlsx:Sheet1:C3:test
test/testdata2.xlsx:Sheet1:B2:test
test/testdata2.xlsx:Sheet1:C3:test"

# regex match
assert "./xlgrep e[s|p] test/testdata.xlsx" 0 "test/testdata.xlsx:Sheet1:A1:grep
test/testdata.xlsx:Sheet1:B2:test
test/testdata.xlsx:Sheet1:C3:test"
