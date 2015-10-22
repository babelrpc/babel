REST
====

This library reads attributes (i.e. annotations) in Babel files that represent RESTful service metadata.

Scope
-----

All of the supported attributes must use the `@rest` scope.

Operation
---------

The `Op` attribute may be used on a method to describe the RESTful operation.

Option     | Type   | Default | Description
-----------|--------|---------|------------
Path       | string | "/"     | The URL path.
Method     | string | "GET"   | The HTTP method. (GET, PUT, POST, DELETE, OPTIONS, HEAD, or PATCH).
Deprecated | bool   | false   | Whether this method is deprecated.
Hide       | bool   | false   | Whether to hide this operation from tools like babel2swagger.

### Example

	@rest [Op(Path="/test/{a}/parms", Method="POST", Deprecated=true)]

Response
--------

Additional responses may be defined with the `Response` attribute. It is only valid on methods.

Option  | Type   | Default | Description
--------|--------|---------|------------
Code    | int    | 0       | The HTTP response code.
Type    | string | ""      | The Babel data type of the response. Defaults to the return type of the method.
Desc    | string | ""      | Description of the response
Headers | string | ""      | A comma-separated list of header names. There must be a `Header` attribute with the same name.

### Example

	@rest [Response(Code=0, Type="string", Headers="Foo,Bar")]

Header
------

The `Header` attribute defines an HTTP headers that is returned.

Option  | Type   | Default | Description
--------|--------|---------|------------
Name    | string | ""      | The HTTP header name.
Type    | string | ""      | The Babel data type of the header.
Desc    | string | ""      | Description of the header.
Format  | string | ""      | The collection format of the data (applicable to `list` types): csv, ssv, tsv, pipes, multi.

### Example

	@rest [Header(Name="Foo", Type="string", Desc="Some header"),
	       Header(Name="Bar", Type="list<string>", Format="csv" Desc="Some other header")]

Parameter
---------

The `Parm` attribute defines where to pull variables from in the REST invocation.

Option   | Type   | Default | Description
---------|--------|---------|------------
In       | string | "query" | Where the field is located: query, header, path, formdata, body.
Required | bool   | false   | Whether the field is required.
Format   | string | ""      |The collection format of the data (applicable to `list` types): csv, ssv, tsv, pipes, multi.
Name     | string | ""      | Used to rename a parameter (usually for headers)


### Example

	@rest [Parm(In="query", Required=false, Format="pipes")]
