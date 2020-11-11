package example.authz

default allowed = false

allowed {
    input.user.preferred_username == "user"
}
