# SwaggerAPI example

## Getting Started

A high-productivity web framework for the [Go language](http://www.golang.org/).

The World's Most Popular Framework for APIs. [Swagger API](http://swagger.io)

### Start the web server:

```
go get github.com/revel/cmd/revel
go get github.com/waiteb3/revel-swagger/...
revel run github.com/waiteb3/revel-swagger/examples/swaggerapi
```

Then go to [http://localhost:9000/@api/](http://localhost:9000/@api/) and you'll see
the example SwaggerAPI's UI for this projects swagger definition.

### Custom UI

You can use Static.Serve (the offical Revel module) to serve your UI.

Example from the routes file
```
# Custom Swagger UI
GET /@custom/swagger.json                       Static.Serve("conf/swagger","swagger.json"
GET /@custom/                                   Static.Serve("public/swagger-ui/dist","index.html")
GET /@custom/*filepath                          Static.Serve("public/swagger-ui/dist")
```

Then go to [http://localhost:9000/@custom/](http://localhost:9000/@custom/) to see
the custom UI.

Thanks to https://github.com/jensoleg/swagger-ui for a modified swagger-ui to demo with.

### Description of Contents

The default directory structure of a generated Revel application:

    swaggerapi            App root
      app                 App sources
        init.go           Interceptor registration
        controllers       App controllers
          api.go          API controller-action definitions
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
        swagger-ui        Custom UI file

app

    The app directory contains the source code and templates for your application.

conf

    The conf directory contains the application’s configuration files. There are two main configuration files:

    * app.conf, the main configuration file for the application, which contains standard configuration parameters
    * routes, the routes definition file.

swagger

    The swagger subdirectory contains the application's swagger specification files. This project has one file:

	* swagger.yml, the swagger-specification definition file.

controllers

    The implementations for controller-actions that can be invoked by a pathItem's 'x-revel-controller-action' extension definition.

    * api.go, the App.Endpoint implementation for [examples/swaggerapi/conf/swagger.yml#L22](the swagger spec)

public

    Resources stored in the public directory are static assets that are served directly by the Web server. Typically it is split into three standard sub-directories for images, CSS stylesheets and JavaScript files.

    The names of these directories may be anything; the developer need only update the routes.

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
