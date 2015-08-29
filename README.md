# webwatch
Sample Go app that watches web resources for HTTP status changes.

This will watch a remote resource at a specified URL, printing the status of it as it polls the resource.

## Config

This app uses [viper](https://godoc.org/github.com/spf13/viper) for the config. The name and URLs are expected to be set within the config.

Sample Config (yaml) - `test.yml`:

```yaml
  name: Test Name
  urls:
    - http://localhost:8000/
    - http://localhost:8001/
```

### Purpose

Created to try out Golang.
