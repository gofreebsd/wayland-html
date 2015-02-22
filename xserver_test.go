package main

import (
	"testing"
	"bufio"
	"os"
	"fmt"
	"strings"
	"strconv"
	. "gopkg.in/check.v1"
)
// Hook up gocheck into the "go test" runner.
func Test(t *testing.T) { TestingT(t) }

type XserverSuite struct{}

func readLines(path string) ([]string, error) {
  file, err := os.Open(path)
  if err != nil {
    return nil, err
  }
  defer file.Close()

  var lines []string
  scanner := bufio.NewScanner(file)
  for scanner.Scan() {
    lines = append(lines, scanner.Text())
  }
  return lines, scanner.Err()
}

func (s *XserverSuite) TestTryLock(c *C) {
	LOCK_FMT := "/tmp/.X%d-lock"
	num, lockFile := TryLock()
	fileName := fmt.Sprintf(LOCK_FMT, num)
	lines, _ := readLines(fileName)

	filePid, _ := strconv.Atoi(strings.TrimSpace(lines[0]))
	pid := os.Getpid()
	c.Assert(filePid, Equals, pid)

	defer lockFile.Unlock()
}

func (s *XserverSuite) TestXServerInit(c *C) {
	display := initDisplay()
	xserverInit(display)
}

var _ = Suite(&XserverSuite{})
