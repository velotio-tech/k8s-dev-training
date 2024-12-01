This code adds a CRD named CodeSanity. CodeSanity represents a sanity instance with a regex
to select relevant code files and a minimum coverage criteria.

This CRD also adds status subresource. Some observations regarding the status subresource.

 - `create`, `update` calls on the CodeSanity CRD will drop the status part entirely.
 - Similarly, a call on the status subresource will drop the other fields of the CodeSanity payload.

The reason for this is that status should not be affected by the changes in the CRDs itself and should
ideally be changed or updated by a controller. That's why a direct `create` call on the `/status`
is also not allowed.

There is also the `scale` subresource available to CRDs. Builtin RDs have some more subresources like `exec`.