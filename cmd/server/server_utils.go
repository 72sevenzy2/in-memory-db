package main

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"time"

	"github.com/72sevenzy2/in-memory-database/db"
)

// utility functions for server/main.go

// todo: add ttl existence validation to cli for setting ttls.

func Set(parts []string, conn net.Conn, b *db.DB) bool {
	if len(parts) < 4 || len(parts) > 4 {
		conn.Write(db.StringToByte("invalid SET format:\n" +
			"SET <KeyName> <value> <TTL Expiration (in minutes, eg: 5)>\n",
		))
		return false
	}

	n, err := strconv.Atoi(parts[3]) // parse ttl expiration
	if err != nil {
		conn.Write(db.StringToByte("please include a valid TTL, (eg: 10, will be in minutes)\n" +
			".\n",
		))
		return false
	}

	f, err := strconv.ParseUint(parts[2], 10, 32) // returns uint64, err.

	if err == nil { // its a int.

		// prevent f from overflowing if number entered is too big
		if f > math.MaxUint32 {
			conn.Write(db.StringToByte("please include a number value within range of unsigned int32.\n" +
				".\n",
			))
			return false
		}

		err2 := b.SetInt(parts[1], uint32(f), time.Minute*time.Duration(n))
		if err2 != nil {
			fmt.Println(err2.Error()) // print on server side
			conn.Write(db.StringToByte("error:\n" +
				err2.Error() + "\n" +
				".\n",
			))
			return false
		}
		conn.Write(db.StringToByte("successful.\n" +
			".\n"),
		)
		return true
	}

	// its a string if unable to parse to uint.
	err2 := b.SetString(parts[1], parts[2], time.Minute*time.Duration(n))
	conn.Write(db.StringToByte("successful.\n" +
		".\n"),
	)
	if err2 != nil {
		fmt.Println(err2.Error())
		conn.Write(db.StringToByte("error:\n" +
			err2.Error() + "\n" +
			".\n",
		))
		return false
	}
	return true
}

func Get(parts []string, b *db.DB, conn net.Conn) bool {
	if len(parts) < 2 || len(parts) > 2 {
		conn.Write(db.StringToByte("invalid GET format:\n"))
		conn.Write(db.StringToByte("GET <KeyName>\n"))
		return false
	}

	val, ok := b.GetInt(parts[1])
	if !ok {
		val2, ok2 := b.GetString(parts[1])
		if !ok2 {
			conn.Write(db.StringToByte("data does not exist\n"))
			return false
		}
		conn.Write(db.StringToByte(val2 + "\n"))
		return false
	}
	// fmt.Println(val)
	conn.Write(db.StringToByte(strconv.FormatUint(uint64(val), 10) + "\n")) // convert uint32 to readable format
	return true
}

func Fetch(b *db.DB, conn net.Conn) {
	vals, ok := b.GetAllInt()      // returns map[string]uint32, bool
	vals2, ok2 := b.GetAllString() // returns map[string]string, bool

	if !ok && !ok2 {
		conn.Write(db.StringToByte("no available key-values\n"))
		return
	}

	for k, v := range vals {
		fmt.Fprintln(conn, k, int(v))
	}
	for k, v := range vals2 {
		fmt.Fprintln(conn, k, v)
	}
}

func Del(parts []string, conn net.Conn, b *db.DB) bool {
	if len(parts) < 2 || len(parts) > 2 {
		conn.Write(db.StringToByte("usage: DEL <KeyName>\n"))
	}
	if _, ok := b.GetInt(parts[1]); ok {
		b.Del(parts[1])
		conn.Write(db.StringToByte("successfully deleted key\n"))
	}
	if _, ok := b.GetString(parts[1]); ok {
		b.Del(parts[1])
		conn.Write(db.StringToByte("successfuly deleted key\n"))
		return false
	}
	fmt.Fprintln(conn, "key does not exit.", parts[1])
	return false

}
