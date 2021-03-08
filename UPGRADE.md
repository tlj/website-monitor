This file includes only the most important items that should
be addressed before attempting to upgrade or during the 
upgrade of a website monitor.

## FROM 0.2.x TO 0.3.0

The main change in this version is the handling of config files
which has been simplified in code. This means some shortcuts have
been removed, and some things have been moved around a bit.

### Configuration file changes

The "simplified" checks are no longer valid ("regex_expected", etc).
Instead the checks need to be added through the verbose method.

