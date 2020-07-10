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
	"net/http/httptest"
	"testing"
	"time"
)

func TestUser(t *testing.T) {
	user := NewUser("testuser", "testtest", "test@test.test")

	if !user.Match(user) {
		t.Fatal("basic equality check failed")
	}

	if !user.GetLock().Valid(map[string]interface{}{
		"pass": "testtest",
	}) {
		t.Fatal("failed to unlock")
	}

	cookies := user.NewCookies(time.Now().Add(12*time.Hour), time.Now().Add(24*time.Hour))

	if user.GetID().String() != cookies[1].Value {
		t.Fatalf("incorrect cookie value")
	}

	testrequest := httptest.NewRequest("GET", "/test/test", &bytes.Buffer{})

	for _, c := range cookies {
		testrequest.AddCookie(c)
	}

	if !user.CheckCookie(testrequest) {
		t.Fatalf("Failed to validate cookies")
	}

	cookieOID, err := GetCookieUID(testrequest)
	if err != nil {
		t.Fatalf("unable to get oid: %s", err)
	}
	if !cookieOID.Equal(user.GetID()) {
		t.Fatalf("mismatched cookie id: %v", cookieOID)
	}
}
