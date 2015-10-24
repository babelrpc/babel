babel2swagger
=============

babel2swagger is a tool to generate a [Swagger 2.0](https://github.com/swagger-api/swagger-spec/blob/master/versions/2.0.md) definition from Babel files. No, this does not suddenly produce a [RESTful API](http://en.wikipedia.org/wiki/Representational_state_transfer) - what it gives you is a Swagger definition that enables [Swagger-UI](https://github.com/swagger-api/swagger-ui), [Swagger-Codegen](https://github.com/swagger-api/swagger-codegen), and other tools to seamlessly interoperate with your Babel-based service. The view in Swagger-UI also makes it very easy to understand how Babel services work over HTTP.

The tool copies over all documentation comments and generates schemas for all the types you defined. It also defines POST operations for all of the services and methods.

The tool has several options to customize the output:

Option | Default | Description
-------|---------|------------
-basepath | /              | Specifies the base path to include in the file, for example /foo/bar
-error    | false          | When -rest is enabled, still include the Babel error definition
-flat     | false          | Flatten composed objects into a single object definition
-format   | json           | Specifies output format - can be json or yaml
-host     | localhost      | Specifies the host to include in the file, for example localhost:8080
-int64    | false          | When -rest is enabled, format int64 Swagger-style instead of Babel-style
-out      |                | Specifies the file to write to
-rest     | false          | Process @rest annotations (resulting Swagger won't be able to invoke Babel services)
-title    | My Application | Sets the application title
-version  | 1.0            | Sets the application version

In pactice, you always need to use the -flat option until [Swagger-JS](https://github.com/swagger-api/swagger-js) Issue [188](https://github.com/swagger-api/swagger-js/issues/188) is resolved.

Examples
--------

Babel files are processed exactly like the babel tool processes them, so you can specify multiple Babel files on the command line to assemble them into one Swagger file.

	babel2swagger -out service.json -host babelrpc.io -basepath /service -flat -title "My Service" *.babel

Limitations
-----------

babel2swagger is mainly limited by Swagger:

* Swagger apparently doesn't support includes or imports, so the API has to be defined in a single file.
* Enumerated types cannot be defined in the `definitions` section of the Swagger file, so they have to be expanded everwhere they are used.
* Swagger doesn't handle composition/inheritance correctly, so using `-flat` is required. Swagger's `allOf` keyword isn't supported in many of its tools.
* Babel quotes `int64` to avoid data loss with JavaScript `number` types. Swagger does not, instead reducing the useful range of an `int64`. Thus, we treat `int64` as strings in Swagger.
* Swagger treats map keys as strings. Babel quotes all map keys, but allows them to be defined as other primitive types.
* babel2swagger uses HTML escaping for some areas that Swagger-UI does not handle correctly.
* Swagger-UI will not show return types that are primitives, even though they are allowed.
* Swagger-UI will not display the full types of items returned in arrays or maps.
* Swagger cannot represent initializers, which are supported by Babel.
