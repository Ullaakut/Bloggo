# Functional tests

This folder contains functional tests that are ran by the CI. They are not stateless, and thus will only work on a fresh DB once. Then, they won't function properly again unless the database is purged and the automatically incrementing indexes are reset.

## Purpose of the tests

Unit tests are good, but I feel like it's safer to validate that the whole API is still behaving like it should. These tests ensure that the software itself responds properly, given a set of inputs.
