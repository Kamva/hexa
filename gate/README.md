`gate` package implements  `hexa.Gate` interface.

Default permission check in order contains:
- Deny guest.
- Deny every user without `activated_account` permission.
- Allows every user with `root` permission to do anything.
- check user has permission by following expression:
```
(user has specified manager paermission) || (user has specifeid user permission && policy returns true)
```

