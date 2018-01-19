package rpp

import (
	"testing"

	"github.com/garyburd/redigo/redis"
	"github.com/rafaeljusto/redigomock"
)

func TestRPP(t *testing.T) {
	testConn := redigomock.NewConn()

	// Register our commands
	authCmd := testConn.Command("AUTH", "foo").ExpectStringSlice("ok")
	selectCmd := testConn.Command("SELECT", 5).ExpectStringSlice("ok")

	dialFn = func(_ string) (redis.Conn, error) {
		return testConn, nil
	}

	var (
		err error
		rpp *redis.Pool
	)
	if rpp, err = RPP("redis://x:foo@10.10.10.10:8443/5", 5, 5); err != nil {
		t.Error("failed to create rpp: ", err)
		t.Fail()
	}

	_ = rpp.Get()

	if testConn.Stats(authCmd) != 1 {
		t.Error("connection was not successfully authenticated")
	}

	if testConn.Stats(selectCmd) != 1 {
		t.Error("connection did not successfully select desired DB")
	}
}
