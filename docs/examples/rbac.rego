package rbac

# user-role assignments
# user_roles := dynamodb.policy("foo/bar", "alice")
# user_roles := {
#     "alice": ["engineering", "webdev"],
#     "bob": ["hr"]
# }

# role-permissions assignments
role_permissions := {
    "engineering": [{"action": "read",  "object": "server123"}],
    "webdev":      [{"action": "read",  "object": "server123"},
                    {"action": "write", "object": "server123"}],
    "hr":          [{"action": "read",  "object": "database456"}]
}
# lookup the list of roles for the user
policy := dynamodb.policy(input.namespace, input.principal)
# logic that implements RBAC.
default allow = false
allow {
    # for each role in that list
    r := policy.roles[_]
    # lookup the permissions list for role r
    permissions := role_permissions[r]
    # for each permission
    p := permissions[_]
    # check if the permission granted to r matches the user's request
    p == {"action": input.action, "object": input.object}
}
