// Copyright 2015 Reborndb Org. All Rights Reserved.
// Licensed under the MIT (MIT-LICENSE.txt) license.

package handler

import (
	"bufio"
	"bytes"
	"testing"

	"github.com/reborndb/go/redis/resp"
	. "gopkg.in/check.v1"
)

func TestT(t *testing.T) {
	TestingT(t)
}

var _ = Suite(&testRedisHandlerSuite{})

type testRedisHandlerSuite struct {
}

func (s *testRedisHandlerSuite) SetUpSuite(c *C) {
}

func (s *testRedisHandlerSuite) TearDownSuite(c *C) {
}

type testHandler struct {
	c map[string]int
}

func (h *testHandler) count(args ...[]byte) (resp.Resp, error) {
	for _, arg := range args {
		h.c[string(arg)]++
	}
	return nil, nil
}

func (h *testHandler) Get(arg0 interface{}, args ...[]byte) (resp.Resp, error) {
	return h.count(args...)
}

func (h *testHandler) Set(arg0 interface{}, args [][]byte) (resp.Resp, error) {
	return h.count(args...)
}

func (s *testRedisHandlerSuite) testmapcount(c *C, m1, m2 map[string]int) {
	c.Assert(len(m1), Equals, len(m2))

	for k, _ := range m1 {
		c.Assert(m1[k], Equals, m2[k])
	}
}

func (s *testRedisHandlerSuite) TestHandlerFunc(c *C) {
	h := &testHandler{make(map[string]int)}
	ss, err := NewServer(h)
	c.Assert(err, IsNil)

	key1, key2, key3, key4 := "key1", "key2", "key3", "key4"
	ss.t["get"](nil)
	s.testmapcount(c, h.c, map[string]int{})

	ss.t["get"](nil, []byte(key1), []byte(key2))
	s.testmapcount(c, h.c, map[string]int{key1: 1, key2: 1})

	ss.t["get"](nil, [][]byte{[]byte(key1), []byte(key3)}...)
	s.testmapcount(c, h.c, map[string]int{key1: 2, key2: 1, key3: 1})

	ss.t["set"](nil)
	s.testmapcount(c, h.c, map[string]int{key1: 2, key2: 1, key3: 1})

	ss.t["set"](nil, []byte(key1), []byte(key4))
	s.testmapcount(c, h.c, map[string]int{key1: 3, key2: 1, key3: 1, key4: 1})

	ss.t["set"](nil, [][]byte{[]byte(key1), []byte(key2), []byte(key3)}...)
	s.testmapcount(c, h.c, map[string]int{key1: 4, key2: 2, key3: 2, key4: 1})
}

func (s *testRedisHandlerSuite) TestServerServe(c *C) {
	h := &testHandler{make(map[string]int)}
	ss, err := NewServer(h)
	c.Assert(err, IsNil)

	resp, err := resp.Decode(bufio.NewReader(bytes.NewReader([]byte("*2\r\n$3\r\nset\r\n$3\r\nfoo\r\n"))))
	c.Assert(err, IsNil)

	_, err = ss.Dispatch(nil, resp)
	c.Assert(err, IsNil)

	s.testmapcount(c, h.c, map[string]int{"foo": 1})
}
