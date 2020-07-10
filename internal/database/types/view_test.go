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

package types

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestViewStore(t *testing.T) {
	contentString := "This is the view content! It's like the file store content, but it should be a PDF version of the file."
	inBytes := []byte(contentString)

	mockStoreID := StoreID{
		Hash:  12345,
		Stamp: 6789,
	}
	vs, err := NewViewStore(mockStoreID, bytes.NewReader(inBytes))
	if err != nil {
		t.Fatalf("error creating viewstore: %s", err)
	}

	rdr, err := vs.Reader()
	if err != nil {
		t.Fatalf("unable to create reader from viewstore: %s", err)
	}

	sb := new(strings.Builder)
	if _, err := io.Copy(sb, rdr); err != nil {
		t.Fatalf("unable to copy from reader: %s", err)
	}

	if s := sb.String(); s != contentString {
		t.Fatalf("incorrect read string: expected '%s', got '%s'", contentString, s)
	}
}
