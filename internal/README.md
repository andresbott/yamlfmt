# internal

This folder contains packages that support the main business logic of the application,
e.g. a package that writes to a database.

### Rules
Rules applicable for packages in this folder:
* don't instantiate shared runtime dependencies, e.g. a logger or a db connection
* they are only instantiated by packages in _apps_ and with few exceptions from _internal_
