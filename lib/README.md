# lib

this folder contains library packages, think of a package that would make sense in more than one application that do
perform the same task as your application.

packages in lib might also be specific to your service but be used by third party as a library, e.g. API go clients

### Rules

* none of the packages will log any message error, this is the duty of the implementation
* dependency between packages should be avoided
