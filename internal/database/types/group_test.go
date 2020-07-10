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

import "testing"

func TestGroup(t *testing.T) {
	testuser := NewUser("testuser", "testtest", "test@test.test")
	testuser2 := NewUser("user2", "testest", "test@test.test")
	group := NewGroup("testgroup", testuser)
	group.AddMember(testuser2)

	if !group.Match(testuser) {
		t.Fatal("owner does not match group")
	}
	if !group.Match(testuser2) {
		t.Fatalf("member does not match group")
	}

	{
		members := group.GetMembers()
		if len(members) != 1 || !members[0].Equal(testuser2) {
			t.Fatal("incorrect member list")
		}
	}
	{
		gjson, err := group.MarshalJSON()
		if err != nil {
			t.Fatalf("unable to MarshalJSON group: %s", err)
		}
		ng := new(Group)
		err = ng.UnmarshalJSON(gjson)
		if err != nil {
			t.Fatalf("unable to UnmarshalJSON group: %s", err)
		}
		if ng.GetName() != group.GetName() {
			t.Fatal("incorrect decoded group")
		}
	}
	{
		gbson, err := group.MarshalBSON()
		if err != nil {
			t.Fatalf("unable to MarshalBSON group: %s", err)
		}
		ng := new(Group)
		err = ng.UnmarshalBSON(gbson)
		if err != nil {
			t.Fatalf("unable to UnmarshalBSON group: %s", err)
		}
		if ng.GetName() != group.GetName() {
			t.Fatal("incorrect decoded group")
		}
	}
}
