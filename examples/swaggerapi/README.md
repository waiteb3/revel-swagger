# SwaggerAPI example

## Getting Started

A high-productivity web framework for the [Go language](http://www.golang.org/).

The World's Most Popular Framework for APIs. [Swagger](http://swagger.io)

For all of the module's more complete details and instructions, go to the [module's README.md](../../modules/swaggerapi/README.md).

### Start the web server:

```
go get github.com/revel/cmd/revel
go get github.com/waiteb3/revel-swagger/examples/swaggerapi/...
go get github.com/go-swagger/go-swagger/spec
revel run github.com/waiteb3/revel-swagger/examples/swaggerapi
```

Then go to [http://localhost:9000/@api/](http://localhost:9000/@api/) and you'll see
the example SwaggerAPI's UI for this projects swagger definition.

Or you can go to [http://localhost:9000/@custom/](http://localhost:9000/@custom/) to see an example of a custom Swagger UI.
You can read also [instructions](../../modules/swaggerapi/README.md#custom-ui) on how to add your own custom Swagger UI distribution.

Thanks to [https://github.com/jensoleg/swagger-ui](https://github.com/jensoleg/swagger-ui) for a modified swagger-ui to demo with.

**Note**: If you get a compilation error from a package not being found,
try refreshing once or twice since the revel tool may not have cloned the package yet.

### Description of Contents

The default directory structure of a generated Revel application:

    swaggerapi            App root
      app                 App sources
        init.go           Interceptor registration
        controllers       App controllers
          app.go          Default API controller-action definitions
          items.go        Items API controller-action definitions
        models            App domain models
        routes            Reverse routes (generated code)
        views             Templates
      tests               Test suites
      conf                Configuration files
        app.conf          Main configuration file
        routes            Routes definition
        swagger           Swagger files
          swagger.yml     Swagger Spec file
      messages            Message files
      public              Public assets
        css               CSS files
        js                Javascript files
        images            Image files
        swagger-ui        Custom UI files

app

    The app directory contains the source code and templates for your application.

conf

    The conf directory contains the application’s configuration files.
    There are two main configuration files:

    * app.conf, the main configuration file for the application, which contains standard configuration parameters
    * routes, the routes definition file.

conf/swagger

    The swagger subdirectory contains the application's swagger specification files.
    This project has one file:

	* swagger.yml, the swagger-specification definition file.

app/controllers
* [App.Endpoint extension declaration](conf/swagger/swagger.yml#L19)
* [Items.{Action} example extension declaration](conf/swagger/swagger.yml#L29)

```
    The implementations for controller and their actions that are invoked by
    a pathItem's 'x-revel-controller-action'.

    * api.go, the App.Endpoint implementation
    * items.go, the Items.{Action} implementations
```

public

    Resources stored in the public directory are static assets that are served
    directly by the Web server. Typically it is split into three standard sub-directories
    for images, CSS stylesheets and JavaScript files.

    The names of these directories may be anything; the developer need only update the routes.

public/swagger-ui

   Swagger UI assets are served from this directory.

test

    Tests are kept in the tests directory. Revel provides a testing framework that makes it easy to write and run functional tests against your application.

### Follow the guidelines to start developing your application:

* The README file created within your application.
* The [Getting Started with Revel](http://revel.github.io/tutorial/index.html).
* The [Revel guides](http://revel.github.io/manual/index.html).
* The [Revel sample apps](http://revel.github.io/samples/index.html).
* The [API documentation](http://revel.github.io/docs/godoc/index.html).

## Contributing
We encourage you to contribute to Revel! Please check out the [Contributing to Revel
guide](https://github.com/revel/revel/blob/master/CONTRIBUTING.md) for guidelines about how
to proceed. [Join us](https://groups.google.com/forum/#!forum/revel-framework)!
