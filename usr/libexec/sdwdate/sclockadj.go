// Copyright (C) 2016 - 2021 ENCRYPTED SUPPORT LP <adrelanos@whonix.org>
// Copyright (C) 2021 qua3k <78624738+qua3k@users.noreply.github.com>
// See the file COPYING for copying conditions.

package main

import (
	"log"
	"os"
	"strconv"
	"time"

	"golang.org/x/sys/unix"
)

// receive time adjustment, negative or positive, in nanoseconds
// int64 is at least -9,223,372,036,854,775,807 to +9,223,372,036,854,775,807
// exits program upon failure
func change_time_by_nanoseconds(t int64) {
	ns := unix.NsecToTimeval(time.Now().UnixNano() + t)
	// set time with new values
	unix.Settimeofday(&ns)
}

func Abs(x int64) int64 {
	if x < 0 {
		return -x
	}
	return x
}

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Too few args!")
	} else if len(os.Args) > 2 {
		log.Fatal("Too many args!")
	}
	ns_time_change, err := strconv.ParseInt(os.Args[1], 10, 0)
	if err != nil {
		log.Fatal("Failed to supply an integer!")
	}

	// since nanosecond jump is fixed, we can count the number of complete jumps.
	const full_jump int64 = 5000000
	number_of_full_jumps := Abs(ns_time_change / full_jump)
	last_jump_nanoseconds := Abs(ns_time_change % full_jump)

	if ns_time_change > 0 { // positive nanosecond change
		for i := int64(0); i < number_of_full_jumps; i++ {
			time.Sleep(1 * time.Second)           // a 1 second wait imitates ntpdate
			change_time_by_nanoseconds(full_jump) // 5,000,000 ns imitates ntpdate
		}
		time.Sleep(1 * time.Second)
		change_time_by_nanoseconds(last_jump_nanoseconds)
	} else { // negative nanosecond change
		for i := int64(0); i < number_of_full_jumps; i++ {
			time.Sleep(1 * time.Second)
			change_time_by_nanoseconds(-full_jump)
		}
		time.Sleep(1 * time.Second)
		change_time_by_nanoseconds(-last_jump_nanoseconds) // negative of absolute value imitates Euclidean modulo
	}
	os.Exit(0)
}
