/*************************************************************************
 *
 * MAXSET CONFIDENTIAL
 * __________________
 *
 *  [2019] - [2020] Maxset WorldWide Inc.
 *  All Rights Reserved.
 *
 * NOTICE:  All information contained herein is, and remains
 * the property of Maxset WorldWide Inc. and its suppliers,
 * if any.  The intellectual and technical concepts contained
 * herein are proprietary to Maxset WorldWide Inc.
 * and its suppliers and may be covered by U.S. and Foreign Patents,
 * patents in process, and are protected by trade secret or copyright law.
 * Dissemination of this information or reproduction of this material
 * is strictly forbidden unless prior written permission is obtained
 * from Maxset WorldWide Inc.
 */

package brand

import (
	"math/rand"
	"sync"
	"time"
)

// ByteGen generates each possible byte in a random order
type ByteGen struct {
	r   *rand.Rand
	buf []byte
	i   int
	m   sync.Mutex
}

func (b *ByteGen) shuffle() {
	b.r.Shuffle(len(b.buf), b.getSwap())
	b.i = 0
}

func (b *ByteGen) getSwap() func(i, j int) {
	return func(i, j int) {
		b.buf[i], b.buf[j] = b.buf[j], b.buf[i]
	}
}

var bg = New(time.Now().Unix())

// New returns a new ByteGen using the given seed
func New(seed int64) *ByteGen {
	n := new(ByteGen)

	n.r = rand.New(rand.NewSource(seed))

	n.buf = make([]byte, 256)
	for i := range n.buf {
		n.buf[i] = byte(i)
	}

	n.shuffle()
	return n
}

// Next returns the next byte
func (b *ByteGen) Next() byte {
	b.m.Lock()
	defer b.m.Unlock()
	defer func() {
		b.i++
		if b.i >= len(bg.buf) {
			b.shuffle()
		}
	}()
	return b.buf[bg.i]
}

// Next returns a random byte without repeating until every possible byte has
// been returned
func Next() byte {
	return bg.Next()
}
