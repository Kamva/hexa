`gate` package implements  `hexa.Gate` interface.

Default permission check contains:
- Deny guest.
- Deny every user with without `activated_account` permission.
- Allows every user with `root` permission to do anything.
- check user has permission by following expression:
```
(user has specified manager paermission) || (user has specifeid permission && policy returns true)
```

