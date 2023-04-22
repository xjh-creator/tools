package xdb

import (
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/log"
	"github.com/gogf/gf/test/gtest"
)

func TestNewAdapter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		a, err := NewAdapterByDB(Default().DB)
		t.AssertEQ(err, nil)
		t.AssertNE(a, nil)
		e, err := casbin.NewEnforcer("rbac_model.conf", a)
		t.AssertEQ(err, nil)
		t.AssertNE(e, nil)
		// Load the policy from DB.
		err = e.LoadPolicy()
		t.AssertEQ(err, nil)
		// Check the permission.
		ok, err := e.Enforce("alice", "data1", "read")
		t.Log(ok)
		t.AssertEQ(err, nil)
		e.AddRoleForUser("alice", "data1_admin")
		// Modify the policy.
		res, err := e.GetRolesForUser("alice")
		t.Log(res)
		res, err = e.GetUsersForRole("data1_admin")
		t.Log(res)
		rm := e.GetRoleManager()
		logger := log.DefaultLogger{}
		logger.EnableLog(true)
		rm.SetLogger(&logger)
		//rm.AddLink("a1","a1","domain")
		rm.PrintRoles()

		// e.RemovePolicy(...)

		// Save the policy back to DB.
		err = e.SavePolicy()
		t.AssertEQ(err, nil)
	})
}
